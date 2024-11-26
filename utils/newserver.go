package utils

// NewServer creates a new instance of the Server struct, initializing its fields
// and setting up the necessary channels and data structures to manage client
// connections and message broadcasting.
func NewServer() (*Server, error) {
	logo, err := LoadLogo()
	if err != nil {
		return nil, err
	}
	return &Server{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan string),
		History:    make([]string, 0, 100),
		maxClients: 10,
		logo:       logo,
	}, nil
}
