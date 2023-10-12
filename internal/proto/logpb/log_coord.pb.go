// Code generated by protoc-gen-go. DO NOT EDIT.
// source: log_coord.proto

package logpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	commonpb "github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	milvuspb "github.com/milvus-io/milvus-proto/go-api/v2/milvuspb"
	msgpb "github.com/milvus-io/milvus-proto/go-api/v2/msgpb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PChannelState int32

const (
	PChannelState_Unused   PChannelState = 0
	PChannelState_Waitting PChannelState = 1
	PChannelState_Watching PChannelState = 2
)

var PChannelState_name = map[int32]string{
	0: "Unused",
	1: "Waitting",
	2: "Watching",
}

var PChannelState_value = map[string]int32{
	"Unused":   0,
	"Waitting": 1,
	"Watching": 2,
}

func (x PChannelState) String() string {
	return proto.EnumName(PChannelState_name, int32(x))
}

func (PChannelState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{0}
}

type PChannelInfo struct {
	Name                 string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	State                PChannelState `protobuf:"varint,2,opt,name=state,proto3,enum=milvus.proto.log.PChannelState" json:"state,omitempty"`
	NodeID               int64         `protobuf:"varint,3,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *PChannelInfo) Reset()         { *m = PChannelInfo{} }
func (m *PChannelInfo) String() string { return proto.CompactTextString(m) }
func (*PChannelInfo) ProtoMessage()    {}
func (*PChannelInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{0}
}

func (m *PChannelInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PChannelInfo.Unmarshal(m, b)
}
func (m *PChannelInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PChannelInfo.Marshal(b, m, deterministic)
}
func (m *PChannelInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PChannelInfo.Merge(m, src)
}
func (m *PChannelInfo) XXX_Size() int {
	return xxx_messageInfo_PChannelInfo.Size(m)
}
func (m *PChannelInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PChannelInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PChannelInfo proto.InternalMessageInfo

func (m *PChannelInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PChannelInfo) GetState() PChannelState {
	if m != nil {
		return m.State
	}
	return PChannelState_Unused
}

func (m *PChannelInfo) GetNodeID() int64 {
	if m != nil {
		return m.NodeID
	}
	return 0
}

type LogNodeInfo struct {
	NodeID               int64    `protobuf:"varint,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogNodeInfo) Reset()         { *m = LogNodeInfo{} }
func (m *LogNodeInfo) String() string { return proto.CompactTextString(m) }
func (*LogNodeInfo) ProtoMessage()    {}
func (*LogNodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{1}
}

func (m *LogNodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogNodeInfo.Unmarshal(m, b)
}
func (m *LogNodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogNodeInfo.Marshal(b, m, deterministic)
}
func (m *LogNodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogNodeInfo.Merge(m, src)
}
func (m *LogNodeInfo) XXX_Size() int {
	return xxx_messageInfo_LogNodeInfo.Size(m)
}
func (m *LogNodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_LogNodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_LogNodeInfo proto.InternalMessageInfo

func (m *LogNodeInfo) GetNodeID() int64 {
	if m != nil {
		return m.NodeID
	}
	return 0
}

func (m *LogNodeInfo) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type WatchChannelRequest struct {
	Base                 *commonpb.MsgBase `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	PChannel             string            `protobuf:"bytes,2,opt,name=pChannel,proto3" json:"pChannel,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *WatchChannelRequest) Reset()         { *m = WatchChannelRequest{} }
func (m *WatchChannelRequest) String() string { return proto.CompactTextString(m) }
func (*WatchChannelRequest) ProtoMessage()    {}
func (*WatchChannelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{2}
}

func (m *WatchChannelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WatchChannelRequest.Unmarshal(m, b)
}
func (m *WatchChannelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WatchChannelRequest.Marshal(b, m, deterministic)
}
func (m *WatchChannelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WatchChannelRequest.Merge(m, src)
}
func (m *WatchChannelRequest) XXX_Size() int {
	return xxx_messageInfo_WatchChannelRequest.Size(m)
}
func (m *WatchChannelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WatchChannelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WatchChannelRequest proto.InternalMessageInfo

func (m *WatchChannelRequest) GetBase() *commonpb.MsgBase {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *WatchChannelRequest) GetPChannel() string {
	if m != nil {
		return m.PChannel
	}
	return ""
}

type UnwatchChannelRequest struct {
	Base                 *commonpb.MsgBase `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	PChannel             string            `protobuf:"bytes,2,opt,name=pChannel,proto3" json:"pChannel,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UnwatchChannelRequest) Reset()         { *m = UnwatchChannelRequest{} }
func (m *UnwatchChannelRequest) String() string { return proto.CompactTextString(m) }
func (*UnwatchChannelRequest) ProtoMessage()    {}
func (*UnwatchChannelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{3}
}

func (m *UnwatchChannelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnwatchChannelRequest.Unmarshal(m, b)
}
func (m *UnwatchChannelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnwatchChannelRequest.Marshal(b, m, deterministic)
}
func (m *UnwatchChannelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnwatchChannelRequest.Merge(m, src)
}
func (m *UnwatchChannelRequest) XXX_Size() int {
	return xxx_messageInfo_UnwatchChannelRequest.Size(m)
}
func (m *UnwatchChannelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UnwatchChannelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UnwatchChannelRequest proto.InternalMessageInfo

func (m *UnwatchChannelRequest) GetBase() *commonpb.MsgBase {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *UnwatchChannelRequest) GetPChannel() string {
	if m != nil {
		return m.PChannel
	}
	return ""
}

type InsertRequest struct {
	Base                 *commonpb.MsgBase      `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Msgs                 []*msgpb.InsertRequest `protobuf:"bytes,2,rep,name=msgs,proto3" json:"msgs,omitempty"`
	PChannels            []string               `protobuf:"bytes,3,rep,name=pChannels,proto3" json:"pChannels,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *InsertRequest) Reset()         { *m = InsertRequest{} }
func (m *InsertRequest) String() string { return proto.CompactTextString(m) }
func (*InsertRequest) ProtoMessage()    {}
func (*InsertRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{4}
}

func (m *InsertRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InsertRequest.Unmarshal(m, b)
}
func (m *InsertRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InsertRequest.Marshal(b, m, deterministic)
}
func (m *InsertRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InsertRequest.Merge(m, src)
}
func (m *InsertRequest) XXX_Size() int {
	return xxx_messageInfo_InsertRequest.Size(m)
}
func (m *InsertRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InsertRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InsertRequest proto.InternalMessageInfo

func (m *InsertRequest) GetBase() *commonpb.MsgBase {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *InsertRequest) GetMsgs() []*msgpb.InsertRequest {
	if m != nil {
		return m.Msgs
	}
	return nil
}

func (m *InsertRequest) GetPChannels() []string {
	if m != nil {
		return m.PChannels
	}
	return nil
}

func init() {
	proto.RegisterEnum("milvus.proto.log.PChannelState", PChannelState_name, PChannelState_value)
	proto.RegisterType((*PChannelInfo)(nil), "milvus.proto.log.PChannelInfo")
	proto.RegisterType((*LogNodeInfo)(nil), "milvus.proto.log.LogNodeInfo")
	proto.RegisterType((*WatchChannelRequest)(nil), "milvus.proto.log.WatchChannelRequest")
	proto.RegisterType((*UnwatchChannelRequest)(nil), "milvus.proto.log.UnwatchChannelRequest")
	proto.RegisterType((*InsertRequest)(nil), "milvus.proto.log.InsertRequest")
}

func init() { proto.RegisterFile("log_coord.proto", fileDescriptor_5674412121b98d7f) }

var fileDescriptor_5674412121b98d7f = []byte{
	// 455 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0x61, 0x8b, 0xd3, 0x40,
	0x14, 0xbc, 0x34, 0xb5, 0x77, 0x79, 0xcd, 0x9d, 0xe5, 0x89, 0x12, 0xe2, 0x81, 0x25, 0x28, 0x06,
	0xc1, 0xf4, 0xe8, 0x21, 0x7e, 0x14, 0xee, 0x04, 0x2d, 0xa8, 0x48, 0xf4, 0x38, 0xf1, 0x8b, 0x6c,
	0x93, 0x75, 0x1b, 0x48, 0x76, 0x7b, 0xd9, 0x8d, 0xfe, 0x11, 0xff, 0xa8, 0xff, 0x40, 0x36, 0x9b,
	0x7a, 0x4d, 0x1a, 0x28, 0x08, 0x7e, 0xdb, 0xd9, 0x9d, 0xcc, 0x3c, 0xde, 0x4c, 0xe0, 0x6e, 0x2e,
	0xd8, 0xb7, 0x44, 0x88, 0x32, 0x8d, 0xd6, 0xa5, 0x50, 0x02, 0x27, 0x45, 0x96, 0xff, 0xa8, 0xa4,
	0x41, 0x51, 0x2e, 0x98, 0xef, 0x26, 0xa2, 0x28, 0x04, 0x37, 0x37, 0xbe, 0xbb, 0xfd, 0xee, 0x3b,
	0x85, 0x64, 0xe6, 0x18, 0xdc, 0x80, 0xfb, 0xf1, 0x72, 0x45, 0x38, 0xa7, 0xf9, 0x82, 0x7f, 0x17,
	0x88, 0x30, 0xe4, 0xa4, 0xa0, 0x9e, 0x35, 0xb5, 0x42, 0x27, 0xae, 0xcf, 0xf8, 0x02, 0xee, 0x48,
	0x45, 0x14, 0xf5, 0x06, 0x53, 0x2b, 0x3c, 0x99, 0x3f, 0x8a, 0xba, 0x66, 0xd1, 0x46, 0xe2, 0x93,
	0xa6, 0xc5, 0x86, 0x8d, 0x0f, 0x60, 0xc4, 0x45, 0x4a, 0x17, 0xaf, 0x3d, 0x7b, 0x6a, 0x85, 0x76,
	0xdc, 0xa0, 0xe0, 0x15, 0x8c, 0xdf, 0x09, 0xf6, 0x41, 0x03, 0xed, 0x78, 0x4b, 0xb3, 0xb6, 0x69,
	0xe8, 0xc1, 0x21, 0x49, 0xd3, 0x92, 0x4a, 0x59, 0xfb, 0x3a, 0xf1, 0x06, 0x06, 0x09, 0xdc, 0xbb,
	0x26, 0x2a, 0x59, 0x35, 0xa6, 0x31, 0xbd, 0xa9, 0xa8, 0x54, 0x78, 0x06, 0xc3, 0x25, 0x91, 0x66,
	0xf4, 0xf1, 0xfc, 0xb4, 0x3d, 0x65, 0xb3, 0x8d, 0xf7, 0x92, 0x5d, 0x10, 0x49, 0xe3, 0x9a, 0x89,
	0x3e, 0x1c, 0xad, 0x1b, 0x91, 0xc6, 0xe3, 0x2f, 0x0e, 0x28, 0xdc, 0xbf, 0xe2, 0x3f, 0xff, 0xbb,
	0xcd, 0x2f, 0x0b, 0x8e, 0x17, 0x5c, 0xd2, 0x52, 0xfd, 0xbb, 0xfe, 0x39, 0x0c, 0x0b, 0xc9, 0xf4,
	0x9a, 0xec, 0x70, 0xdc, 0x8d, 0x47, 0x47, 0xdd, 0x32, 0x88, 0x6b, 0x32, 0x9e, 0x82, 0xb3, 0x19,
	0x42, 0x7a, 0xf6, 0xd4, 0x0e, 0x9d, 0xf8, 0xf6, 0xe2, 0xd9, 0x4b, 0x38, 0x6e, 0x65, 0x8a, 0x00,
	0xa3, 0x2b, 0x5e, 0x49, 0x9a, 0x4e, 0x0e, 0xd0, 0x85, 0xa3, 0x6b, 0x92, 0x29, 0x95, 0x71, 0x36,
	0xb1, 0x0c, 0x52, 0xc9, 0x4a, 0xa3, 0xc1, 0xfc, 0xf7, 0x00, 0x0e, 0x9b, 0x74, 0x31, 0x07, 0x7c,
	0x43, 0xd5, 0xa5, 0x28, 0xd6, 0x82, 0x53, 0xae, 0x6a, 0x21, 0x89, 0x51, 0x67, 0x3e, 0x03, 0x76,
	0x89, 0xcd, 0xb8, 0xfe, 0xe3, 0x5e, 0x7e, 0x87, 0x1c, 0x1c, 0xe0, 0x67, 0x70, 0xb7, 0x5b, 0x81,
	0x4f, 0x76, 0x6b, 0xda, 0xd3, 0x1a, 0xff, 0x61, 0xef, 0x82, 0xb5, 0x6a, 0xa5, 0x55, 0xbf, 0xc0,
	0x49, 0xbb, 0x06, 0xf8, 0x74, 0x57, 0xb7, 0xb7, 0x28, 0xfb, 0x94, 0xdf, 0xc2, 0xc8, 0xe4, 0x82,
	0x3d, 0x3f, 0x54, 0x2b, 0xb1, 0x3d, 0x4a, 0x17, 0xf3, 0xaf, 0x67, 0x2c, 0x53, 0xab, 0x6a, 0xa9,
	0x5f, 0x66, 0x86, 0xfa, 0x3c, 0x13, 0xcd, 0x69, 0x96, 0x71, 0x45, 0x4b, 0x4e, 0xf2, 0x59, 0xfd,
	0xf5, 0x2c, 0x17, 0x6c, 0xbd, 0x5c, 0x8e, 0x6a, 0x70, 0xfe, 0x27, 0x00, 0x00, 0xff, 0xff, 0xc0,
	0x53, 0xa0, 0x86, 0x4a, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LogNodeClient is the client API for LogNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LogNodeClient interface {
	GetComponentStates(ctx context.Context, in *milvuspb.GetComponentStatesRequest, opts ...grpc.CallOption) (*milvuspb.ComponentStates, error)
	WatchChannel(ctx context.Context, in *WatchChannelRequest, opts ...grpc.CallOption) (*commonpb.Status, error)
	UnwatchChannel(ctx context.Context, in *UnwatchChannelRequest, opts ...grpc.CallOption) (*commonpb.Status, error)
	Insert(ctx context.Context, in *InsertRequest, opts ...grpc.CallOption) (*commonpb.Status, error)
}

type logNodeClient struct {
	cc *grpc.ClientConn
}

func NewLogNodeClient(cc *grpc.ClientConn) LogNodeClient {
	return &logNodeClient{cc}
}

func (c *logNodeClient) GetComponentStates(ctx context.Context, in *milvuspb.GetComponentStatesRequest, opts ...grpc.CallOption) (*milvuspb.ComponentStates, error) {
	out := new(milvuspb.ComponentStates)
	err := c.cc.Invoke(ctx, "/milvus.proto.log.LogNode/GetComponentStates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logNodeClient) WatchChannel(ctx context.Context, in *WatchChannelRequest, opts ...grpc.CallOption) (*commonpb.Status, error) {
	out := new(commonpb.Status)
	err := c.cc.Invoke(ctx, "/milvus.proto.log.LogNode/WatchChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logNodeClient) UnwatchChannel(ctx context.Context, in *UnwatchChannelRequest, opts ...grpc.CallOption) (*commonpb.Status, error) {
	out := new(commonpb.Status)
	err := c.cc.Invoke(ctx, "/milvus.proto.log.LogNode/UnwatchChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logNodeClient) Insert(ctx context.Context, in *InsertRequest, opts ...grpc.CallOption) (*commonpb.Status, error) {
	out := new(commonpb.Status)
	err := c.cc.Invoke(ctx, "/milvus.proto.log.LogNode/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogNodeServer is the server API for LogNode service.
type LogNodeServer interface {
	GetComponentStates(context.Context, *milvuspb.GetComponentStatesRequest) (*milvuspb.ComponentStates, error)
	WatchChannel(context.Context, *WatchChannelRequest) (*commonpb.Status, error)
	UnwatchChannel(context.Context, *UnwatchChannelRequest) (*commonpb.Status, error)
	Insert(context.Context, *InsertRequest) (*commonpb.Status, error)
}

// UnimplementedLogNodeServer can be embedded to have forward compatible implementations.
type UnimplementedLogNodeServer struct {
}

func (*UnimplementedLogNodeServer) GetComponentStates(ctx context.Context, req *milvuspb.GetComponentStatesRequest) (*milvuspb.ComponentStates, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetComponentStates not implemented")
}
func (*UnimplementedLogNodeServer) WatchChannel(ctx context.Context, req *WatchChannelRequest) (*commonpb.Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WatchChannel not implemented")
}
func (*UnimplementedLogNodeServer) UnwatchChannel(ctx context.Context, req *UnwatchChannelRequest) (*commonpb.Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnwatchChannel not implemented")
}
func (*UnimplementedLogNodeServer) Insert(ctx context.Context, req *InsertRequest) (*commonpb.Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}

func RegisterLogNodeServer(s *grpc.Server, srv LogNodeServer) {
	s.RegisterService(&_LogNode_serviceDesc, srv)
}

func _LogNode_GetComponentStates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(milvuspb.GetComponentStatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogNodeServer).GetComponentStates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/milvus.proto.log.LogNode/GetComponentStates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogNodeServer).GetComponentStates(ctx, req.(*milvuspb.GetComponentStatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogNode_WatchChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WatchChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogNodeServer).WatchChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/milvus.proto.log.LogNode/WatchChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogNodeServer).WatchChannel(ctx, req.(*WatchChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogNode_UnwatchChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnwatchChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogNodeServer).UnwatchChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/milvus.proto.log.LogNode/UnwatchChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogNodeServer).UnwatchChannel(ctx, req.(*UnwatchChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogNode_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogNodeServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/milvus.proto.log.LogNode/Insert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogNodeServer).Insert(ctx, req.(*InsertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LogNode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "milvus.proto.log.LogNode",
	HandlerType: (*LogNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetComponentStates",
			Handler:    _LogNode_GetComponentStates_Handler,
		},
		{
			MethodName: "WatchChannel",
			Handler:    _LogNode_WatchChannel_Handler,
		},
		{
			MethodName: "UnwatchChannel",
			Handler:    _LogNode_UnwatchChannel_Handler,
		},
		{
			MethodName: "Insert",
			Handler:    _LogNode_Insert_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "log_coord.proto",
}
