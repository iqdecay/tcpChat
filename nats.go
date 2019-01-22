package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// Client is a basic client to the NATS server
type Client struct {
	conn net.Conn
	w    *bufio.Writer
	sync.Mutex
}

func NewClient() *Client {
	return &Client{}
}

//	Connect establishes a connection to the NAT server
func (c *Client) Connect(netloc string) error {
	conn, err := net.Dial("tcp", netloc)
	if err != nil {
		return err
	}
	c.conn = conn
	c.w = bufio.NewWriter(conn) // conn implements io.Writer
	return nil
}

// Publish takes a subject as an immutable string and payload in bytes,
// then sends the message to the server
func (c *Client) Publish(subject string, payload []byte) error {
	c.Lock()
	pub := fmt.Sprintf("PUB %s %d\r\n", subject, len(payload))
	_, err := c.w.WriteString(pub)
	if err == nil {
		_, err = c.w.Write(payload)
	}
	if err == nil {
		_, err = c.w.WriteString("\r\n")
	}
	if err == nil {
		err = c.w.Flush()
	}
	c.Unlock()
	if err != nil {
		return err
	}
	return nil

}
