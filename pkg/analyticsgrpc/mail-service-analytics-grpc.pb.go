// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.13.0
// source: proto/mail-service-analytics-grpc.proto

package analyticsgrpc

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId string               `protobuf:"bytes,1,opt,name=EventId,proto3" json:"EventId,omitempty"`
	TaskId  string               `protobuf:"bytes,2,opt,name=TaskId,proto3" json:"TaskId,omitempty"`
	Time    *timestamp.Timestamp `protobuf:"bytes,3,opt,name=Time,proto3" json:"Time,omitempty"`
	Type    string               `protobuf:"bytes,4,opt,name=Type,proto3" json:"Type,omitempty"`
	Status  string               `protobuf:"bytes,5,opt,name=Status,proto3" json:"Status,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mail_service_analytics_grpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mail_service_analytics_grpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_proto_mail_service_analytics_grpc_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

func (x *Event) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *Event) GetTime() *timestamp.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Event) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Event) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_proto_mail_service_analytics_grpc_proto protoreflect.FileDescriptor

var file_proto_mail_service_analytics_grpc_proto_rawDesc = []byte{
	0x0a, 0x27, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x2d, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2d, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x2d, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x61, 0x6e, 0x61, 0x6c, 0x79,
	0x74, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x61,
	0x73, 0x6b, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x54, 0x61, 0x73, 0x6b,
	0x49, 0x64, 0x12, 0x2e, 0x0a, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x49,
	0x0a, 0x09, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x12, 0x3c, 0x0a, 0x0a, 0x53,
	0x74, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x2e, 0x61, 0x6e, 0x61, 0x6c,
	0x79, 0x74, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x1f, 0x5a, 0x1d, 0x2e, 0x2f, 0x61,
	0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x3b, 0x61, 0x6e, 0x61,
	0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_mail_service_analytics_grpc_proto_rawDescOnce sync.Once
	file_proto_mail_service_analytics_grpc_proto_rawDescData = file_proto_mail_service_analytics_grpc_proto_rawDesc
)

func file_proto_mail_service_analytics_grpc_proto_rawDescGZIP() []byte {
	file_proto_mail_service_analytics_grpc_proto_rawDescOnce.Do(func() {
		file_proto_mail_service_analytics_grpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_mail_service_analytics_grpc_proto_rawDescData)
	})
	return file_proto_mail_service_analytics_grpc_proto_rawDescData
}

var file_proto_mail_service_analytics_grpc_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_mail_service_analytics_grpc_proto_goTypes = []interface{}{
	(*Event)(nil),               // 0: analyticsgrpc.Event
	(*timestamp.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*empty.Empty)(nil),         // 2: google.protobuf.Empty
}
var file_proto_mail_service_analytics_grpc_proto_depIdxs = []int32{
	1, // 0: analyticsgrpc.Event.Time:type_name -> google.protobuf.Timestamp
	0, // 1: analyticsgrpc.Analytics.StoreEvent:input_type -> analyticsgrpc.Event
	2, // 2: analyticsgrpc.Analytics.StoreEvent:output_type -> google.protobuf.Empty
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_mail_service_analytics_grpc_proto_init() }
func file_proto_mail_service_analytics_grpc_proto_init() {
	if File_proto_mail_service_analytics_grpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_mail_service_analytics_grpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
			RawDescriptor: file_proto_mail_service_analytics_grpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_mail_service_analytics_grpc_proto_goTypes,
		DependencyIndexes: file_proto_mail_service_analytics_grpc_proto_depIdxs,
		MessageInfos:      file_proto_mail_service_analytics_grpc_proto_msgTypes,
	}.Build()
	File_proto_mail_service_analytics_grpc_proto = out.File
	file_proto_mail_service_analytics_grpc_proto_rawDesc = nil
	file_proto_mail_service_analytics_grpc_proto_goTypes = nil
	file_proto_mail_service_analytics_grpc_proto_depIdxs = nil
}
