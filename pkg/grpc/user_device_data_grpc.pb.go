// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: pkg/grpc/user_device_data.proto

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

// UseDeviceDataServiceClient is the client API for UseDeviceDataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UseDeviceDataServiceClient interface {
	GetUserDeviceData(ctx context.Context, in *UserDeviceDataRequest, opts ...grpc.CallOption) (*UserDeviceDataResponse, error)
}

type useDeviceDataServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUseDeviceDataServiceClient(cc grpc.ClientConnInterface) UseDeviceDataServiceClient {
	return &useDeviceDataServiceClient{cc}
}

func (c *useDeviceDataServiceClient) GetUserDeviceData(ctx context.Context, in *UserDeviceDataRequest, opts ...grpc.CallOption) (*UserDeviceDataResponse, error) {
	out := new(UserDeviceDataResponse)
	err := c.cc.Invoke(ctx, "/grpc.UseDeviceDataService/GetUserDeviceData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UseDeviceDataServiceServer is the server API for UseDeviceDataService service.
// All implementations must embed UnimplementedUseDeviceDataServiceServer
// for forward compatibility
type UseDeviceDataServiceServer interface {
	GetUserDeviceData(context.Context, *UserDeviceDataRequest) (*UserDeviceDataResponse, error)
	mustEmbedUnimplementedUseDeviceDataServiceServer()
}

// UnimplementedUseDeviceDataServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUseDeviceDataServiceServer struct {
}

func (UnimplementedUseDeviceDataServiceServer) GetUserDeviceData(context.Context, *UserDeviceDataRequest) (*UserDeviceDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserDeviceData not implemented")
}
func (UnimplementedUseDeviceDataServiceServer) mustEmbedUnimplementedUseDeviceDataServiceServer() {}

// UnsafeUseDeviceDataServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UseDeviceDataServiceServer will
// result in compilation errors.
type UnsafeUseDeviceDataServiceServer interface {
	mustEmbedUnimplementedUseDeviceDataServiceServer()
}

func RegisterUseDeviceDataServiceServer(s grpc.ServiceRegistrar, srv UseDeviceDataServiceServer) {
	s.RegisterService(&UseDeviceDataService_ServiceDesc, srv)
}

func _UseDeviceDataService_GetUserDeviceData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDeviceDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UseDeviceDataServiceServer).GetUserDeviceData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.UseDeviceDataService/GetUserDeviceData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UseDeviceDataServiceServer).GetUserDeviceData(ctx, req.(*UserDeviceDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UseDeviceDataService_ServiceDesc is the grpc.ServiceDesc for UseDeviceDataService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UseDeviceDataService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.UseDeviceDataService",
	HandlerType: (*UseDeviceDataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserDeviceData",
			Handler:    _UseDeviceDataService_GetUserDeviceData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/grpc/user_device_data.proto",
}
