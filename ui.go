package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
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
			b := make([]byte, 100)
			_, err := v.Read(b)
			fmt.Println("\n", string(b))
			if err != nil && err != io.EOF {
				return err
			}
			g.SetCurrentView("chat")
			v = g.CurrentView()
			v.Write(b)
			g.SetCurrentView("input")
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
