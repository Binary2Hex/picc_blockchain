// Code generated by protoc-gen-go.
// source: insurance.proto
// DO NOT EDIT!

package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Insurance struct {
	Farm          string `protobuf:"bytes,1,opt,name=farm" json:"farm,omitempty"`
	Number        string `protobuf:"bytes,2,opt,name=number" json:"number,omitempty"`
	StartDate     string `protobuf:"bytes,3,opt,name=startDate" json:"startDate,omitempty"`
	EndDate       string `protobuf:"bytes,4,opt,name=endDate" json:"endDate,omitempty"`
	Checked       bool   `protobuf:"varint,5,opt,name=checked" json:"checked,omitempty"`
	Details       string `protobuf:"bytes,6,opt,name=details" json:"details,omitempty"`
	AmountInsured int64  `protobuf:"varint,7,opt,name=amountInsured" json:"amountInsured,omitempty"`
	AmountLoss    int64  `protobuf:"varint,8,opt,name=amountLoss" json:"amountLoss,omitempty"`
}

func (m *Insurance) Reset()                    { *m = Insurance{} }
func (m *Insurance) String() string            { return proto.CompactTextString(m) }
func (*Insurance) ProtoMessage()               {}
func (*Insurance) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func init() {
	proto.RegisterType((*Insurance)(nil), "main.Insurance")
}

func init() { proto.RegisterFile("insurance.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x54, 0x8f, 0xcd, 0xaa, 0xc2, 0x30,
	0x10, 0x46, 0xe9, 0x6d, 0x6f, 0x7f, 0x06, 0x2e, 0x17, 0x66, 0x21, 0x59, 0x88, 0x88, 0xb8, 0x70,
	0xe5, 0xc6, 0x57, 0x70, 0x23, 0xb8, 0xea, 0x1b, 0x4c, 0xdb, 0x11, 0x8b, 0x26, 0x91, 0x24, 0x7d,
	0x68, 0xdf, 0xc2, 0x76, 0xda, 0xa2, 0xee, 0x72, 0xce, 0xc9, 0x07, 0x09, 0xfc, 0xb7, 0xc6, 0x77,
	0x8e, 0x4c, 0xcd, 0xfb, 0x87, 0xb3, 0xc1, 0x62, 0xa2, 0xa9, 0x35, 0x9b, 0x67, 0x04, 0xc5, 0x69,
	0x2e, 0x88, 0x90, 0x5c, 0xc8, 0x69, 0x15, 0xad, 0xa3, 0x5d, 0x51, 0xca, 0x19, 0x17, 0x90, 0x9a,
	0x4e, 0x57, 0xec, 0xd4, 0x8f, 0xd8, 0x89, 0x70, 0x09, 0x85, 0x0f, 0xe4, 0xc2, 0x91, 0x02, 0xab,
	0x58, 0xd2, 0x5b, 0xa0, 0x82, 0x8c, 0x4d, 0x23, 0x2d, 0x91, 0x36, 0xe3, 0x50, 0xea, 0x2b, 0xd7,
	0x37, 0x6e, 0xd4, 0x6f, 0x5f, 0xf2, 0x72, 0xc6, 0xa1, 0x34, 0x1c, 0xa8, 0xbd, 0x7b, 0x95, 0x8e,
	0x9b, 0x09, 0x71, 0x0b, 0x7f, 0xa4, 0x6d, 0x67, 0x82, 0x3c, 0xb5, 0x5f, 0x66, 0x7d, 0x8f, 0xcb,
	0x6f, 0x89, 0x2b, 0x80, 0x51, 0x9c, 0xad, 0xf7, 0x2a, 0x97, 0x2b, 0x1f, 0xa6, 0x4a, 0xe5, 0xe3,
	0x87, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x04, 0x68, 0x85, 0xba, 0x0b, 0x01, 0x00, 0x00,
}
