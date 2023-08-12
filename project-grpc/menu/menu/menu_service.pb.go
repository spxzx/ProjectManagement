// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: menu_service.proto

package menu

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Menu struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Pid        int64   `protobuf:"varint,2,opt,name=pid,proto3" json:"pid,omitempty"`
	Title      string  `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Icon       string  `protobuf:"bytes,4,opt,name=icon,proto3" json:"icon,omitempty"`
	Url        string  `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
	FilePath   string  `protobuf:"bytes,6,opt,name=filePath,proto3" json:"filePath,omitempty"`
	Params     string  `protobuf:"bytes,7,opt,name=params,proto3" json:"params,omitempty"`
	Node       string  `protobuf:"bytes,8,opt,name=node,proto3" json:"node,omitempty"`
	Sort       int32   `protobuf:"varint,9,opt,name=sort,proto3" json:"sort,omitempty"`
	Status     int32   `protobuf:"varint,10,opt,name=status,proto3" json:"status,omitempty"`
	CreateBy   int64   `protobuf:"varint,11,opt,name=createBy,proto3" json:"createBy,omitempty"`
	IsInner    int32   `protobuf:"varint,12,opt,name=isInner,proto3" json:"isInner,omitempty"`
	Values     string  `protobuf:"bytes,13,opt,name=values,proto3" json:"values,omitempty"`
	ShowSlider int32   `protobuf:"varint,14,opt,name=showSlider,proto3" json:"showSlider,omitempty"`
	StatusText string  `protobuf:"bytes,15,opt,name=statusText,proto3" json:"statusText,omitempty"`
	InnerText  string  `protobuf:"bytes,16,opt,name=innerText,proto3" json:"innerText,omitempty"`
	FullUrl    string  `protobuf:"bytes,17,opt,name=fullUrl,proto3" json:"fullUrl,omitempty"`
	Children   []*Menu `protobuf:"bytes,18,rep,name=children,proto3" json:"children,omitempty"`
}

func (x *Menu) Reset() {
	*x = Menu{}
	if protoimpl.UnsafeEnabled {
		mi := &file_menu_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Menu) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Menu) ProtoMessage() {}

func (x *Menu) ProtoReflect() protoreflect.Message {
	mi := &file_menu_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Menu.ProtoReflect.Descriptor instead.
func (*Menu) Descriptor() ([]byte, []int) {
	return file_menu_service_proto_rawDescGZIP(), []int{0}
}

func (x *Menu) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Menu) GetPid() int64 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *Menu) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Menu) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

func (x *Menu) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Menu) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *Menu) GetParams() string {
	if x != nil {
		return x.Params
	}
	return ""
}

func (x *Menu) GetNode() string {
	if x != nil {
		return x.Node
	}
	return ""
}

func (x *Menu) GetSort() int32 {
	if x != nil {
		return x.Sort
	}
	return 0
}

func (x *Menu) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Menu) GetCreateBy() int64 {
	if x != nil {
		return x.CreateBy
	}
	return 0
}

func (x *Menu) GetIsInner() int32 {
	if x != nil {
		return x.IsInner
	}
	return 0
}

func (x *Menu) GetValues() string {
	if x != nil {
		return x.Values
	}
	return ""
}

func (x *Menu) GetShowSlider() int32 {
	if x != nil {
		return x.ShowSlider
	}
	return 0
}

func (x *Menu) GetStatusText() string {
	if x != nil {
		return x.StatusText
	}
	return ""
}

func (x *Menu) GetInnerText() string {
	if x != nil {
		return x.InnerText
	}
	return ""
}

func (x *Menu) GetFullUrl() string {
	if x != nil {
		return x.FullUrl
	}
	return ""
}

func (x *Menu) GetChildren() []*Menu {
	if x != nil {
		return x.Children
	}
	return nil
}

type MenuResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*Menu `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *MenuResponse) Reset() {
	*x = MenuResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_menu_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MenuResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MenuResponse) ProtoMessage() {}

func (x *MenuResponse) ProtoReflect() protoreflect.Message {
	mi := &file_menu_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MenuResponse.ProtoReflect.Descriptor instead.
func (*MenuResponse) Descriptor() ([]byte, []int) {
	return file_menu_service_proto_rawDescGZIP(), []int{1}
}

func (x *MenuResponse) GetList() []*Menu {
	if x != nil {
		return x.List
	}
	return nil
}

var File_menu_service_proto protoreflect.FileDescriptor

var file_menu_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6d, 0x65, 0x6e, 0x75, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6d, 0x65, 0x6e, 0x75, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xce, 0x03, 0x0a, 0x04, 0x4d, 0x65, 0x6e, 0x75, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x69, 0x63, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61,
	0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61,
	0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f,
	0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x6f,
	0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x42, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x42, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x49, 0x6e, 0x6e, 0x65,
	0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x69, 0x73, 0x49, 0x6e, 0x6e, 0x65, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x68, 0x6f, 0x77,
	0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x68,
	0x6f, 0x77, 0x53, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x54, 0x65, 0x78, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x54, 0x65, 0x78, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x6e, 0x65,
	0x72, 0x54, 0x65, 0x78, 0x74, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x6e,
	0x65, 0x72, 0x54, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x66, 0x75, 0x6c, 0x6c, 0x55, 0x72,
	0x6c, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x66, 0x75, 0x6c, 0x6c, 0x55, 0x72, 0x6c,
	0x12, 0x2e, 0x0a, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x18, 0x12, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x6e, 0x75, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x4d, 0x65, 0x6e, 0x75, 0x52, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e,
	0x22, 0x36, 0x0a, 0x0c, 0x4d, 0x65, 0x6e, 0x75, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x26, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x6d, 0x65, 0x6e, 0x75, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65,
	0x6e, 0x75, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x32, 0x52, 0x0a, 0x0b, 0x4d, 0x65, 0x6e, 0x75,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4d, 0x65,
	0x6e, 0x75, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1a,
	0x2e, 0x6d, 0x65, 0x6e, 0x75, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65,
	0x6e, 0x75, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x22, 0x5a, 0x20,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x6d, 0x65, 0x6e, 0x75,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_menu_service_proto_rawDescOnce sync.Once
	file_menu_service_proto_rawDescData = file_menu_service_proto_rawDesc
)

func file_menu_service_proto_rawDescGZIP() []byte {
	file_menu_service_proto_rawDescOnce.Do(func() {
		file_menu_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_menu_service_proto_rawDescData)
	})
	return file_menu_service_proto_rawDescData
}

var file_menu_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_menu_service_proto_goTypes = []interface{}{
	(*Menu)(nil),          // 0: menu.service.Menu
	(*MenuResponse)(nil),  // 1: menu.service.MenuResponse
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_menu_service_proto_depIdxs = []int32{
	0, // 0: menu.service.Menu.children:type_name -> menu.service.Menu
	0, // 1: menu.service.MenuResponse.list:type_name -> menu.service.Menu
	2, // 2: menu.service.MenuService.GetMenuList:input_type -> google.protobuf.Empty
	1, // 3: menu.service.MenuService.GetMenuList:output_type -> menu.service.MenuResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_menu_service_proto_init() }
func file_menu_service_proto_init() {
	if File_menu_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_menu_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Menu); i {
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
		file_menu_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MenuResponse); i {
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
			RawDescriptor: file_menu_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_menu_service_proto_goTypes,
		DependencyIndexes: file_menu_service_proto_depIdxs,
		MessageInfos:      file_menu_service_proto_msgTypes,
	}.Build()
	File_menu_service_proto = out.File
	file_menu_service_proto_rawDesc = nil
	file_menu_service_proto_goTypes = nil
	file_menu_service_proto_depIdxs = nil
}
