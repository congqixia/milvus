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

type MessageType int32

const (
	MessageType_TIME_TICK         MessageType = 0
	MessageType_INSERT            MessageType = 1
	MessageType_DELETE            MessageType = 2
	MessageType_CREATE_COLLECTION MessageType = 3
	MessageType_DROP_COLLECTION   MessageType = 4
	MessageType_CREATE_PARTITION  MessageType = 5
	MessageType_DROP_PARTITION    MessageType = 6
)

var MessageType_name = map[int32]string{
	0: "TIME_TICK",
	1: "INSERT",
	2: "DELETE",
	3: "CREATE_COLLECTION",
	4: "DROP_COLLECTION",
	5: "CREATE_PARTITION",
	6: "DROP_PARTITION",
}

var MessageType_value = map[string]int32{
	"TIME_TICK":         0,
	"INSERT":            1,
	"DELETE":            2,
	"CREATE_COLLECTION": 3,
	"DROP_COLLECTION":   4,
	"CREATE_PARTITION":  5,
	"DROP_PARTITION":    6,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}

func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{1}
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

type SendRequest struct {
	ChannelName          string            `protobuf:"bytes,1,opt,name=channel_name,json=channelName,proto3" json:"channel_name,omitempty"`
	Payloads             [][]byte          `protobuf:"bytes,2,rep,name=payloads,proto3" json:"payloads,omitempty"`
	MessageType          MessageType       `protobuf:"varint,3,opt,name=message_type,json=messageType,proto3,enum=milvus.proto.log.MessageType" json:"message_type,omitempty"`
	Header               map[string]string `protobuf:"bytes,4,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SendRequest) Reset()         { *m = SendRequest{} }
func (m *SendRequest) String() string { return proto.CompactTextString(m) }
func (*SendRequest) ProtoMessage()    {}
func (*SendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{5}
}

func (m *SendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendRequest.Unmarshal(m, b)
}
func (m *SendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendRequest.Marshal(b, m, deterministic)
}
func (m *SendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendRequest.Merge(m, src)
}
func (m *SendRequest) XXX_Size() int {
	return xxx_messageInfo_SendRequest.Size(m)
}
func (m *SendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendRequest proto.InternalMessageInfo

func (m *SendRequest) GetChannelName() string {
	if m != nil {
		return m.ChannelName
	}
	return ""
}

func (m *SendRequest) GetPayloads() [][]byte {
	if m != nil {
		return m.Payloads
	}
	return nil
}

func (m *SendRequest) GetMessageType() MessageType {
	if m != nil {
		return m.MessageType
	}
	return MessageType_TIME_TICK
}

func (m *SendRequest) GetHeader() map[string]string {
	if m != nil {
		return m.Header
	}
	return nil
}

type SendResponse struct {
	Status               *commonpb.Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Offset               int64            `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Timestamp            int64            `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SendResponse) Reset()         { *m = SendResponse{} }
func (m *SendResponse) String() string { return proto.CompactTextString(m) }
func (*SendResponse) ProtoMessage()    {}
func (*SendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5674412121b98d7f, []int{6}
}

func (m *SendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendResponse.Unmarshal(m, b)
}
func (m *SendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendResponse.Marshal(b, m, deterministic)
}
func (m *SendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendResponse.Merge(m, src)
}
func (m *SendResponse) XXX_Size() int {
	return xxx_messageInfo_SendResponse.Size(m)
}
func (m *SendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SendResponse proto.InternalMessageInfo

func (m *SendResponse) GetStatus() *commonpb.Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *SendResponse) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *SendResponse) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func init() {
	proto.RegisterEnum("milvus.proto.log.PChannelState", PChannelState_name, PChannelState_value)
	proto.RegisterEnum("milvus.proto.log.MessageType", MessageType_name, MessageType_value)
	proto.RegisterType((*PChannelInfo)(nil), "milvus.proto.log.PChannelInfo")
	proto.RegisterType((*LogNodeInfo)(nil), "milvus.proto.log.LogNodeInfo")
	proto.RegisterType((*WatchChannelRequest)(nil), "milvus.proto.log.WatchChannelRequest")
	proto.RegisterType((*UnwatchChannelRequest)(nil), "milvus.proto.log.UnwatchChannelRequest")
	proto.RegisterType((*InsertRequest)(nil), "milvus.proto.log.InsertRequest")
	proto.RegisterType((*SendRequest)(nil), "milvus.proto.log.SendRequest")
	proto.RegisterMapType((map[string]string)(nil), "milvus.proto.log.SendRequest.HeaderEntry")
	proto.RegisterType((*SendResponse)(nil), "milvus.proto.log.SendResponse")
}

func init() { proto.RegisterFile("log_coord.proto", fileDescriptor_5674412121b98d7f) }

var fileDescriptor_5674412121b98d7f = []byte{
	// 739 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x5d, 0x6f, 0xe3, 0x44,
	0x14, 0x8d, 0xe3, 0x34, 0xbb, 0xb9, 0x76, 0xb3, 0x66, 0x76, 0x17, 0x45, 0x66, 0x17, 0x82, 0x05,
	0x22, 0xac, 0x44, 0xb2, 0x4a, 0x85, 0xf8, 0x78, 0x81, 0x36, 0xb5, 0x5a, 0x8b, 0x34, 0xad, 0x26,
	0xae, 0x8a, 0x78, 0x89, 0x9c, 0x78, 0xea, 0x44, 0xd8, 0x33, 0x6e, 0x66, 0x52, 0x94, 0xbf, 0x80,
	0x78, 0xe4, 0xd7, 0xf0, 0xeb, 0xd0, 0xd8, 0x93, 0xc6, 0xf9, 0xa0, 0x95, 0x90, 0xf6, 0x6d, 0xce,
	0xf5, 0xf1, 0xb9, 0x77, 0xe6, 0x9e, 0x7b, 0xe1, 0x45, 0xcc, 0xa2, 0xd1, 0x84, 0xb1, 0x79, 0xd8,
	0x4e, 0xe7, 0x4c, 0x30, 0x64, 0x25, 0xb3, 0xf8, 0x7e, 0xc1, 0x73, 0xd4, 0x8e, 0x59, 0x64, 0x9b,
	0x13, 0x96, 0x24, 0x8c, 0xe6, 0x11, 0xdb, 0x2c, 0x7e, 0xb7, 0x6b, 0x09, 0x8f, 0xf2, 0xa3, 0x73,
	0x07, 0xe6, 0x55, 0x6f, 0x1a, 0x50, 0x4a, 0x62, 0x8f, 0xde, 0x32, 0x84, 0xa0, 0x42, 0x83, 0x84,
	0x34, 0xb4, 0xa6, 0xd6, 0xaa, 0xe1, 0xec, 0x8c, 0xbe, 0x85, 0x03, 0x2e, 0x02, 0x41, 0x1a, 0xe5,
	0xa6, 0xd6, 0xaa, 0x77, 0x3f, 0x6b, 0x6f, 0x27, 0x6b, 0xaf, 0x24, 0x86, 0x92, 0x86, 0x73, 0x36,
	0xfa, 0x18, 0xaa, 0x94, 0x85, 0xc4, 0x3b, 0x6d, 0xe8, 0x4d, 0xad, 0xa5, 0x63, 0x85, 0x9c, 0x9f,
	0xc0, 0xe8, 0xb3, 0x68, 0x20, 0x81, 0xcc, 0xb8, 0xa6, 0x69, 0x45, 0x1a, 0x6a, 0xc0, 0xb3, 0x20,
	0x0c, 0xe7, 0x84, 0xf3, 0x2c, 0x6f, 0x0d, 0xaf, 0xa0, 0x33, 0x81, 0x97, 0x37, 0x81, 0x98, 0x4c,
	0x55, 0x52, 0x4c, 0xee, 0x16, 0x84, 0x0b, 0xf4, 0x1e, 0x2a, 0xe3, 0x80, 0xe7, 0xa5, 0x1b, 0xdd,
	0x37, 0x9b, 0x55, 0xaa, 0xd7, 0xb8, 0xe0, 0xd1, 0x49, 0xc0, 0x09, 0xce, 0x98, 0xc8, 0x86, 0xe7,
	0xa9, 0x12, 0x51, 0x39, 0x1e, 0xb0, 0x43, 0xe0, 0xf5, 0x35, 0xfd, 0xe3, 0x83, 0xa7, 0xf9, 0x5b,
	0x83, 0x43, 0x8f, 0x72, 0x32, 0x17, 0xff, 0x5f, 0xff, 0x08, 0x2a, 0x09, 0x8f, 0xe4, 0x33, 0xe9,
	0x2d, 0x63, 0xbb, 0x3d, 0xb2, 0xd5, 0x1b, 0x09, 0x70, 0x46, 0x46, 0x6f, 0xa0, 0xb6, 0x2a, 0x82,
	0x37, 0xf4, 0xa6, 0xde, 0xaa, 0xe1, 0x75, 0xc0, 0xf9, 0xab, 0x0c, 0xc6, 0x90, 0xd0, 0x70, 0x55,
	0xd4, 0xe7, 0x60, 0x4e, 0xf2, 0x6f, 0xa3, 0x82, 0x3d, 0x0c, 0x15, 0x1b, 0x48, 0x97, 0xc8, 0x5b,
	0x06, 0xcb, 0x98, 0x05, 0x61, 0x5e, 0x89, 0x89, 0x1f, 0x30, 0xfa, 0x19, 0xcc, 0x84, 0x70, 0x1e,
	0x44, 0x64, 0x24, 0x96, 0x29, 0xc9, 0x0c, 0x51, 0xef, 0xbe, 0xdd, 0x35, 0xd2, 0x45, 0xce, 0xf2,
	0x97, 0x29, 0xc1, 0x46, 0xb2, 0x06, 0xe8, 0x18, 0xaa, 0x53, 0x12, 0x84, 0x64, 0xde, 0xa8, 0x64,
	0xb7, 0xfc, 0x7a, 0xf7, 0xdf, 0x42, 0xbd, 0xed, 0xf3, 0x8c, 0xeb, 0x52, 0x31, 0x5f, 0x62, 0xf5,
	0xa3, 0xfd, 0x03, 0x18, 0x85, 0x30, 0xb2, 0x40, 0xff, 0x9d, 0x2c, 0xd5, 0x4d, 0xe4, 0x11, 0xbd,
	0x82, 0x83, 0xfb, 0x20, 0x5e, 0x10, 0xd5, 0xa4, 0x1c, 0xfc, 0x58, 0xfe, 0x5e, 0x73, 0x96, 0x60,
	0xe6, 0xea, 0x3c, 0x65, 0x34, 0x7b, 0xf1, 0xaa, 0xf4, 0xf8, 0x82, 0xab, 0x2e, 0x7d, 0xb2, 0xb7,
	0x4b, 0xc3, 0x8c, 0x82, 0x15, 0x55, 0x1a, 0x9d, 0xdd, 0xde, 0x72, 0x22, 0x32, 0x7d, 0x1d, 0x2b,
	0x24, 0x3b, 0x21, 0x66, 0x09, 0xe1, 0x22, 0x48, 0x52, 0x35, 0x2a, 0xeb, 0xc0, 0xbb, 0xef, 0xe0,
	0x70, 0x63, 0xba, 0x10, 0x40, 0xf5, 0x9a, 0x2e, 0x38, 0x09, 0xad, 0x12, 0x32, 0xe1, 0xf9, 0x4d,
	0x30, 0x13, 0x62, 0x46, 0x23, 0x4b, 0xcb, 0x91, 0x98, 0x4c, 0x25, 0x2a, 0xbf, 0xfb, 0x53, 0x03,
	0xa3, 0xf0, 0x9c, 0xe8, 0x10, 0x6a, 0xbe, 0x77, 0xe1, 0x8e, 0x7c, 0xaf, 0xf7, 0x8b, 0x55, 0x92,
	0x32, 0xde, 0x60, 0xe8, 0x62, 0xdf, 0xd2, 0xe4, 0xf9, 0xd4, 0xed, 0xbb, 0xbe, 0x6b, 0x95, 0xd1,
	0x6b, 0xf8, 0xa8, 0x87, 0xdd, 0x63, 0xdf, 0x1d, 0xf5, 0x2e, 0xfb, 0x7d, 0xb7, 0xe7, 0x7b, 0x97,
	0x03, 0x4b, 0x47, 0x2f, 0xe1, 0xc5, 0x29, 0xbe, 0xbc, 0x2a, 0x06, 0x2b, 0xe8, 0x15, 0x58, 0x8a,
	0x7b, 0x75, 0x8c, 0x7d, 0x2f, 0x8b, 0x1e, 0x20, 0x04, 0xf5, 0x8c, 0xba, 0x8e, 0x55, 0xbb, 0xff,
	0xe8, 0xf0, 0x4c, 0x0d, 0x3d, 0x8a, 0x01, 0x9d, 0x11, 0xd1, 0x63, 0x49, 0xca, 0x28, 0xa1, 0x22,
	0xbb, 0x15, 0x47, 0xed, 0x2d, 0xdb, 0xe6, 0x60, 0x97, 0xa8, 0x3a, 0x6c, 0x7f, 0xb1, 0x97, 0xbf,
	0x45, 0x76, 0x4a, 0xc8, 0x07, 0xb3, 0xb8, 0x2c, 0xd0, 0x97, 0xbb, 0xc6, 0xd9, 0xb3, 0x4c, 0xec,
	0xc7, 0x3a, 0xea, 0x94, 0xd0, 0xaf, 0x50, 0xdf, 0xdc, 0x0e, 0xe8, 0xab, 0x5d, 0xdd, 0xbd, 0xfb,
	0xe3, 0x29, 0xe5, 0x73, 0xa8, 0xe6, 0xe3, 0x8a, 0xf6, 0xec, 0xd9, 0x8d, 0x41, 0x7e, 0x4a, 0xe9,
	0x0c, 0x2a, 0xd2, 0xb4, 0xe8, 0xed, 0xa3, 0xa3, 0x62, 0x7f, 0xfa, 0x5f, 0x9f, 0x73, 0xaf, 0x3b,
	0xa5, 0x93, 0xee, 0x6f, 0xef, 0xa3, 0x99, 0x98, 0x2e, 0xc6, 0x32, 0x45, 0x27, 0x67, 0x7f, 0x33,
	0x63, 0xea, 0xd4, 0x99, 0x51, 0x41, 0xe6, 0x34, 0x88, 0x3b, 0x99, 0x40, 0x27, 0x66, 0x51, 0x3a,
	0x1e, 0x57, 0x33, 0x70, 0xf4, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6f, 0xd5, 0xb4, 0x84, 0xaa,
	0x06, 0x00, 0x00,
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
	Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error)
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

func (c *logNodeClient) Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error) {
	out := new(SendResponse)
	err := c.cc.Invoke(ctx, "/milvus.proto.log.LogNode/Send", in, out, opts...)
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
	Send(context.Context, *SendRequest) (*SendResponse, error)
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
func (*UnimplementedLogNodeServer) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
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

func _LogNode_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogNodeServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/milvus.proto.log.LogNode/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogNodeServer).Send(ctx, req.(*SendRequest))
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
		{
			MethodName: "Send",
			Handler:    _LogNode_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "log_coord.proto",
}
