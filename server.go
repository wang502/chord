package chord

import (
	"bytes"
	"fmt"
	"log"
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
	return server.Stabilize()
}

// Stabilize is called periodically to verify this server's immediate successor and tells the successor about this server
func (server *Server) Stabilize() error {
	//return nil
	if server.node.Successor() == nil {
		return fmt.Errorf("no need to stabilize.no sucessor")
	}

	successor := server.node.Successor()
	predResp, err := server.transporter.SendGetPredecessorRequest(server, successor.host)
	if predResp == nil {
		log.Printf("Chord.stabilize.error.%s", err)
	} else if err != nil {
		return fmt.Errorf("Chord.stabilize.notify.%s", err)
	} else {
		ID := []byte(predResp.ID)
		host := predResp.host

		// verifies server's immediate successor
		// if the successor's predecessor has an ID bigger than this server, then it means this server's immediate successor
		// should be updated to the one contained in the response
		if between(server.node.ID, successor.ID, ID) {
			server.node.SetSuccessor(NewRemoteNode(ID, host))
		}
	}

	// notify the immediate successor about the server
	_, err = server.transporter.SendNotifyRequest(server, NewNotifyRequest(server.node.ID, server.config.Host, server.node.Successor().host))
	if err != nil {
		return fmt.Errorf("Chord.stabilize.notify.%s", err)
	}

	return nil
}

// Notify handles the NotifyRequest sent from another server
func (server *Server) Notify(req *NotifyRequest) (*NotifyResponse, error) {
	possiblePredID := []byte(req.ID)
	possiblePredHost := req.host
	currentPredecessor := server.node.Predecessor()

	// when this node haven't set its predecessor, then new incoming notify request is from a node that should be a predecessor
	if currentPredecessor == nil {
		server.node.SetPredecessor(NewRemoteNode(possiblePredID, possiblePredHost))
		return NewNotifyResponse(server.node.ID, server.config.Host), nil
	}
	// update the predecessor if the notify request is from a node that has bigger byte value than the current predecessor
	if between(currentPredecessor.ID, server.node.ID, possiblePredID) {
		server.node.SetPredecessor(NewRemoteNode(possiblePredID, possiblePredHost))
		return NewNotifyResponse(server.node.ID, server.config.Host), nil
	}
	return nil, nil
}

// FindSuccessor handles a incoming request sent from other server to help find successor
func (server *Server) FindSuccessor(req *FindSuccessorRequest) (*FindSuccessorResponse, error) {
	id := []byte(req.ID)
	localNode := server.node
	resp := &FindSuccessorResponse{}

	// if this local node does not have successor yet, compare local server's bytes id with incoming id
	if localNode.Successor() == nil {
		resp.ID = string(localNode.ID)
		resp.host = server.config.Host
		return resp, nil
	}
	successor := localNode.Successor()
	if between(localNode.ID, successor.ID, id) {
		resp.ID = string(successor.ID)
		resp.host = successor.host
		return resp, nil
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

// HandleGetPredecessorRequest returns the predecessor of this local node
func (server *Server) HandleGetPredecessorRequest() (*GetPredecessorResponse, error) {
	pred := server.node.Predecessor()
	if pred == nil {
		return nil, fmt.Errorf("this node has no predecessor")
	}
	resp := NewGetPredecessorResponse(pred.ID, pred.host)
	return resp, nil
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
