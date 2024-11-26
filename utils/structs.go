package utils

import (
	"net"
	"sync"
)

type Client struct {
	Conn     net.Conn
	Name     string
	Messages chan string
}

type Server struct {
	Clients    map[*Client]bool
	Broadcast  chan string
	History    []string
	mu         sync.RWMutex
	maxClients int
	logo       string
}
