// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: rpc/rpc.proto

package rpc

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

type PeerPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcId    string   `protobuf:"bytes,1,opt,name=src_id,json=srcId,proto3" json:"src_id,omitempty"`
	SrcIp    string   `protobuf:"bytes,2,opt,name=src_ip,json=srcIp,proto3" json:"src_ip,omitempty"`
	Start    string   `protobuf:"bytes,3,opt,name=start,proto3" json:"start,omitempty"`
	Interval []string `protobuf:"bytes,4,rep,name=interval,proto3" json:"interval,omitempty"`
}

func (x *PeerPacket) Reset() {
	*x = PeerPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerPacket) ProtoMessage() {}

func (x *PeerPacket) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerPacket.ProtoReflect.Descriptor instead.
func (*PeerPacket) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *PeerPacket) GetSrcId() string {
	if x != nil {
		return x.SrcId
	}
	return ""
}

func (x *PeerPacket) GetSrcIp() string {
	if x != nil {
		return x.SrcIp
	}
	return ""
}

func (x *PeerPacket) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

func (x *PeerPacket) GetInterval() []string {
	if x != nil {
		return x.Interval
	}
	return nil
}

type LookupPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcId   string `protobuf:"bytes,1,opt,name=src_id,json=srcId,proto3" json:"src_id,omitempty"`
	SrcIp   string `protobuf:"bytes,2,opt,name=src_ip,json=srcIp,proto3" json:"src_ip,omitempty"`
	Id      string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	IdShift string `protobuf:"bytes,4,opt,name=idShift,proto3" json:"idShift,omitempty"`
	Iid     string `protobuf:"bytes,5,opt,name=iid,proto3" json:"iid,omitempty"`
}

func (x *LookupPacket) Reset() {
	*x = LookupPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupPacket) ProtoMessage() {}

func (x *LookupPacket) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupPacket.ProtoReflect.Descriptor instead.
func (*LookupPacket) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *LookupPacket) GetSrcId() string {
	if x != nil {
		return x.SrcId
	}
	return ""
}

func (x *LookupPacket) GetSrcIp() string {
	if x != nil {
		return x.SrcIp
	}
	return ""
}

func (x *LookupPacket) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LookupPacket) GetIdShift() string {
	if x != nil {
		return x.IdShift
	}
	return ""
}

func (x *LookupPacket) GetIid() string {
	if x != nil {
		return x.Iid
	}
	return ""
}

type BlockPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BlockPacket) Reset() {
	*x = BlockPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockPacket) ProtoMessage() {}

func (x *BlockPacket) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockPacket.ProtoReflect.Descriptor instead.
func (*BlockPacket) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{2}
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{3}
}

type BootStrapPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcId string `protobuf:"bytes,1,opt,name=src_id,json=srcId,proto3" json:"src_id,omitempty"`
	SrcIp string `protobuf:"bytes,2,opt,name=src_ip,json=srcIp,proto3" json:"src_ip,omitempty"`
}

func (x *BootStrapPacket) Reset() {
	*x = BootStrapPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootStrapPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootStrapPacket) ProtoMessage() {}

func (x *BootStrapPacket) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootStrapPacket.ProtoReflect.Descriptor instead.
func (*BootStrapPacket) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{4}
}

func (x *BootStrapPacket) GetSrcId() string {
	if x != nil {
		return x.SrcId
	}
	return ""
}

func (x *BootStrapPacket) GetSrcIp() string {
	if x != nil {
		return x.SrcIp
	}
	return ""
}

type BootStrapReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Successor *PeerPacket `protobuf:"bytes,1,opt,name=Successor,proto3" json:"Successor,omitempty"`
	D         *PeerPacket `protobuf:"bytes,2,opt,name=D,proto3" json:"D,omitempty"`
}

func (x *BootStrapReply) Reset() {
	*x = BootStrapReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootStrapReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootStrapReply) ProtoMessage() {}

func (x *BootStrapReply) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootStrapReply.ProtoReflect.Descriptor instead.
func (*BootStrapReply) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{5}
}

func (x *BootStrapReply) GetSuccessor() *PeerPacket {
	if x != nil {
		return x.Successor
	}
	return nil
}

func (x *BootStrapReply) GetD() *PeerPacket {
	if x != nil {
		return x.D
	}
	return nil
}

var File_rpc_rpc_proto protoreflect.FileDescriptor

