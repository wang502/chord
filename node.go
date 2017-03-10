package chord

// Node represents a Node node involved in Chord protocol
type Node struct {
	ID         []byte
	successors []*Node
	finger     []*Node
}

// NewNode initializes a Node server involved in Chord protocol
func NewNode(config *Config) *Node {
	return &Node{
		ID:         generateID(config),
		successors: make([]*Node, config.NumSuccessors),
		finger:     make([]*Node, config.HashBits),
	}
}

// generateId is helper function that generates Id for a Node server based index number
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

// Successors returns slice of successors inside the Node
func (p *Node) Successors() []*Node {
	return p.successors
}

// Finger returns finger table inside the Node
func (p *Node) Finger() []*Node {
	return p.finger
}
