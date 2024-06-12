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
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

type TransferIntent struct {
	SourceAddress      string                                  `protobuf:"bytes,1,opt,name=source_address,json=sourceAddress,proto3" json:"source_address,omitempty"`
	DestinationAddress string                                  `protobuf:"bytes,2,opt,name=destination_address,json=destinationAddress,proto3" json:"destination_address,omitempty"`
	ClientId           string                                  `protobuf:"bytes,3,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	TimeoutHeight      int64                                   `protobuf:"varint,4,opt,name=timeout_height,json=timeoutHeight,proto3" json:"timeout_height,omitempty"`
	TransferTokens     *TransferTokens                         `protobuf:"bytes,5,opt,name=transfer_tokens,json=transferTokens,proto3" json:"transfer_tokens,omitempty"`
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

func (m *TransferIntent) GetTimeoutHeight() int64 {
	if m != nil {
		return m.TimeoutHeight
	}
	return 0
}

func (m *TransferIntent) GetTransferTokens() *TransferTokens {
	if m != nil {
		return m.TransferTokens
	}
	return nil
}

type TransferTokens struct {
	// The ERC20 address of the token to transfer
	Erc20Address string `protobuf:"bytes,1,opt,name=erc20_address,json=erc20Address,proto3" json:"erc20_address,omitempty"`
	// The amount of tokens to transfer
	Amount cosmossdk_io_math.Uint `protobuf:"bytes,2,opt,name=amount,proto3,customtype=cosmossdk.io/math.Uint" json:"amount"`
}

func (m *TransferTokens) Reset()         { *m = TransferTokens{} }
func (m *TransferTokens) String() string { return proto.CompactTextString(m) }
func (*TransferTokens) ProtoMessage()    {}
func (*TransferTokens) Descriptor() ([]byte, []int) {
	return fileDescriptor_5047e26d000a53a3, []int{1}
}
func (m *TransferTokens) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferTokens) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferTokens.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferTokens) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferTokens.Merge(m, src)
}
func (m *TransferTokens) XXX_Size() int {
	return m.Size()
}
func (m *TransferTokens) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferTokens.DiscardUnknown(m)
}

var xxx_messageInfo_TransferTokens proto.InternalMessageInfo

func (m *TransferTokens) GetErc20Address() string {
	if m != nil {
		return m.Erc20Address
	}
	return ""
}

func init() {
	proto.RegisterType((*TransferIntent)(nil), "composable.xcvm.v1beta1.TransferIntent")
	proto.RegisterType((*TransferTokens)(nil), "composable.xcvm.v1beta1.TransferTokens")
}

func init() {
	proto.RegisterFile("composable/xcvm/v1beta1/intent.proto", fileDescriptor_5047e26d000a53a3)
}

var fileDescriptor_5047e26d000a53a3 = []byte{
	// 413 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x6e, 0xd4, 0x30,
	0x10, 0xc6, 0x93, 0x16, 0x22, 0x6a, 0xba, 0x8b, 0x14, 0x10, 0x84, 0x22, 0x65, 0x57, 0xe5, 0x4f,
	0xf7, 0x82, 0x4d, 0x8b, 0xc4, 0x9d, 0xe5, 0x00, 0xbd, 0xa1, 0xa8, 0x5c, 0xb8, 0x44, 0x4e, 0xe2,
	0x26, 0xd6, 0xae, 0x3d, 0x51, 0x3c, 0xa9, 0xda, 0xb7, 0xe0, 0xad, 0xe8, 0xb1, 0x47, 0xc4, 0xa1,
	0x42, 0xbb, 0x2f, 0x82, 0x62, 0x7b, 0x29, 0x20, 0xf5, 0x94, 0xcc, 0x97, 0xdf, 0xcc, 0x7c, 0xf9,
	0x6c, 0xf2, 0xa2, 0x04, 0xd5, 0x82, 0xe1, 0xc5, 0x52, 0xb0, 0xf3, 0xf2, 0x4c, 0xb1, 0xb3, 0xc3,
	0x42, 0x20, 0x3f, 0x64, 0x52, 0xa3, 0xd0, 0x48, 0xdb, 0x0e, 0x10, 0xe2, 0x27, 0x37, 0x14, 0x1d,
	0x28, 0xea, 0xa9, 0xbd, 0xa7, 0x25, 0x18, 0x05, 0x26, 0xb7, 0x18, 0x73, 0x85, 0xeb, 0xd9, 0x7b,
	0x54, 0x43, 0x0d, 0x4e, 0x1f, 0xde, 0xbc, 0x3a, 0xa9, 0x01, 0xea, 0xa5, 0x60, 0xb6, 0x2a, 0xfa,
	0x53, 0x86, 0x52, 0x09, 0x83, 0x5c, 0xb5, 0x0e, 0xd8, 0xff, 0xbe, 0x45, 0xc6, 0x27, 0x1d, 0xd7,
	0xe6, 0x54, 0x74, 0xc7, 0xd6, 0x43, 0xfc, 0x92, 0x8c, 0x0d, 0xf4, 0x5d, 0x29, 0x72, 0x5e, 0x55,
	0x9d, 0x30, 0x26, 0x09, 0xa7, 0xe1, 0x6c, 0x27, 0x1b, 0x39, 0xf5, 0xbd, 0x13, 0x63, 0x46, 0x1e,
	0x56, 0xc2, 0xa0, 0xd4, 0x1c, 0x25, 0xe8, 0x3f, 0xec, 0x96, 0x65, 0xe3, 0xbf, 0x3e, 0x6d, 0x1a,
	0x9e, 0x91, 0x9d, 0x72, 0x29, 0x85, 0xc6, 0x5c, 0x56, 0xc9, 0xb6, 0xc5, 0xee, 0x39, 0xe1, 0xb8,
	0x1a, 0x96, 0x0e, 0xd6, 0xa0, 0xc7, 0xbc, 0x11, 0xb2, 0x6e, 0x30, 0xb9, 0x33, 0x0d, 0x67, 0xdb,
	0xd9, 0xc8, 0xab, 0x9f, 0xac, 0x18, 0x7f, 0x26, 0x0f, 0xd0, 0xbb, 0xcd, 0x11, 0x16, 0x42, 0x9b,
	0xe4, 0xee, 0x34, 0x9c, 0xdd, 0x3f, 0x3a, 0xa0, 0xb7, 0x64, 0x46, 0x37, 0x7f, 0x77, 0x62, 0xf1,
	0x6c, 0x8c, 0xff, 0xd4, 0xf1, 0x47, 0x12, 0x15, 0xd0, 0x6b, 0xbc, 0x48, 0xa2, 0xc1, 0xd2, 0x9c,
	0x5d, 0x5e, 0x4f, 0x82, 0x9f, 0xd7, 0x93, 0x83, 0x5a, 0x62, 0xd3, 0x17, 0xc3, 0x58, 0x1f, 0xb4,
	0x7f, 0xbc, 0x36, 0xd5, 0x82, 0xe1, 0x45, 0x2b, 0x0c, 0xfd, 0x00, 0x52, 0x67, 0xbe, 0x7d, 0x5f,
	0xdd, 0x04, 0xe9, 0x47, 0x3f, 0x27, 0x23, 0xd1, 0x95, 0x47, 0x6f, 0xfe, 0xcb, 0x71, 0xd7, 0x8a,
	0x9b, 0x54, 0xde, 0x91, 0x88, 0xab, 0x61, 0x82, 0x4b, 0x6e, 0x9e, 0xfa, 0xfd, 0x8f, 0xdd, 0x36,
	0x53, 0x2d, 0xa8, 0x04, 0xa6, 0x38, 0x36, 0xf4, 0x8b, 0xd4, 0x98, 0x79, 0x7a, 0xfe, 0xea, 0x72,
	0x95, 0x86, 0x57, 0xab, 0x34, 0xfc, 0xb5, 0x4a, 0xc3, 0x6f, 0xeb, 0x34, 0xb8, 0x5a, 0xa7, 0xc1,
	0x8f, 0x75, 0x1a, 0x7c, 0xdd, 0x3d, 0x77, 0x57, 0xcb, 0x7a, 0x2c, 0x22, 0x7b, 0xce, 0x6f, 0x7f,
	0x07, 0x00, 0x00, 0xff, 0xff, 0x8d, 0xb3, 0x66, 0x4e, 0x7a, 0x02, 0x00, 0x00,
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
	if m.TransferTokens != nil {
		{
			size, err := m.TransferTokens.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIntent(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.TimeoutHeight != 0 {
		i = encodeVarintIntent(dAtA, i, uint64(m.TimeoutHeight))
		i--
		dAtA[i] = 0x20
	}
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

func (m *TransferTokens) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferTokens) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintIntent(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Erc20Address) > 0 {
		i -= len(m.Erc20Address)
		copy(dAtA[i:], m.Erc20Address)
		i = encodeVarintIntent(dAtA, i, uint64(len(m.Erc20Address)))
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
	if m.TimeoutHeight != 0 {
		n += 1 + sovIntent(uint64(m.TimeoutHeight))
	}
	if m.TransferTokens != nil {
		l = m.TransferTokens.Size()
		n += 1 + l + sovIntent(uint64(l))
	}
	l = m.Bounty.Size()
	n += 1 + l + sovIntent(uint64(l))
	return n
}

func (m *TransferTokens) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Erc20Address)
	if l > 0 {
		n += 1 + l + sovIntent(uint64(l))
	}
	l = m.Amount.Size()
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TimeoutHeight", wireType)
			}
			m.TimeoutHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TimeoutHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferTokens", wireType)
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
			if m.TransferTokens == nil {
				m.TransferTokens = &TransferTokens{}
			}
			if err := m.TransferTokens.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *TransferTokens) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TransferTokens: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferTokens: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Erc20Address", wireType)
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
			m.Erc20Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
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
