// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: proto/roundupload.proto

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

type RoundUploadMessage struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	RoomId          string                 `protobuf:"bytes,1,opt,name=RoomId,proto3" json:"RoomId,omitempty"`                   //主播房间id
	AnchorOpenId    string                 `protobuf:"bytes,2,opt,name=AnchorOpenId,proto3" json:"AnchorOpenId,omitempty"`       //主播id
	RoundId         int64                  `protobuf:"varint,3,opt,name=RoundId,proto3" json:"RoundId,omitempty"`                //对局id
	GroupResultList []*GroupResult         `protobuf:"bytes,4,rep,name=GroupResultList,proto3" json:"GroupResultList,omitempty"` //对局结果列表
	GroupUserList   []*GroupUser           `protobuf:"bytes,5,rep,name=GroupUserList,proto3" json:"GroupUserList,omitempty"`     //玩家列表
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *RoundUploadMessage) Reset() {
	*x = RoundUploadMessage{}
	mi := &file_proto_roundupload_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RoundUploadMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoundUploadMessage) ProtoMessage() {}

func (x *RoundUploadMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_roundupload_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoundUploadMessage.ProtoReflect.Descriptor instead.
func (*RoundUploadMessage) Descriptor() ([]byte, []int) {
	return file_proto_roundupload_proto_rawDescGZIP(), []int{0}
}

func (x *RoundUploadMessage) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *RoundUploadMessage) GetAnchorOpenId() string {
	if x != nil {
		return x.AnchorOpenId
	}
	return ""
}

func (x *RoundUploadMessage) GetRoundId() int64 {
	if x != nil {
		return x.RoundId
	}
	return 0
}

func (x *RoundUploadMessage) GetGroupResultList() []*GroupResult {
	if x != nil {
		return x.GroupResultList
	}
	return nil
}

func (x *RoundUploadMessage) GetGroupUserList() []*GroupUser {
	if x != nil {
		return x.GroupUserList
	}
	return nil
}

var File_proto_roundupload_proto protoreflect.FileDescriptor

const file_proto_roundupload_proto_rawDesc = "" +
	"\n" +
	"\x17proto/roundupload.proto\x12\x04pmsg\x1a\x17proto/groupresult.proto\x1a\x15proto/groupuser.proto\"\xde\x01\n" +
	"\x12RoundUploadMessage\x12\x16\n" +
	"\x06RoomId\x18\x01 \x01(\tR\x06RoomId\x12\"\n" +
	"\fAnchorOpenId\x18\x02 \x01(\tR\fAnchorOpenId\x12\x18\n" +
	"\aRoundId\x18\x03 \x01(\x03R\aRoundId\x12;\n" +
	"\x0fGroupResultList\x18\x04 \x03(\v2\x11.pmsg.GroupResultR\x0fGroupResultList\x125\n" +
	"\rGroupUserList\x18\x05 \x03(\v2\x0f.pmsg.GroupUserR\rGroupUserListB\bZ\x06./pmsgb\x06proto3"

var (
	file_proto_roundupload_proto_rawDescOnce sync.Once
	file_proto_roundupload_proto_rawDescData []byte
)

func file_proto_roundupload_proto_rawDescGZIP() []byte {
	file_proto_roundupload_proto_rawDescOnce.Do(func() {
		file_proto_roundupload_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_roundupload_proto_rawDesc), len(file_proto_roundupload_proto_rawDesc)))
	})
	return file_proto_roundupload_proto_rawDescData
}

var file_proto_roundupload_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_roundupload_proto_goTypes = []any{
	(*RoundUploadMessage)(nil), // 0: pmsg.RoundUploadMessage
	(*GroupResult)(nil),        // 1: pmsg.GroupResult
	(*GroupUser)(nil),          // 2: pmsg.GroupUser
}
var file_proto_roundupload_proto_depIdxs = []int32{
	1, // 0: pmsg.RoundUploadMessage.GroupResultList:type_name -> pmsg.GroupResult
	2, // 1: pmsg.RoundUploadMessage.GroupUserList:type_name -> pmsg.GroupUser
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_roundupload_proto_init() }
func file_proto_roundupload_proto_init() {
	if File_proto_roundupload_proto != nil {
		return
	}
	file_proto_groupresult_proto_init()
	file_proto_groupuser_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_roundupload_proto_rawDesc), len(file_proto_roundupload_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_roundupload_proto_goTypes,
		DependencyIndexes: file_proto_roundupload_proto_depIdxs,
		MessageInfos:      file_proto_roundupload_proto_msgTypes,
	}.Build()
	File_proto_roundupload_proto = out.File
	file_proto_roundupload_proto_goTypes = nil
	file_proto_roundupload_proto_depIdxs = nil
}
