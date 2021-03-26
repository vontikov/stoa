// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: raft.proto

package pb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Type is a Raft claster command type.
type ClusterCommand_Type int32

const (
	ClusterCommand_RESERVED                 ClusterCommand_Type = 0
	ClusterCommand_QUEUE_SIZE               ClusterCommand_Type = 1
	ClusterCommand_QUEUE_CLEAR              ClusterCommand_Type = 2
	ClusterCommand_QUEUE_OFFER              ClusterCommand_Type = 3
	ClusterCommand_QUEUE_POLL               ClusterCommand_Type = 4
	ClusterCommand_QUEUE_PEEK               ClusterCommand_Type = 5
	ClusterCommand_DICTIONARY_SIZE          ClusterCommand_Type = 6
	ClusterCommand_DICTIONARY_CLEAR         ClusterCommand_Type = 7
	ClusterCommand_DICTIONARY_PUT_IF_ABSENT ClusterCommand_Type = 8
	ClusterCommand_DICTIONARY_PUT           ClusterCommand_Type = 9
	ClusterCommand_DICTIONARY_GET           ClusterCommand_Type = 10
	ClusterCommand_DICTIONARY_REMOVE        ClusterCommand_Type = 11
	ClusterCommand_DICTIONARY_RANGE         ClusterCommand_Type = 12
	ClusterCommand_MUTEX_TRY_LOCK           ClusterCommand_Type = 13
	ClusterCommand_MUTEX_UNLOCK             ClusterCommand_Type = 14
	ClusterCommand_SERVICE_PING             ClusterCommand_Type = 15
	ClusterCommand_MAX_INDEX                ClusterCommand_Type = 16
)

// Enum value maps for ClusterCommand_Type.
var (
	ClusterCommand_Type_name = map[int32]string{
		0:  "RESERVED",
		1:  "QUEUE_SIZE",
		2:  "QUEUE_CLEAR",
		3:  "QUEUE_OFFER",
		4:  "QUEUE_POLL",
		5:  "QUEUE_PEEK",
		6:  "DICTIONARY_SIZE",
		7:  "DICTIONARY_CLEAR",
		8:  "DICTIONARY_PUT_IF_ABSENT",
		9:  "DICTIONARY_PUT",
		10: "DICTIONARY_GET",
		11: "DICTIONARY_REMOVE",
		12: "DICTIONARY_RANGE",
		13: "MUTEX_TRY_LOCK",
		14: "MUTEX_UNLOCK",
		15: "SERVICE_PING",
		16: "MAX_INDEX",
	}
	ClusterCommand_Type_value = map[string]int32{
		"RESERVED":                 0,
		"QUEUE_SIZE":               1,
		"QUEUE_CLEAR":              2,
		"QUEUE_OFFER":              3,
		"QUEUE_POLL":               4,
		"QUEUE_PEEK":               5,
		"DICTIONARY_SIZE":          6,
		"DICTIONARY_CLEAR":         7,
		"DICTIONARY_PUT_IF_ABSENT": 8,
		"DICTIONARY_PUT":           9,
		"DICTIONARY_GET":           10,
		"DICTIONARY_REMOVE":        11,
		"DICTIONARY_RANGE":         12,
		"MUTEX_TRY_LOCK":           13,
		"MUTEX_UNLOCK":             14,
		"SERVICE_PING":             15,
		"MAX_INDEX":                16,
	}
)

func (x ClusterCommand_Type) Enum() *ClusterCommand_Type {
	p := new(ClusterCommand_Type)
	*p = x
	return p
}

func (x ClusterCommand_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClusterCommand_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_raft_proto_enumTypes[0].Descriptor()
}

func (ClusterCommand_Type) Type() protoreflect.EnumType {
	return &file_raft_proto_enumTypes[0]
}

func (x ClusterCommand_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ClusterCommand_Type.Descriptor instead.
func (ClusterCommand_Type) EnumDescriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{1, 0}
}

// Discovery contains Raft Peer information.
type Discovery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                              // Peer ID
	BindIp   string `protobuf:"bytes,2,opt,name=bind_ip,json=bindIp,proto3" json:"bind_ip,omitempty"`        // bind IP
	BindPort int32  `protobuf:"varint,3,opt,name=bind_port,json=bindPort,proto3" json:"bind_port,omitempty"` // bind port
	GrpcIp   string `protobuf:"bytes,4,opt,name=grpc_ip,json=grpcIp,proto3" json:"grpc_ip,omitempty"`        // gRPC IP
	GrpcPort int32  `protobuf:"varint,5,opt,name=grpc_port,json=grpcPort,proto3" json:"grpc_port,omitempty"` // gRPC port
	Leader   bool   `protobuf:"varint,6,opt,name=leader,proto3" json:"leader,omitempty"`                     // leader status
}

