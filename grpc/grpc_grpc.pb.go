// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.0--rc3
// source: grpc/grpc.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GrpcClient is the client API for Grpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GrpcClient interface {
	SendBatch(ctx context.Context, in *Batch, opts ...grpc.CallOption) (*BatchResponse, error)
}

type grpcClient struct {
	cc grpc.ClientConnInterface
}

func NewGrpcClient(cc grpc.ClientConnInterface) GrpcClient {
	return &grpcClient{cc}
}

func (c *grpcClient) SendBatch(ctx context.Context, in *Batch, opts ...grpc.CallOption) (*BatchResponse, error) {
	out := new(BatchResponse)
	err := c.cc.Invoke(ctx, "/grpc.Grpc/SendBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrpcServer is the server API for Grpc service.
// All implementations must embed UnimplementedGrpcServer
// for forward compatibility
type GrpcServer interface {
	SendBatch(context.Context, *Batch) (*BatchResponse, error)
	mustEmbedUnimplementedGrpcServer()
}

// UnimplementedGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedGrpcServer struct {
}

func (UnimplementedGrpcServer) SendBatch(context.Context, *Batch) (*BatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBatch not implemented")
}
func (UnimplementedGrpcServer) mustEmbedUnimplementedGrpcServer() {}

// UnsafeGrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GrpcServer will
// result in compilation errors.
type UnsafeGrpcServer interface {
	mustEmbedUnimplementedGrpcServer()
}

func RegisterGrpcServer(s grpc.ServiceRegistrar, srv GrpcServer) {
	s.RegisterService(&Grpc_ServiceDesc, srv)
}

func _Grpc_SendBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Batch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).SendBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Grpc/SendBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).SendBatch(ctx, req.(*Batch))
	}
	return interceptor(ctx, in, info, handler)
}

// Grpc_ServiceDesc is the grpc.ServiceDesc for Grpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Grpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Grpc",
	HandlerType: (*GrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendBatch",
			Handler:    _Grpc_SendBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/grpc.proto",
}
