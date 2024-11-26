package utils

func (s *Server) HandleClientMessages(client *Client) {
	client.Listen(s.Broadcast)
}
