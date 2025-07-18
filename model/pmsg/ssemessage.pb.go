// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: proto/ssemessage.proto

package pmsg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SseMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UidList       []string               `protobuf:"bytes,1,rep,name=uid_list,json=uidList,proto3" json:"uid_list,omitempty"`                            // uid列表
	Data          []byte                 `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`                                                 // 数据
	MessageId     MessageId              `protobuf:"varint,3,opt,name=message_id,json=messageId,proto3,enum=pmsg.MessageId" json:"message_id,omitempty"` // 消息类型
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SseMessage) Reset() {
	*x = SseMessage{}
	mi := &file_proto_ssemessage_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SseMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SseMessage) ProtoMessage() {}

func (x *SseMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_ssemessage_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SseMessage.ProtoReflect.Descriptor instead.
func (*SseMessage) Descriptor() ([]byte, []int) {
	return file_proto_ssemessage_proto_rawDescGZIP(), []int{0}
}

func (x *SseMessage) GetUidList() []string {
	if x != nil {
		return x.UidList
	}
	return nil
}

func (x *SseMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *SseMessage) GetMessageId() MessageId {
	if x != nil {
		return x.MessageId
	}
	return MessageId_Unknown
}

var File_proto_ssemessage_proto protoreflect.FileDescriptor

const file_proto_ssemessage_proto_rawDesc = "" +
	"\n" +
	"\x16proto/ssemessage.proto\x12\x04pmsg\x1a\x13proto/constid.proto\"k\n" +
	"\n" +
	"SseMessage\x12\x19\n" +
	"\buid_list\x18\x01 \x03(\tR\auidList\x12\x12\n" +
	"\x04data\x18\x02 \x01(\fR\x04data\x12.\n" +
	"\n" +
	"message_id\x18\x03 \x01(\x0e2\x0f.pmsg.MessageIdR\tmessageIdB\bZ\x06./pmsgb\x06proto3"

var (
	file_proto_ssemessage_proto_rawDescOnce sync.Once
	file_proto_ssemessage_proto_rawDescData []byte
)

func file_proto_ssemessage_proto_rawDescGZIP() []byte {
	file_proto_ssemessage_proto_rawDescOnce.Do(func() {
		file_proto_ssemessage_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_ssemessage_proto_rawDesc), len(file_proto_ssemessage_proto_rawDesc)))
	})
	return file_proto_ssemessage_proto_rawDescData
}

var file_proto_ssemessage_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_ssemessage_proto_goTypes = []any{
	(*SseMessage)(nil), // 0: pmsg.SseMessage
	(MessageId)(0),     // 1: pmsg.MessageId
}
var file_proto_ssemessage_proto_depIdxs = []int32{
	1, // 0: pmsg.SseMessage.message_id:type_name -> pmsg.MessageId
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_ssemessage_proto_init() }
func file_proto_ssemessage_proto_init() {
	if File_proto_ssemessage_proto != nil {
		return
	}
	file_proto_constid_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_ssemessage_proto_rawDesc), len(file_proto_ssemessage_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_ssemessage_proto_goTypes,
		DependencyIndexes: file_proto_ssemessage_proto_depIdxs,
		MessageInfos:      file_proto_ssemessage_proto_msgTypes,
	}.Build()
	File_proto_ssemessage_proto = out.File
	file_proto_ssemessage_proto_goTypes = nil
	file_proto_ssemessage_proto_depIdxs = nil
}
