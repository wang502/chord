package chord

import (
    "sync"
)

// represents a single node in Chord protocol
type server struct {
    name string
    ring *Ring
    sync.RWMutex
}
