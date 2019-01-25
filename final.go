package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"strings"
)

type Client struct {
	name string
	id   int
}

type Server struct {
	allClients map[net.Conn]Client
	allPseudo  []string
	totalClients int
	connectedClients int
	listener net.Listener
}

var validPseudo = regexp.MustCompile(`([A-Z]|[a-z]|[0-9]){4,12}`)

func getValidPseudo(conn net.Conn) string {
	conn.Write([]byte("\n Please enter a new pseudo : "))
	reader := bufio.NewReader(conn)
	pseudo, _ := reader.ReadString('\n')
	pseudo = strings.Trim(pseudo, "\n")
	for !validPseudo.MatchString(pseudo) {
		conn.Write([]byte("Pseudo are alphanumerical and of length in [4,12]"))
		conn.Write([]byte("Please enter a new pseudo : "))
		pseudo, _ := reader.ReadString('\n')
		pseudo = strings.Trim(pseudo,"\n")
	}
	return pseudo
}

func disconnect(conn net.Conn) {
	client := allClients[conn]
	pseudo := client.name
	i := find(allPseudo, pseudo)
	allPseudo = append(allPseudo[:i],allPseudo[i+1:]...)
	conn.Close()
	delete(allClients, conn)
	log.Printf("Client with pseudo %s disconnected", pseudo)
}

func find(s interface{}, elem interface{}) int {
	// Return -1 if elem is in s, its index in s otherwise
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == elem {
				return i
			}
		}
	}
	return -1

}
func contains(s interface{}, elem interface{}) bool {
	// Return true if elem is in s, false otherwise
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {

			// panics if slice element points to an unexported struct field
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

func main() {
	server := new(Server)
	// TCP will push new connections to it
	newConnections := make(chan net.Conn)
	// We will remove those clients from allClients
	deadConnections := make(chan net.Conn)
	// Channel which will contain message from connected clients
	messages := make(chan string)
	// Start TCP server
	var err error
	server.listener, err = net.Listen("tcp", ":6060")
	if err != nil {
		err = fmt.Errorf("error launching the server : %e", err)
		fmt.Println(err)
	}

	// Server accepts connections forever and pushes new ones
	// to the dedicated channel
	go func() {
		for {
			conn, err := server.listener.Accept()
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
			clientCount++
			// Read all incoming messages from this client into a goroutine
			// and push them to the message chan
			go func(conn net.Conn) {
				reader := bufio.NewReader(conn)
				conn.Write([]byte("Welcome to the server ! \n"))
				pseudo := getValidPseudo(conn)
				for contains(allPseudo, pseudo) {
					conn.Write([]byte("Pseudo already in use, please choose a new one"))
					pseudo = getValidPseudo(conn)
				}
				conn.Write([]byte(fmt.Sprintf("Your pseudo is now %s \n", pseudo)))
				client := Client{pseudo, clientCount}
				allClients[conn] = client
				allPseudo = append(allPseudo, pseudo)
				reader = bufio.NewReader(conn)
				for {
					incoming, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					messages <- fmt.Sprintf("%s > %s", client.name, incoming)
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
			disconnect(conn)
		}
	}
}
