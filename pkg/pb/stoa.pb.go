// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.13.0
// source: stoa.proto

package pb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Entity is a Stoa entity.
type Entity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityName string `protobuf:"bytes,1,opt,name=entity_name,json=entityName,proto3" json:"entity_name,omitempty"`
}

func (x *Entity) Reset() {
	*x = Entity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stoa_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entity) ProtoMessage() {}

func (x *Entity) ProtoReflect() protoreflect.Message {
	mi := &file_stoa_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entity.ProtoReflect.Descriptor instead.
func (*Entity) Descriptor() ([]byte, []int) {
	return file_stoa_proto_rawDescGZIP(), []int{0}
}

func (x *Entity) GetEntityName() string {
	if x != nil {
		return x.EntityName
	}
	return ""
}

// ClientId contains Client ID.
type ClientId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityName string `protobuf:"bytes,1,opt,name=entity_name,json=entityName,proto3" json:"entity_name,omitempty"`
	Id         []byte `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Payload    []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ClientId) Reset() {
	*x = ClientId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stoa_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientId) ProtoMessage() {}

func (x *ClientId) ProtoReflect() protoreflect.Message {
	mi := &file_stoa_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientId.ProtoReflect.Descriptor instead.
func (*ClientId) Descriptor() ([]byte, []int) {
	return file_stoa_proto_rawDescGZIP(), []int{1}
}

func (x *ClientId) GetEntityName() string {
	if x != nil {
		return x.EntityName
	}
	return ""
}

func (x *ClientId) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ClientId) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

// Value contains an arbitrary data.
type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityName string `protobuf:"bytes,1,opt,name=entity_name,json=entityName,proto3" json:"entity_name,omitempty"`
	Value      []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stoa_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_stoa_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_stoa_proto_rawDescGZIP(), []int{2}
}

func (x *Value) GetEntityName() string {
	if x != nil {
		return x.EntityName
	}
	return ""
}

func (x *Value) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// Key contains Dictionary key.
type Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityName string `protobuf:"bytes,1,opt,name=entity_name,json=entityName,proto3" json:"entity_name,omitempty"`
	Key        []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *Key) Reset() {
	*x = Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stoa_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Key) ProtoMessage() {}

func (x *Key) ProtoReflect() protoreflect.Message {
	mi := &file_stoa_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Key.ProtoReflect.Descriptor instead.
func (*Key) Descriptor() ([]byte, []int) {
	return file_stoa_proto_rawDescGZIP(), []int{3}
}

func (x *Key) GetEntityName() string {
	if x != nil {
		return x.EntityName
	}
	return ""
}

func (x *Key) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

// KeyValue contains Dictionary key-value pair.
type KeyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityName string `protobuf:"bytes,1,opt,name=entity_name,json=entityName,proto3" json:"entity_name,omitempty"`
	Key        []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value      []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *KeyValue) Reset() {
	*x = KeyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stoa_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValue) ProtoMessage() {}

func (x *KeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_stoa_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValue.ProtoReflect.Descriptor instead.
func (*KeyValue) Descriptor() ([]byte, []int) {
	return file_stoa_proto_rawDescGZIP(), []int{4}
}

func (x *KeyValue) GetEntityName() string {
	if x != nil {
		return x.EntityName
	}
	return ""
}

func (x *KeyValue) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *KeyValue) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_stoa_proto protoreflect.FileDescriptor

var file_stoa_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f,
	0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a, 0x06, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x12, 0x28, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01,
	0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x5e, 0x0a, 0x08,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x7a, 0x02, 0x10, 0x01, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x50, 0x0a, 0x05,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x28, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x7a, 0x02, 0x10, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x4a,
	0x0a, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x19, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x7a, 0x02, 0x10, 0x01, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x6e, 0x0a, 0x08, 0x4b, 0x65,
	0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x28, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x19, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x7a, 0x02, 0x10, 0x01, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1d, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x7a,
	0x02, 0x10, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0xd5, 0x0f, 0x0a, 0x04, 0x53,
	0x74, 0x6f, 0x61, 0x12, 0x7a, 0x0a, 0x09, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f,
	0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a, 0x22, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61,
	0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1e, 0x12, 0x1c, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2f, 0x73, 0x69, 0x7a,
	0x65, 0x2f, 0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12,
	0x7c, 0x0a, 0x0a, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x12, 0x23, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69,
	0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x1a, 0x22, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x25, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x22, 0x1d,
	0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2f, 0x63, 0x6c, 0x65, 0x61, 0x72, 0x2f,
	0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0x7f, 0x0a,
	0x0a, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4f, 0x66, 0x66, 0x65, 0x72, 0x12, 0x22, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f,
	0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a,
	0x23, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e,
	0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x22, 0x1d, 0x2f, 0x76,
	0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2f, 0x6f, 0x66, 0x66, 0x65, 0x72, 0x2f, 0x7b, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x7a,
	0x0a, 0x09, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x12, 0x23, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f,
	0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x1a, 0x22, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f,
	0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x12, 0x1c, 0x2f, 0x76,
	0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x6c, 0x2f, 0x7b, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0x7a, 0x0a, 0x09, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x50, 0x65, 0x65, 0x6b, 0x12, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74,
	0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a, 0x22, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b,
	0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x12, 0x1c, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x2f, 0x70, 0x65, 0x65, 0x6b, 0x2f, 0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0x84, 0x01, 0x0a, 0x0e, 0x44, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x72, 0x79, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e,
	0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a, 0x22,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74,
	0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x29, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x23, 0x12, 0x21, 0x2f, 0x76, 0x31, 0x2f,
	0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x2f, 0x73, 0x69, 0x7a, 0x65, 0x2f,
	0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0x86, 0x01,
	0x0a, 0x0f, 0x44, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x43, 0x6c, 0x65, 0x61,
	0x72, 0x12, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76,
	0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a, 0x22, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f,
	0x61, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x24, 0x22, 0x22, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61,
	0x72, 0x79, 0x2f, 0x63, 0x6c, 0x65, 0x61, 0x72, 0x2f, 0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0x98, 0x01, 0x0a, 0x15, 0x44, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x72, 0x79, 0x50, 0x75, 0x74, 0x49, 0x66, 0x41, 0x62, 0x73, 0x65, 0x6e, 0x74,
	0x12, 0x25, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f,
	0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4b,
	0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74,
	0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x33, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x2d, 0x22, 0x28, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x61, 0x72, 0x79, 0x2f, 0x70, 0x75, 0x74, 0x49, 0x66, 0x41, 0x62, 0x73, 0x65, 0x6e, 0x74,
	0x2f, 0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x3a, 0x01,
	0x2a, 0x12, 0x87, 0x01, 0x0a, 0x0d, 0x44, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79,
	0x50, 0x75, 0x74, 0x12, 0x25, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76,
	0x31, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x22, 0x2e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76,
	0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x2b,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x22, 0x20, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x2f, 0x70, 0x75, 0x74, 0x2f, 0x7b, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x85, 0x01, 0x0a, 0x0d,
	0x44, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x47, 0x65, 0x74, 0x12, 0x20, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69,
	0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x1a,
	0x22, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e,
	0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x2e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x28, 0x12, 0x26, 0x2f, 0x76, 0x31,
	0x2f, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x2f, 0x67, 0x65, 0x74, 0x2f,
	0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x2f, 0x7b, 0x6b,
	0x65, 0x79, 0x7d, 0x12, 0x91, 0x01, 0x0a, 0x10, 0x44, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61,
	0x72, 0x79, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x20, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73,
	0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x1a, 0x23, 0x2e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76,
	0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22,
	0x36, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x30, 0x22, 0x2e, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x2f, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x2f, 0x64, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x2f, 0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x7d, 0x2f, 0x7b, 0x6b, 0x65, 0x79, 0x7d, 0x12, 0x8b, 0x01, 0x0a, 0x0f, 0x44, 0x69, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x23, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f,
	0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x1a, 0x25, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f,
	0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4b,
	0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x12,
	0x22, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x2f,
	0x72, 0x61, 0x6e, 0x67, 0x65, 0x2f, 0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x7d, 0x30, 0x01, 0x12, 0x86, 0x01, 0x0a, 0x0c, 0x4d, 0x75, 0x74, 0x65, 0x78, 0x54,
	0x72, 0x79, 0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x25, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f,
	0x61, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x1a, 0x23, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69,
	0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x22, 0x1f, 0x2f, 0x76, 0x31, 0x2f,
	0x6d, 0x75, 0x74, 0x65, 0x78, 0x2f, 0x74, 0x72, 0x79, 0x6c, 0x6f, 0x63, 0x6b, 0x2f, 0x7b, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x84,
	0x01, 0x0a, 0x0b, 0x4d, 0x75, 0x74, 0x65, 0x78, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x25,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74,
	0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x1a, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x29, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x23, 0x22, 0x1e, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x2f, 0x75, 0x6e,
	0x6c, 0x6f, 0x63, 0x6b, 0x2f, 0x7b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x6b, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x25, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69,
	0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x1a, 0x22, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12,
	0x22, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x69,
	0x6e, 0x67, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_stoa_proto_rawDescOnce sync.Once
	file_stoa_proto_rawDescData = file_stoa_proto_rawDesc
)

func file_stoa_proto_rawDescGZIP() []byte {
	file_stoa_proto_rawDescOnce.Do(func() {
		file_stoa_proto_rawDescData = protoimpl.X.CompressGZIP(file_stoa_proto_rawDescData)
	})
	return file_stoa_proto_rawDescData
}

var file_stoa_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_stoa_proto_goTypes = []interface{}{
	(*Entity)(nil),   // 0: github.com.vontikov.stoa.v1.Entity
	(*ClientId)(nil), // 1: github.com.vontikov.stoa.v1.ClientId
	(*Value)(nil),    // 2: github.com.vontikov.stoa.v1.Value
	(*Key)(nil),      // 3: github.com.vontikov.stoa.v1.Key
	(*KeyValue)(nil), // 4: github.com.vontikov.stoa.v1.KeyValue
	(*Empty)(nil),    // 5: github.com.vontikov.stoa.v1.Empty
	(*Result)(nil),   // 6: github.com.vontikov.stoa.v1.Result
}
var file_stoa_proto_depIdxs = []int32{
	0,  // 0: github.com.vontikov.stoa.v1.Stoa.QueueSize:input_type -> github.com.vontikov.stoa.v1.Entity
	0,  // 1: github.com.vontikov.stoa.v1.Stoa.QueueClear:input_type -> github.com.vontikov.stoa.v1.Entity
	2,  // 2: github.com.vontikov.stoa.v1.Stoa.QueueOffer:input_type -> github.com.vontikov.stoa.v1.Value
	0,  // 3: github.com.vontikov.stoa.v1.Stoa.QueuePoll:input_type -> github.com.vontikov.stoa.v1.Entity
	0,  // 4: github.com.vontikov.stoa.v1.Stoa.QueuePeek:input_type -> github.com.vontikov.stoa.v1.Entity
	0,  // 5: github.com.vontikov.stoa.v1.Stoa.DictionarySize:input_type -> github.com.vontikov.stoa.v1.Entity
	0,  // 6: github.com.vontikov.stoa.v1.Stoa.DictionaryClear:input_type -> github.com.vontikov.stoa.v1.Entity
	4,  // 7: github.com.vontikov.stoa.v1.Stoa.DictionaryPutIfAbsent:input_type -> github.com.vontikov.stoa.v1.KeyValue
	4,  // 8: github.com.vontikov.stoa.v1.Stoa.DictionaryPut:input_type -> github.com.vontikov.stoa.v1.KeyValue
	3,  // 9: github.com.vontikov.stoa.v1.Stoa.DictionaryGet:input_type -> github.com.vontikov.stoa.v1.Key
	3,  // 10: github.com.vontikov.stoa.v1.Stoa.DictionaryRemove:input_type -> github.com.vontikov.stoa.v1.Key
	0,  // 11: github.com.vontikov.stoa.v1.Stoa.DictionaryRange:input_type -> github.com.vontikov.stoa.v1.Entity
	1,  // 12: github.com.vontikov.stoa.v1.Stoa.MutexTryLock:input_type -> github.com.vontikov.stoa.v1.ClientId
	1,  // 13: github.com.vontikov.stoa.v1.Stoa.MutexUnlock:input_type -> github.com.vontikov.stoa.v1.ClientId
	1,  // 14: github.com.vontikov.stoa.v1.Stoa.Ping:input_type -> github.com.vontikov.stoa.v1.ClientId
	2,  // 15: github.com.vontikov.stoa.v1.Stoa.QueueSize:output_type -> github.com.vontikov.stoa.v1.Value
	5,  // 16: github.com.vontikov.stoa.v1.Stoa.QueueClear:output_type -> github.com.vontikov.stoa.v1.Empty
	6,  // 17: github.com.vontikov.stoa.v1.Stoa.QueueOffer:output_type -> github.com.vontikov.stoa.v1.Result
	2,  // 18: github.com.vontikov.stoa.v1.Stoa.QueuePoll:output_type -> github.com.vontikov.stoa.v1.Value
	2,  // 19: github.com.vontikov.stoa.v1.Stoa.QueuePeek:output_type -> github.com.vontikov.stoa.v1.Value
	2,  // 20: github.com.vontikov.stoa.v1.Stoa.DictionarySize:output_type -> github.com.vontikov.stoa.v1.Value
	5,  // 21: github.com.vontikov.stoa.v1.Stoa.DictionaryClear:output_type -> github.com.vontikov.stoa.v1.Empty
	6,  // 22: github.com.vontikov.stoa.v1.Stoa.DictionaryPutIfAbsent:output_type -> github.com.vontikov.stoa.v1.Result
	2,  // 23: github.com.vontikov.stoa.v1.Stoa.DictionaryPut:output_type -> github.com.vontikov.stoa.v1.Value
	2,  // 24: github.com.vontikov.stoa.v1.Stoa.DictionaryGet:output_type -> github.com.vontikov.stoa.v1.Value
	6,  // 25: github.com.vontikov.stoa.v1.Stoa.DictionaryRemove:output_type -> github.com.vontikov.stoa.v1.Result
	4,  // 26: github.com.vontikov.stoa.v1.Stoa.DictionaryRange:output_type -> github.com.vontikov.stoa.v1.KeyValue
	6,  // 27: github.com.vontikov.stoa.v1.Stoa.MutexTryLock:output_type -> github.com.vontikov.stoa.v1.Result
	6,  // 28: github.com.vontikov.stoa.v1.Stoa.MutexUnlock:output_type -> github.com.vontikov.stoa.v1.Result
	5,  // 29: github.com.vontikov.stoa.v1.Stoa.Ping:output_type -> github.com.vontikov.stoa.v1.Empty
	15, // [15:30] is the sub-list for method output_type
	0,  // [0:15] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_stoa_proto_init() }
func file_stoa_proto_init() {
	if File_stoa_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_stoa_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entity); i {
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
		file_stoa_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientId); i {
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
		file_stoa_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
		file_stoa_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Key); i {
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
		file_stoa_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValue); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stoa_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stoa_proto_goTypes,
		DependencyIndexes: file_stoa_proto_depIdxs,
		MessageInfos:      file_stoa_proto_msgTypes,
	}.Build()
	File_stoa_proto = out.File
	file_stoa_proto_rawDesc = nil
	file_stoa_proto_goTypes = nil
	file_stoa_proto_depIdxs = nil
}
