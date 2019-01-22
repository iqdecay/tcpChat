package main

import (
	"bufio"
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
