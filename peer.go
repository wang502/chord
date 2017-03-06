package chord

import (
    "encoding/binary"
)

// represents a peer node involved in Chord protocol
type Peer struct {
    Id []byte
    ring *Ring
    successors []*Peer
    finger []*Peer
}

// initialize a peer server involved in Chord protocol
func NewPeer(ring *Ring, idx int) (*Peer) {
    return &Peer{
        Id         : generateId(ring.config, uint16(idx)),
        ring       : ring,
        successors : make([]*Peer, ring.config.NumSuccessors),
        finger     : make([]*Peer, ring.config.HashBits),
    }
}

// generate Id for a peer server
func generateId(config *Config, idx uint16) ([]byte) {
    hash := config.HashFunc
    hash.Write([]byte(config.Host))
    binary.Write(hash, binary.BigEndian, idx)
    return hash.Sum(nil)
}
