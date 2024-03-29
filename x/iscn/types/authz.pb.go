// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: likechain/iscn/authz.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type UpdateAuthorization struct {
	IscnIdPrefix string `protobuf:"bytes,1,opt,name=iscn_id_prefix,json=iscnIdPrefix,proto3" json:"iscn_id_prefix,omitempty"`
}

func (m *UpdateAuthorization) Reset()         { *m = UpdateAuthorization{} }
func (m *UpdateAuthorization) String() string { return proto.CompactTextString(m) }
func (*UpdateAuthorization) ProtoMessage()    {}
func (*UpdateAuthorization) Descriptor() ([]byte, []int) {
	return fileDescriptor_69559c01192448b8, []int{0}
}
func (m *UpdateAuthorization) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpdateAuthorization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpdateAuthorization.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UpdateAuthorization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateAuthorization.Merge(m, src)
}
func (m *UpdateAuthorization) XXX_Size() int {
	return m.Size()
}
func (m *UpdateAuthorization) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateAuthorization.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateAuthorization proto.InternalMessageInfo

func (m *UpdateAuthorization) GetIscnIdPrefix() string {
	if m != nil {
		return m.IscnIdPrefix
	}
	return ""
}

func init() {
	proto.RegisterType((*UpdateAuthorization)(nil), "likechain.iscn.UpdateAuthorization")
}

func init() { proto.RegisterFile("likechain/iscn/authz.proto", fileDescriptor_69559c01192448b8) }

var fileDescriptor_69559c01192448b8 = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xca, 0xc9, 0xcc, 0x4e,
	0x4d, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0xcf, 0x2c, 0x4e, 0xce, 0xd3, 0x4f, 0x2c, 0x2d, 0xc9, 0xa8,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x83, 0xcb, 0xe9, 0x81, 0xe4, 0xa4, 0x24, 0x93,
	0xf3, 0x8b, 0x73, 0xf3, 0x8b, 0xe3, 0xc1, 0xb2, 0xfa, 0x10, 0x0e, 0x44, 0xa9, 0x94, 0x48, 0x7a,
	0x7e, 0x7a, 0x3e, 0x44, 0x1c, 0xc4, 0x82, 0x88, 0x2a, 0xf9, 0x71, 0x09, 0x87, 0x16, 0xa4, 0x24,
	0x96, 0xa4, 0x3a, 0x96, 0x96, 0x64, 0xe4, 0x17, 0x65, 0x56, 0x25, 0x96, 0x64, 0xe6, 0xe7, 0x09,
	0xa9, 0x70, 0xf1, 0x81, 0xcc, 0x8b, 0xcf, 0x4c, 0x89, 0x2f, 0x28, 0x4a, 0x4d, 0xcb, 0xac, 0x90,
	0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0xe2, 0x01, 0x89, 0x7a, 0xa6, 0x04, 0x80, 0xc5, 0xac, 0x04,
	0x2f, 0x6d, 0xd1, 0xe5, 0x45, 0xd1, 0xe8, 0xe4, 0x73, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72,
	0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7,
	0x72, 0x0c, 0x51, 0x46, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x60,
	0x57, 0xe7, 0x67, 0xe6, 0xc1, 0x19, 0xba, 0x10, 0xff, 0x95, 0x99, 0xe8, 0x57, 0x40, 0x3c, 0x59,
	0x52, 0x59, 0x90, 0x5a, 0x9c, 0xc4, 0x06, 0x76, 0xa4, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x1f,
	0x41, 0x59, 0x6b, 0x03, 0x01, 0x00, 0x00,
}

func (m *UpdateAuthorization) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateAuthorization) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UpdateAuthorization) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.IscnIdPrefix) > 0 {
		i -= len(m.IscnIdPrefix)
		copy(dAtA[i:], m.IscnIdPrefix)
		i = encodeVarintAuthz(dAtA, i, uint64(len(m.IscnIdPrefix)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAuthz(dAtA []byte, offset int, v uint64) int {
	offset -= sovAuthz(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *UpdateAuthorization) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.IscnIdPrefix)
	if l > 0 {
		n += 1 + l + sovAuthz(uint64(l))
	}
	return n
}

func sovAuthz(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAuthz(x uint64) (n int) {
	return sovAuthz(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UpdateAuthorization) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthz
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UpdateAuthorization: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateAuthorization: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IscnIdPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuthz
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthz
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IscnIdPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAuthz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipAuthz(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAuthz
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAuthz
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAuthz
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthAuthz
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAuthz
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAuthz
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAuthz        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAuthz          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAuthz = fmt.Errorf("proto: unexpected end of group")
)
