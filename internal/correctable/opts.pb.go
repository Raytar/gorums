// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: internal/correctable/opts.proto

package correctable

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
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

var file_internal_correctable_opts_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         6200,
		Name:          "correctable.correctable",
		Tag:           "varint,6200,opt,name=correctable",
		Filename:      "internal/correctable/opts.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         6201,
		Name:          "correctable.correctable_stream",
		Tag:           "varint,6201,opt,name=correctable_stream",
		Filename:      "internal/correctable/opts.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional bool correctable = 6200;
	E_Correctable = &file_internal_correctable_opts_proto_extTypes[0]
	// optional bool correctable_stream = 6201;
	E_CorrectableStream = &file_internal_correctable_opts_proto_extTypes[1]
)

var File_internal_correctable_opts_proto protoreflect.FileDescriptor

var file_internal_correctable_opts_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x72, 0x72, 0x65,
	0x63, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2f, 0x6f, 0x70, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0b, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x1a, 0x20,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x3a, 0x41, 0x0a, 0x0b, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xb8, 0x30, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x3a, 0x4e, 0x0a, 0x12, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xb9, 0x30, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x11, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x62, 0x2f, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_internal_correctable_opts_proto_goTypes = []interface{}{
	(*descriptorpb.MethodOptions)(nil), // 0: google.protobuf.MethodOptions
}
var file_internal_correctable_opts_proto_depIdxs = []int32{
	0, // 0: correctable.correctable:extendee -> google.protobuf.MethodOptions
	0, // 1: correctable.correctable_stream:extendee -> google.protobuf.MethodOptions
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	0, // [0:2] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_correctable_opts_proto_init() }
func file_internal_correctable_opts_proto_init() {
	if File_internal_correctable_opts_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_correctable_opts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_internal_correctable_opts_proto_goTypes,
		DependencyIndexes: file_internal_correctable_opts_proto_depIdxs,
		ExtensionInfos:    file_internal_correctable_opts_proto_extTypes,
	}.Build()
	File_internal_correctable_opts_proto = out.File
	file_internal_correctable_opts_proto_rawDesc = nil
	file_internal_correctable_opts_proto_goTypes = nil
	file_internal_correctable_opts_proto_depIdxs = nil
}
