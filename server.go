package chord

import (
	"sync"
)

// Server represents a single node in Chord protocol
type Server struct {
	name string
	ring *Ring
	sync.RWMutex
}
