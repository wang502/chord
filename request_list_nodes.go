package chord

import (
	"io"
)

// ListNodesRequest represents a request sent to another server to return all current nodes in the ring
type ListNodesRequest struct {
	serverName string
}

// ListNodesResponse represents a response returned containing info about all current nodes in the ring
type ListNodesResponse struct {
}

// NewListNodesRequest initilizes a new request to list all nodes
func NewListNodesRequest(serverName string) *ListNodesRequest {
	return &ListNodesRequest{
		serverName: serverName,
	}
}

func (req *ListNodesRequest) Encode(buf io.Writer) (int, error) {

}

func (req *ListNodesRequest) Decode(r io.Reader) (int, error) {

}

func (resp *ListNodesResponse) Encode(buf io.Writer) (int, error) {

}

func (resp *ListNodesResponse) Decode(r io.Reader) (int, error) {

}
