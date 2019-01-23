package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	clientCount := 0
	// int serves as unique id
	allClients := make(map[net.Conn]int)
	// TCP will push new connections to it
	newConnections := make(chan net.Conn)
	// We will remove those clients from allClients
	deadConnections := make(chan net.Conn)
	// Channel which will contain message from connected clients
	messages := make(chan string)
	// Start TCP server
	server, err := net.Listen("tcp", ":6060")
	if err != nil {
		err = fmt.Errorf("Error launching the server :", err)
		fmt.Println(err)
	}

	// Server accepts connections forever and pushes new ones
	// to the dedicated channel
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println(err)
			}
			newConnections <- conn
		}

	}()

	// Infinite loop
	for {
		select {

		// Continuously accept new clients
		case conn := <-newConnections:
			log.Printf("Accepted new client with id %d", clientCount)
			allClients[conn] = clientCount
			clientCount++
			// Read all incoming messages from this client into a goroutine
			// and push them to the message chan
			go func(conn net.Conn) {
				clientId := allClients[conn]
				reader := bufio.NewReader(conn)
				for {
					incoming, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					messages <- fmt.Sprintf("Client %d > %s", clientId, incoming)
				}
				// When encounter an error, the client will be removed
				deadConnections <- conn

			}(conn)


		// Continuously read incoming messages and broadcast them
		case message := <-messages:
			for conn := range allClients {
				//Send the message in a go routine
				go func(conn net.Conn, message string) {
					_, err := conn.Write([]byte(message))
					// If it doesn't work the connection is dead
					if err != nil {
						deadConnections <- conn
					}
				}(conn, message)
			}
			log.Printf("New message : %s", message)

		//Remove dead clients
		case conn := <-deadConnections:
			log.Printf("Client %d disconnected", allClients[conn])
			delete(allClients, conn)
		}

	}

}
