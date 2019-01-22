package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	for {
		// Read input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send : \n")
		text, _ := reader.ReadString('\n')
		// Send to socket
		fmt.Fprintf(conn, text+"\n")
		// Listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Message from server : ", message)
	}

}
