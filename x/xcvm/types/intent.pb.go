// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: composable/xcvm/v1beta1/intent.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type TransferIntent struct {
	SourceAddress      string                                  `protobuf:"bytes,1,opt,name=source_address,json=sourceAddress,proto3" json:"source_address,omitempty"`
	DestinationAddress string                                  `protobuf:"bytes,2,opt,name=destination_address,json=destinationAddress,proto3" json:"destination_address,omitempty"`
	ClientId           string                                  `protobuf:"bytes,3,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Timeout            time.Time                               `protobuf:"bytes,4,opt,name=timeout,proto3,stdtime" json:"timeout"`
	Amount             cosmossdk_io_math.Uint                  `protobuf:"bytes,5,opt,name=amount,proto3,customtype=cosmossdk.io/math.Uint" json:"amount"`
	Bounty             github_com_cosmos_cosmos_sdk_types.Coin `protobuf:"bytes,6,opt,name=bounty,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin" json:"bounty"`
}

func (m *TransferIntent) Reset()         { *m = TransferIntent{} }
func (m *TransferIntent) String() string { return proto.CompactTextString(m) }
func (*TransferIntent) ProtoMessage()    {}
func (*TransferIntent) Descriptor() ([]byte, []int) {
	return fileDescriptor_5047e26d000a53a3, []int{0}
}
func (m *TransferIntent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferIntent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferIntent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferIntent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferIntent.Merge(m, src)
}
func (m *TransferIntent) XXX_Size() int {
	return m.Size()
}
func (m *TransferIntent) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferIntent.DiscardUnknown(m)
}

var xxx_messageInfo_TransferIntent proto.InternalMessageInfo

func (m *TransferIntent) GetSourceAddress() string {
	if m != nil {
		return m.SourceAddress
	}
	return ""
}

func (m *TransferIntent) GetDestinationAddress() string {
	if m != nil {
		return m.DestinationAddress
	}
	return ""
}

func (m *TransferIntent) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *TransferIntent) GetTimeout() time.Time {
	if m != nil {
		return m.Timeout
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*TransferIntent)(nil), "composable.xcvm.v1beta1.TransferIntent")
}

func init() {
	proto.RegisterFile("composable/xcvm/v1beta1/intent.proto", fileDescriptor_5047e26d000a53a3)
}

var fileDescriptor_5047e26d000a53a3 = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0x4f, 0x8f, 0xd2, 0x40,
	0x18, 0xc6, 0x5b, 0xd4, 0x0a, 0xa3, 0x72, 0xa8, 0x46, 0x2b, 0x26, 0x2d, 0x31, 0xfe, 0xe1, 0xe2,
	0x4c, 0xd0, 0xc4, 0xa3, 0x89, 0x78, 0x30, 0x5c, 0x1b, 0xbc, 0x78, 0x21, 0xd3, 0x76, 0x28, 0x13,
	0x98, 0x79, 0x9b, 0xce, 0x94, 0xc0, 0xb7, 0xe0, 0x53, 0x19, 0x8e, 0x1c, 0xcd, 0x1e, 0xd8, 0x0d,
	0x7c, 0x91, 0x4d, 0x67, 0xda, 0xdd, 0x3d, 0x75, 0xde, 0xa7, 0xbf, 0x5f, 0xf3, 0xf4, 0x1d, 0xf4,
	0x21, 0x05, 0x51, 0x80, 0xa2, 0xc9, 0x9a, 0x91, 0x6d, 0xba, 0x11, 0x64, 0x33, 0x4e, 0x98, 0xa6,
	0x63, 0xc2, 0xa5, 0x66, 0x52, 0xe3, 0xa2, 0x04, 0x0d, 0xfe, 0x9b, 0x7b, 0x0a, 0xd7, 0x14, 0x6e,
	0xa8, 0xc1, 0xdb, 0x14, 0x94, 0x00, 0x35, 0x37, 0x18, 0xb1, 0x83, 0x75, 0x06, 0xaf, 0x72, 0xc8,
	0xc1, 0xe6, 0xf5, 0xa9, 0x49, 0xa3, 0x1c, 0x20, 0x5f, 0x33, 0x62, 0xa6, 0xa4, 0x5a, 0x10, 0xcd,
	0x05, 0x53, 0x9a, 0x8a, 0xc2, 0x02, 0xef, 0xff, 0x75, 0x50, 0x7f, 0x56, 0x52, 0xa9, 0x16, 0xac,
	0x9c, 0x9a, 0x0e, 0xfe, 0x47, 0xd4, 0x57, 0x50, 0x95, 0x29, 0x9b, 0xd3, 0x2c, 0x2b, 0x99, 0x52,
	0x81, 0x3b, 0x74, 0x47, 0xbd, 0xf8, 0x85, 0x4d, 0x7f, 0xda, 0xd0, 0x27, 0xe8, 0x65, 0xc6, 0x94,
	0xe6, 0x92, 0x6a, 0x0e, 0xf2, 0x8e, 0xed, 0x18, 0xd6, 0x7f, 0xf0, 0xaa, 0x15, 0xde, 0xa1, 0x5e,
	0xba, 0xe6, 0x4c, 0xea, 0x39, 0xcf, 0x82, 0x47, 0x06, 0xeb, 0xda, 0x60, 0x9a, 0xf9, 0x3f, 0xd0,
	0xd3, 0xba, 0x1a, 0x54, 0x3a, 0x78, 0x3c, 0x74, 0x47, 0xcf, 0xbe, 0x0e, 0xb0, 0xad, 0x8e, 0xdb,
	0xea, 0x78, 0xd6, 0x56, 0x9f, 0x74, 0x0f, 0xa7, 0xc8, 0xd9, 0x5f, 0x47, 0x6e, 0xdc, 0x4a, 0xfe,
	0x77, 0xe4, 0x51, 0x01, 0x95, 0xd4, 0xc1, 0x93, 0xfa, 0xcb, 0x93, 0xb0, 0x46, 0xae, 0x4e, 0xd1,
	0x6b, 0xbb, 0x24, 0x95, 0xad, 0x30, 0x07, 0x22, 0xa8, 0x5e, 0xe2, 0x3f, 0x5c, 0xea, 0xb8, 0xa1,
	0xfd, 0xdf, 0xc8, 0x4b, 0xea, 0xc3, 0x2e, 0xf0, 0x8c, 0x47, 0x1a, 0xef, 0x73, 0xce, 0xf5, 0xb2,
	0x4a, 0x70, 0x0a, 0xa2, 0xd9, 0x73, 0xf3, 0xf8, 0xa2, 0xb2, 0x15, 0xd1, 0xbb, 0x82, 0x29, 0xfc,
	0x0b, 0xb8, 0x8c, 0x1b, 0x7d, 0xf2, 0xe9, 0x70, 0x0e, 0xdd, 0xe3, 0x39, 0x74, 0x6f, 0xce, 0xa1,
	0xbb, 0xbf, 0x84, 0xce, 0xf1, 0x12, 0x3a, 0xff, 0x2f, 0xa1, 0xf3, 0xf7, 0xf9, 0xd6, 0x5e, 0xb5,
	0x91, 0x12, 0xcf, 0xfc, 0xcf, 0xb7, 0xdb, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa8, 0xae, 0xdf, 0xb2,
	0x0a, 0x02, 0x00, 0x00,
}

func (m *TransferIntent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferIntent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferIntent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Bounty.Size()
		i -= size
		if _, err := m.Bounty.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintIntent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintIntent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Timeout, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timeout):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintIntent(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	if len(m.ClientId) > 0 {
		i -= len(m.ClientId)
		copy(dAtA[i:], m.ClientId)
		i = encodeVarintIntent(dAtA, i, uint64(len(m.ClientId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.DestinationAddress) > 0 {
		i -= len(m.DestinationAddress)
		copy(dAtA[i:], m.DestinationAddress)
		i = encodeVarintIntent(dAtA, i, uint64(len(m.DestinationAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.SourceAddress) > 0 {
		i -= len(m.SourceAddress)
		copy(dAtA[i:], m.SourceAddress)
		i = encodeVarintIntent(dAtA, i, uint64(len(m.SourceAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintIntent(dAtA []byte, offset int, v uint64) int {
	offset -= sovIntent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TransferIntent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SourceAddress)
	if l > 0 {
		n += 1 + l + sovIntent(uint64(l))
	}
	l = len(m.DestinationAddress)
	if l > 0 {
		n += 1 + l + sovIntent(uint64(l))
	}
	l = len(m.ClientId)
	if l > 0 {
		n += 1 + l + sovIntent(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timeout)
	n += 1 + l + sovIntent(uint64(l))
	l = m.Amount.Size()
	n += 1 + l + sovIntent(uint64(l))
	l = m.Bounty.Size()
	n += 1 + l + sovIntent(uint64(l))
	return n
}

func sovIntent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozIntent(x uint64) (n int) {
	return sovIntent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TransferIntent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIntent
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
			return fmt.Errorf("proto: TransferIntent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferIntent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntent
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
				return ErrInvalidLengthIntent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIntent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntent
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
				return ErrInvalidLengthIntent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIntent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntent
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
				return ErrInvalidLengthIntent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIntent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timeout", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIntent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIntent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Timeout, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntent
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
				return ErrInvalidLengthIntent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIntent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bounty", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntent
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
				return ErrInvalidLengthIntent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIntent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Bounty.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIntent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIntent
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
func skipIntent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIntent
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
					return 0, ErrIntOverflowIntent
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
					return 0, ErrIntOverflowIntent
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
				return 0, ErrInvalidLengthIntent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupIntent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthIntent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthIntent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIntent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupIntent = fmt.Errorf("proto: unexpected end of group")
)
