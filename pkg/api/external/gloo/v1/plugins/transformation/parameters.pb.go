// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: plugins/transformation/parameters.proto

package transformation // import "github.com/solo-io/supergloo/pkg/api/external/gloo/v1/plugins/transformation"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import types "github.com/gogo/protobuf/types"

import bytes "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Parameters struct {
	// headers that will be used to extract data for processing output templates
	// Gloo will search for parameters by their name in header value strings, enclosed in single
	// curly braces
	// Example:
	//   extensions:
	//     parameters:
	//         headers:
	//           x-user-id: { userId }
	Headers map[string]string `protobuf:"bytes,1,rep,name=headers" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// part of the (or the entire) path that will be used extract data for processing output templates
	// Gloo will search for parameters by their name in header value strings, enclosed in single
	// curly braces
	// Example:
	//   extensions:
	//     parameters:
	//         path: /users/{ userId }
	Path                 *types.StringValue `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Parameters) Reset()         { *m = Parameters{} }
func (m *Parameters) String() string { return proto.CompactTextString(m) }
func (*Parameters) ProtoMessage()    {}
func (*Parameters) Descriptor() ([]byte, []int) {
	return fileDescriptor_parameters_168485b2fd641fe3, []int{0}
}
func (m *Parameters) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Parameters.Unmarshal(m, b)
}
func (m *Parameters) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Parameters.Marshal(b, m, deterministic)
}
func (dst *Parameters) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Parameters.Merge(dst, src)
}
func (m *Parameters) XXX_Size() int {
	return xxx_messageInfo_Parameters.Size(m)
}
func (m *Parameters) XXX_DiscardUnknown() {
	xxx_messageInfo_Parameters.DiscardUnknown(m)
}

var xxx_messageInfo_Parameters proto.InternalMessageInfo

func (m *Parameters) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *Parameters) GetPath() *types.StringValue {
	if m != nil {
		return m.Path
	}
	return nil
}

func init() {
	proto.RegisterType((*Parameters)(nil), "transformation.plugins.gloo.solo.io.Parameters")
	proto.RegisterMapType((map[string]string)(nil), "transformation.plugins.gloo.solo.io.Parameters.HeadersEntry")
}
func (this *Parameters) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Parameters)
	if !ok {
		that2, ok := that.(Parameters)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Headers) != len(that1.Headers) {
		return false
	}
	for i := range this.Headers {
		if this.Headers[i] != that1.Headers[i] {
			return false
		}
	}
	if !this.Path.Equal(that1.Path) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}

func init() {
	proto.RegisterFile("plugins/transformation/parameters.proto", fileDescriptor_parameters_168485b2fd641fe3)
}

var fileDescriptor_parameters_168485b2fd641fe3 = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xc1, 0x4a, 0x33, 0x31,
	0x10, 0xc7, 0x49, 0xfb, 0x7d, 0x4a, 0x53, 0x0f, 0xb2, 0xf4, 0x50, 0x8a, 0x94, 0xa2, 0x07, 0x7b,
	0x71, 0x46, 0xeb, 0x45, 0x8a, 0x27, 0x41, 0xf0, 0xe0, 0x41, 0x56, 0xe8, 0xc1, 0x5b, 0xaa, 0x69,
	0x36, 0x34, 0x9b, 0x09, 0x49, 0xb6, 0xda, 0x37, 0xf2, 0x79, 0x7c, 0x04, 0x9f, 0x44, 0x36, 0x6b,
	0xab, 0x82, 0x07, 0x6f, 0xff, 0x0c, 0xf3, 0xff, 0xcd, 0x0f, 0xc2, 0x8f, 0x9d, 0xa9, 0x94, 0xb6,
	0x01, 0xa3, 0x17, 0x36, 0x2c, 0xc8, 0x97, 0x22, 0x6a, 0xb2, 0xe8, 0x84, 0x17, 0xa5, 0x8c, 0xd2,
	0x07, 0x70, 0x9e, 0x22, 0x65, 0x47, 0x3f, 0x17, 0xe0, 0xb3, 0x07, 0xca, 0x10, 0x41, 0x20, 0x43,
	0xa0, 0x69, 0x30, 0x54, 0x44, 0xca, 0x48, 0x4c, 0x95, 0x79, 0xb5, 0xc0, 0x67, 0x2f, 0x9c, 0xdb,
	0x42, 0x06, 0x3d, 0x45, 0x8a, 0x52, 0xc4, 0x3a, 0x35, 0xd3, 0xc3, 0x37, 0xc6, 0xf9, 0xdd, 0xf6,
	0x5e, 0x36, 0xe3, 0xbb, 0x85, 0x14, 0x4f, 0xd2, 0x87, 0x3e, 0x1b, 0xb5, 0xc7, 0xdd, 0xc9, 0x25,
	0xfc, 0xe1, 0x36, 0x7c, 0x11, 0xe0, 0xa6, 0xa9, 0x5f, 0xdb, 0xe8, 0xd7, 0xf9, 0x06, 0x96, 0x9d,
	0xf2, 0x7f, 0x4e, 0xc4, 0xa2, 0xdf, 0x1a, 0xb1, 0x71, 0x77, 0x72, 0x00, 0x8d, 0x2b, 0x6c, 0x5c,
	0xe1, 0x3e, 0x7a, 0x6d, 0xd5, 0x4c, 0x98, 0x4a, 0xe6, 0x69, 0x73, 0x30, 0xe5, 0x7b, 0xdf, 0x51,
	0xd9, 0x3e, 0x6f, 0x2f, 0xe5, 0xba, 0xcf, 0x46, 0x6c, 0xdc, 0xc9, 0xeb, 0x98, 0xf5, 0xf8, 0xff,
	0x55, 0x5d, 0x48, 0xd0, 0x4e, 0xde, 0x3c, 0xa6, 0xad, 0x0b, 0x76, 0x95, 0xbf, 0xbe, 0x0f, 0xd9,
	0xc3, 0xad, 0xd2, 0xb1, 0xa8, 0xe6, 0xf0, 0x48, 0x25, 0xd6, 0xa2, 0x27, 0x9a, 0x30, 0x54, 0x4e,
	0xfa, 0x5a, 0x1d, 0xdd, 0x52, 0xa1, 0x70, 0x1a, 0xe5, 0x4b, 0x94, 0xde, 0x0a, 0x83, 0x69, 0xba,
	0x3a, 0xc3, 0xdf, 0x3f, 0x65, 0xbe, 0x93, 0x5c, 0xcf, 0x3f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc5,
	0xa5, 0x9e, 0x2d, 0xb5, 0x01, 0x00, 0x00,
}
