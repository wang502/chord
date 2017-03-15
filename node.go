package chord

import (
	"sync"
)

// Node represents a Node node involved in Chord protocol
type Node struct {
	ID          []byte
	successor   *RemoteNode
	finger      []*RemoteNode
	predecessor *RemoteNode
	fingerIndex int
	sync.RWMutex
}

// RemoteNode represents a virtual remote Node involved in Chord protocol, containing hashed ID and host
type RemoteNode struct {
	ID   []byte
	host string
}

// NewNode initializes a Node server involved in Chord protocol
func NewNode(config *Config) *Node {
	return &Node{
		ID:          generateID(config),
		successor:   defaultSuccessor(config), // the successor is the node itself at the beginning
		finger:      make([]*RemoteNode, config.HashBits),
		predecessor: nil,
		fingerIndex: -1,
	}
}

// NewRemoteNode initializes a remote Node server involved in Chord protocol
func NewRemoteNode(id []byte, host string) *RemoteNode {
	return &RemoteNode{
		ID:   id,
		host: host,
	}
}

// generateId is helper function that uses configured hash function to generates Id for a Node server
func generateID(config *Config) []byte {
	hash := config.HashFunc
	hash.Write([]byte(config.Host))
	return hash.Sum(nil)
}

/*
   Getter
*/

// GetID returns the ID of the Node
func (n *Node) GetID() []byte {
	return n.ID
}

// Successor returns the successor inside this local Node
func (n *Node) Successor() *RemoteNode {
	n.Lock()
	defer n.Unlock()
	return n.successor
}

// Finger returns finger table inside the Node
func (n *Node) Finger() []*RemoteNode {
	n.Lock()
	defer n.Unlock()
	return n.finger
}

// Predecessor returns the predecessor
func (n *Node) Predecessor() *RemoteNode {
	n.Lock()
	defer n.Unlock()
	return n.predecessor
}

/*
  Setter
*/

// SetID sets node's id
func (n *Node) SetID(id []byte) {
	n.Lock()
	defer n.Unlock()
	n.ID = id
}

// SetSuccessor sets node's successor
func (n *Node) SetSuccessor(succ *RemoteNode) {
	n.Lock()
	defer n.Unlock()
	n.successor = succ
}

// SetPredecessor sets node's predecessor
func (n *Node) SetPredecessor(pred *RemoteNode) {
	n.Lock()
	defer n.Unlock()
	n.predecessor = pred
}

func defaultSuccessor(config *Config) *RemoteNode {
	return NewRemoteNode(generateID(config), config.Host)
}
