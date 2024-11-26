package utils

import (
	"fmt"
	"net"
)

// The HandleConnection() function is responsible for the initial setup and
// integration of a new client into the server's message broadcasting system.
// It ensures that the new client is properly added, notified to other clients,
// and able to both send and receive messages.
func (s *Server) HandleConnection(conn net.Conn) {
	// Write the server's logo to the new connection
	conn.Write([]byte(s.logo))

	// Create a new Client instance for the connection
	client := NewClient(conn, s)

	if !s.addClient(client) {

		conn.Write([]byte("Server is full. Please try again later.\n"))
		conn.Close()
		return
	}

	s.Broadcast <- fmt.Sprintf("%s has joined our chat...\n", client.Name)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, msg := range s.History {
		client.Messages <- msg
	}

	// Start a new goroutine to handle processing of messages sent by the new client
	go s.HandleClientMessages(client)

	// Start a new goroutine to handle sending of messages from the new client's Messages channel to the connection
	go client.SendMessages()
}
