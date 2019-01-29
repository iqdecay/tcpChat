package main

import (
	"bufio"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"net"
)

func navigator(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
	maxX, maxY := g.Size()
	// Panel presenting connected users on the left
	if v, err := g.SetView("users", 0, 0, maxX/4, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Users"
		v.Write([]byte("victor"))
	}
	// The chat history up to the connection of the user
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
		v.Editor = gocui.EditorFunc(navigator)
		v.Editable = true
		v.Wrap = true

	}

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

func initKeyBindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("chat", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.Autoscroll = true
			g.SetCurrentView("input")
			g.Update(update)
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.SetCurrentView("chat")
			v, _ = g.View("chat")
			v.Autoscroll = false
			g.Update(update)
			return nil
		}); err != nil {
		return err
	}

	return nil
}
func receiveMessage(conn net.Conn, g *gocui.Gui) {
	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message := []byte(incoming)
		displayMessage(g, "chat", message)
		g.Update(func(g *gocui.Gui) error {
			return nil
		})
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
	reader := new(bufio.Reader)
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
