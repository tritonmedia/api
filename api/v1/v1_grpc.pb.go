// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIClient interface {
	GetMedia(ctx context.Context, in *GetMediaRequest, opts ...grpc.CallOption) (*Media, error)
	CreateMedia(ctx context.Context, in *CreateMediaRequest, opts ...grpc.CallOption) (*Media, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) GetMedia(ctx context.Context, in *GetMediaRequest, opts ...grpc.CallOption) (*Media, error) {
	out := new(Media)
	err := c.cc.Invoke(ctx, "/api.v1.API/GetMedia", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) CreateMedia(ctx context.Context, in *CreateMediaRequest, opts ...grpc.CallOption) (*Media, error) {
	out := new(Media)
	err := c.cc.Invoke(ctx, "/api.v1.API/CreateMedia", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// APIServer is the server API for API service.
// All implementations must embed UnimplementedAPIServer
// for forward compatibility
type APIServer interface {
	GetMedia(context.Context, *GetMediaRequest) (*Media, error)
	CreateMedia(context.Context, *CreateMediaRequest) (*Media, error)
	mustEmbedUnimplementedAPIServer()
}

// UnimplementedAPIServer must be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (UnimplementedAPIServer) GetMedia(context.Context, *GetMediaRequest) (*Media, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMedia not implemented")
}
func (UnimplementedAPIServer) CreateMedia(context.Context, *CreateMediaRequest) (*Media, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMedia not implemented")
}
func (UnimplementedAPIServer) mustEmbedUnimplementedAPIServer() {}

// UnsafeAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIServer will
// result in compilation errors.
type UnsafeAPIServer interface {
	mustEmbedUnimplementedAPIServer()
}

func RegisterAPIServer(s grpc.ServiceRegistrar, srv APIServer) {
	s.RegisterService(&_API_serviceDesc, srv)
}

func _API_GetMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).GetMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.API/GetMedia",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).GetMedia(ctx, req.(*GetMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_CreateMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).CreateMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.API/CreateMedia",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).CreateMedia(ctx, req.(*CreateMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _API_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMedia",
			Handler:    _API_GetMedia_Handler,
		},
		{
			MethodName: "CreateMedia",
			Handler:    _API_CreateMedia_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1.proto",
}
