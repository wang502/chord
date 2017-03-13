package chord

import (
	"bytes"
	"fmt"
	"sync"
)

// Server represents a single node in Chord protocol
type Server struct {
	name string
	node *Node
	sync.RWMutex
	config      *Config
	transporter *Transporter
}

// NewServer initializes a new local server involved in Chord protocol
func NewServer(name string, config *Config, transporter *Transporter) *Server {
	server := &Server{
		name:        name,
		node:        NewNode(config),
		config:      config,
		transporter: transporter,
	}
	return server
}

// Join joins an existing chord ring, given existingHost is one of the node in the ring
func (server *Server) Join(existingHost string) error {
	localNode := server.node
	findSuccessorReq := NewFindSuccessorRequest(localNode.ID, existingHost)
	findSuccessorResp, err := server.transporter.SendFindSuccessorRequest(server, findSuccessorReq)
	if err != nil {
		return fmt.Errorf("Chord.join.error.%s", err)
	}

	successorNode := NewRemoteNode([]byte(findSuccessorResp.ID), findSuccessorResp.host)
	localNode.SetSuccessor(successorNode)
	return nil
}

// Stabilize is called periodically to verify this server's immediate successor and tells the successor about this server
func (server *Server) Stabilize() error {
	return nil
}

// Notify handles the NotifyRequest sent from another server
func (server *Server) Notify(req *NotifyRequest) (*NotifyResponse, error) {
	return nil, nil
}

// FindSuccessor handles a incoming request sent from other server to help find successor
func (server *Server) FindSuccessor(req *FindSuccessorRequest) (*FindSuccessorResponse, error) {
	id := []byte(req.ID)
	localNode := server.node
	resp := &FindSuccessorResponse{}

	// if this local node does not have successor yet, compare local server's bytes id with incoming id
	if localNode.successor == nil {
		if bytes.Compare(id, localNode.ID) == -1 {
			resp.ID = string(localNode.ID)
			resp.host = server.config.Host
			return resp, nil
		}
	} else {
		successor := localNode.Successor()
		if between(localNode.ID, []byte(successor.ID), id) {
			resp.ID = string(successor.ID)
			resp.host = successor.host
			return resp, nil
		}
	}

	closestPre := server.closestPreceedingNode(id)
	if closestPre == nil {
		resp.ID = string(localNode.ID)
		resp.host = server.config.Host
		return resp, nil
	}
	findSuccRequest := NewFindSuccessorRequest(id, closestPre.host)
	return server.transporter.SendFindSuccessorRequest(server, findSuccRequest)
}

// closestPreceedingNode is a helper function to find the cloest preceeding node of the node with given hashed id from finger table
func (server *Server) closestPreceedingNode(id []byte) *RemoteNode {
	localNode := server.node
	finger := localNode.Finger()
	for i := server.config.HashBits - 1; i >= 0; i-- {
		if finger[i] != nil {
			if between(localNode.ID, id, finger[i].ID) {
				return finger[i]
			}
		}
	}
	//return &RemoteNode{ID: localNode.ID, host: server.config.Host}
	return nil
}

// utils

// Checks if a key is STRICTLY between two ID's exclusively
func between(id1, id2, key []byte) bool {
	// Check for ring wrap around
	if bytes.Compare(id1, id2) == 1 {
		return bytes.Compare(id1, key) == -1 ||
			bytes.Compare(id2, key) == 1
	}

	// Handle the normal case
	return bytes.Compare(id1, key) == -1 &&
		bytes.Compare(id2, key) == 1
}