var file_rpc_rpc_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x72, 0x70, 0x63, 0x22, 0x6c, 0x0a, 0x0a, 0x50, 0x65, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x72, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x72, 0x63, 0x49, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x72, 0x63,
	0x5f, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x72, 0x63, 0x49, 0x70,
	0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x22, 0x78, 0x0a, 0x0c, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x72, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x72, 0x63, 0x49, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x72, 0x63,
	0x5f, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x72, 0x63, 0x49, 0x70,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x69, 0x64, 0x53, 0x68, 0x69, 0x66, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x69, 0x64, 0x53, 0x68, 0x69, 0x66, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x69,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x69, 0x69, 0x64, 0x22, 0x0d, 0x0a, 0x0b,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x07, 0x0a, 0x05, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x3f, 0x0a, 0x0f, 0x42, 0x6f, 0x6f, 0x74, 0x53, 0x74, 0x72, 0x61,
	0x70, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x72, 0x63, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x72, 0x63, 0x49, 0x64, 0x12, 0x15,
	0x0a, 0x06, 0x73, 0x72, 0x63, 0x5f, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x73, 0x72, 0x63, 0x49, 0x70, 0x22, 0x5e, 0x0a, 0x0e, 0x42, 0x6f, 0x6f, 0x74, 0x53, 0x74, 0x72,
	0x61, 0x70, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x2d, 0x0a, 0x09, 0x53, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x50, 0x65, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x09, 0x53, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x1d, 0x0a, 0x01, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x52, 0x01, 0x44, 0x32, 0x85, 0x02, 0x0a, 0x06, 0x4b, 0x6f, 0x6f, 0x72, 0x64, 0x65,
	0x12, 0x3b, 0x0a, 0x0c, 0x42, 0x6f, 0x6f, 0x74, 0x53, 0x74, 0x61, 0x72, 0x70, 0x52, 0x50, 0x43,
	0x12, 0x14, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x53, 0x74, 0x72, 0x61, 0x70,
	0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x1a, 0x13, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x42, 0x6f, 0x6f,
	0x74, 0x53, 0x74, 0x72, 0x61, 0x70, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x31, 0x0a,
	0x09, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x50, 0x43, 0x12, 0x11, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x1a, 0x0f, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x00,
	0x12, 0x2d, 0x0a, 0x0c, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x52, 0x50, 0x43,
	0x12, 0x0a, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x00, 0x12,
	0x2c, 0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x69, 0x67, 0x62, 0x6f, 0x72,
	0x52, 0x50, 0x43, 0x12, 0x0a, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x0a, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2e, 0x0a,
	0x0c, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x43, 0x61, 0x73, 0x74, 0x52, 0x50, 0x43, 0x12, 0x10, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x1a,
	0x0a, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x2e, 0x5a,
	0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x68, 0x6d,
	0x65, 0x64, 0x6e, 0x61, 0x62, 0x69, 0x6c, 0x31, 0x30, 0x39, 0x2f, 0x6b, 0x6f, 0x6f, 0x72, 0x64,
	0x65, 0x2d, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_rpc_proto_rawDescOnce sync.Once
	file_rpc_rpc_proto_rawDescData = file_rpc_rpc_proto_rawDesc
)

func file_rpc_rpc_proto_rawDescGZIP() []byte {
	file_rpc_rpc_proto_rawDescOnce.Do(func() {
		file_rpc_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_rpc_proto_rawDescData)
	})
	return file_rpc_rpc_proto_rawDescData
}

var file_rpc_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_rpc_rpc_proto_goTypes = []interface{}{
	(*PeerPacket)(nil),      // 0: rpc.PeerPacket
	(*LookupPacket)(nil),    // 1: rpc.LookupPacket
	(*BlockPacket)(nil),     // 2: rpc.BlockPacket
	(*Empty)(nil),           // 3: rpc.Empty
	(*BootStrapPacket)(nil), // 4: rpc.BootStrapPacket
	(*BootStrapReply)(nil),  // 5: rpc.BootStrapReply
}
var file_rpc_rpc_proto_depIdxs = []int32{
	0, // 0: rpc.BootStrapReply.Successor:type_name -> rpc.PeerPacket
	0, // 1: rpc.BootStrapReply.D:type_name -> rpc.PeerPacket
	4, // 2: rpc.Koorde.BootStarpRPC:input_type -> rpc.BootStrapPacket
	1, // 3: rpc.Koorde.LookupRPC:input_type -> rpc.LookupPacket
	3, // 4: rpc.Koorde.SuccessorRPC:input_type -> rpc.Empty
	3, // 5: rpc.Koorde.UpdateNeigborRPC:input_type -> rpc.Empty
	2, // 6: rpc.Koorde.BroadCastRPC:input_type -> rpc.BlockPacket
	5, // 7: rpc.Koorde.BootStarpRPC:output_type -> rpc.BootStrapReply
	0, // 8: rpc.Koorde.LookupRPC:output_type -> rpc.PeerPacket
	0, // 9: rpc.Koorde.SuccessorRPC:output_type -> rpc.PeerPacket
	3, // 10: rpc.Koorde.UpdateNeigborRPC:output_type -> rpc.Empty
	3, // 11: rpc.Koorde.BroadCastRPC:output_type -> rpc.Empty
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rpc_rpc_proto_init() }
func file_rpc_rpc_proto_init() {
	if File_rpc_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerPacket); i {
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
		file_rpc_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupPacket); i {
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
		file_rpc_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockPacket); i {
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
		file_rpc_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_rpc_rpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootStrapPacket); i {
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
		file_rpc_rpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootStrapReply); i {
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
			RawDescriptor: file_rpc_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_rpc_proto_goTypes,
		DependencyIndexes: file_rpc_rpc_proto_depIdxs,
		MessageInfos:      file_rpc_rpc_proto_msgTypes,
	}.Build()
	File_rpc_rpc_proto = out.File
	file_rpc_rpc_proto_rawDesc = nil
	file_rpc_rpc_proto_goTypes = nil
	file_rpc_rpc_proto_depIdxs = nil
}
