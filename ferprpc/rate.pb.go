// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ferprpc/rate.proto

package ferprpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type SubscribeRatesRequest struct {
	Currency             string   `protobuf:"bytes,1,opt,name=currency,proto3" json:"currency,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscribeRatesRequest) Reset()         { *m = SubscribeRatesRequest{} }
func (m *SubscribeRatesRequest) String() string { return proto.CompactTextString(m) }
func (*SubscribeRatesRequest) ProtoMessage()    {}
func (*SubscribeRatesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e260cc4981fd8bef, []int{0}
}

func (m *SubscribeRatesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeRatesRequest.Unmarshal(m, b)
}
func (m *SubscribeRatesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeRatesRequest.Marshal(b, m, deterministic)
}
func (m *SubscribeRatesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeRatesRequest.Merge(m, src)
}
func (m *SubscribeRatesRequest) XXX_Size() int {
	return xxx_messageInfo_SubscribeRatesRequest.Size(m)
}
func (m *SubscribeRatesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeRatesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeRatesRequest proto.InternalMessageInfo

func (m *SubscribeRatesRequest) GetCurrency() string {
	if m != nil {
		return m.Currency
	}
	return ""
}

type SubscribeRatesResponse struct {
	Currency             string          `protobuf:"bytes,1,opt,name=currency,proto3" json:"currency,omitempty"`
	Rate                 int64           `protobuf:"varint,2,opt,name=rate,proto3" json:"rate,omitempty"`
	RateMsat             int64           `protobuf:"varint,3,opt,name=rate_msat,json=rateMsat,proto3" json:"rate_msat,omitempty"`
	ConversionRate       *ConversionRate `protobuf:"bytes,4,opt,name=conversion_rate,json=conversionRate,proto3" json:"conversion_rate,omitempty"`
	LastUpdated          int64           `protobuf:"varint,5,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SubscribeRatesResponse) Reset()         { *m = SubscribeRatesResponse{} }
func (m *SubscribeRatesResponse) String() string { return proto.CompactTextString(m) }
func (*SubscribeRatesResponse) ProtoMessage()    {}
func (*SubscribeRatesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e260cc4981fd8bef, []int{1}
}

func (m *SubscribeRatesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeRatesResponse.Unmarshal(m, b)
}
func (m *SubscribeRatesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeRatesResponse.Marshal(b, m, deterministic)
}
func (m *SubscribeRatesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeRatesResponse.Merge(m, src)
}
func (m *SubscribeRatesResponse) XXX_Size() int {
	return xxx_messageInfo_SubscribeRatesResponse.Size(m)
}
func (m *SubscribeRatesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeRatesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeRatesResponse proto.InternalMessageInfo

func (m *SubscribeRatesResponse) GetCurrency() string {
	if m != nil {
		return m.Currency
	}
	return ""
}

func (m *SubscribeRatesResponse) GetRate() int64 {
	if m != nil {
		return m.Rate
	}
	return 0
}

func (m *SubscribeRatesResponse) GetRateMsat() int64 {
	if m != nil {
		return m.RateMsat
	}
	return 0
}

func (m *SubscribeRatesResponse) GetConversionRate() *ConversionRate {
	if m != nil {
		return m.ConversionRate
	}
	return nil
}

func (m *SubscribeRatesResponse) GetLastUpdated() int64 {
	if m != nil {
		return m.LastUpdated
	}
	return 0
}

type ConversionRate struct {
	Currency             string   `protobuf:"bytes,1,opt,name=currency,proto3" json:"currency,omitempty"`
	Rate                 float32  `protobuf:"fixed32,2,opt,name=rate,proto3" json:"rate,omitempty"`
	LastUpdated          int64    `protobuf:"varint,5,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConversionRate) Reset()         { *m = ConversionRate{} }
func (m *ConversionRate) String() string { return proto.CompactTextString(m) }
func (*ConversionRate) ProtoMessage()    {}
func (*ConversionRate) Descriptor() ([]byte, []int) {
	return fileDescriptor_e260cc4981fd8bef, []int{2}
}

func (m *ConversionRate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConversionRate.Unmarshal(m, b)
}
func (m *ConversionRate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConversionRate.Marshal(b, m, deterministic)
}
func (m *ConversionRate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConversionRate.Merge(m, src)
}
func (m *ConversionRate) XXX_Size() int {
	return xxx_messageInfo_ConversionRate.Size(m)
}
func (m *ConversionRate) XXX_DiscardUnknown() {
	xxx_messageInfo_ConversionRate.DiscardUnknown(m)
}

var xxx_messageInfo_ConversionRate proto.InternalMessageInfo

func (m *ConversionRate) GetCurrency() string {
	if m != nil {
		return m.Currency
	}
	return ""
}

func (m *ConversionRate) GetRate() float32 {
	if m != nil {
		return m.Rate
	}
	return 0
}

func (m *ConversionRate) GetLastUpdated() int64 {
	if m != nil {
		return m.LastUpdated
	}
	return 0
}

func init() {
	proto.RegisterType((*SubscribeRatesRequest)(nil), "session.SubscribeRatesRequest")
	proto.RegisterType((*SubscribeRatesResponse)(nil), "session.SubscribeRatesResponse")
	proto.RegisterType((*ConversionRate)(nil), "session.ConversionRate")
}

func init() { proto.RegisterFile("ferprpc/rate.proto", fileDescriptor_e260cc4981fd8bef) }

var fileDescriptor_e260cc4981fd8bef = []byte{
	// 290 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x65, 0xdb, 0xaa, 0xed, 0x46, 0x22, 0x2c, 0xa8, 0xa1, 0x82, 0xc6, 0x88, 0x90, 0x8b, 0x89,
	0xb4, 0x7f, 0x40, 0xf4, 0xec, 0x25, 0xc1, 0x8b, 0x97, 0xb0, 0xd9, 0x8e, 0x35, 0x60, 0xb2, 0x71,
	0x67, 0x52, 0xf0, 0x2f, 0xfa, 0xab, 0x64, 0xd3, 0x10, 0x88, 0xf8, 0xd1, 0x53, 0x32, 0xf3, 0x78,
	0x33, 0xef, 0xcd, 0x3e, 0x2e, 0x5e, 0xc0, 0xd4, 0xa6, 0x56, 0xb1, 0x91, 0x04, 0x51, 0x6d, 0x34,
	0x69, 0x71, 0x80, 0x80, 0x58, 0xe8, 0x2a, 0x58, 0xf2, 0xe3, 0xb4, 0xc9, 0x51, 0x99, 0x22, 0x87,
	0x44, 0x12, 0x60, 0x02, 0xef, 0x0d, 0x20, 0x89, 0x39, 0x9f, 0xaa, 0xc6, 0x18, 0xa8, 0xd4, 0x87,
	0xc7, 0x7c, 0x16, 0xce, 0x92, 0xbe, 0x0e, 0x3e, 0x19, 0x3f, 0xf9, 0xce, 0xc2, 0x5a, 0x57, 0x08,
	0x7f, 0xd1, 0x84, 0xe0, 0x13, 0x2b, 0xc1, 0x1b, 0xf9, 0x2c, 0x1c, 0x27, 0xed, 0xbf, 0x38, 0xe3,
	0x33, 0xfb, 0xcd, 0x4a, 0x94, 0xe4, 0x8d, 0x5b, 0x60, 0x6a, 0x1b, 0x8f, 0x28, 0x49, 0xdc, 0xf1,
	0x23, 0xa5, 0xab, 0x0d, 0x18, 0x2b, 0x35, 0x6b, 0xb9, 0x13, 0x9f, 0x85, 0xce, 0xe2, 0x34, 0xea,
	0xf4, 0x47, 0x0f, 0x3d, 0x6e, 0x75, 0x24, 0xae, 0x1a, 0xd4, 0xe2, 0x92, 0x1f, 0xbe, 0x49, 0xa4,
	0xac, 0xa9, 0x57, 0x92, 0x60, 0xe5, 0xed, 0xb5, 0x1b, 0x1c, 0xdb, 0x7b, 0xda, 0xb6, 0x02, 0xc5,
	0xdd, 0xe1, 0x90, 0x9d, 0x3d, 0x8c, 0x3a, 0x0f, 0xff, 0x2f, 0x59, 0xe4, 0xdc, 0xb1, 0xa3, 0x53,
	0x30, 0x9b, 0x42, 0x81, 0x48, 0xb9, 0x3b, 0xbc, 0x9f, 0x38, 0xef, 0x1d, 0xfd, 0xf8, 0x1c, 0xf3,
	0x8b, 0x5f, 0xf1, 0xed, 0xe1, 0x6f, 0xd9, 0xfd, 0xf5, 0xf3, 0xd5, 0xba, 0xa0, 0xd7, 0x26, 0x8f,
	0x94, 0x2e, 0x63, 0x94, 0x54, 0x94, 0x9a, 0x74, 0xbc, 0xd6, 0x37, 0x36, 0x00, 0x71, 0x97, 0x82,
	0x7c, 0xbf, 0x4d, 0xc0, 0xf2, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xb2, 0x13, 0xd9, 0xcb, 0x17, 0x02,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RateServiceClient is the client API for RateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RateServiceClient interface {
	SubscribeRates(ctx context.Context, in *SubscribeRatesRequest, opts ...grpc.CallOption) (RateService_SubscribeRatesClient, error)
}

type rateServiceClient struct {
	cc *grpc.ClientConn
}

func NewRateServiceClient(cc *grpc.ClientConn) RateServiceClient {
	return &rateServiceClient{cc}
}

func (c *rateServiceClient) SubscribeRates(ctx context.Context, in *SubscribeRatesRequest, opts ...grpc.CallOption) (RateService_SubscribeRatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RateService_serviceDesc.Streams[0], "/session.RateService/SubscribeRates", opts...)
	if err != nil {
		return nil, err
	}
	x := &rateServiceSubscribeRatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RateService_SubscribeRatesClient interface {
	Recv() (*SubscribeRatesResponse, error)
	grpc.ClientStream
}

type rateServiceSubscribeRatesClient struct {
	grpc.ClientStream
}

func (x *rateServiceSubscribeRatesClient) Recv() (*SubscribeRatesResponse, error) {
	m := new(SubscribeRatesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RateServiceServer is the server API for RateService service.
type RateServiceServer interface {
	SubscribeRates(*SubscribeRatesRequest, RateService_SubscribeRatesServer) error
}

// UnimplementedRateServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRateServiceServer struct {
}

func (*UnimplementedRateServiceServer) SubscribeRates(req *SubscribeRatesRequest, srv RateService_SubscribeRatesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeRates not implemented")
}

func RegisterRateServiceServer(s *grpc.Server, srv RateServiceServer) {
	s.RegisterService(&_RateService_serviceDesc, srv)
}

func _RateService_SubscribeRates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRatesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RateServiceServer).SubscribeRates(m, &rateServiceSubscribeRatesServer{stream})
}

type RateService_SubscribeRatesServer interface {
	Send(*SubscribeRatesResponse) error
	grpc.ServerStream
}

type rateServiceSubscribeRatesServer struct {
	grpc.ServerStream
}

func (x *rateServiceSubscribeRatesServer) Send(m *SubscribeRatesResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _RateService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "session.RateService",
	HandlerType: (*RateServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeRates",
			Handler:       _RateService_SubscribeRates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ferprpc/rate.proto",
}