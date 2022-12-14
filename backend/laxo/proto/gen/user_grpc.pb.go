// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: user.proto

package gen

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

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	GetNotificationUpdate(ctx context.Context, in *NotificationUpdateRequest, opts ...grpc.CallOption) (UserService_GetNotificationUpdateClient, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) GetNotificationUpdate(ctx context.Context, in *NotificationUpdateRequest, opts ...grpc.CallOption) (UserService_GetNotificationUpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserService_ServiceDesc.Streams[0], "/user.UserService/GetNotificationUpdate", opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceGetNotificationUpdateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserService_GetNotificationUpdateClient interface {
	Recv() (*NotificationUpdateReply, error)
	grpc.ClientStream
}

type userServiceGetNotificationUpdateClient struct {
	grpc.ClientStream
}

func (x *userServiceGetNotificationUpdateClient) Recv() (*NotificationUpdateReply, error) {
	m := new(NotificationUpdateReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	GetNotificationUpdate(*NotificationUpdateRequest, UserService_GetNotificationUpdateServer) error
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) GetNotificationUpdate(*NotificationUpdateRequest, UserService_GetNotificationUpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method GetNotificationUpdate not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_GetNotificationUpdate_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(NotificationUpdateRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServiceServer).GetNotificationUpdate(m, &userServiceGetNotificationUpdateServer{stream})
}

type UserService_GetNotificationUpdateServer interface {
	Send(*NotificationUpdateReply) error
	grpc.ServerStream
}

type userServiceGetNotificationUpdateServer struct {
	grpc.ServerStream
}

func (x *userServiceGetNotificationUpdateServer) Send(m *NotificationUpdateReply) error {
	return x.ServerStream.SendMsg(m)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetNotificationUpdate",
			Handler:       _UserService_GetNotificationUpdate_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "user.proto",
}
