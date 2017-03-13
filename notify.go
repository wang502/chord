package chord

import (
	"io"
	"log"

	"fmt"

	"io/ioutil"

	"github.com/golang/protobuf/proto"
	pb "github.com/wang502/chord/protobuf"
)

// NotifyRequest represents a request sent to successor to notify it about local node
type NotifyRequest struct {
	ID   string
	host string
}

// NotifyResponse represents a response to a NotifyRequest
type NotifyResponse struct {
	ID   string
	host string
}

// NewNotifyRequest initializes a new notify request
func NewNotifyRequest(id []byte, host string) *NotifyRequest {
	return &NotifyRequest{
		ID:   string(id),
		host: host,
	}
}

// NewNotifyResponse initializes a new notify response
func NewNotifyResponse(id []byte, host string) *NotifyResponse {
	return &NotifyResponse{
		ID:   string(id),
		host: host,
	}
}

// Encode encodes NotifyRequest into data buffer
func (req *NotifyRequest) Encode(w io.Writer) (int, error) {
	pb := &pb.NotifyRequest{
		ID:   req.ID,
		Host: req.host,
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		log.Println("chord.NotifyRequest.encode.error")
		return -1, fmt.Errorf("chord.NotifyRequest.encode.error.%s", err)
	}

	return w.Write(data)
}

// Decode decodes data from buffer and stores it in NotifyRequest
func (req *NotifyRequest) Decode(r io.Reader) (int, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("chord.NotifyRequest.decode.error")
		return -1, fmt.Errorf("chord.NotifyRequest.decode.error.%s", err)
	}

	pb := &pb.NotifyRequest{}
	if err = proto.Unmarshal(data, pb); err != nil {
		return -1, fmt.Errorf("chord.NotifyRequest.decode.error.%s", err)
	}

	req.ID = pb.ID
	req.host = pb.Host
	return len(data), nil
}

// Encode encodes NotifyResponse into data buffer
func (resp *NotifyResponse) Encode(w io.Writer) (int, error) {
	pb := &pb.NotifyResponse{
		ID:   resp.ID,
		Host: resp.host,
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		log.Println("chord.NotifyResponse.encode.error")
		return -1, fmt.Errorf("chord.NotifyResponse.encode.error.%s", err)
	}

	return w.Write(data)
}

// Decode decodes data from buffer and stores it in NotifyResponse
func (resp *NotifyResponse) Decode(r io.Reader) (int, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("chord.NotifyResponse.decode.error")
		return -1, fmt.Errorf("chord.NotifyResponse.decode.error.%s", err)
	}

	pb := &pb.NotifyResponse{}
	if err = proto.Unmarshal(data, pb); err != nil {
		return -1, fmt.Errorf("chord.NotifyResponse.decode.error.%s", err)
	}

	resp.ID = pb.ID
	resp.host = pb.Host
	return len(data), nil
}
