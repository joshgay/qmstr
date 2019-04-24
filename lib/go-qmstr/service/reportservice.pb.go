// Code generated by protoc-gen-go. DO NOT EDIT.
// source: reportservice.proto

package service

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ReporterConfigRequest struct {
	ReporterID           int32    `protobuf:"varint,1,opt,name=reporterID,proto3" json:"reporterID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReporterConfigRequest) Reset()         { *m = ReporterConfigRequest{} }
func (m *ReporterConfigRequest) String() string { return proto.CompactTextString(m) }
func (*ReporterConfigRequest) ProtoMessage()    {}
func (*ReporterConfigRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_414fd7029dd18bbe, []int{0}
}

func (m *ReporterConfigRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReporterConfigRequest.Unmarshal(m, b)
}
func (m *ReporterConfigRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReporterConfigRequest.Marshal(b, m, deterministic)
}
func (m *ReporterConfigRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReporterConfigRequest.Merge(m, src)
}
func (m *ReporterConfigRequest) XXX_Size() int {
	return xxx_messageInfo_ReporterConfigRequest.Size(m)
}
func (m *ReporterConfigRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReporterConfigRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReporterConfigRequest proto.InternalMessageInfo

func (m *ReporterConfigRequest) GetReporterID() int32 {
	if m != nil {
		return m.ReporterID
	}
	return 0
}

type ReporterConfigResponse struct {
	ConfigMap            map[string]string `protobuf:"bytes,1,rep,name=configMap,proto3" json:"configMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Name                 string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ReporterConfigResponse) Reset()         { *m = ReporterConfigResponse{} }
func (m *ReporterConfigResponse) String() string { return proto.CompactTextString(m) }
func (*ReporterConfigResponse) ProtoMessage()    {}
func (*ReporterConfigResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_414fd7029dd18bbe, []int{1}
}

func (m *ReporterConfigResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReporterConfigResponse.Unmarshal(m, b)
}
func (m *ReporterConfigResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReporterConfigResponse.Marshal(b, m, deterministic)
}
func (m *ReporterConfigResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReporterConfigResponse.Merge(m, src)
}
func (m *ReporterConfigResponse) XXX_Size() int {
	return xxx_messageInfo_ReporterConfigResponse.Size(m)
}
func (m *ReporterConfigResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ReporterConfigResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReporterConfigResponse proto.InternalMessageInfo

func (m *ReporterConfigResponse) GetConfigMap() map[string]string {
	if m != nil {
		return m.ConfigMap
	}
	return nil
}

func (m *ReporterConfigResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type InfoDataRequest struct {
	RootID               string   `protobuf:"bytes,1,opt,name=rootID,proto3" json:"rootID,omitempty"`
	Infotype             string   `protobuf:"bytes,2,opt,name=infotype,proto3" json:"infotype,omitempty"`
	Datatype             string   `protobuf:"bytes,3,opt,name=datatype,proto3" json:"datatype,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoDataRequest) Reset()         { *m = InfoDataRequest{} }
func (m *InfoDataRequest) String() string { return proto.CompactTextString(m) }
func (*InfoDataRequest) ProtoMessage()    {}
func (*InfoDataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_414fd7029dd18bbe, []int{2}
}

func (m *InfoDataRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoDataRequest.Unmarshal(m, b)
}
func (m *InfoDataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoDataRequest.Marshal(b, m, deterministic)
}
func (m *InfoDataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoDataRequest.Merge(m, src)
}
func (m *InfoDataRequest) XXX_Size() int {
	return xxx_messageInfo_InfoDataRequest.Size(m)
}
func (m *InfoDataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoDataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InfoDataRequest proto.InternalMessageInfo

func (m *InfoDataRequest) GetRootID() string {
	if m != nil {
		return m.RootID
	}
	return ""
}

func (m *InfoDataRequest) GetInfotype() string {
	if m != nil {
		return m.Infotype
	}
	return ""
}

func (m *InfoDataRequest) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

type InfoDataResponse struct {
	Data                 []string `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoDataResponse) Reset()         { *m = InfoDataResponse{} }
func (m *InfoDataResponse) String() string { return proto.CompactTextString(m) }
func (*InfoDataResponse) ProtoMessage()    {}
func (*InfoDataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_414fd7029dd18bbe, []int{3}
}

func (m *InfoDataResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoDataResponse.Unmarshal(m, b)
}
func (m *InfoDataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoDataResponse.Marshal(b, m, deterministic)
}
func (m *InfoDataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoDataResponse.Merge(m, src)
}
func (m *InfoDataResponse) XXX_Size() int {
	return xxx_messageInfo_InfoDataResponse.Size(m)
}
func (m *InfoDataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoDataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InfoDataResponse proto.InternalMessageInfo

func (m *InfoDataResponse) GetData() []string {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*ReporterConfigRequest)(nil), "service.ReporterConfigRequest")
	proto.RegisterType((*ReporterConfigResponse)(nil), "service.ReporterConfigResponse")
	proto.RegisterMapType((map[string]string)(nil), "service.ReporterConfigResponse.ConfigMapEntry")
	proto.RegisterType((*InfoDataRequest)(nil), "service.InfoDataRequest")
	proto.RegisterType((*InfoDataResponse)(nil), "service.InfoDataResponse")
}

func init() { proto.RegisterFile("reportservice.proto", fileDescriptor_414fd7029dd18bbe) }

var fileDescriptor_414fd7029dd18bbe = []byte{
	// 368 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x4f, 0x4f, 0xf2, 0x40,
	0x10, 0xc6, 0x59, 0xfe, 0xbd, 0x6f, 0x87, 0xbc, 0xc0, 0xbb, 0x22, 0xa9, 0x3d, 0x20, 0xe9, 0xc1,
	0x70, 0xea, 0x01, 0x0f, 0x1a, 0x63, 0x3c, 0x08, 0x86, 0x90, 0xa8, 0x31, 0x6b, 0x62, 0xbc, 0xae,
	0x65, 0x20, 0x28, 0x74, 0xcb, 0x76, 0x21, 0xe1, 0x9b, 0xf9, 0xa5, 0xfc, 0x0e, 0xa6, 0xdb, 0x6e,
	0xf9, 0x13, 0x8c, 0xb7, 0x79, 0xe6, 0xd9, 0x67, 0x3a, 0xf3, 0x4b, 0xe1, 0x48, 0x62, 0x28, 0xa4,
	0x8a, 0x50, 0xae, 0xa6, 0x3e, 0x7a, 0xa1, 0x14, 0x4a, 0xd0, 0x3f, 0xa9, 0x74, 0x6a, 0x23, 0xae,
	0xf8, 0x5c, 0x8c, 0x70, 0x96, 0x38, 0xee, 0x05, 0x1c, 0x33, 0x1d, 0x40, 0xd9, 0x13, 0xc1, 0x78,
	0x3a, 0x61, 0xb8, 0x58, 0x62, 0xa4, 0x68, 0x0b, 0x40, 0xa6, 0xc6, 0xb0, 0x6f, 0x93, 0x36, 0xe9,
	0x94, 0xd8, 0x56, 0xc7, 0xfd, 0x24, 0xd0, 0xdc, 0x4f, 0x46, 0xa1, 0x08, 0x22, 0xa4, 0xf7, 0x60,
	0xf9, 0xba, 0xf3, 0xc0, 0x43, 0x9b, 0xb4, 0x0b, 0x9d, 0x4a, 0xd7, 0xf3, 0xcc, 0x42, 0x87, 0x33,
	0x5e, 0xcf, 0x04, 0xee, 0x02, 0x25, 0xd7, 0x6c, 0x33, 0x80, 0x52, 0x28, 0x06, 0x7c, 0x8e, 0x76,
	0xbe, 0x4d, 0x3a, 0x16, 0xd3, 0xb5, 0x73, 0x0d, 0xd5, 0xdd, 0x00, 0xad, 0x43, 0xe1, 0x03, 0xd7,
	0x7a, 0x4f, 0x8b, 0xc5, 0x25, 0x6d, 0x40, 0x69, 0xc5, 0x67, 0x4b, 0x13, 0x4c, 0xc4, 0x55, 0xfe,
	0x92, 0xb8, 0x1c, 0x6a, 0xc3, 0x60, 0x2c, 0xfa, 0x5c, 0x71, 0x73, 0x6d, 0x13, 0xca, 0x52, 0x08,
	0x95, 0x5e, 0x6a, 0xb1, 0x54, 0x51, 0x07, 0xfe, 0x4e, 0x83, 0xb1, 0x50, 0xeb, 0xd0, 0xcc, 0xc9,
	0x74, 0xec, 0xc5, 0x34, 0xb5, 0x57, 0x48, 0x3c, 0xa3, 0xdd, 0x33, 0xa8, 0x6f, 0x3e, 0x91, 0x62,
	0xa1, 0x50, 0x8c, 0x7d, 0x4d, 0xc4, 0x62, 0xba, 0xee, 0x7e, 0x11, 0xf8, 0x97, 0x10, 0x79, 0x4e,
	0xf8, 0xd0, 0x17, 0xf8, 0x3f, 0x40, 0xb5, 0x4b, 0x89, 0xb6, 0x7e, 0xc4, 0xa7, 0xd7, 0x77, 0x4e,
	0x7f, 0xc1, 0xeb, 0xe6, 0x68, 0x1f, 0x2a, 0x03, 0x54, 0x66, 0x29, 0x6a, 0x67, 0x89, 0x3d, 0x14,
	0xce, 0xc9, 0x01, 0x27, 0x9b, 0x72, 0x03, 0xd5, 0x01, 0xaa, 0x27, 0x29, 0xde, 0xd1, 0x57, 0x8f,
	0x62, 0x84, 0xb4, 0x91, 0x3d, 0xdf, 0xea, 0x3a, 0x07, 0xbb, 0x6e, 0xee, 0xd6, 0x86, 0xa6, 0x90,
	0x13, 0x6f, 0x31, 0x8f, 0x94, 0xf4, 0x26, 0x32, 0xf4, 0xcd, 0xbb, 0xd7, 0xdc, 0x5b, 0x59, 0xff,
	0x91, 0xe7, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb1, 0x4b, 0x53, 0x52, 0xc2, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ReportServiceClient is the client API for ReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReportServiceClient interface {
	GetReporterConfig(ctx context.Context, in *ReporterConfigRequest, opts ...grpc.CallOption) (*ReporterConfigResponse, error)
	GetInfoData(ctx context.Context, in *InfoDataRequest, opts ...grpc.CallOption) (*InfoDataResponse, error)
	GetProjectNode(ctx context.Context, in *ProjectNode, opts ...grpc.CallOption) (*ProjectNode, error)
}

type reportServiceClient struct {
	cc *grpc.ClientConn
}

func NewReportServiceClient(cc *grpc.ClientConn) ReportServiceClient {
	return &reportServiceClient{cc}
}

func (c *reportServiceClient) GetReporterConfig(ctx context.Context, in *ReporterConfigRequest, opts ...grpc.CallOption) (*ReporterConfigResponse, error) {
	out := new(ReporterConfigResponse)
	err := c.cc.Invoke(ctx, "/service.ReportService/GetReporterConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportServiceClient) GetInfoData(ctx context.Context, in *InfoDataRequest, opts ...grpc.CallOption) (*InfoDataResponse, error) {
	out := new(InfoDataResponse)
	err := c.cc.Invoke(ctx, "/service.ReportService/GetInfoData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportServiceClient) GetProjectNode(ctx context.Context, in *ProjectNode, opts ...grpc.CallOption) (*ProjectNode, error) {
	out := new(ProjectNode)
	err := c.cc.Invoke(ctx, "/service.ReportService/GetProjectNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReportServiceServer is the server API for ReportService service.
type ReportServiceServer interface {
	GetReporterConfig(context.Context, *ReporterConfigRequest) (*ReporterConfigResponse, error)
	GetInfoData(context.Context, *InfoDataRequest) (*InfoDataResponse, error)
	GetProjectNode(context.Context, *ProjectNode) (*ProjectNode, error)
}

func RegisterReportServiceServer(s *grpc.Server, srv ReportServiceServer) {
	s.RegisterService(&_ReportService_serviceDesc, srv)
}

func _ReportService_GetReporterConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReporterConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).GetReporterConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ReportService/GetReporterConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).GetReporterConfig(ctx, req.(*ReporterConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReportService_GetInfoData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).GetInfoData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ReportService/GetInfoData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).GetInfoData(ctx, req.(*InfoDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReportService_GetProjectNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).GetProjectNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ReportService/GetProjectNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).GetProjectNode(ctx, req.(*ProjectNode))
	}
	return interceptor(ctx, in, info, handler)
}

var _ReportService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.ReportService",
	HandlerType: (*ReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetReporterConfig",
			Handler:    _ReportService_GetReporterConfig_Handler,
		},
		{
			MethodName: "GetInfoData",
			Handler:    _ReportService_GetInfoData_Handler,
		},
		{
			MethodName: "GetProjectNode",
			Handler:    _ReportService_GetProjectNode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "reportservice.proto",
}