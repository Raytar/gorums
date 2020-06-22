// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.11.2
// source: nodestream.proto

package dev

import (
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/relab/gorums"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//	*Message_Request
	//	*Message_Response
	//	*Message_Empty
	Data     isMessage_Data `protobuf_oneof:"Data"`
	ID       uint64         `protobuf:"varint,20000,opt,name=ID,proto3" json:"ID,omitempty"`
	MethodID int32          `protobuf:"varint,20001,opt,name=MethodID,proto3" json:"MethodID,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodestream_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_nodestream_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_nodestream_proto_rawDescGZIP(), []int{0}
}

func (m *Message) GetData() isMessage_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *Message) GetRequest() *Request {
	if x, ok := x.GetData().(*Message_Request); ok {
		return x.Request
	}
	return nil
}

func (x *Message) GetResponse() *Response {
	if x, ok := x.GetData().(*Message_Response); ok {
		return x.Response
	}
	return nil
}

func (x *Message) GetEmpty() *empty.Empty {
	if x, ok := x.GetData().(*Message_Empty); ok {
		return x.Empty
	}
	return nil
}

func (x *Message) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Message) GetMethodID() int32 {
	if x != nil {
		return x.MethodID
	}
	return 0
}

type isMessage_Data interface {
	isMessage_Data()
}

type Message_Request struct {
	Request *Request `protobuf:"bytes,1,opt,name=Request,proto3,oneof"`
}

type Message_Response struct {
	Response *Response `protobuf:"bytes,2,opt,name=Response,proto3,oneof"`
}

type Message_Empty struct {
	Empty *empty.Empty `protobuf:"bytes,3,opt,name=Empty,proto3,oneof"`
}

func (*Message_Request) isMessage_Data() {}

func (*Message_Response) isMessage_Data() {}

func (*Message_Empty) isMessage_Data() {}

var File_nodestream_proto protoreflect.FileDescriptor

var file_nodestream_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0d, 0x64, 0x65, 0x76, 0x2e, 0x67, 0x65, 0x6e, 0x67, 0x6f, 0x72, 0x75, 0x6d,
	0x73, 0x1a, 0x0c, 0x7a, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0c, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8, 0x01, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x2b, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x48, 0x00, 0x52, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a,
	0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x48, 0x00, 0x52, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x10, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0xa0, 0x9c, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12,
	0x1c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x44, 0x18, 0xa1, 0x9c, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x44, 0x42, 0x06, 0x0a,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x32, 0x4a, 0x0a, 0x06, 0x47, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x12,
	0x40, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x16, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x67, 0x65, 0x6e, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x67, 0x65, 0x6e, 0x67,
	0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x28, 0x01, 0x30,
	0x01, 0x42, 0x1b, 0x5a, 0x19, 0x63, 0x6d, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d,
	0x67, 0x65, 0x6e, 0x2d, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2f, 0x64, 0x65, 0x76, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nodestream_proto_rawDescOnce sync.Once
	file_nodestream_proto_rawDescData = file_nodestream_proto_rawDesc
)

func file_nodestream_proto_rawDescGZIP() []byte {
	file_nodestream_proto_rawDescOnce.Do(func() {
		file_nodestream_proto_rawDescData = protoimpl.X.CompressGZIP(file_nodestream_proto_rawDescData)
	})
	return file_nodestream_proto_rawDescData
}

var file_nodestream_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_nodestream_proto_goTypes = []interface{}{
	(*Message)(nil),     // 0: dev.gengorums.Message
	(*Request)(nil),     // 1: dev.Request
	(*Response)(nil),    // 2: dev.Response
	(*empty.Empty)(nil), // 3: google.protobuf.Empty
}
var file_nodestream_proto_depIdxs = []int32{
	1, // 0: dev.gengorums.Message.Request:type_name -> dev.Request
	2, // 1: dev.gengorums.Message.Response:type_name -> dev.Response
	3, // 2: dev.gengorums.Message.Empty:type_name -> google.protobuf.Empty
	0, // 3: dev.gengorums.Gorums.NodeStream:input_type -> dev.gengorums.Message
	0, // 4: dev.gengorums.Gorums.NodeStream:output_type -> dev.gengorums.Message
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_nodestream_proto_init() }
func file_nodestream_proto_init() {
	if File_nodestream_proto != nil {
		return
	}
	file_zorums_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_nodestream_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_nodestream_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_Request)(nil),
		(*Message_Response)(nil),
		(*Message_Empty)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_nodestream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_nodestream_proto_goTypes,
		DependencyIndexes: file_nodestream_proto_depIdxs,
		MessageInfos:      file_nodestream_proto_msgTypes,
	}.Build()
	File_nodestream_proto = out.File
	file_nodestream_proto_rawDesc = nil
	file_nodestream_proto_goTypes = nil
	file_nodestream_proto_depIdxs = nil
}
