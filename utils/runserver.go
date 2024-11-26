package utils

// The Run() function is designed to run in a separate goroutine, as it will
// continue to process messages indefinitely until the server is stopped.
func (s *Server) Run() {
	for message := range s.Broadcast {
		// Acquire a lock on the server's mutex to ensure thread-safety
		s.mu.Lock()

		// Append the message to the message history
		s.History = append(s.History, message)

		// Release the lock on the server's mutex
		s.mu.Unlock()

		// Iterate over the connected clients and attempt to send the message to each one
		for client := range s.Clients {
			select {
			case client.Messages <- message:
				// If the send to the client's Messages channel is successful, continue to the next client
			default:
				// If the send to the client's Messages channel would block, remove the client
				s.removeClient(client)
			}
		}
	}
}
