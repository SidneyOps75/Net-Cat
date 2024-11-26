package utils

/*
removeClient removes the specified client from the server's list of connected clients.
The function first acquires a lock on the server's mutex to ensure thread-safety when
modifying the Clients map. Once the lock is obtained, it checks if the client is
still present in the Clients map. If so, it removes the client from the map, closes
the client's Messages channel, and closes the client's connection.
*/
func (s *Server) removeClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Clients[client]; ok {
		delete(s.Clients, client)
		close(client.Messages)
		client.Conn.Close()
	}
}
