package chord

import (
	"io"
	"io/ioutil"
	"log"

	"fmt"

	"github.com/golang/protobuf/proto"
	pb "github.com/wang502/chord/protobuf"
)

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
		ID:   string(bytes),
		host: host,
	}
}

// Encode encodes the FindSuccessorRequest into a buffer
// returns the number of bytes written to the buffer, and error if occurred
func (req *FindSuccessorRequest) Encode(buf io.Writer) (int, error) {
	pb := &pb.FindSuccessorRequest{
		ID:   req.ID,
		Host: req.host,
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		log.Println("[ERROR]chord.FindSuccessorRequest.encode.error")
		return -1, fmt.Errorf("encode FindSuccessorRequest failed: %s", err)
	}

	return buf.Write(data)
}

// Decode decodes the bytes read from buffer and store data into FindSuccessorRequest entry
func (req *FindSuccessorRequest) Decode(r io.Reader) (int, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		//log.Println("[ERROR]chord.FindSuccessorRequest.decode.error")
		return -1, fmt.Errorf("decode FindSuccessorRequest failed: %s", err)
	}

	pb := &pb.FindSuccessorRequest{}
	if err = proto.Unmarshal(data, pb); err != nil {
		//log.Println("[ERROR]chord.FindSuccessorRequest.decode.error")
		return -1, fmt.Errorf("decode FindSuccessorRequest failed: %s", err)
	}

	req.ID = pb.ID
	req.host = pb.Host
	return len(data), nil
}

// Encode encodes the FindSuccessorResponse into a buffer
// returns the number of bytes written to the buffer, and error if occurred
func (resp *FindSuccessorResponse) Encode(buf io.Writer) (int, error) {
	pb := &pb.FindSuccessorResponse{
		ID:   resp.ID,
		Host: resp.host,
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		log.Println("[ERROR]chord.FindSuccessorResponse.encode.error")
		return -1, fmt.Errorf("encode FindSuccessorReponse failed: %s", err)
	}

	return buf.Write(data)
}

// Decode decodes the bytes read from buffer and store data into FindSuccessorResponse entry
func (resp *FindSuccessorResponse) Decode(r io.Reader) (int, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("[ERROR]chord.FindSuccessorRequest.decode.error")
		return -1, fmt.Errorf("decode FindSuccessorResponse failed: %s", err)
	}

	pb := &pb.FindSuccessorResponse{}
	if err = proto.Unmarshal(data, pb); err != nil {
		log.Println("[ERROR]chord.FindSuccessorResponse.decode.error")
		return -1, fmt.Errorf("decode FindSuccessorResponse failed: %s", err)
	}

	resp.ID = pb.ID
	resp.host = pb.Host
	return len(data), nil
}
