package chord

import (
	"math/big"
	"sync"
)

// Node represents a Node node involved in Chord protocol
type Node struct {
	ID          []byte
	successor   *RemoteNode
	finger      []*FingerEntry
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
	node := &Node{
		ID: generateID(config),
		//successor:   defaultSuccessor(config), // the successor is the node itself at the beginning
		finger:      make([]*FingerEntry, config.HashBits),
		predecessor: nil,
		fingerIndex: -1,
	}
	node.successor = defaultSuccessor(node.ID, config.Host)

	return node
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
	b := hash.Sum(nil)

	idInt := big.Int{}
	idInt.SetBytes(b)

	// Get the ceiling
	two := big.NewInt(2)
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(config.HashBits)), nil)

	// Apply the mod
	idInt.Mod(&idInt, &ceil)

	// Add together
	return idInt.Bytes()
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
func (n *Node) Finger() []*FingerEntry {
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

func defaultSuccessor(id []byte, host string) *RemoteNode {
	return NewRemoteNode(id, host)
}

func defaultFingerEntry(id []byte, exp int, config *Config) *FingerEntry {
	return &FingerEntry{
		start: powerOffset(id, exp, config.HashBits),
		node:  id,
		host:  config.Host,
	}
}

func defaultFingerTable(config *Config) []*FingerEntry {
	hb := config.HashBits
	id := generateID(config)
	fingerTable := make([]*FingerEntry, hb)

	for i := 0; i < hb; i++ {
		fingerTable[i] = defaultFingerEntry(id, i, config)
	}

	return fingerTable
}
