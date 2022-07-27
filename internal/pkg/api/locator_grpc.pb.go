// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CourierLocatorClient is the client API for CourierLocator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CourierLocatorClient interface {
	// use this procedure to push the latest locations of drivers.
	// guarantees to return an ack message for each location push but not in the same order.
	PushGpsPoints(ctx context.Context, opts ...grpc.CallOption) (CourierLocator_PushGpsPointsClient, error)
	// use this procedure to find all the drivers within radius of a given point.
	// there obviously are some limits on the radius.
	FindNearbyCouriers(ctx context.Context, in *FindNearbyCouriersQuery, opts ...grpc.CallOption) (*CourierList, error)
}

type courierLocatorClient struct {
	cc grpc.ClientConnInterface
}

func NewCourierLocatorClient(cc grpc.ClientConnInterface) CourierLocatorClient {
	return &courierLocatorClient{cc}
}

func (c *courierLocatorClient) PushGpsPoints(ctx context.Context, opts ...grpc.CallOption) (CourierLocator_PushGpsPointsClient, error) {
	stream, err := c.cc.NewStream(ctx, &CourierLocator_ServiceDesc.Streams[0], "/dlocator.CourierLocator/PushGpsPoints", opts...)
	if err != nil {
		return nil, err
	}
	x := &courierLocatorPushGpsPointsClient{stream}
	return x, nil
}

type CourierLocator_PushGpsPointsClient interface {
	Send(*CourierGpsPoint) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type courierLocatorPushGpsPointsClient struct {
	grpc.ClientStream
}

func (x *courierLocatorPushGpsPointsClient) Send(m *CourierGpsPoint) error {
	return x.ClientStream.SendMsg(m)
}

func (x *courierLocatorPushGpsPointsClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *courierLocatorClient) FindNearbyCouriers(ctx context.Context, in *FindNearbyCouriersQuery, opts ...grpc.CallOption) (*CourierList, error) {
	out := new(CourierList)
	err := c.cc.Invoke(ctx, "/dlocator.CourierLocator/FindNearbyCouriers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CourierLocatorServer is the server API for CourierLocator service.
// All implementations must embed UnimplementedCourierLocatorServer
// for forward compatibility
type CourierLocatorServer interface {
	// use this procedure to push the latest locations of drivers.
	// guarantees to return an ack message for each location push but not in the same order.
	PushGpsPoints(CourierLocator_PushGpsPointsServer) error
	// use this procedure to find all the drivers within radius of a given point.
	// there obviously are some limits on the radius.
	FindNearbyCouriers(context.Context, *FindNearbyCouriersQuery) (*CourierList, error)
	mustEmbedUnimplementedCourierLocatorServer()
}

// UnimplementedCourierLocatorServer must be embedded to have forward compatible implementations.
type UnimplementedCourierLocatorServer struct {
}

func (UnimplementedCourierLocatorServer) PushGpsPoints(CourierLocator_PushGpsPointsServer) error {
	return status.Errorf(codes.Unimplemented, "method PushGpsPoints not implemented")
}
func (UnimplementedCourierLocatorServer) FindNearbyCouriers(context.Context, *FindNearbyCouriersQuery) (*CourierList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindNearbyCouriers not implemented")
}
func (UnimplementedCourierLocatorServer) mustEmbedUnimplementedCourierLocatorServer() {}

// UnsafeCourierLocatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CourierLocatorServer will
// result in compilation errors.
type UnsafeCourierLocatorServer interface {
	mustEmbedUnimplementedCourierLocatorServer()
}

func RegisterCourierLocatorServer(s grpc.ServiceRegistrar, srv CourierLocatorServer) {
	s.RegisterService(&CourierLocator_ServiceDesc, srv)
}

func _CourierLocator_PushGpsPoints_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CourierLocatorServer).PushGpsPoints(&courierLocatorPushGpsPointsServer{stream})
}

type CourierLocator_PushGpsPointsServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*CourierGpsPoint, error)
	grpc.ServerStream
}

type courierLocatorPushGpsPointsServer struct {
	grpc.ServerStream
}

func (x *courierLocatorPushGpsPointsServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *courierLocatorPushGpsPointsServer) Recv() (*CourierGpsPoint, error) {
	m := new(CourierGpsPoint)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _CourierLocator_FindNearbyCouriers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindNearbyCouriersQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourierLocatorServer).FindNearbyCouriers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dlocator.CourierLocator/FindNearbyCouriers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourierLocatorServer).FindNearbyCouriers(ctx, req.(*FindNearbyCouriersQuery))
	}
	return interceptor(ctx, in, info, handler)
}

// CourierLocator_ServiceDesc is the grpc.ServiceDesc for CourierLocator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CourierLocator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dlocator.CourierLocator",
	HandlerType: (*CourierLocatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindNearbyCouriers",
			Handler:    _CourierLocator_FindNearbyCouriers_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushGpsPoints",
			Handler:       _CourierLocator_PushGpsPoints_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "api/proto/v1/locator.proto",
}