package main

import (
	"bufio"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"net"
	"strings"
)

func navigator(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// Rules for navigating the chat history
	switch {
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}

func layout(g *gocui.Gui) error {
	// Generate the UI and its rules
	maxX, maxY := g.Size()
	// Panel presenting connected users on the left
	if v, err := g.SetView("users", 0, 0, maxX/4, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Users"
	}
	// The chat history
	if v, err := g.SetView("chat", maxX/4+1, 0, maxX-1, 3*maxY/4-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("chat"); err != nil {
			return err
		}
		v.Title = "Chat"
		v.Autoscroll = true
		v.Overwrite = false
		// We use the navigator as an editor, but the text displayed will not change
		v.Editor = gocui.EditorFunc(navigator)
		v.Editable = true
		v.Wrap = true

	}

	// The input box
	if v, err := g.SetView("input", maxX/4+1, 3*maxY/4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
		v.Editable = true
		v.Wrap = true
	}
	// No error occured during any initialization
	return nil

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func displayMessage(g *gocui.Gui, viewName string, message []byte) error {
	originalView := g.CurrentView()
	g.SetCurrentView(viewName)
	v := g.CurrentView()
	v.Write(message)
	g.SetCurrentView(originalView.Name())
	return nil
}

func extractUserList(s string) []string {
	// From a string containing comma-separated alphanumerical pseudos,
	// return a list of said pseudos
	var users []string
	var user strings.Builder
	s = s[2 : len(s)-1] // remove the leading ## and trailing newline
	for _, char := range s {
		if char != rune(',') {
			user.WriteRune(char)
		} else {
			users = append(users, user.String())
			user.Reset()
		}
	}
	return users
}

func displayUserList(g *gocui.Gui, userList string) error {
	originalView := g.CurrentView()
	users := extractUserList(userList)
	g.SetCurrentView("users")
	v := g.CurrentView()
	for _, pseudo := range users {
		v.Write([]byte(pseudo + "\n"))
	}
	g.SetCurrentView(originalView.Name())
	return nil
}

func initKeyBindings(g *gocui.Gui) error {
	// Initialize the keybindings that depend only on the gui

	// When in chat history, ENTER puts you back in the input box,
	// where you left of
	if err := g.SetKeybinding("chat", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.Autoscroll = true       // Show the latest messages
			g.SetCurrentView("input") // Get back to input box
			g.Update(update)
			return nil
		}); err != nil {
		return err
	}

	// Ctrl-C leaves the application, whatever the focused view
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// When in input box, ArrowUp puts you in chat history so you can navigate
	// all the history
	if err := g.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.SetCurrentView("chat")
			v, _ = g.View("chat")
			v.Autoscroll = false // Earlier messages can be displayed
			g.Update(update)
			return nil
		}); err != nil {
		return err
	}

	return nil
}

var c int
func receiveMessage(conn net.Conn, g *gocui.Gui) {
	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// if the message is the userList
		//if incoming[0:2] == "##" {
		//	fmt.Println(c)
		//	displayUserList(g, incoming)
		//	g.Update(update)
		//	c++
		//	return
		//} else {
			message := []byte(incoming)
			displayMessage(g, "chat", message)
			g.Update(update)
		//}
	}

}

func update(g *gocui.Gui) error {
	return nil
}
func main() {
	//Initialize connection to chat server
	conn, err := net.Dial("tcp", ":6060")
	if err != nil {
		err = fmt.Errorf("error connecting to server : #{err}")
		fmt.Println(err)
	}
	defer conn.Close()
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = false
	g.SetManagerFunc(layout)
	initKeyBindings(g)
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			b := make([]byte, len(v.Buffer()))
			n, err := v.Read(b)
			if n != len(v.Buffer()) {
				return fmt.Errorf("mismatch between buffer length and bytes read")
			}
			if err != nil && err != io.EOF {
				return err
			}
			if n != 0 {
				conn.Write(b)
			}
			v = g.CurrentView()
			v.Clear()
			x0, y0 := v.Origin()
			x, y := v.Cursor()
			dx, dy := x0-x, y0-y
			v.MoveCursor(dx, dy, true)
			return nil

		}); err != nil {
		panic(err)
	}

	go receiveMessage(conn, g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}
