// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: proto/singeroomaddgroup.proto

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

type AddGroupMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GroupId       int32                  `protobuf:"varint,1,opt,name=GroupId,proto3" json:"GroupId,omitempty"`    //组Id
	GroupName     string                 `protobuf:"bytes,2,opt,name=GroupName,proto3" json:"GroupName,omitempty"` //组名
	UserId        string                 `protobuf:"bytes,3,opt,name=UserId,proto3" json:"UserId,omitempty"`       //加入的用户Id
	UdpAddr       string                 `protobuf:"bytes,4,opt,name=UdpAddr,proto3" json:"UdpAddr,omitempty"`     //加入组的玩家的udpAddr地址，暂时不用
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddGroupMessage) Reset() {
	*x = AddGroupMessage{}
	mi := &file_proto_singeroomaddgroup_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddGroupMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddGroupMessage) ProtoMessage() {}

func (x *AddGroupMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_singeroomaddgroup_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddGroupMessage.ProtoReflect.Descriptor instead.
func (*AddGroupMessage) Descriptor() ([]byte, []int) {
	return file_proto_singeroomaddgroup_proto_rawDescGZIP(), []int{0}
}

func (x *AddGroupMessage) GetGroupId() int32 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *AddGroupMessage) GetGroupName() string {
	if x != nil {
		return x.GroupName
	}
	return ""
}

func (x *AddGroupMessage) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *AddGroupMessage) GetUdpAddr() string {
	if x != nil {
		return x.UdpAddr
	}
	return ""
}

var File_proto_singeroomaddgroup_proto protoreflect.FileDescriptor

const file_proto_singeroomaddgroup_proto_rawDesc = "" +
	"\n" +
	"\x1dproto/singeroomaddgroup.proto\x12\x04pmsg\"{\n" +
	"\x0fAddGroupMessage\x12\x18\n" +
	"\aGroupId\x18\x01 \x01(\x05R\aGroupId\x12\x1c\n" +
	"\tGroupName\x18\x02 \x01(\tR\tGroupName\x12\x16\n" +
	"\x06UserId\x18\x03 \x01(\tR\x06UserId\x12\x18\n" +
	"\aUdpAddr\x18\x04 \x01(\tR\aUdpAddrB\bZ\x06./pmsgb\x06proto3"

var (
	file_proto_singeroomaddgroup_proto_rawDescOnce sync.Once
	file_proto_singeroomaddgroup_proto_rawDescData []byte
)

func file_proto_singeroomaddgroup_proto_rawDescGZIP() []byte {
	file_proto_singeroomaddgroup_proto_rawDescOnce.Do(func() {
		file_proto_singeroomaddgroup_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_singeroomaddgroup_proto_rawDesc), len(file_proto_singeroomaddgroup_proto_rawDesc)))
	})
	return file_proto_singeroomaddgroup_proto_rawDescData
}

var file_proto_singeroomaddgroup_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_singeroomaddgroup_proto_goTypes = []any{
	(*AddGroupMessage)(nil), // 0: pmsg.AddGroupMessage
}
var file_proto_singeroomaddgroup_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_singeroomaddgroup_proto_init() }
func file_proto_singeroomaddgroup_proto_init() {
	if File_proto_singeroomaddgroup_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_singeroomaddgroup_proto_rawDesc), len(file_proto_singeroomaddgroup_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_singeroomaddgroup_proto_goTypes,
		DependencyIndexes: file_proto_singeroomaddgroup_proto_depIdxs,
		MessageInfos:      file_proto_singeroomaddgroup_proto_msgTypes,
	}.Build()
	File_proto_singeroomaddgroup_proto = out.File
	file_proto_singeroomaddgroup_proto_goTypes = nil
	file_proto_singeroomaddgroup_proto_depIdxs = nil
}
