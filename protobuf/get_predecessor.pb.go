// Code generated by protoc-gen-go.
// source: get_predecessor.proto
// DO NOT EDIT!

package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type GetPredecessorResponse struct {
	ID   string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Host string `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
}

func (m *GetPredecessorResponse) Reset()                    { *m = GetPredecessorResponse{} }
func (m *GetPredecessorResponse) String() string            { return proto.CompactTextString(m) }
func (*GetPredecessorResponse) ProtoMessage()               {}
func (*GetPredecessorResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *GetPredecessorResponse) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *GetPredecessorResponse) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func init() {
	proto.RegisterType((*GetPredecessorResponse)(nil), "protobuf.GetPredecessorResponse")
}

func init() { proto.RegisterFile("get_predecessor.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 106 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x4d, 0x4f, 0x2d, 0x89,
	0x2f, 0x28, 0x4a, 0x4d, 0x49, 0x4d, 0x4e, 0x2d, 0x2e, 0xce, 0x2f, 0xd2, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0xe2, 0x00, 0x53, 0x49, 0xa5, 0x69, 0x4a, 0x36, 0x5c, 0x62, 0xee, 0xa9, 0x25, 0x01,
	0x08, 0x15, 0x41, 0xa9, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x7c, 0x5c, 0x4c, 0x9e, 0x2e,
	0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x4c, 0x9e, 0x2e, 0x42, 0x42, 0x5c, 0x2c, 0x19, 0xf9,
	0xc5, 0x25, 0x12, 0x4c, 0x60, 0x11, 0x30, 0x3b, 0x89, 0x0d, 0x6c, 0x8e, 0x31, 0x20, 0x00, 0x00,
	0xff, 0xff, 0x38, 0xef, 0xd8, 0x93, 0x67, 0x00, 0x00, 0x00,
}