func (x *Discovery) Reset() {
	*x = Discovery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Discovery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Discovery) ProtoMessage() {}

func (x *Discovery) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Discovery.ProtoReflect.Descriptor instead.
func (*Discovery) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{0}
}

func (x *Discovery) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Discovery) GetBindIp() string {
	if x != nil {
		return x.BindIp
	}
	return ""
}

func (x *Discovery) GetBindPort() int32 {
	if x != nil {
		return x.BindPort
	}
	return 0
}

func (x *Discovery) GetGrpcIp() string {
	if x != nil {
		return x.GrpcIp
	}
	return ""
}

func (x *Discovery) GetGrpcPort() int32 {
	if x != nil {
		return x.GrpcPort
	}
	return 0
}

func (x *Discovery) GetLeader() bool {
	if x != nil {
		return x.Leader
	}
	return false
}

// ClusterCommand contains Raft claster log command.
type ClusterCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command        ClusterCommand_Type `protobuf:"varint,1,opt,name=command,proto3,enum=github.com.vontikov.stoa.v1.ClusterCommand_Type" json:"command,omitempty"`
	TtlEnabled     bool                `protobuf:"varint,2,opt,name=ttl_enabled,json=ttlEnabled,proto3" json:"ttl_enabled,omitempty"`
	TtlMillis      int64               `protobuf:"varint,3,opt,name=ttl_millis,json=ttlMillis,proto3" json:"ttl_millis,omitempty"`
	BarrierEnabled bool                `protobuf:"varint,4,opt,name=barrier_enabled,json=barrierEnabled,proto3" json:"barrier_enabled,omitempty"`
	BarrierMillis  int64               `protobuf:"varint,5,opt,name=barrier_millis,json=barrierMillis,proto3" json:"barrier_millis,omitempty"`
	// Types that are assignable to Payload:
	//	*ClusterCommand_Entity
	//	*ClusterCommand_Value
	//	*ClusterCommand_Key
	//	*ClusterCommand_KeyValue
	//	*ClusterCommand_ClientId
	Payload isClusterCommand_Payload `protobuf_oneof:"payload"`
}

func (x *ClusterCommand) Reset() {
	*x = ClusterCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterCommand) ProtoMessage() {}

func (x *ClusterCommand) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterCommand.ProtoReflect.Descriptor instead.
func (*ClusterCommand) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{1}
}

func (x *ClusterCommand) GetCommand() ClusterCommand_Type {
	if x != nil {
		return x.Command
	}
	return ClusterCommand_RESERVED
}

func (x *ClusterCommand) GetTtlEnabled() bool {
	if x != nil {
		return x.TtlEnabled
	}
	return false
}

func (x *ClusterCommand) GetTtlMillis() int64 {
	if x != nil {
		return x.TtlMillis
	}
	return 0
}

func (x *ClusterCommand) GetBarrierEnabled() bool {
	if x != nil {
		return x.BarrierEnabled
	}
	return false
}

func (x *ClusterCommand) GetBarrierMillis() int64 {
	if x != nil {
		return x.BarrierMillis
	}
	return 0
}

func (m *ClusterCommand) GetPayload() isClusterCommand_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *ClusterCommand) GetEntity() *Entity {
	if x, ok := x.GetPayload().(*ClusterCommand_Entity); ok {
		return x.Entity
	}
	return nil
}

func (x *ClusterCommand) GetValue() *Value {
	if x, ok := x.GetPayload().(*ClusterCommand_Value); ok {
		return x.Value
	}
	return nil
}

func (x *ClusterCommand) GetKey() *Key {
	if x, ok := x.GetPayload().(*ClusterCommand_Key); ok {
		return x.Key
	}
	return nil
}

func (x *ClusterCommand) GetKeyValue() *KeyValue {
	if x, ok := x.GetPayload().(*ClusterCommand_KeyValue); ok {
		return x.KeyValue
	}
	return nil
}

func (x *ClusterCommand) GetClientId() *ClientId {
	if x, ok := x.GetPayload().(*ClusterCommand_ClientId); ok {
		return x.ClientId
	}
	return nil
}

type isClusterCommand_Payload interface {
	isClusterCommand_Payload()
}

type ClusterCommand_Entity struct {
	Entity *Entity `protobuf:"bytes,6,opt,name=entity,proto3,oneof"`
}

type ClusterCommand_Value struct {
	Value *Value `protobuf:"bytes,7,opt,name=value,proto3,oneof"`
}

