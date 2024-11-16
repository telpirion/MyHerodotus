// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: conversation.proto

package main

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConversationBit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BotResponse string `protobuf:"bytes,1,opt,name=bot_response,json=botResponse,proto3" json:"bot_response,omitempty"`
	UserQuery   string `protobuf:"bytes,2,opt,name=user_query,json=userQuery,proto3" json:"user_query,omitempty"`
	Model       string `protobuf:"bytes,3,opt,name=model,proto3" json:"model,omitempty"`
	Prompt      string `protobuf:"bytes,4,opt,name=prompt,proto3" json:"prompt,omitempty"`
	Created     int64  `protobuf:"varint,5,opt,name=created,proto3" json:"created,omitempty"`
	TokenCount  int32  `protobuf:"varint,6,opt,name=token_count,json=tokenCount,proto3" json:"token_count,omitempty"`
	Rating      string `protobuf:"bytes,7,opt,name=rating,proto3" json:"rating,omitempty"`
}

func (x *ConversationBit) Reset() {
	*x = ConversationBit{}
	mi := &file_conversation_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConversationBit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConversationBit) ProtoMessage() {}

func (x *ConversationBit) ProtoReflect() protoreflect.Message {
	mi := &file_conversation_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConversationBit.ProtoReflect.Descriptor instead.
func (*ConversationBit) Descriptor() ([]byte, []int) {
	return file_conversation_proto_rawDescGZIP(), []int{0}
}

func (x *ConversationBit) GetBotResponse() string {
	if x != nil {
		return x.BotResponse
	}
	return ""
}

func (x *ConversationBit) GetUserQuery() string {
	if x != nil {
		return x.UserQuery
	}
	return ""
}

func (x *ConversationBit) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *ConversationBit) GetPrompt() string {
	if x != nil {
		return x.Prompt
	}
	return ""
}

func (x *ConversationBit) GetCreated() int64 {
	if x != nil {
		return x.Created
	}
	return 0
}

func (x *ConversationBit) GetTokenCount() int32 {
	if x != nil {
		return x.TokenCount
	}
	return 0
}

func (x *ConversationBit) GetRating() string {
	if x != nil {
		return x.Rating
	}
	return ""
}

var File_conversation_proto protoreflect.FileDescriptor

var file_conversation_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6d, 0x79, 0x68, 0x65, 0x72, 0x6f, 0x64, 0x6f, 0x74, 0x75,
	0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xd4, 0x01, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x42, 0x69, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6f, 0x74, 0x5f, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x6f,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75,
	0x73, 0x65, 0x72, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x16,
	0x0a, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x42, 0x16, 0x5a, 0x14, 0x6d, 0x79, 0x68,
	0x65, 0x72, 0x6f, 0x64, 0x6f, 0x74, 0x75, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x69,
	0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_conversation_proto_rawDescOnce sync.Once
	file_conversation_proto_rawDescData = file_conversation_proto_rawDesc
)

func file_conversation_proto_rawDescGZIP() []byte {
	file_conversation_proto_rawDescOnce.Do(func() {
		file_conversation_proto_rawDescData = protoimpl.X.CompressGZIP(file_conversation_proto_rawDescData)
	})
	return file_conversation_proto_rawDescData
}

var file_conversation_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_conversation_proto_goTypes = []any{
	(*ConversationBit)(nil), // 0: myherodotus.ConversationBit
}
var file_conversation_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_conversation_proto_init() }
func file_conversation_proto_init() {
	if File_conversation_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_conversation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_conversation_proto_goTypes,
		DependencyIndexes: file_conversation_proto_depIdxs,
		MessageInfos:      file_conversation_proto_msgTypes,
	}.Build()
	File_conversation_proto = out.File
	file_conversation_proto_rawDesc = nil
	file_conversation_proto_goTypes = nil
	file_conversation_proto_depIdxs = nil
}
