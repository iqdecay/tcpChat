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
	allClients       map[net.Conn]Client
	allPseudo        []string
	totalClients     int
	connectedClients int
	listener         net.Listener
}

var validPseudo = regexp.MustCompile(`([A-Z]|[a-z]|[0-9]){4,12}`)

func getValidPseudo(conn net.Conn) string {
	// Get pseudo from client
	conn.Write([]byte("Please enter a new pseudo : \n"))
	reader := bufio.NewReader(conn)
	pseudo, _ := reader.ReadString('\n')
	pseudo = strings.Trim(pseudo, "\n")

	// Until it has correct format
	for !validPseudo.MatchString(pseudo) {
		conn.Write([]byte("Pseudo are alphanumerical and of length in [4,12]\n"))
		conn.Write([]byte("Please enter a new pseudo : \n"))
		pseudo, _ := reader.ReadString('\n')
		pseudo = strings.Trim(pseudo, "\n")
	}
	return pseudo
}

func disconnect(conn net.Conn, server *Server) {
	// Properly close the connection and delete the client from the server
	client := server.allClients[conn]
	pseudo := client.name
	i := find(server.allPseudo, pseudo)
	server.allPseudo = append(server.allPseudo[:i], server.allPseudo[i+1:]...)
	server.connectedClients--
	conn.Close()
	delete(server.allClients, conn)
	log.Printf("Client with pseudo %s disconnected", pseudo)
	log.Printf("There are %d connected clients", server.connectedClients)
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
	// Initialization

	server := new(Server)
	server.allClients = make(map[net.Conn]Client)
	// Server will push new connections to it
	newConnections := make(chan net.Conn)
	// Clients that will be remove from allClients
	deadConnections := make(chan net.Conn)
	// Receives messages from connected clients
	messages := make(chan string)

	// Start TCP server

	var err error
	server.listener, err = net.Listen("tcp", ":6060")
	if err != nil {
		err = fmt.Errorf("error launching the server : %e", err)
		fmt.Println(err)
	}

	// Server accepts connections forever and pushes new ones to the channel
	go func() {
		for {
			conn, err := server.listener.Accept()
			if err != nil {
				fmt.Println(err)
			}
			newConnections <- conn
		}
	}()

	for {
		select {

		// Continuously accept new clients
		case conn := <-newConnections:
			server.totalClients++
			log.Printf("Accepted new client with id %d", server.totalClients)

			// Read all incoming messages from this client and push them to the chan
			go func(conn net.Conn, server *Server) {
				conn.Write([]byte("Welcome to the server ! \n"))

				// Get a pseudo in valid format
				pseudo := getValidPseudo(conn)
				// Get a pseudo not used
				for contains(server.allPseudo, pseudo) {
					conn.Write([]byte("Pseudo already in use, please choose a new one"))
					pseudo = getValidPseudo(conn)
				}

				messages <- fmt.Sprintf("User %s joined the room !\n", pseudo)
				conn.Write([]byte(fmt.Sprintf("Your pseudo is now %s \n", pseudo)))

				// Add client to server
				client := Client{pseudo, server.totalClients}
				server.allClients[conn] = client
				server.allPseudo = append(server.allPseudo, pseudo)
				server.connectedClients++
				reader := bufio.NewReader(conn)

				// Read all his incoming messages
				for {
					incoming, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					messages <- fmt.Sprintf("%s > %s", client.name, incoming)
				}

				// If there was an error, we delete the client
				deadConnections <- conn
			}(conn, server)


		// Continuously read incoming messages and broadcast them
		case message := <-messages:
			log.Printf("New message : %s", message)

			// Send the message to all connected clients
			for conn := range server.allClients {
				//Send the message in a go routine
				go func(conn net.Conn, message string) {
					_, err := conn.Write([]byte(message))
					// If it doesn't work the connection is dead
					if err != nil {
						deadConnections <- conn
					}
				}(conn, message)
			}

		//Remove dead clients
		case conn := <-deadConnections:
			disconnect(conn, server)
		}
	}
}
