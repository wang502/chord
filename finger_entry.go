package chord

import (
	"fmt"
)

// FingerEntry represents a entry in the finger table
type FingerEntry struct {
	start []byte
	node  []byte
	host  string
}

func (entry *FingerEntry) String() string {
	return fmt.Sprintf("start byte: %x \n node byte: %x \n host: %s", entry.start, entry.node, entry.host)
}
