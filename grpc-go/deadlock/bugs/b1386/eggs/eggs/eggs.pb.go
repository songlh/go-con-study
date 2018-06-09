// Code generated by protoc-gen-go. DO NOT EDIT.
// source: eggs.proto

/*
Package main is a generated protocol buffer package.

It is generated from these files:
	eggs.proto

It has these top-level messages:
*/
package eggs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/wrappers"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Eggs service

type EggsClient interface {
	Echo(ctx context.Context, in *google_protobuf.StringValue, opts ...grpc.CallOption) (*google_protobuf.StringValue, error)
}

type eggsClient struct {
	cc *grpc.ClientConn
}

func NewEggsClient(cc *grpc.ClientConn) EggsClient {
	return &eggsClient{cc}
}

func (c *eggsClient) Echo(ctx context.Context, in *google_protobuf.StringValue, opts ...grpc.CallOption) (*google_protobuf.StringValue, error) {
	out := new(google_protobuf.StringValue)
	err := grpc.Invoke(ctx, "/eggs.Eggs/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Eggs service

type EggsServer interface {
	Echo(context.Context, *google_protobuf.StringValue) (*google_protobuf.StringValue, error)
}

func RegisterEggsServer(s *grpc.Server, srv EggsServer) {
	s.RegisterService(&_Eggs_serviceDesc, srv)
}

func _Eggs_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EggsServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eggs.Eggs/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EggsServer).Echo(ctx, req.(*google_protobuf.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

var _Eggs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eggs.Eggs",
	HandlerType: (*EggsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Eggs_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "eggs.proto",
}

func init() { proto.RegisterFile("eggs.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 113 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x4d, 0x4f, 0x2f,
	0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0xa5, 0xe4, 0xd2, 0xf3, 0xf3, 0xd3,
	0x73, 0x52, 0xf5, 0xc1, 0x62, 0x49, 0xa5, 0x69, 0xfa, 0xe5, 0x45, 0x89, 0x05, 0x05, 0xa9, 0x45,
	0x50, 0x55, 0x46, 0x5e, 0x5c, 0x2c, 0xae, 0x40, 0x75, 0x42, 0x4e, 0x40, 0x3a, 0x39, 0x23, 0x5f,
	0x48, 0x46, 0x0f, 0xa2, 0x41, 0x0f, 0xa6, 0x41, 0x2f, 0xb8, 0xa4, 0x28, 0x33, 0x2f, 0x3d, 0x2c,
	0x31, 0xa7, 0x34, 0x55, 0x0a, 0xaf, 0xac, 0x13, 0x5b, 0x14, 0x4b, 0x6e, 0x62, 0x66, 0x5e, 0x12,
	0x1b, 0x58, 0xd6, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xa0, 0x1b, 0x9a, 0x96, 0x8e, 0x00, 0x00,
	0x00,
}