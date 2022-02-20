// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: pkg/proto/ledger.proto

package proto

import (
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

type ReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WriterId string `protobuf:"bytes,1,opt,name=writer_id,json=writerId,proto3" json:"writer_id,omitempty"`
	ReaderId string `protobuf:"bytes,2,opt,name=reader_id,json=readerId,proto3" json:"reader_id,omitempty"`
}

func (x *ReadRequest) Reset() {
	*x = ReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRequest) ProtoMessage() {}

func (x *ReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRequest.ProtoReflect.Descriptor instead.
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{0}
}

func (x *ReadRequest) GetWriterId() string {
	if x != nil {
		return x.WriterId
	}
	return ""
}

func (x *ReadRequest) GetReaderId() string {
	if x != nil {
		return x.ReaderId
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset uint64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Data   []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *Message) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type CommitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WriterId string `protobuf:"bytes,1,opt,name=writer_id,json=writerId,proto3" json:"writer_id,omitempty"`
	ReaderId string `protobuf:"bytes,2,opt,name=reader_id,json=readerId,proto3" json:"reader_id,omitempty"`
	Offset   uint64 `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *CommitRequest) Reset() {
	*x = CommitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitRequest) ProtoMessage() {}

func (x *CommitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitRequest.ProtoReflect.Descriptor instead.
func (*CommitRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{2}
}

func (x *CommitRequest) GetWriterId() string {
	if x != nil {
		return x.WriterId
	}
	return ""
}

func (x *CommitRequest) GetReaderId() string {
	if x != nil {
		return x.ReaderId
	}
	return ""
}

func (x *CommitRequest) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type CommitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CommitResponse) Reset() {
	*x = CommitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitResponse) ProtoMessage() {}

func (x *CommitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitResponse.ProtoReflect.Descriptor instead.
func (*CommitResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{3}
}

type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WriterId string `protobuf:"bytes,1,opt,name=writer_id,json=writerId,proto3" json:"writer_id,omitempty"`
	Data     []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{4}
}

func (x *WriteRequest) GetWriterId() string {
	if x != nil {
		return x.WriterId
	}
	return ""
}

func (x *WriteRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type WriteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset uint64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *WriteResponse) Reset() {
	*x = WriteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteResponse) ProtoMessage() {}

func (x *WriteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteResponse.ProtoReflect.Descriptor instead.
func (*WriteResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{5}
}

func (x *WriteResponse) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type WriterCloseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WriterId string `protobuf:"bytes,1,opt,name=writer_id,json=writerId,proto3" json:"writer_id,omitempty"`
}

func (x *WriterCloseRequest) Reset() {
	*x = WriterCloseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriterCloseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriterCloseRequest) ProtoMessage() {}

func (x *WriterCloseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriterCloseRequest.ProtoReflect.Descriptor instead.
func (*WriterCloseRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{6}
}

func (x *WriterCloseRequest) GetWriterId() string {
	if x != nil {
		return x.WriterId
	}
	return ""
}

type WriterCloseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WriterCloseResponse) Reset() {
	*x = WriterCloseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriterCloseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriterCloseResponse) ProtoMessage() {}

func (x *WriterCloseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriterCloseResponse.ProtoReflect.Descriptor instead.
func (*WriterCloseResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{7}
}

type ReaderCloseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WriterId string `protobuf:"bytes,1,opt,name=writer_id,json=writerId,proto3" json:"writer_id,omitempty"`
	ReaderId string `protobuf:"bytes,2,opt,name=reader_id,json=readerId,proto3" json:"reader_id,omitempty"`
}

func (x *ReaderCloseRequest) Reset() {
	*x = ReaderCloseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReaderCloseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReaderCloseRequest) ProtoMessage() {}

func (x *ReaderCloseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReaderCloseRequest.ProtoReflect.Descriptor instead.
func (*ReaderCloseRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{8}
}

func (x *ReaderCloseRequest) GetWriterId() string {
	if x != nil {
		return x.WriterId
	}
	return ""
}

func (x *ReaderCloseRequest) GetReaderId() string {
	if x != nil {
		return x.ReaderId
	}
	return ""
}

type ReaderCloseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReaderCloseResponse) Reset() {
	*x = ReaderCloseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_ledger_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReaderCloseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReaderCloseResponse) ProtoMessage() {}

func (x *ReaderCloseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_ledger_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReaderCloseResponse.ProtoReflect.Descriptor instead.
func (*ReaderCloseResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_ledger_proto_rawDescGZIP(), []int{9}
}

var File_pkg_proto_ledger_proto protoreflect.FileDescriptor

var file_pkg_proto_ledger_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x65, 0x64, 0x67,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x0b, 0x52, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x72, 0x69, 0x74, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x35, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x61, 0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x72, 0x69,
	0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x72,
	0x69, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x10, 0x0a, 0x0e, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3f, 0x0a,
	0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x27,
	0x0a, 0x0d, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x31, 0x0a, 0x12, 0x57, 0x72, 0x69, 0x74, 0x65,
	0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x49, 0x64, 0x22, 0x15, 0x0a, 0x13, 0x57, 0x72,
	0x69, 0x74, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x4e, 0x0a, 0x12, 0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x72, 0x69, 0x74, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x15, 0x0a, 0x13, 0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xfb, 0x01, 0x0a, 0x06, 0x4c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12, 0x0c, 0x2e, 0x52, 0x65,
	0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x2b, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x12, 0x0e, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x28, 0x0a, 0x05, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x0d, 0x2e,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3a,
	0x0a, 0x0b, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x72, 0x12, 0x13, 0x2e,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x14, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0b, 0x43, 0x6c,
	0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x52, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0b, 0x5a, 0x09, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_ledger_proto_rawDescOnce sync.Once
	file_pkg_proto_ledger_proto_rawDescData = file_pkg_proto_ledger_proto_rawDesc
)

func file_pkg_proto_ledger_proto_rawDescGZIP() []byte {
	file_pkg_proto_ledger_proto_rawDescOnce.Do(func() {
		file_pkg_proto_ledger_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_ledger_proto_rawDescData)
	})
	return file_pkg_proto_ledger_proto_rawDescData
}

var file_pkg_proto_ledger_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_pkg_proto_ledger_proto_goTypes = []interface{}{
	(*ReadRequest)(nil),         // 0: ReadRequest
	(*Message)(nil),             // 1: Message
	(*CommitRequest)(nil),       // 2: CommitRequest
	(*CommitResponse)(nil),      // 3: CommitResponse
	(*WriteRequest)(nil),        // 4: WriteRequest
	(*WriteResponse)(nil),       // 5: WriteResponse
	(*WriterCloseRequest)(nil),  // 6: WriterCloseRequest
	(*WriterCloseResponse)(nil), // 7: WriterCloseResponse
	(*ReaderCloseRequest)(nil),  // 8: ReaderCloseRequest
	(*ReaderCloseResponse)(nil), // 9: ReaderCloseResponse
}
var file_pkg_proto_ledger_proto_depIdxs = []int32{
	0, // 0: Ledger.Read:input_type -> ReadRequest
	2, // 1: Ledger.Commit:input_type -> CommitRequest
	4, // 2: Ledger.Write:input_type -> WriteRequest
	6, // 3: Ledger.CloseWriter:input_type -> WriterCloseRequest
	8, // 4: Ledger.CloseReader:input_type -> ReaderCloseRequest
	1, // 5: Ledger.Read:output_type -> Message
	3, // 6: Ledger.Commit:output_type -> CommitResponse
	5, // 7: Ledger.Write:output_type -> WriteResponse
	7, // 8: Ledger.CloseWriter:output_type -> WriterCloseResponse
	9, // 9: Ledger.CloseReader:output_type -> ReaderCloseResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_proto_ledger_proto_init() }
func file_pkg_proto_ledger_proto_init() {
	if File_pkg_proto_ledger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_ledger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadRequest); i {
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
		file_pkg_proto_ledger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_pkg_proto_ledger_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitRequest); i {
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
		file_pkg_proto_ledger_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitResponse); i {
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
		file_pkg_proto_ledger_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteRequest); i {
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
		file_pkg_proto_ledger_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteResponse); i {
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
		file_pkg_proto_ledger_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriterCloseRequest); i {
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
		file_pkg_proto_ledger_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriterCloseResponse); i {
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
		file_pkg_proto_ledger_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReaderCloseRequest); i {
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
		file_pkg_proto_ledger_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReaderCloseResponse); i {
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
			RawDescriptor: file_pkg_proto_ledger_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_ledger_proto_goTypes,
		DependencyIndexes: file_pkg_proto_ledger_proto_depIdxs,
		MessageInfos:      file_pkg_proto_ledger_proto_msgTypes,
	}.Build()
	File_pkg_proto_ledger_proto = out.File
	file_pkg_proto_ledger_proto_rawDesc = nil
	file_pkg_proto_ledger_proto_goTypes = nil
	file_pkg_proto_ledger_proto_depIdxs = nil
}