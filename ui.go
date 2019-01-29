package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"time"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("users", 0, 0, maxX/4, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Users"
		fmt.Println("")
		fmt.Println(" victor")
		fmt.Println(" jean")
		fmt.Println(" mano")
	}

	if v, err := g.SetView("chat", maxX/4+1, 0, maxX-1, 3*maxY/4-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("chat"); err != nil {
			return err
		}
		v.Title = "Chat"
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
	message = append([]byte(time.Now().Format("15:04")+">"), message...)
	originalView := g.CurrentView()
	g.SetCurrentView(viewName)
	v := g.CurrentView()
	v.Write(message)
	g.SetCurrentView(originalView.Name())
	return nil
}

func initKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	g.Cursor = true
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
			displayMessage(g, "chat", b)
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

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}

}
