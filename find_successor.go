package chord

import "io"
import "fmt"

//FindSuccessorRequest represents a request entry sent to other server to find successor of this local node
type FindSuccessorRequest struct {
	ID   string
	host string
}

//FindSuccessorResponse represents a response entry sent back to other server to help find successor
type FindSuccessorResponse struct {
	ID   string
	host string
}

//NewFindSuccessorRequest initializes a new request to find successor
func NewFindSuccessorRequest(bytes []byte, host string) *FindSuccessorRequest {
	return &FindSuccessorRequest{
		ID:   fmt.Sprintf("%s", bytes),
		host: host,
	}
}

func (req *FindSuccessorRequest) Encode(buf io.Writer) (int, error) {

}

func (req *FindSuccessorRequest) Decode(r io.Reader) (int, error) {

}

func (resp *FindSuccessorResponse) Encode(buf io.Writer) (int, error) {

}

func (resp *FindSuccessorResponse) Decode(r io.Reader) (int, error) {

}
