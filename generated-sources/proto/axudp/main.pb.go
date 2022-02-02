// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.2
// source: main.proto

package axudp

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

type PacketMode int32

const (
	PacketMode_PM_OPTIONAL_CONSISTENTLY  PacketMode = 0  // Optional last id
	PacketMode_PM_OPTIONAL               PacketMode = 1  // Optional
	PacketMode_PM_MANDATORY_CONSISTENTLY PacketMode = 3  // Mandatory ordered
	PacketMode_PM_MANDATORY              PacketMode = 4  // Mandatory unordered
	PacketMode_PM_SERVICE                PacketMode = 10 // Service layer
)

// Enum value maps for PacketMode.
var (
	PacketMode_name = map[int32]string{
		0:  "PM_OPTIONAL_CONSISTENTLY",
		1:  "PM_OPTIONAL",
		3:  "PM_MANDATORY_CONSISTENTLY",
		4:  "PM_MANDATORY",
		10: "PM_SERVICE",
	}
	PacketMode_value = map[string]int32{
		"PM_OPTIONAL_CONSISTENTLY":  0,
		"PM_OPTIONAL":               1,
		"PM_MANDATORY_CONSISTENTLY": 3,
		"PM_MANDATORY":              4,
		"PM_SERVICE":                10,
	}
)

func (x PacketMode) Enum() *PacketMode {
	p := new(PacketMode)
	*p = x
	return p
}

func (x PacketMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PacketMode) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[0].Descriptor()
}

func (PacketMode) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[0]
}

func (x PacketMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketMode.Descriptor instead.
func (PacketMode) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{0}
}

type PacketTag int32

const (
	PacketTag_PT_PAYLOAD      PacketTag = 0 // Plain payload
	PacketTag_PT_GZIP_PAYLOAD PacketTag = 1 // GZip payload
	PacketTag_PT_PING         PacketTag = 7 // Ping
	PacketTag_PT_DONE         PacketTag = 8 // Packet response
)

// Enum value maps for PacketTag.
var (
	PacketTag_name = map[int32]string{
		0: "PT_PAYLOAD",
		1: "PT_GZIP_PAYLOAD",
		7: "PT_PING",
		8: "PT_DONE",
	}
	PacketTag_value = map[string]int32{
		"PT_PAYLOAD":      0,
		"PT_GZIP_PAYLOAD": 1,
		"PT_PING":         7,
		"PT_DONE":         8,
	}
)

func (x PacketTag) Enum() *PacketTag {
	p := new(PacketTag)
	*p = x
	return p
}

func (x PacketTag) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PacketTag) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[1].Descriptor()
}

func (PacketTag) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[1]
}

func (x PacketTag) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketTag.Descriptor instead.
func (PacketTag) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{1}
}

type Packet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload    []byte     `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`                           // Payload
	Id         uint64     `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`                                    // ID
	Tag        PacketTag  `protobuf:"varint,3,opt,name=tag,proto3,enum=com.axgrid.udp.PacketTag" json:"tag,omitempty"`    // TAG
	Parts      []byte     `protobuf:"bytes,4,opt,name=parts,proto3" json:"parts,omitempty"`                               // PARTS in Bytes
	PartsCount uint32     `protobuf:"varint,5,opt,name=parts_count,json=partsCount,proto3" json:"parts_count,omitempty"`  // PARTS_COUNT
	Mode       PacketMode `protobuf:"varint,6,opt,name=mode,proto3,enum=com.axgrid.udp.PacketMode" json:"mode,omitempty"` // MODE
	Mtu        uint32     `protobuf:"varint,7,opt,name=mtu,proto3" json:"mtu,omitempty"`                                  // Current MTU
}

func (x *Packet) Reset() {
	*x = Packet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Packet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Packet) ProtoMessage() {}

func (x *Packet) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Packet.ProtoReflect.Descriptor instead.
func (*Packet) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{0}
}

func (x *Packet) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Packet) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Packet) GetTag() PacketTag {
	if x != nil {
		return x.Tag
	}
	return PacketTag_PT_PAYLOAD
}

func (x *Packet) GetParts() []byte {
	if x != nil {
		return x.Parts
	}
	return nil
}

func (x *Packet) GetPartsCount() uint32 {
	if x != nil {
		return x.PartsCount
	}
	return 0
}

func (x *Packet) GetMode() PacketMode {
	if x != nil {
		return x.Mode
	}
	return PacketMode_PM_OPTIONAL_CONSISTENTLY
}

