// Code generated by protoc-gen-go.
// source: notify.proto
// DO NOT EDIT!

package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type NotifyRequest struct {
	ID         string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Host       string `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
	TargetHost string `protobuf:"bytes,3,opt,name=targetHost" json:"targetHost,omitempty"`
}

func (m *NotifyRequest) Reset()                    { *m = NotifyRequest{} }
func (m *NotifyRequest) String() string            { return proto.CompactTextString(m) }
func (*NotifyRequest) ProtoMessage()               {}
func (*NotifyRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *NotifyRequest) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *NotifyRequest) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *NotifyRequest) GetTargetHost() string {
	if m != nil {
		return m.TargetHost
	}
	return ""
}

type NotifyResponse struct {
	ID   string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Host string `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
}

func (m *NotifyResponse) Reset()                    { *m = NotifyResponse{} }
func (m *NotifyResponse) String() string            { return proto.CompactTextString(m) }
func (*NotifyResponse) ProtoMessage()               {}
func (*NotifyResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *NotifyResponse) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *NotifyResponse) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func init() {
	proto.RegisterType((*NotifyRequest)(nil), "protobuf.NotifyRequest")
	proto.RegisterType((*NotifyResponse)(nil), "protobuf.NotifyResponse")
}

func init() { proto.RegisterFile("notify.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 126 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xc9, 0xcb, 0x2f, 0xc9,
	0x4c, 0xab, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0x49, 0xa5, 0x69, 0x4a,
	0xc1, 0x5c, 0xbc, 0x7e, 0x60, 0x99, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x3e, 0x2e,
	0x26, 0x4f, 0x17, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x26, 0x4f, 0x17, 0x21, 0x21, 0x2e,
	0x96, 0x8c, 0xfc, 0xe2, 0x12, 0x09, 0x26, 0xb0, 0x08, 0x98, 0x2d, 0x24, 0xc7, 0xc5, 0x55, 0x92,
	0x58, 0x94, 0x9e, 0x5a, 0xe2, 0x01, 0x92, 0x61, 0x06, 0xcb, 0x20, 0x89, 0x28, 0x99, 0x70, 0xf1,
	0xc1, 0x0c, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x25, 0xc6, 0xd4, 0x24, 0x36, 0xb0, 0xa3, 0x8c,
	0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x92, 0x4f, 0x71, 0x08, 0xab, 0x00, 0x00, 0x00,
}
