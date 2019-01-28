package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	g.Cursor = true
	g.SetManagerFunc(layout)

	if err := initKeyBindings(g); err != nil {
		log.Fatalln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

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

	if v, err := g.SetView("chat", maxX/4+1, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("chat"); err != nil {
			return err
		}
		v.Title = "Chat"
		v.Editable = true
		v.Wrap = true

	}
	return nil

}
func view(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func initKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, view); err != nil {
		return err
	}
	return nil



}
