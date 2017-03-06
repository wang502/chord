package chord

// represents a ring consistes of nodes in Chord protocol
type Ring struct {
    nodes []*Peer
    config *Config
}

// initialize a ring in local Chord server
func NewRing(conf *Config) (*Ring){
    ring := &Ring{}
    ring.nodes = make([]*Peer, conf.NumNodes)
    conf.HashBits = conf.HashFunc.Size() * 8
    ring.config = conf

    for i:=0; i < conf.NumNodes; i++ {
        ring.nodes[i] = NewPeer(ring, i)
    }
    return ring
}
