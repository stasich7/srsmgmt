// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.1
// source: srsmgmt.proto

package srsmgmt_proto

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Stream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	App       string                 `protobuf:"bytes,2,opt,name=app,proto3" json:"app,omitempty"`
	Password  string                 `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Status    int32                  `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
	Hls       string                 `protobuf:"bytes,5,opt,name=hls,proto3" json:"hls,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
	ClientId  string                 `protobuf:"bytes,8,opt,name=clientId,proto3" json:"clientId,omitempty"`
	StartedAt *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=startedAt,proto3" json:"startedAt,omitempty"`
	StopedAt  *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=stopedAt,proto3" json:"stopedAt,omitempty"`
}

func (x *Stream) Reset() {
	*x = Stream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_srsmgmt_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stream) ProtoMessage() {}

func (x *Stream) ProtoReflect() protoreflect.Message {
	mi := &file_srsmgmt_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stream.ProtoReflect.Descriptor instead.
func (*Stream) Descriptor() ([]byte, []int) {
	return file_srsmgmt_proto_rawDescGZIP(), []int{0}
}

func (x *Stream) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Stream) GetApp() string {
	if x != nil {
		return x.App
	}
	return ""
}

func (x *Stream) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Stream) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Stream) GetHls() string {
	if x != nil {
		return x.Hls
	}
	return ""
}

func (x *Stream) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Stream) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Stream) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Stream) GetStartedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StartedAt
	}
	return nil
}

func (x *Stream) GetStopedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StopedAt
	}
	return nil
}

type GetStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetStreamRequest) Reset() {
	*x = GetStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_srsmgmt_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStreamRequest) ProtoMessage() {}

