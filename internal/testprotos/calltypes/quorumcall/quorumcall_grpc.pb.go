// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package quorumcall

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// QuorumCallClient is the client API for QuorumCall service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuorumCallClient interface {
	QuorumCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type quorumCallClient struct {
	cc grpc.ClientConnInterface
}

func NewQuorumCallClient(cc grpc.ClientConnInterface) QuorumCallClient {
	return &quorumCallClient{cc}
}

func (c *quorumCallClient) QuorumCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/gorums.testprotos.calltypes.quorumcall.QuorumCall/QuorumCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuorumCallServer is the server API for QuorumCall service.
type QuorumCallServer interface {
	QuorumCall(context.Context, *Request) (*Response, error)
}

// UnimplementedQuorumCallServer can be embedded to have forward compatible implementations.
type UnimplementedQuorumCallServer struct {
}

func (*UnimplementedQuorumCallServer) QuorumCall(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuorumCall not implemented")
}

func RegisterQuorumCallServer(s *grpc.Server, srv QuorumCallServer) {
	s.RegisterService(&_QuorumCall_serviceDesc, srv)
}

func _QuorumCall_QuorumCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuorumCallServer).QuorumCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gorums.testprotos.calltypes.quorumcall.QuorumCall/QuorumCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuorumCallServer).QuorumCall(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _QuorumCall_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gorums.testprotos.calltypes.quorumcall.QuorumCall",
	HandlerType: (*QuorumCallServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QuorumCall",
			Handler:    _QuorumCall_QuorumCall_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/testprotos/calltypes/quorumcall/quorumcall.proto",
}
