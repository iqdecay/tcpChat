package main

import (
	"bufio"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func trial(conn net.Conn, h *tui.Box) {
	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		fmt.Printf("Message received : %s", incoming)
		if err != nil {
			fmt.Println("There is an error")
		}
		h.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
			tui.NewLabel(incoming),
			tui.NewSpacer(),
		))

	}

}

func main() {
	//Initialize connection to chat server
	conn, err := net.Dial("tcp", ":6060")
	if err != nil {
		err = fmt.Errorf("error connecting to server : #{err}")
		fmt.Println(err)
	}
	defer conn.Close()

	userList := tui.NewList()
	sidebar := tui.NewVBox(
		tui.NewLabel("Users"),
		userList)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	input.OnSubmit(func(e *tui.Entry) {
		conn.Write([]byte(e.Text() + "\n"))
		input.SetText("")
	})

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)
	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	go trial(conn, history)
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
	go mustCopy(os.Stdout, conn)
	mustCopy(conn, os.Stdin)
}
