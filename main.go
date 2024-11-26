package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"nc/utils"
)

// The main() function is responsible for setting up the server, starting the message processing loop,
// and handling new client connections. It provides the entry point for the TCP chat server application.
func main() {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	// Attempt to create a TCP listener on the specified port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	// Attempt to create a new instance of the Server struct
	srv, err := utils.NewServer()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	fmt.Printf("Listening on the port :%s\n", port)

	// Start the server's main message processing loop in a new goroutine
	go srv.Run()

	// Enter a loop to continuously accept new client connections
	for {
		// Accept a new client connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Start a new goroutine to handle the new client connection
		go srv.HandleConnection(conn)
	}
}
