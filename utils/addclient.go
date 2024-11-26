package utils

func (s *Server) addClient(client *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the client name already exists
	for existingClient := range s.Clients {
		if existingClient.Name == client.Name {
			client.Conn.Write([]byte("Name already in use, please choose another.\n"))
			return false
		}
	}

	// Check if the server has reached its max client capacity
	if len(s.Clients) >= s.maxClients {
		return false
	}

	// Add the client if the name is unique
	s.Clients[client] = true
	return true
}
