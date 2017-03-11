package chord

// Node represents a Node node involved in Chord protocol
type Node struct {
	ID          []byte
	successor   *RemoteNode
	finger      []*RemoteNode
	predecessor *Node
	fingerIndex int
}

type RemoteNode struct {
	ID   []byte
	host string
}

// NewNode initializes a Node server involved in Chord protocol
func NewNode(config *Config) *Node {
	return &Node{
		ID:          generateID(config),
		successor:   nil,
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
func (p *Node) GetID() []byte {
	return p.ID
}

// Successor returns the successor inside this local Node
func (p *Node) Successor() *RemoteNode {
	return p.successor
}

// Finger returns finger table inside the Node
func (p *Node) Finger() []*RemoteNode {
	return p.finger
}

// Predecessor returns the predecessor
func (p *Node) Predecessor() *Node {
	return p.predecessor
}