func (x *GetStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_srsmgmt_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStreamRequest.ProtoReflect.Descriptor instead.
func (*GetStreamRequest) Descriptor() ([]byte, []int) {
	return file_srsmgmt_proto_rawDescGZIP(), []int{1}
}

func (x *GetStreamRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetStreamReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stream *Stream `protobuf:"bytes,1,opt,name=stream,proto3" json:"stream,omitempty"`
}

func (x *GetStreamReply) Reset() {
	*x = GetStreamReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_srsmgmt_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStreamReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStreamReply) ProtoMessage() {}

func (x *GetStreamReply) ProtoReflect() protoreflect.Message {
	mi := &file_srsmgmt_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStreamReply.ProtoReflect.Descriptor instead.
func (*GetStreamReply) Descriptor() ([]byte, []int) {
	return file_srsmgmt_proto_rawDescGZIP(), []int{2}
}

func (x *GetStreamReply) GetStream() *Stream {
	if x != nil {
		return x.Stream
	}
	return nil
}

type CreateStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stream *Stream `protobuf:"bytes,1,opt,name=stream,proto3" json:"stream,omitempty"`
}

func (x *CreateStreamRequest) Reset() {
	*x = CreateStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_srsmgmt_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStreamRequest) ProtoMessage() {}

func (x *CreateStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_srsmgmt_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStreamRequest.ProtoReflect.Descriptor instead.
func (*CreateStreamRequest) Descriptor() ([]byte, []int) {
	return file_srsmgmt_proto_rawDescGZIP(), []int{3}
}

func (x *CreateStreamRequest) GetStream() *Stream {
	if x != nil {
		return x.Stream
	}
	return nil
}

type CreateStreamReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stream *Stream `protobuf:"bytes,1,opt,name=stream,proto3" json:"stream,omitempty"`
}

func (x *CreateStreamReply) Reset() {
	*x = CreateStreamReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_srsmgmt_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateStreamReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStreamReply) ProtoMessage() {}

func (x *CreateStreamReply) ProtoReflect() protoreflect.Message {
	mi := &file_srsmgmt_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStreamReply.ProtoReflect.Descriptor instead.
func (*CreateStreamReply) Descriptor() ([]byte, []int) {
	return file_srsmgmt_proto_rawDescGZIP(), []int{4}
}

func (x *CreateStreamReply) GetStream() *Stream {
	if x != nil {
		return x.Stream
	}
	return nil
}

type DeleteStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteStreamRequest) Reset() {
	*x = DeleteStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_srsmgmt_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStreamRequest) ProtoMessage() {}

func (x *DeleteStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_srsmgmt_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStreamRequest.ProtoReflect.Descriptor instead.
func (*DeleteStreamRequest) Descriptor() ([]byte, []int) {
	return file_srsmgmt_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteStreamRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteStreamReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteStreamReply) Reset() {
	*x = DeleteStreamReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_srsmgmt_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteStreamReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStreamReply) ProtoMessage() {}

func (x *DeleteStreamReply) ProtoReflect() protoreflect.Message {
	mi := &file_srsmgmt_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStreamReply.ProtoReflect.Descriptor instead.
func (*DeleteStreamReply) Descriptor() ([]byte, []int) {
	return file_srsmgmt_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteStreamReply) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_srsmgmt_proto protoreflect.FileDescriptor

var file_srsmgmt_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x72, 0x73, 0x6d, 0x67, 0x6d, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xf2, 0x02, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x61, 0x70, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x70, 0x70, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x68, 0x6c, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x68, 0x6c, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x38, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x36, 0x0a, 0x08, 0x73, 0x74, 0x6f, 0x70, 0x65, 0x64, 0x41, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x73,
	0x74, 0x6f, 0x70, 0x65, 0x64, 0x41, 0x74, 0x22, 0x22, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x34, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x22, 0x0a,
	0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x70, 0x62, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x22, 0x39, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x37, 0x0a, 0x11,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x22, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x06, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x25, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x23, 0x0a, 0x11,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x32, 0xdc, 0x01, 0x0a, 0x07, 0x53, 0x72, 0x73, 0x4d, 0x67, 0x6d, 0x74, 0x12, 0x4d, 0x0a,
	0x09, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x14, 0x2e, 0x70, 0x62, 0x2e,
	0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x40, 0x0a, 0x0c,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x17, 0x2e, 0x70,
	0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x40,
	0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x17,
	0x2e, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x42, 0x0f, 0x5a, 0x0d, 0x73, 0x72, 0x73, 0x6d, 0x67, 0x6d, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_srsmgmt_proto_rawDescOnce sync.Once
	file_srsmgmt_proto_rawDescData = file_srsmgmt_proto_rawDesc
)

func file_srsmgmt_proto_rawDescGZIP() []byte {
	file_srsmgmt_proto_rawDescOnce.Do(func() {
		file_srsmgmt_proto_rawDescData = protoimpl.X.CompressGZIP(file_srsmgmt_proto_rawDescData)
	})
	return file_srsmgmt_proto_rawDescData
}

var file_srsmgmt_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_srsmgmt_proto_goTypes = []interface{}{
	(*Stream)(nil),                // 0: pb.Stream
	(*GetStreamRequest)(nil),      // 1: pb.GetStreamRequest
	(*GetStreamReply)(nil),        // 2: pb.GetStreamReply
	(*CreateStreamRequest)(nil),   // 3: pb.CreateStreamRequest
	(*CreateStreamReply)(nil),     // 4: pb.CreateStreamReply
	(*DeleteStreamRequest)(nil),   // 5: pb.DeleteStreamRequest
	(*DeleteStreamReply)(nil),     // 6: pb.DeleteStreamReply
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_srsmgmt_proto_depIdxs = []int32{
	7,  // 0: pb.Stream.createdAt:type_name -> google.protobuf.Timestamp
	7,  // 1: pb.Stream.updatedAt:type_name -> google.protobuf.Timestamp
	7,  // 2: pb.Stream.startedAt:type_name -> google.protobuf.Timestamp
	7,  // 3: pb.Stream.stopedAt:type_name -> google.protobuf.Timestamp
	0,  // 4: pb.GetStreamReply.stream:type_name -> pb.Stream
	0,  // 5: pb.CreateStreamRequest.stream:type_name -> pb.Stream
	0,  // 6: pb.CreateStreamReply.stream:type_name -> pb.Stream
	1,  // 7: pb.SrsMgmt.GetStream:input_type -> pb.GetStreamRequest
	3,  // 8: pb.SrsMgmt.CreateStream:input_type -> pb.CreateStreamRequest
	5,  // 9: pb.SrsMgmt.DeleteStream:input_type -> pb.DeleteStreamRequest
	2,  // 10: pb.SrsMgmt.GetStream:output_type -> pb.GetStreamReply
	4,  // 11: pb.SrsMgmt.CreateStream:output_type -> pb.CreateStreamReply
	6,  // 12: pb.SrsMgmt.DeleteStream:output_type -> pb.DeleteStreamReply
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_srsmgmt_proto_init() }
func file_srsmgmt_proto_init() {
	if File_srsmgmt_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_srsmgmt_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stream); i {
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
		file_srsmgmt_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStreamRequest); i {
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
		file_srsmgmt_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStreamReply); i {
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
		file_srsmgmt_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateStreamRequest); i {
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
		file_srsmgmt_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateStreamReply); i {
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
		file_srsmgmt_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteStreamRequest); i {
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
		file_srsmgmt_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteStreamReply); i {
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
			RawDescriptor: file_srsmgmt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_srsmgmt_proto_goTypes,
		DependencyIndexes: file_srsmgmt_proto_depIdxs,
		MessageInfos:      file_srsmgmt_proto_msgTypes,
	}.Build()
	File_srsmgmt_proto = out.File
	file_srsmgmt_proto_rawDesc = nil
	file_srsmgmt_proto_goTypes = nil
	file_srsmgmt_proto_depIdxs = nil
}