func (x *Packet) GetMtu() uint32 {
	if x != nil {
		return x.Mtu
	}
	return 0
}

var File_main_proto protoreflect.FileDescriptor

var file_main_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x63, 0x6f,
	0x6d, 0x2e, 0x61, 0x78, 0x67, 0x72, 0x69, 0x64, 0x2e, 0x75, 0x64, 0x70, 0x22, 0xd8, 0x01, 0x0a,
	0x06, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x2b, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x78, 0x67, 0x72, 0x69, 0x64, 0x2e, 0x75, 0x64, 0x70, 0x2e,
	0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54, 0x61, 0x67, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x14,
	0x0a, 0x05, 0x70, 0x61, 0x72, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x70,
	0x61, 0x72, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x72, 0x74, 0x73, 0x5f, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x73,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x78, 0x67, 0x72, 0x69, 0x64,
	0x2e, 0x75, 0x64, 0x70, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x52,
	0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x74, 0x75, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x03, 0x6d, 0x74, 0x75, 0x2a, 0x7c, 0x0a, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x18, 0x50, 0x4d, 0x5f, 0x4f, 0x50, 0x54, 0x49,
	0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x43, 0x4f, 0x4e, 0x53, 0x49, 0x53, 0x54, 0x45, 0x4e, 0x54, 0x4c,
	0x59, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x4d, 0x5f, 0x4f, 0x50, 0x54, 0x49, 0x4f, 0x4e,
	0x41, 0x4c, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x50, 0x4d, 0x5f, 0x4d, 0x41, 0x4e, 0x44, 0x41,
	0x54, 0x4f, 0x52, 0x59, 0x5f, 0x43, 0x4f, 0x4e, 0x53, 0x49, 0x53, 0x54, 0x45, 0x4e, 0x54, 0x4c,
	0x59, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x50, 0x4d, 0x5f, 0x4d, 0x41, 0x4e, 0x44, 0x41, 0x54,
	0x4f, 0x52, 0x59, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x4d, 0x5f, 0x53, 0x45, 0x52, 0x56,
	0x49, 0x43, 0x45, 0x10, 0x0a, 0x2a, 0x4a, 0x0a, 0x09, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54,
	0x61, 0x67, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x54, 0x5f, 0x50, 0x41, 0x59, 0x4c, 0x4f, 0x41, 0x44,
	0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x54, 0x5f, 0x47, 0x5a, 0x49, 0x50, 0x5f, 0x50, 0x41,
	0x59, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x54, 0x5f, 0x50, 0x49,
	0x4e, 0x47, 0x10, 0x07, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x54, 0x5f, 0x44, 0x4f, 0x4e, 0x45, 0x10,
	0x08, 0x42, 0x2a, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x78, 0x67, 0x72, 0x69, 0x64, 0x2e,
	0x75, 0x64, 0x70, 0x50, 0x01, 0xaa, 0x02, 0x15, 0x41, 0x78, 0x47, 0x72, 0x69, 0x64, 0x2e, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_main_proto_rawDescOnce sync.Once
	file_main_proto_rawDescData = file_main_proto_rawDesc
)

func file_main_proto_rawDescGZIP() []byte {
	file_main_proto_rawDescOnce.Do(func() {
		file_main_proto_rawDescData = protoimpl.X.CompressGZIP(file_main_proto_rawDescData)
	})
	return file_main_proto_rawDescData
}

var file_main_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_main_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_main_proto_goTypes = []interface{}{
	(PacketMode)(0), // 0: com.axgrid.udp.PacketMode
	(PacketTag)(0),  // 1: com.axgrid.udp.PacketTag
	(*Packet)(nil),  // 2: com.axgrid.udp.Packet
}
var file_main_proto_depIdxs = []int32{
	1, // 0: com.axgrid.udp.Packet.tag:type_name -> com.axgrid.udp.PacketTag
	0, // 1: com.axgrid.udp.Packet.mode:type_name -> com.axgrid.udp.PacketMode
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_main_proto_init() }
func file_main_proto_init() {
	if File_main_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_main_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Packet); i {
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
			RawDescriptor: file_main_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_main_proto_goTypes,
		DependencyIndexes: file_main_proto_depIdxs,
		EnumInfos:         file_main_proto_enumTypes,
		MessageInfos:      file_main_proto_msgTypes,
	}.Build()
	File_main_proto = out.File
	file_main_proto_rawDesc = nil
	file_main_proto_goTypes = nil
	file_main_proto_depIdxs = nil
}
