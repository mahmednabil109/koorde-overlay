// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package drpc

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

// DServerClient is the client API for DServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DServerClient interface {
	UpdatePointers(ctx context.Context, in *DPointers, opts ...grpc.CallOption) (*DEmpty, error)
	BroadcastBlock(ctx context.Context, in *DBlock, opts ...grpc.CallOption) (*DEmpty, error)
	Connect(ctx context.Context, in *DPeer, opts ...grpc.CallOption) (*DEmpty, error)
	GetLocalBlock(ctx context.Context, in *DEmpty, opts ...grpc.CallOption) (*DBlocks, error)
}

type dServerClient struct {
	cc grpc.ClientConnInterface
}

func NewDServerClient(cc grpc.ClientConnInterface) DServerClient {
	return &dServerClient{cc}
}

func (c *dServerClient) UpdatePointers(ctx context.Context, in *DPointers, opts ...grpc.CallOption) (*DEmpty, error) {
	out := new(DEmpty)
	err := c.cc.Invoke(ctx, "/rpc.DServer/UpdatePointers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dServerClient) BroadcastBlock(ctx context.Context, in *DBlock, opts ...grpc.CallOption) (*DEmpty, error) {
	out := new(DEmpty)
	err := c.cc.Invoke(ctx, "/rpc.DServer/BroadcastBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dServerClient) Connect(ctx context.Context, in *DPeer, opts ...grpc.CallOption) (*DEmpty, error) {
	out := new(DEmpty)
	err := c.cc.Invoke(ctx, "/rpc.DServer/Connect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dServerClient) GetLocalBlock(ctx context.Context, in *DEmpty, opts ...grpc.CallOption) (*DBlocks, error) {
	out := new(DBlocks)
	err := c.cc.Invoke(ctx, "/rpc.DServer/GetLocalBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DServerServer is the server API for DServer service.
// All implementations must embed UnimplementedDServerServer
// for forward compatibility
type DServerServer interface {
	UpdatePointers(context.Context, *DPointers) (*DEmpty, error)
	BroadcastBlock(context.Context, *DBlock) (*DEmpty, error)
	Connect(context.Context, *DPeer) (*DEmpty, error)
	GetLocalBlock(context.Context, *DEmpty) (*DBlocks, error)
	mustEmbedUnimplementedDServerServer()
}

// UnimplementedDServerServer must be embedded to have forward compatible implementations.
type UnimplementedDServerServer struct {
}

func (UnimplementedDServerServer) UpdatePointers(context.Context, *DPointers) (*DEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePointers not implemented")
}
func (UnimplementedDServerServer) BroadcastBlock(context.Context, *DBlock) (*DEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastBlock not implemented")
}
func (UnimplementedDServerServer) Connect(context.Context, *DPeer) (*DEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedDServerServer) GetLocalBlock(context.Context, *DEmpty) (*DBlocks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLocalBlock not implemented")
}
func (UnimplementedDServerServer) mustEmbedUnimplementedDServerServer() {}

// UnsafeDServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DServerServer will
// result in compilation errors.
type UnsafeDServerServer interface {
	mustEmbedUnimplementedDServerServer()
}

func RegisterDServerServer(s grpc.ServiceRegistrar, srv DServerServer) {
	s.RegisterService(&DServer_ServiceDesc, srv)
}

func _DServer_UpdatePointers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DPointers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DServerServer).UpdatePointers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.DServer/UpdatePointers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DServerServer).UpdatePointers(ctx, req.(*DPointers))
	}
	return interceptor(ctx, in, info, handler)
}

func _DServer_BroadcastBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DBlock)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DServerServer).BroadcastBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.DServer/BroadcastBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DServerServer).BroadcastBlock(ctx, req.(*DBlock))
	}
	return interceptor(ctx, in, info, handler)
}

func _DServer_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DPeer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DServerServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.DServer/Connect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DServerServer).Connect(ctx, req.(*DPeer))
	}
	return interceptor(ctx, in, info, handler)
}

func _DServer_GetLocalBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DServerServer).GetLocalBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.DServer/GetLocalBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DServerServer).GetLocalBlock(ctx, req.(*DEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

// DServer_ServiceDesc is the grpc.ServiceDesc for DServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.DServer",
	HandlerType: (*DServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdatePointers",
			Handler:    _DServer_UpdatePointers_Handler,
		},
		{
			MethodName: "BroadcastBlock",
			Handler:    _DServer_BroadcastBlock_Handler,
		},
		{
			MethodName: "Connect",
			Handler:    _DServer_Connect_Handler,
		},
		{
			MethodName: "GetLocalBlock",
			Handler:    _DServer_GetLocalBlock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "drpc/drpc.proto",
}
