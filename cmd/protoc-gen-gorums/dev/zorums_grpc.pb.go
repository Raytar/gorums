// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package dev

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ReaderServiceClient is the client API for ReaderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReaderServiceClient interface {
	ReadGrpc(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	ReadQuorumCall(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	ReadQuorumCallPerNodeArg(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	ReadQuorumCallQFWithRequestArg(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	ReadQuorumCallCustomReturnType(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	// ReadQuorumCallCombo does it all. Comment testing.
	ReadQuorumCallCombo(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	// ReadMulticast is testing a comment.
	ReadMulticast(ctx context.Context, opts ...grpc.CallOption) (ReaderService_ReadMulticastClient, error)
	ReadQuorumCallFuture(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	ReadCorrectable(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	ReadCorrectableStream(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (ReaderService_ReadCorrectableStreamClient, error)
}

type readerServiceClient struct {
	cc *grpc.ClientConn
}

func NewReaderServiceClient(cc *grpc.ClientConn) ReaderServiceClient {
	return &readerServiceClient{cc}
}

func (c *readerServiceClient) ReadGrpc(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dev.ReaderService/ReadGrpc", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readerServiceClient) ReadQuorumCall(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dev.ReaderService/ReadQuorumCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readerServiceClient) ReadQuorumCallPerNodeArg(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallPerNodeArg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readerServiceClient) ReadQuorumCallQFWithRequestArg(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallQFWithRequestArg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readerServiceClient) ReadQuorumCallCustomReturnType(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallCustomReturnType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readerServiceClient) ReadQuorumCallCombo(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallCombo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readerServiceClient) ReadMulticast(ctx context.Context, opts ...grpc.CallOption) (ReaderService_ReadMulticastClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ReaderService_serviceDesc.Streams[0], "/dev.ReaderService/ReadMulticast", opts...)
	if err != nil {
		return nil, err
	}
	x := &readerServiceReadMulticastClient{stream}
	return x, nil
}

type ReaderService_ReadMulticastClient interface {
	Send(*ReadRequest) error
	CloseAndRecv() (*ReadResponse, error)
	grpc.ClientStream
}

type readerServiceReadMulticastClient struct {
	grpc.ClientStream
}

func (x *readerServiceReadMulticastClient) Send(m *ReadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *readerServiceReadMulticastClient) CloseAndRecv() (*ReadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ReadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *readerServiceClient) ReadQuorumCallFuture(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallFuture", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readerServiceClient) ReadCorrectable(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dev.ReaderService/ReadCorrectable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readerServiceClient) ReadCorrectableStream(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (ReaderService_ReadCorrectableStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ReaderService_serviceDesc.Streams[1], "/dev.ReaderService/ReadCorrectableStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &readerServiceReadCorrectableStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ReaderService_ReadCorrectableStreamClient interface {
	Recv() (*ReadResponse, error)
	grpc.ClientStream
}

type readerServiceReadCorrectableStreamClient struct {
	grpc.ClientStream
}

func (x *readerServiceReadCorrectableStreamClient) Recv() (*ReadResponse, error) {
	m := new(ReadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReaderServiceServer is the server API for ReaderService service.
type ReaderServiceServer interface {
	ReadGrpc(context.Context, *ReadRequest) (*ReadResponse, error)
	ReadQuorumCall(context.Context, *ReadRequest) (*ReadResponse, error)
	ReadQuorumCallPerNodeArg(context.Context, *ReadRequest) (*ReadResponse, error)
	ReadQuorumCallQFWithRequestArg(context.Context, *ReadRequest) (*ReadResponse, error)
	ReadQuorumCallCustomReturnType(context.Context, *ReadRequest) (*ReadResponse, error)
	// ReadQuorumCallCombo does it all. Comment testing.
	ReadQuorumCallCombo(context.Context, *ReadRequest) (*ReadResponse, error)
	// ReadMulticast is testing a comment.
	ReadMulticast(ReaderService_ReadMulticastServer) error
	ReadQuorumCallFuture(context.Context, *ReadRequest) (*ReadResponse, error)
	ReadCorrectable(context.Context, *ReadRequest) (*ReadResponse, error)
	ReadCorrectableStream(*ReadRequest, ReaderService_ReadCorrectableStreamServer) error
}

// UnimplementedReaderServiceServer can be embedded to have forward compatible implementations.
type UnimplementedReaderServiceServer struct {
}

func (*UnimplementedReaderServiceServer) ReadGrpc(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadGrpc not implemented")
}
func (*UnimplementedReaderServiceServer) ReadQuorumCall(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadQuorumCall not implemented")
}
func (*UnimplementedReaderServiceServer) ReadQuorumCallPerNodeArg(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadQuorumCallPerNodeArg not implemented")
}
func (*UnimplementedReaderServiceServer) ReadQuorumCallQFWithRequestArg(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadQuorumCallQFWithRequestArg not implemented")
}
func (*UnimplementedReaderServiceServer) ReadQuorumCallCustomReturnType(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadQuorumCallCustomReturnType not implemented")
}
func (*UnimplementedReaderServiceServer) ReadQuorumCallCombo(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadQuorumCallCombo not implemented")
}
func (*UnimplementedReaderServiceServer) ReadMulticast(ReaderService_ReadMulticastServer) error {
	return status.Errorf(codes.Unimplemented, "method ReadMulticast not implemented")
}
func (*UnimplementedReaderServiceServer) ReadQuorumCallFuture(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadQuorumCallFuture not implemented")
}
func (*UnimplementedReaderServiceServer) ReadCorrectable(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadCorrectable not implemented")
}
func (*UnimplementedReaderServiceServer) ReadCorrectableStream(*ReadRequest, ReaderService_ReadCorrectableStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ReadCorrectableStream not implemented")
}

func RegisterReaderServiceServer(s *grpc.Server, srv ReaderServiceServer) {
	s.RegisterService(&_ReaderService_serviceDesc, srv)
}

func _ReaderService_ReadGrpc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReaderServiceServer).ReadGrpc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ReaderService/ReadGrpc",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReaderServiceServer).ReadGrpc(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReaderService_ReadQuorumCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReaderServiceServer).ReadQuorumCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ReaderService/ReadQuorumCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReaderServiceServer).ReadQuorumCall(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReaderService_ReadQuorumCallPerNodeArg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReaderServiceServer).ReadQuorumCallPerNodeArg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ReaderService/ReadQuorumCallPerNodeArg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReaderServiceServer).ReadQuorumCallPerNodeArg(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReaderService_ReadQuorumCallQFWithRequestArg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReaderServiceServer).ReadQuorumCallQFWithRequestArg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ReaderService/ReadQuorumCallQFWithRequestArg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReaderServiceServer).ReadQuorumCallQFWithRequestArg(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReaderService_ReadQuorumCallCustomReturnType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReaderServiceServer).ReadQuorumCallCustomReturnType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ReaderService/ReadQuorumCallCustomReturnType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReaderServiceServer).ReadQuorumCallCustomReturnType(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReaderService_ReadQuorumCallCombo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReaderServiceServer).ReadQuorumCallCombo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ReaderService/ReadQuorumCallCombo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReaderServiceServer).ReadQuorumCallCombo(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReaderService_ReadMulticast_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ReaderServiceServer).ReadMulticast(&readerServiceReadMulticastServer{stream})
}

type ReaderService_ReadMulticastServer interface {
	SendAndClose(*ReadResponse) error
	Recv() (*ReadRequest, error)
	grpc.ServerStream
}

type readerServiceReadMulticastServer struct {
	grpc.ServerStream
}

func (x *readerServiceReadMulticastServer) SendAndClose(m *ReadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *readerServiceReadMulticastServer) Recv() (*ReadRequest, error) {
	m := new(ReadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ReaderService_ReadQuorumCallFuture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReaderServiceServer).ReadQuorumCallFuture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ReaderService/ReadQuorumCallFuture",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReaderServiceServer).ReadQuorumCallFuture(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReaderService_ReadCorrectable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReaderServiceServer).ReadCorrectable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ReaderService/ReadCorrectable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReaderServiceServer).ReadCorrectable(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReaderService_ReadCorrectableStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReaderServiceServer).ReadCorrectableStream(m, &readerServiceReadCorrectableStreamServer{stream})
}

type ReaderService_ReadCorrectableStreamServer interface {
	Send(*ReadResponse) error
	grpc.ServerStream
}

type readerServiceReadCorrectableStreamServer struct {
	grpc.ServerStream
}

func (x *readerServiceReadCorrectableStreamServer) Send(m *ReadResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _ReaderService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dev.ReaderService",
	HandlerType: (*ReaderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReadGrpc",
			Handler:    _ReaderService_ReadGrpc_Handler,
		},
		{
			MethodName: "ReadQuorumCall",
			Handler:    _ReaderService_ReadQuorumCall_Handler,
		},
		{
			MethodName: "ReadQuorumCallPerNodeArg",
			Handler:    _ReaderService_ReadQuorumCallPerNodeArg_Handler,
		},
		{
			MethodName: "ReadQuorumCallQFWithRequestArg",
			Handler:    _ReaderService_ReadQuorumCallQFWithRequestArg_Handler,
		},
		{
			MethodName: "ReadQuorumCallCustomReturnType",
			Handler:    _ReaderService_ReadQuorumCallCustomReturnType_Handler,
		},
		{
			MethodName: "ReadQuorumCallCombo",
			Handler:    _ReaderService_ReadQuorumCallCombo_Handler,
		},
		{
			MethodName: "ReadQuorumCallFuture",
			Handler:    _ReaderService_ReadQuorumCallFuture_Handler,
		},
		{
			MethodName: "ReadCorrectable",
			Handler:    _ReaderService_ReadCorrectable_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReadMulticast",
			Handler:       _ReaderService_ReadMulticast_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ReadCorrectableStream",
			Handler:       _ReaderService_ReadCorrectableStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "zorums.proto",
}
