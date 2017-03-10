package chord

import (
	"sync"
)

// Server represents a single node in Chord protocol
type Server struct {
	name string
	node *Node
	sync.RWMutex
	config *Config
}

// NewServer initializes a new local server involved in Chord protocol
func NewServer(name string, config *Config) *Server {
	server := &Server{
		name:   name,
		node:   NewNode(config),
		config: config,
	}
}

// ListNodes handles a incoming request sent from other server to list existing nodes
func (server *Server) ListNodes(req *ListNodesRequest) *ListNodesResponse {

}

// FindSuccessor handles a incoming request sent from other server to help find successor
func (server *Server) FindSuccessor(req *FindSuccessorRequest) *FindSuccessorResponse {

}
