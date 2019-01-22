package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Launching server")

	// Listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081909")
	// Accept connections on port
	conn, _ := ln.Accept()
	// Run loop forever
	for {
		// Listen for a message ending in \n
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Message received :", string(message))
		// Example process of string
		newMessage := strings.ToUpper(message)
		// Send back new message to client
		conn.Write([]byte(newMessage + "\n"))

	}

}
