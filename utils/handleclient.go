package utils

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// NewClient creates a new client with a connection and a reference to the server
func NewClient(conn net.Conn, s *Server) *Client {
	// Create a buffered reader to efficiently read from the connection
	reader := bufio.NewReader(conn)

	// Write an empty byte slice to the connection to ensure it's ready for reading
	conn.Write([]byte(""))

	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	client := &Client{
		Conn:     conn,
		Name:     name,
		Messages: make(chan string, 10),
	}

	// Return the newly created client
	return client
}

// SendMessages continuously sends any messages in the client's message channel to the connection
func (c *Client) SendMessages() {
	// Iterate over the messages in the channel
	for msg := range c.Messages {
		// Write the message to the connection
		_, err := c.Conn.Write([]byte(msg))
		if err != nil {
			break
		}
	}
}

// Listen reads messages from the client's connection, processes them, and broadcasts them to the server
func (c *Client) Listen(broadcast chan<- string) {
	// Create a buffered reader to efficiently read from the connection
	reader := bufio.NewReader(c.Conn)

	for {
		// Read a line of text from the connection
		message, err := reader.ReadString('\n')
		if err != nil {
			// If there's an error, broadcast that the client has left and break out of the loop
			broadcast <- fmt.Sprintf("%s has left our chat...\n", c.Name)
			break
		}

		message = strings.TrimSpace(message)

		if message == "" {
			continue
		}

		formattedMsg := fmt.Sprintf("[%s][%s]: %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Name,
			message)

		// Broadcast the formatted message to the server
		broadcast <- formattedMsg
	}
}
