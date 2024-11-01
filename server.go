package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	maxClients  = 10 // Max number of clients
	defaultPort = "8989"
)

var (
	clients       = make(map[net.Conn]string) // Store client connections with names
	messages      []string                    // Chat history
	mu            sync.Mutex                  // Mutex to protect access to shared resources
	clientCounter int                         // Count of active clients
)

// Function to handle individual client connections
func handleClient(conn net.Conn) {
	defer conn.Close()

	// Request client's name
	conn.Write([]byte("Enter your name: "))
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	mu.Lock()
	if clientCounter >= maxClients {
		conn.Write([]byte("Chatroom is full. Try again later.\n"))
		mu.Unlock()
		return
	}

	clientCounter++
	clients[conn] = name
	mu.Unlock()

	broadcast(fmt.Sprintf("%s joined the chat", name), conn)

	// Send chat history to the new client
	mu.Lock()
	for _, msg := range messages {
		conn.Write([]byte(msg + "\n"))
	}
	mu.Unlock()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)

		// Ignore empty messages
		if message == "" {
			continue
		}

		// Format message with timestamp
		formattedMessage := fmt.Sprintf("[%s][%s]: %s",
			time.Now().Format("2006-01-02 15:04:05"), name, message)

		// Store the message in chat history
		mu.Lock()
		messages = append(messages, formattedMessage)
		mu.Unlock()

		// Broadcast the message to all clients
		broadcast(formattedMessage, conn)
	}

	// Handle client disconnection
	mu.Lock()
	delete(clients, conn)
	clientCounter--
	mu.Unlock()

	broadcast(fmt.Sprintf("%s left the chat", name), conn)
}

// Function to broadcast a message to all clients
func broadcast(message string, sender net.Conn) {
	fmt.Println(message) // Server log
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		if conn != sender { // Avoid sending message back to sender
			conn.Write([]byte(message + "\n"))
		}
	}
}

func main() {
	port := defaultPort
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	defer listener.Close()
	fmt.Printf("Server is listening on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}

		go handleClient(conn)
	}
}
