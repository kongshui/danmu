// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: proto/errorstatus.proto

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

// 错误状态
type ErrorStatus int32

const (
	ErrorStatus_UnknownError                   ErrorStatus = 0  // 未知错误
	ErrorStatus_MatchError                     ErrorStatus = 1  // 匹配错误
	ErrorStatus_MatchRoomIdError               ErrorStatus = 2  // 房间ID错误
	ErrorStatus_MatchOpenIdError               ErrorStatus = 3  // 用户ID错误
	ErrorStatus_MatchBattleRoomIdError         ErrorStatus = 4  // 战斗房间ID错误
	ErrorStatus_SetMatchStatusError            ErrorStatus = 5  // 设置匹配状态错误
	ErrorStatus_ProtoMarshalError              ErrorStatus = 6  // 协议序列化错误
	ErrorStatus_GetUuidError                   ErrorStatus = 7  // 获取UUID错误
	ErrorStatus_GetRedisError                  ErrorStatus = 8  // 获取Redis错误
	ErrorStatus_SetRedisError                  ErrorStatus = 9  // 设置Redis错误
	ErrorStatus_GetMatchStatusError            ErrorStatus = 10 // 获取匹配状态错误
	ErrorStatus_SendMessageError               ErrorStatus = 11 // 发送消息错误
	ErrorStatus_ProtoUnmarshalError            ErrorStatus = 12 // 协议反序列化错误
	ErrorStatus_MatchBattleV1SetStartTimeError ErrorStatus = 13 // 设置超时错误
)

// Enum value maps for ErrorStatus.
var (
	ErrorStatus_name = map[int32]string{
		0:  "UnknownError",
		1:  "MatchError",
		2:  "MatchRoomIdError",
		3:  "MatchOpenIdError",
		4:  "MatchBattleRoomIdError",
		5:  "SetMatchStatusError",
		6:  "ProtoMarshalError",
		7:  "GetUuidError",
		8:  "GetRedisError",
		9:  "SetRedisError",
		10: "GetMatchStatusError",
		11: "SendMessageError",
		12: "ProtoUnmarshalError",
		13: "MatchBattleV1SetStartTimeError",
	}
	ErrorStatus_value = map[string]int32{
		"UnknownError":                   0,
		"MatchError":                     1,
		"MatchRoomIdError":               2,
		"MatchOpenIdError":               3,
		"MatchBattleRoomIdError":         4,
		"SetMatchStatusError":            5,
		"ProtoMarshalError":              6,
		"GetUuidError":                   7,
		"GetRedisError":                  8,
		"SetRedisError":                  9,
		"GetMatchStatusError":            10,
		"SendMessageError":               11,
		"ProtoUnmarshalError":            12,
		"MatchBattleV1SetStartTimeError": 13,
	}
)

func (x ErrorStatus) Enum() *ErrorStatus {
	p := new(ErrorStatus)
	*p = x
	return p
}

func (x ErrorStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_errorstatus_proto_enumTypes[0].Descriptor()
}

func (ErrorStatus) Type() protoreflect.EnumType {
	return &file_proto_errorstatus_proto_enumTypes[0]
}

func (x ErrorStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorStatus.Descriptor instead.
func (ErrorStatus) EnumDescriptor() ([]byte, []int) {
	return file_proto_errorstatus_proto_rawDescGZIP(), []int{0}
}

var File_proto_errorstatus_proto protoreflect.FileDescriptor

const file_proto_errorstatus_proto_rawDesc = "" +
	"\n" +
	"\x17proto/errorstatus.proto\x12\x04pmsg*\xcb\x02\n" +
	"\vErrorStatus\x12\x10\n" +
	"\fUnknownError\x10\x00\x12\x0e\n" +
	"\n" +
	"MatchError\x10\x01\x12\x14\n" +
	"\x10MatchRoomIdError\x10\x02\x12\x14\n" +
	"\x10MatchOpenIdError\x10\x03\x12\x1a\n" +
	"\x16MatchBattleRoomIdError\x10\x04\x12\x17\n" +
	"\x13SetMatchStatusError\x10\x05\x12\x15\n" +
	"\x11ProtoMarshalError\x10\x06\x12\x10\n" +
	"\fGetUuidError\x10\a\x12\x11\n" +
	"\rGetRedisError\x10\b\x12\x11\n" +
	"\rSetRedisError\x10\t\x12\x17\n" +
	"\x13GetMatchStatusError\x10\n" +
	"\x12\x14\n" +
	"\x10SendMessageError\x10\v\x12\x17\n" +
	"\x13ProtoUnmarshalError\x10\f\x12\"\n" +
	"\x1eMatchBattleV1SetStartTimeError\x10\rB\bZ\x06./pmsgb\x06proto3"

var (
	file_proto_errorstatus_proto_rawDescOnce sync.Once
	file_proto_errorstatus_proto_rawDescData []byte
)

func file_proto_errorstatus_proto_rawDescGZIP() []byte {
	file_proto_errorstatus_proto_rawDescOnce.Do(func() {
		file_proto_errorstatus_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_errorstatus_proto_rawDesc), len(file_proto_errorstatus_proto_rawDesc)))
	})
	return file_proto_errorstatus_proto_rawDescData
}

var file_proto_errorstatus_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_errorstatus_proto_goTypes = []any{
	(ErrorStatus)(0), // 0: pmsg.ErrorStatus
}
var file_proto_errorstatus_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_errorstatus_proto_init() }
func file_proto_errorstatus_proto_init() {
	if File_proto_errorstatus_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_errorstatus_proto_rawDesc), len(file_proto_errorstatus_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_errorstatus_proto_goTypes,
		DependencyIndexes: file_proto_errorstatus_proto_depIdxs,
		EnumInfos:         file_proto_errorstatus_proto_enumTypes,
	}.Build()
	File_proto_errorstatus_proto = out.File
	file_proto_errorstatus_proto_goTypes = nil
	file_proto_errorstatus_proto_depIdxs = nil
}