type ClusterCommand_Key struct {
	Key *Key `protobuf:"bytes,8,opt,name=key,proto3,oneof"`
}

type ClusterCommand_KeyValue struct {
	KeyValue *KeyValue `protobuf:"bytes,9,opt,name=key_value,json=keyValue,proto3,oneof"`
}

type ClusterCommand_ClientId struct {
	ClientId *ClientId `protobuf:"bytes,10,opt,name=client_id,json=clientId,proto3,oneof"`
}

func (*ClusterCommand_Entity) isClusterCommand_Payload() {}

func (*ClusterCommand_Value) isClusterCommand_Payload() {}

func (*ClusterCommand_Key) isClusterCommand_Payload() {}

func (*ClusterCommand_KeyValue) isClusterCommand_Payload() {}

func (*ClusterCommand_ClientId) isClusterCommand_Payload() {}

var File_raft_proto protoreflect.FileDescriptor

var file_raft_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f,
	0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x1a, 0x0a, 0x73, 0x74, 0x6f, 0x61, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9f, 0x01, 0x0a, 0x09, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x69, 0x6e, 0x64, 0x5f, 0x69, 0x70, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x69, 0x6e, 0x64, 0x49, 0x70, 0x12, 0x1b, 0x0a, 0x09,
	0x62, 0x69, 0x6e, 0x64, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x62, 0x69, 0x6e, 0x64, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x67, 0x72, 0x70,
	0x63, 0x5f, 0x69, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x72, 0x70, 0x63,
	0x49, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x67, 0x72, 0x70, 0x63, 0x50, 0x6f, 0x72, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x22, 0x82, 0x07, 0x0a, 0x0e, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x4a, 0x0a, 0x07, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x30, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f,
	0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x74, 0x6c, 0x5f, 0x65, 0x6e,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x74, 0x74, 0x6c,
	0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x74, 0x6c, 0x5f, 0x6d,
	0x69, 0x6c, 0x6c, 0x69, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x74, 0x6c,
	0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x62, 0x61, 0x72, 0x72, 0x69, 0x65,
	0x72, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0e, 0x62, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12,
	0x25, 0x0a, 0x0e, 0x62, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x5f, 0x6d, 0x69, 0x6c, 0x6c, 0x69,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x62, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72,
	0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x12, 0x3d, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f,
	0x61, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x48, 0x00, 0x52, 0x06, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x3a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e,
	0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x34, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74,
	0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79,
	0x48, 0x00, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x44, 0x0a, 0x09, 0x6b, 0x65, 0x79, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76,
	0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x48, 0x00, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x44, 0x0a,
	0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x25, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f,
	0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x48, 0x00, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x22, 0xcb, 0x02, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08,
	0x52, 0x45, 0x53, 0x45, 0x52, 0x56, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x51, 0x55,
	0x45, 0x55, 0x45, 0x5f, 0x53, 0x49, 0x5a, 0x45, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x51, 0x55,
	0x45, 0x55, 0x45, 0x5f, 0x43, 0x4c, 0x45, 0x41, 0x52, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x51,
	0x55, 0x45, 0x55, 0x45, 0x5f, 0x4f, 0x46, 0x46, 0x45, 0x52, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a,
	0x51, 0x55, 0x45, 0x55, 0x45, 0x5f, 0x50, 0x4f, 0x4c, 0x4c, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a,
	0x51, 0x55, 0x45, 0x55, 0x45, 0x5f, 0x50, 0x45, 0x45, 0x4b, 0x10, 0x05, 0x12, 0x13, 0x0a, 0x0f,
	0x44, 0x49, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x41, 0x52, 0x59, 0x5f, 0x53, 0x49, 0x5a, 0x45, 0x10,
	0x06, 0x12, 0x14, 0x0a, 0x10, 0x44, 0x49, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x41, 0x52, 0x59, 0x5f,
	0x43, 0x4c, 0x45, 0x41, 0x52, 0x10, 0x07, 0x12, 0x1c, 0x0a, 0x18, 0x44, 0x49, 0x43, 0x54, 0x49,
	0x4f, 0x4e, 0x41, 0x52, 0x59, 0x5f, 0x50, 0x55, 0x54, 0x5f, 0x49, 0x46, 0x5f, 0x41, 0x42, 0x53,
	0x45, 0x4e, 0x54, 0x10, 0x08, 0x12, 0x12, 0x0a, 0x0e, 0x44, 0x49, 0x43, 0x54, 0x49, 0x4f, 0x4e,
	0x41, 0x52, 0x59, 0x5f, 0x50, 0x55, 0x54, 0x10, 0x09, 0x12, 0x12, 0x0a, 0x0e, 0x44, 0x49, 0x43,
	0x54, 0x49, 0x4f, 0x4e, 0x41, 0x52, 0x59, 0x5f, 0x47, 0x45, 0x54, 0x10, 0x0a, 0x12, 0x15, 0x0a,
	0x11, 0x44, 0x49, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x41, 0x52, 0x59, 0x5f, 0x52, 0x45, 0x4d, 0x4f,
	0x56, 0x45, 0x10, 0x0b, 0x12, 0x14, 0x0a, 0x10, 0x44, 0x49, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x41,
	0x52, 0x59, 0x5f, 0x52, 0x41, 0x4e, 0x47, 0x45, 0x10, 0x0c, 0x12, 0x12, 0x0a, 0x0e, 0x4d, 0x55,
	0x54, 0x45, 0x58, 0x5f, 0x54, 0x52, 0x59, 0x5f, 0x4c, 0x4f, 0x43, 0x4b, 0x10, 0x0d, 0x12, 0x10,
	0x0a, 0x0c, 0x4d, 0x55, 0x54, 0x45, 0x58, 0x5f, 0x55, 0x4e, 0x4c, 0x4f, 0x43, 0x4b, 0x10, 0x0e,
	0x12, 0x10, 0x0a, 0x0c, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x50, 0x49, 0x4e, 0x47,
	0x10, 0x0f, 0x12, 0x0d, 0x0a, 0x09, 0x4d, 0x41, 0x58, 0x5f, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x10,
	0x10, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x06, 0x5a, 0x04,
	0x2e, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_raft_proto_rawDescOnce sync.Once
	file_raft_proto_rawDescData = file_raft_proto_rawDesc
)

