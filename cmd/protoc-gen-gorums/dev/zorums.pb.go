// Code generated by protoc-gen-go. DO NOT EDIT.
// source: zorums.proto

package dev

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/relab/gorums"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Request struct {
	Value                string   `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	MsgID                uint64   `protobuf:"varint,2,opt,name=MsgID,proto3" json:"MsgID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_eefcc80370fc6ec7, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Request) GetMsgID() uint64 {
	if m != nil {
		return m.MsgID
	}
	return 0
}

type Response struct {
	Result               int64    `protobuf:"varint,1,opt,name=Result,proto3" json:"Result,omitempty"`
	MsgID                uint64   `protobuf:"varint,2,opt,name=MsgID,proto3" json:"MsgID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_eefcc80370fc6ec7, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetResult() int64 {
	if m != nil {
		return m.Result
	}
	return 0
}

func (m *Response) GetMsgID() uint64 {
	if m != nil {
		return m.MsgID
	}
	return 0
}

type MyResponse struct {
	Value                string   `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	Result               int64    `protobuf:"varint,2,opt,name=Result,proto3" json:"Result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MyResponse) Reset()         { *m = MyResponse{} }
func (m *MyResponse) String() string { return proto.CompactTextString(m) }
func (*MyResponse) ProtoMessage()    {}
func (*MyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_eefcc80370fc6ec7, []int{2}
}

func (m *MyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MyResponse.Unmarshal(m, b)
}
func (m *MyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MyResponse.Marshal(b, m, deterministic)
}
func (m *MyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MyResponse.Merge(m, src)
}
func (m *MyResponse) XXX_Size() int {
	return xxx_messageInfo_MyResponse.Size(m)
}
func (m *MyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MyResponse proto.InternalMessageInfo

func (m *MyResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *MyResponse) GetResult() int64 {
	if m != nil {
		return m.Result
	}
	return 0
}

func init() {
	proto.RegisterType((*Request)(nil), "dev.Request")
	proto.RegisterType((*Response)(nil), "dev.Response")
	proto.RegisterType((*MyResponse)(nil), "dev.MyResponse")
}

func init() {
	proto.RegisterFile("zorums.proto", fileDescriptor_eefcc80370fc6ec7)
}

var fileDescriptor_eefcc80370fc6ec7 = []byte{
	// 681 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x96, 0xcf, 0x73, 0xd2, 0x40,
	0x14, 0xc7, 0xa5, 0xc5, 0x0a, 0x4f, 0xaa, 0x6d, 0xda, 0x22, 0xc2, 0xa5, 0x72, 0xc2, 0x43, 0x03,
	0x52, 0xeb, 0x74, 0x3a, 0xda, 0x51, 0xa1, 0x68, 0x55, 0x2a, 0x0d, 0x4e, 0x9d, 0xe9, 0x0d, 0xc2,
	0x1a, 0x99, 0x49, 0x08, 0xe6, 0x47, 0x67, 0xf0, 0xe4, 0x4d, 0x8e, 0x1c, 0x7b, 0xcc, 0xd1, 0x63,
	0x2f, 0xbd, 0x78, 0xf4, 0x5f, 0xf2, 0x1f, 0xe8, 0x26, 0x9b, 0x94, 0x85, 0x05, 0xb2, 0xb9, 0x75,
	0xb7, 0xf9, 0x7e, 0x36, 0xfb, 0xde, 0xe7, 0x01, 0x90, 0xfa, 0xa9, 0x1b, 0xb6, 0x66, 0x8a, 0x7d,
	0x43, 0xb7, 0x74, 0x61, 0xb9, 0x83, 0x2e, 0xb2, 0x29, 0x85, 0xda, 0xca, 0xe6, 0x14, 0x5d, 0x57,
	0x54, 0x54, 0xf4, 0x56, 0x6d, 0xfb, 0x5b, 0x11, 0x69, 0x7d, 0x6b, 0x40, 0xfe, 0x99, 0xdf, 0x83,
	0x7b, 0x12, 0xfa, 0x61, 0x23, 0xd3, 0x12, 0x36, 0xe1, 0xee, 0x59, 0x4b, 0xb5, 0x51, 0x26, 0xb6,
	0x1d, 0x2b, 0x24, 0x25, 0xb2, 0x70, 0x77, 0xeb, 0xa6, 0x72, 0x5c, 0xcd, 0x2c, 0xe1, 0xdd, 0xb8,
	0x44, 0x16, 0xf9, 0x7d, 0x48, 0x48, 0xc8, 0xec, 0xeb, 0x3d, 0x13, 0x09, 0x69, 0x58, 0xc1, 0x7f,
	0xdb, 0xaa, 0xe5, 0x05, 0x97, 0x25, 0x7f, 0x35, 0x27, 0x79, 0x00, 0x50, 0x1f, 0xdc, 0x66, 0x67,
	0x9f, 0x39, 0x26, 0x2e, 0xd1, 0xc4, 0xf2, 0xef, 0x0d, 0x58, 0x3d, 0xf7, 0xae, 0xd6, 0x44, 0xc6,
	0x45, 0x57, 0x46, 0xc2, 0x53, 0x48, 0xbc, 0x93, 0x1a, 0x95, 0x4a, 0x4b, 0x55, 0x85, 0x94, 0x88,
	0xef, 0x2e, 0xfa, 0xb7, 0xc9, 0xae, 0xfa, 0x2b, 0x72, 0x50, 0xfe, 0x8e, 0x50, 0x04, 0x38, 0xb5,
	0xdd, 0x70, 0xf8, 0xc3, 0xf1, 0x5f, 0xd7, 0x99, 0x98, 0xf0, 0x0a, 0x36, 0xc7, 0x81, 0x06, 0x32,
	0x4e, 0xf4, 0x0e, 0x7a, 0x63, 0x28, 0x8b, 0xa3, 0x09, 0x37, 0x7a, 0xe5, 0xc6, 0x2b, 0x90, 0x1d,
	0xc7, 0x4f, 0x6b, 0x5f, 0xbb, 0xd6, 0x77, 0x3f, 0xc3, 0x07, 0xf9, 0xe3, 0x42, 0xea, 0x34, 0xa4,
	0x62, 0x9b, 0x96, 0xae, 0x49, 0xc8, 0xb2, 0x8d, 0xde, 0x97, 0x41, 0x1f, 0x2d, 0x86, 0x08, 0x2e,
	0xe4, 0xdf, 0xff, 0x0c, 0x5d, 0xee, 0x0f, 0xf0, 0x90, 0xc2, 0xe9, 0x5a, 0x5b, 0x5f, 0xcc, 0xc8,
	0x06, 0x2f, 0x72, 0xc5, 0xb2, 0x0e, 0x69, 0xd6, 0x91, 0xab, 0x94, 0x90, 0x16, 0x89, 0x6a, 0x62,
	0xa0, 0x9a, 0xe8, 0xed, 0xcf, 0x2e, 0xef, 0x21, 0xac, 0x4d, 0xe5, 0xcb, 0x53, 0x2f, 0x33, 0x07,
	0xe7, 0xe7, 0x4b, 0x90, 0xac, 0x63, 0x29, 0xba, 0x72, 0x0b, 0xbb, 0xbb, 0xb8, 0x9d, 0x97, 0xf8,
	0xf9, 0x82, 0x7b, 0xe2, 0xc6, 0x6d, 0x82, 0xbb, 0x9f, 0x97, 0xfe, 0xed, 0x71, 0xfe, 0x19, 0x56,
	0x37, 0xc8, 0x97, 0xf9, 0x8e, 0x3c, 0xa0, 0x22, 0xbb, 0xbc, 0xd7, 0xf3, 0xb3, 0x55, 0x2a, 0xfb,
	0x7c, 0x6e, 0x6d, 0xc3, 0x28, 0x7b, 0x74, 0x99, 0x6b, 0x36, 0x76, 0x27, 0xc4, 0x9b, 0xb8, 0xc3,
	0xd8, 0x4b, 0x62, 0xdc, 0x25, 0x73, 0x82, 0x11, 0x38, 0x86, 0xed, 0x69, 0x48, 0xc4, 0x41, 0x70,
	0x82, 0x41, 0x68, 0xb2, 0xa8, 0x88, 0xe3, 0xe0, 0xb0, 0x0a, 0x37, 0x60, 0x8b, 0x81, 0x72, 0x0c,
	0x85, 0x33, 0x7f, 0x28, 0x5e, 0xc0, 0xfa, 0x34, 0xb1, 0xcc, 0x57, 0x6e, 0xe6, 0x4d, 0xc8, 0x48,
	0x71, 0x2a, 0xe3, 0x41, 0x8e, 0x20, 0x3d, 0x13, 0x52, 0xe6, 0x1e, 0x4c, 0x87, 0x0c, 0xd6, 0xfd,
	0x8a, 0x6e, 0x18, 0x48, 0xb6, 0x5a, 0x6d, 0x35, 0x4c, 0x96, 0x21, 0x19, 0xe5, 0x2d, 0x2a, 0xc1,
	0xed, 0xc9, 0x30, 0xf0, 0xa4, 0x0a, 0x39, 0x2a, 0x1f, 0x51, 0x91, 0x61, 0xa0, 0xc8, 0xc9, 0x04,
	0x25, 0xa2, 0x1d, 0x43, 0xb6, 0x97, 0x1f, 0x61, 0x8d, 0xe6, 0x71, 0x88, 0x31, 0x5c, 0xf4, 0x69,
	0x49, 0xc3, 0x22, 0xf5, 0xd6, 0x2b, 0xf1, 0x6b, 0x58, 0x9f, 0xce, 0xf3, 0xb7, 0xd5, 0x23, 0xbc,
	0x84, 0x07, 0x4d, 0xcb, 0xe8, 0xca, 0xd6, 0x67, 0xa3, 0x83, 0x8c, 0x6e, 0x2f, 0xa4, 0xae, 0xc9,
	0xbf, 0xd7, 0x19, 0xf2, 0x95, 0x5d, 0x88, 0x95, 0x62, 0xc2, 0xfe, 0xc4, 0xf9, 0x18, 0x84, 0x5a,
	0x5a, 0x88, 0x1a, 0x23, 0x7c, 0x2a, 0x4e, 0xd6, 0x26, 0xda, 0x42, 0x92, 0xdc, 0x8a, 0x8c, 0xfc,
	0x6a, 0x62, 0xce, 0x27, 0x78, 0xc2, 0x70, 0x22, 0xaa, 0x32, 0xf2, 0xfb, 0x83, 0x69, 0x67, 0x33,
	0x68, 0x11, 0x95, 0x19, 0x31, 0x5d, 0x2e, 0xb9, 0x9f, 0x53, 0x69, 0x96, 0xcb, 0xa1, 0xce, 0x68,
	0xae, 0x3a, 0x5e, 0x09, 0x59, 0x68, 0x24, 0x85, 0xfc, 0x56, 0xbc, 0x87, 0x47, 0xb3, 0x39, 0xfc,
	0x2a, 0x11, 0xd2, 0xdb, 0xdc, 0xf9, 0x63, 0x59, 0xeb, 0x90, 0x9f, 0x94, 0xf2, 0x8e, 0x82, 0x7a,
	0x3b, 0xe4, 0x27, 0x67, 0x11, 0x3f, 0xdd, 0x5e, 0xf1, 0xb6, 0x77, 0x6f, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x0f, 0x1d, 0x60, 0x3b, 0x9a, 0x0a, 0x00, 0x00,
}
