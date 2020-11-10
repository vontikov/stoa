// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: plugin.proto

package pb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type QueueEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value     []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	TtlMillis int64  `protobuf:"varint,3,opt,name=ttl_millis,json=ttlMillis,proto3" json:"ttl_millis,omitempty"`
}

func (x *QueueEntry) Reset() {
	*x = QueueEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueEntry) ProtoMessage() {}

func (x *QueueEntry) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueEntry.ProtoReflect.Descriptor instead.
func (*QueueEntry) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *QueueEntry) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *QueueEntry) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *QueueEntry) GetTtlMillis() int64 {
	if x != nil {
		return x.TtlMillis
	}
	return 0
}

type DictionaryEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Key       []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value     []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	TtlMillis int64  `protobuf:"varint,4,opt,name=ttl_millis,json=ttlMillis,proto3" json:"ttl_millis,omitempty"`
}

func (x *DictionaryEntry) Reset() {
	*x = DictionaryEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DictionaryEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DictionaryEntry) ProtoMessage() {}

func (x *DictionaryEntry) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DictionaryEntry.ProtoReflect.Descriptor instead.
func (*DictionaryEntry) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *DictionaryEntry) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DictionaryEntry) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *DictionaryEntry) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *DictionaryEntry) GetTtlMillis() int64 {
	if x != nil {
		return x.TtlMillis
	}
	return 0
}

var File_plugin_proto protoreflect.FileDescriptor

var file_plugin_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69,
	0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x1a, 0x0c, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x67, 0x0a, 0x0a, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x7a, 0x02, 0x10, 0x01, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x74, 0x6c, 0x5f, 0x6d, 0x69, 0x6c, 0x6c,
	0x69, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x74, 0x6c, 0x4d, 0x69, 0x6c,
	0x6c, 0x69, 0x73, 0x22, 0x87, 0x01, 0x0a, 0x0f, 0x44, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61,
	0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x7a, 0x02, 0x10, 0x01, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x1d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x7a, 0x02, 0x10, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x74, 0x74, 0x6c, 0x5f, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x74, 0x74, 0x6c, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x32, 0xed, 0x01,
	0x0a, 0x08, 0x45, 0x76, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x68, 0x0a, 0x05, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x12, 0x27, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76,
	0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x1a, 0x23, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b,
	0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x12, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x12, 0x77, 0x0a, 0x0a, 0x44, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61,
	0x72, 0x79, 0x12, 0x2c, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x76, 0x6f, 0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31,
	0x2e, 0x44, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x1a, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x6f,
	0x6e, 0x74, 0x69, 0x6b, 0x6f, 0x76, 0x2e, 0x73, 0x74, 0x6f, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f,
	0x76, 0x31, 0x2f, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x42, 0x06, 0x5a,
	0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plugin_proto_rawDescOnce sync.Once
	file_plugin_proto_rawDescData = file_plugin_proto_rawDesc
)

func file_plugin_proto_rawDescGZIP() []byte {
	file_plugin_proto_rawDescOnce.Do(func() {
		file_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_proto_rawDescData)
	})
	return file_plugin_proto_rawDescData
}

var file_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_plugin_proto_goTypes = []interface{}{
	(*QueueEntry)(nil),      // 0: github.com.vontikov.stoa.v1.QueueEntry
	(*DictionaryEntry)(nil), // 1: github.com.vontikov.stoa.v1.DictionaryEntry
	(*Result)(nil),          // 2: github.com.vontikov.stoa.v1.Result
}
var file_plugin_proto_depIdxs = []int32{
	0, // 0: github.com.vontikov.stoa.v1.Eviction.Queue:input_type -> github.com.vontikov.stoa.v1.QueueEntry
	1, // 1: github.com.vontikov.stoa.v1.Eviction.Dictionary:input_type -> github.com.vontikov.stoa.v1.DictionaryEntry
	2, // 2: github.com.vontikov.stoa.v1.Eviction.Queue:output_type -> github.com.vontikov.stoa.v1.Result
	2, // 3: github.com.vontikov.stoa.v1.Eviction.Dictionary:output_type -> github.com.vontikov.stoa.v1.Result
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_plugin_proto_init() }
func file_plugin_proto_init() {
	if File_plugin_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueEntry); i {
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
		file_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DictionaryEntry); i {
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
			RawDescriptor: file_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_proto_goTypes,
		DependencyIndexes: file_plugin_proto_depIdxs,
		MessageInfos:      file_plugin_proto_msgTypes,
	}.Build()
	File_plugin_proto = out.File
	file_plugin_proto_rawDesc = nil
	file_plugin_proto_goTypes = nil
	file_plugin_proto_depIdxs = nil
}
