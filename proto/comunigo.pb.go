// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: proto/comunigo.proto

package proto

import (
	empty "github.com/golang/protobuf/ptypes/empty"
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

type RawMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From string `protobuf:"bytes,1,opt,name=From,proto3" json:"From,omitempty"`
	Body string `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *RawMessage) Reset() {
	*x = RawMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_comunigo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawMessage) ProtoMessage() {}

func (x *RawMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_comunigo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawMessage.ProtoReflect.Descriptor instead.
func (*RawMessage) Descriptor() ([]byte, []int) {
	return file_proto_comunigo_proto_rawDescGZIP(), []int{0}
}

func (x *RawMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *RawMessage) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type SequencerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp uint64 `protobuf:"varint,1,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	From      string `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
	Body      string `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *SequencerMessage) Reset() {
	*x = SequencerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_comunigo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SequencerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SequencerMessage) ProtoMessage() {}

func (x *SequencerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_comunigo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SequencerMessage.ProtoReflect.Descriptor instead.
func (*SequencerMessage) Descriptor() ([]byte, []int) {
	return file_proto_comunigo_proto_rawDescGZIP(), []int{1}
}

func (x *SequencerMessage) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *SequencerMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *SequencerMessage) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type ScalarClockMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp uint64 `protobuf:"varint,1,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	From      string `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
	Body      string `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *ScalarClockMessage) Reset() {
	*x = ScalarClockMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_comunigo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScalarClockMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScalarClockMessage) ProtoMessage() {}

func (x *ScalarClockMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_comunigo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScalarClockMessage.ProtoReflect.Descriptor instead.
func (*ScalarClockMessage) Descriptor() ([]byte, []int) {
	return file_proto_comunigo_proto_rawDescGZIP(), []int{2}
}

func (x *ScalarClockMessage) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *ScalarClockMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *ScalarClockMessage) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type ScalarClockAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp uint64 `protobuf:"varint,1,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	From      string `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
}

func (x *ScalarClockAck) Reset() {
	*x = ScalarClockAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_comunigo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScalarClockAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScalarClockAck) ProtoMessage() {}

func (x *ScalarClockAck) ProtoReflect() protoreflect.Message {
	mi := &file_proto_comunigo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScalarClockAck.ProtoReflect.Descriptor instead.
func (*ScalarClockAck) Descriptor() ([]byte, []int) {
	return file_proto_comunigo_proto_rawDescGZIP(), []int{3}
}

func (x *ScalarClockAck) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *ScalarClockAck) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

type VectorialClockMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp []uint64 `protobuf:"varint,1,rep,packed,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	From      string   `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
	Body      string   `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *VectorialClockMessage) Reset() {
	*x = VectorialClockMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_comunigo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VectorialClockMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorialClockMessage) ProtoMessage() {}

func (x *VectorialClockMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_comunigo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorialClockMessage.ProtoReflect.Descriptor instead.
func (*VectorialClockMessage) Descriptor() ([]byte, []int) {
	return file_proto_comunigo_proto_rawDescGZIP(), []int{4}
}

func (x *VectorialClockMessage) GetTimestamp() []uint64 {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *VectorialClockMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *VectorialClockMessage) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

var File_proto_comunigo_proto protoreflect.FileDescriptor

var file_proto_comunigo_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x75, 0x6e, 0x69, 0x67, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x34, 0x0a, 0x0a, 0x52, 0x61,
	0x77, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04,
	0x42, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79,
	0x22, 0x58, 0x0a, 0x10, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x5a, 0x0a, 0x12, 0x53, 0x63,
	0x61, 0x6c, 0x61, 0x72, 0x43, 0x6c, 0x6f, 0x63, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12,
	0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72,
	0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x42, 0x0a, 0x0e, 0x53, 0x63, 0x61, 0x6c, 0x61, 0x72,
	0x43, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x63, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x22, 0x5d, 0x0a, 0x15, 0x56, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x69, 0x61, 0x6c, 0x43, 0x6c, 0x6f, 0x63, 0x6b, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x32, 0xf9, 0x02, 0x0a, 0x08, 0x43, 0x6f,
	0x6d, 0x75, 0x6e, 0x69, 0x67, 0x6f, 0x12, 0x44, 0x0a, 0x17, 0x53, 0x65, 0x6e, 0x64, 0x46, 0x72,
	0x6f, 0x6d, 0x50, 0x65, 0x65, 0x72, 0x54, 0x6f, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65,
	0x72, 0x12, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x61, 0x77, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4a, 0x0a, 0x17,
	0x53, 0x65, 0x6e, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65,
	0x72, 0x54, 0x6f, 0x50, 0x65, 0x65, 0x72, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x48, 0x0a, 0x13, 0x53, 0x65, 0x6e, 0x64,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x32, 0x50, 0x53, 0x63, 0x61, 0x6c, 0x61, 0x72, 0x12,
	0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x63, 0x61, 0x6c, 0x61, 0x72, 0x43, 0x6c,
	0x6f, 0x63, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x41, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x50, 0x32, 0x50,
	0x53, 0x63, 0x61, 0x6c, 0x61, 0x72, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53,
	0x63, 0x61, 0x6c, 0x61, 0x72, 0x43, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x63, 0x6b, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4e, 0x0a, 0x16, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x50, 0x32, 0x50, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x69, 0x61, 0x6c, 0x12,
	0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x69, 0x61,
	0x6c, 0x43, 0x6c, 0x6f, 0x63, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x62, 0x77, 0x65, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x75,
	0x6e, 0x69, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_comunigo_proto_rawDescOnce sync.Once
	file_proto_comunigo_proto_rawDescData = file_proto_comunigo_proto_rawDesc
)

func file_proto_comunigo_proto_rawDescGZIP() []byte {
	file_proto_comunigo_proto_rawDescOnce.Do(func() {
		file_proto_comunigo_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_comunigo_proto_rawDescData)
	})
	return file_proto_comunigo_proto_rawDescData
}

var file_proto_comunigo_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_comunigo_proto_goTypes = []interface{}{
	(*RawMessage)(nil),            // 0: proto.RawMessage
	(*SequencerMessage)(nil),      // 1: proto.SequencerMessage
	(*ScalarClockMessage)(nil),    // 2: proto.ScalarClockMessage
	(*ScalarClockAck)(nil),        // 3: proto.ScalarClockAck
	(*VectorialClockMessage)(nil), // 4: proto.VectorialClockMessage
	(*empty.Empty)(nil),           // 5: google.protobuf.Empty
}
var file_proto_comunigo_proto_depIdxs = []int32{
	0, // 0: proto.Comunigo.SendFromPeerToSequencer:input_type -> proto.RawMessage
	1, // 1: proto.Comunigo.SendFromSequencerToPeer:input_type -> proto.SequencerMessage
	2, // 2: proto.Comunigo.SendUpdateP2PScalar:input_type -> proto.ScalarClockMessage
	3, // 3: proto.Comunigo.SendAckP2PScalar:input_type -> proto.ScalarClockAck
	4, // 4: proto.Comunigo.SendUpdateP2PVectorial:input_type -> proto.VectorialClockMessage
	5, // 5: proto.Comunigo.SendFromPeerToSequencer:output_type -> google.protobuf.Empty
	5, // 6: proto.Comunigo.SendFromSequencerToPeer:output_type -> google.protobuf.Empty
	5, // 7: proto.Comunigo.SendUpdateP2PScalar:output_type -> google.protobuf.Empty
	5, // 8: proto.Comunigo.SendAckP2PScalar:output_type -> google.protobuf.Empty
	5, // 9: proto.Comunigo.SendUpdateP2PVectorial:output_type -> google.protobuf.Empty
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_comunigo_proto_init() }
func file_proto_comunigo_proto_init() {
	if File_proto_comunigo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_comunigo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RawMessage); i {
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
		file_proto_comunigo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SequencerMessage); i {
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
		file_proto_comunigo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScalarClockMessage); i {
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
		file_proto_comunigo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScalarClockAck); i {
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
		file_proto_comunigo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VectorialClockMessage); i {
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
			RawDescriptor: file_proto_comunigo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_comunigo_proto_goTypes,
		DependencyIndexes: file_proto_comunigo_proto_depIdxs,
		MessageInfos:      file_proto_comunigo_proto_msgTypes,
	}.Build()
	File_proto_comunigo_proto = out.File
	file_proto_comunigo_proto_rawDesc = nil
	file_proto_comunigo_proto_goTypes = nil
	file_proto_comunigo_proto_depIdxs = nil
}
