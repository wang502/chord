package chord

import (
	"encoding/binary"
)

// Peer represents a peer node involved in Chord protocol
type Peer struct {
	ID         []byte
	ring       *Ring
	successors []*Peer
	finger     []*Peer
}

// NewPeer initializes a peer server involved in Chord protocol
func NewPeer(ring *Ring, idx int) *Peer {
	return &Peer{
		ID:         generateID(ring.config, uint16(idx)),
		ring:       ring,
		successors: make([]*Peer, ring.config.NumSuccessors),
		finger:     make([]*Peer, ring.config.HashBits),
	}
}

// generateId is helper function that generates Id for a peer server based index number
func generateID(config *Config, idx uint16) []byte {
	hash := config.HashFunc
	hash.Write([]byte(config.Host))
	binary.Write(hash, binary.BigEndian, idx)
	return hash.Sum(nil)
}

/*
   Getter
*/

// GetID returns the ID of the peer
func (p *Peer) GetID() []byte {
	return p.ID
}

// Ring returns pointer to ring inside the peer
func (p *Peer) Ring() *Ring {
	return p.ring
}

// Successors returns slice of successors inside the peer
func (p *Peer) Successors() []*Peer {
	return p.successors
}

// Finger returns finger table inside the peer
func (p *Peer) Finger() []*Peer {
	return p.finger
}