func file_raft_proto_rawDescGZIP() []byte {
	file_raft_proto_rawDescOnce.Do(func() {
		file_raft_proto_rawDescData = protoimpl.X.CompressGZIP(file_raft_proto_rawDescData)
	})
	return file_raft_proto_rawDescData
}

var file_raft_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_raft_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_raft_proto_goTypes = []interface{}{
	(ClusterCommand_Type)(0), // 0: github.com.vontikov.stoa.v1.ClusterCommand.Type
	(*Discovery)(nil),        // 1: github.com.vontikov.stoa.v1.Discovery
	(*ClusterCommand)(nil),   // 2: github.com.vontikov.stoa.v1.ClusterCommand
	(*Entity)(nil),           // 3: github.com.vontikov.stoa.v1.Entity
	(*Value)(nil),            // 4: github.com.vontikov.stoa.v1.Value
	(*Key)(nil),              // 5: github.com.vontikov.stoa.v1.Key
	(*KeyValue)(nil),         // 6: github.com.vontikov.stoa.v1.KeyValue
	(*ClientId)(nil),         // 7: github.com.vontikov.stoa.v1.ClientId
}
var file_raft_proto_depIdxs = []int32{
	0, // 0: github.com.vontikov.stoa.v1.ClusterCommand.command:type_name -> github.com.vontikov.stoa.v1.ClusterCommand.Type
	3, // 1: github.com.vontikov.stoa.v1.ClusterCommand.entity:type_name -> github.com.vontikov.stoa.v1.Entity
	4, // 2: github.com.vontikov.stoa.v1.ClusterCommand.value:type_name -> github.com.vontikov.stoa.v1.Value
	5, // 3: github.com.vontikov.stoa.v1.ClusterCommand.key:type_name -> github.com.vontikov.stoa.v1.Key
	6, // 4: github.com.vontikov.stoa.v1.ClusterCommand.key_value:type_name -> github.com.vontikov.stoa.v1.KeyValue
	7, // 5: github.com.vontikov.stoa.v1.ClusterCommand.client_id:type_name -> github.com.vontikov.stoa.v1.ClientId
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_raft_proto_init() }
func file_raft_proto_init() {
	if File_raft_proto != nil {
		return
	}
	file_stoa_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_raft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Discovery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_raft_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*ClusterCommand_Entity)(nil),
		(*ClusterCommand_Value)(nil),
		(*ClusterCommand_Key)(nil),
		(*ClusterCommand_KeyValue)(nil),
		(*ClusterCommand_ClientId)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_raft_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_raft_proto_goTypes,
		DependencyIndexes: file_raft_proto_depIdxs,
		EnumInfos:         file_raft_proto_enumTypes,
		MessageInfos:      file_raft_proto_msgTypes,
	}.Build()
	File_raft_proto = out.File
	file_raft_proto_rawDesc = nil
	file_raft_proto_goTypes = nil
	file_raft_proto_depIdxs = nil
}
