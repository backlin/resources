// Code generated by protoc-gen-gogo.
// source: event.proto
// DO NOT EDIT!

/*
	Package orchestrator is a generated protocol buffer package.

	It is generated from these files:
		event.proto

	It has these top-level messages:
		Event
		Work
*/
package orchestrator

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Event_Level int32

const (
	Event_BATCH Event_Level = 0
	Event_JOB   Event_Level = 1
)

var Event_Level_name = map[int32]string{
	0: "BATCH",
	1: "JOB",
}
var Event_Level_value = map[string]int32{
	"BATCH": 0,
	"JOB":   1,
}

func (x Event_Level) String() string {
	return proto.EnumName(Event_Level_name, int32(x))
}
func (Event_Level) EnumDescriptor() ([]byte, []int) { return fileDescriptorEvent, []int{0, 0} }

type Event_Status int32

const (
	Event_MISSING Event_Status = 0
	Event_PENDING Event_Status = 1
	Event_RUNNING Event_Status = 2
	Event_SUCCESS Event_Status = 3
	Event_FAILURE Event_Status = 4
)

var Event_Status_name = map[int32]string{
	0: "MISSING",
	1: "PENDING",
	2: "RUNNING",
	3: "SUCCESS",
	4: "FAILURE",
}
var Event_Status_value = map[string]int32{
	"MISSING": 0,
	"PENDING": 1,
	"RUNNING": 2,
	"SUCCESS": 3,
	"FAILURE": 4,
}

func (x Event_Status) String() string {
	return proto.EnumName(Event_Status_name, int32(x))
}
func (Event_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptorEvent, []int{0, 1} }

type Event struct {
	BatchId     []byte       `protobuf:"bytes,1,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	JobId       int32        `protobuf:"varint,2,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	StatusLevel Event_Level  `protobuf:"varint,3,opt,name=status_level,json=statusLevel,proto3,enum=orchestrator.Event_Level" json:"status_level,omitempty"`
	Status      Event_Status `protobuf:"varint,4,opt,name=status,proto3,enum=orchestrator.Event_Status" json:"status,omitempty"`
	JobCount    int32        `protobuf:"varint,5,opt,name=job_count,json=jobCount,proto3" json:"job_count,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptorEvent, []int{0} }

func (m *Event) GetBatchId() []byte {
	if m != nil {
		return m.BatchId
	}
	return nil
}

func (m *Event) GetJobId() int32 {
	if m != nil {
		return m.JobId
	}
	return 0
}

func (m *Event) GetStatusLevel() Event_Level {
	if m != nil {
		return m.StatusLevel
	}
	return Event_BATCH
}

func (m *Event) GetStatus() Event_Status {
	if m != nil {
		return m.Status
	}
	return Event_MISSING
}

func (m *Event) GetJobCount() int32 {
	if m != nil {
		return m.JobCount
	}
	return 0
}

type Work struct {
	BatchId  []byte `protobuf:"bytes,1,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	JobId    int32  `protobuf:"varint,2,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	Duration int64  `protobuf:"varint,3,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (m *Work) Reset()                    { *m = Work{} }
func (m *Work) String() string            { return proto.CompactTextString(m) }
func (*Work) ProtoMessage()               {}
func (*Work) Descriptor() ([]byte, []int) { return fileDescriptorEvent, []int{1} }

func (m *Work) GetBatchId() []byte {
	if m != nil {
		return m.BatchId
	}
	return nil
}

func (m *Work) GetJobId() int32 {
	if m != nil {
		return m.JobId
	}
	return 0
}

func (m *Work) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func init() {
	proto.RegisterType((*Event)(nil), "orchestrator.Event")
	proto.RegisterType((*Work)(nil), "orchestrator.Work")
	proto.RegisterEnum("orchestrator.Event_Level", Event_Level_name, Event_Level_value)
	proto.RegisterEnum("orchestrator.Event_Status", Event_Status_name, Event_Status_value)
}
func (m *Event) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Event) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.BatchId) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintEvent(dAtA, i, uint64(len(m.BatchId)))
		i += copy(dAtA[i:], m.BatchId)
	}
	if m.JobId != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintEvent(dAtA, i, uint64(m.JobId))
	}
	if m.StatusLevel != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintEvent(dAtA, i, uint64(m.StatusLevel))
	}
	if m.Status != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintEvent(dAtA, i, uint64(m.Status))
	}
	if m.JobCount != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintEvent(dAtA, i, uint64(m.JobCount))
	}
	return i, nil
}

func (m *Work) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Work) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.BatchId) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintEvent(dAtA, i, uint64(len(m.BatchId)))
		i += copy(dAtA[i:], m.BatchId)
	}
	if m.JobId != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintEvent(dAtA, i, uint64(m.JobId))
	}
	if m.Duration != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintEvent(dAtA, i, uint64(m.Duration))
	}
	return i, nil
}

func encodeFixed64Event(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Event(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Event) Size() (n int) {
	var l int
	_ = l
	l = len(m.BatchId)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.JobId != 0 {
		n += 1 + sovEvent(uint64(m.JobId))
	}
	if m.StatusLevel != 0 {
		n += 1 + sovEvent(uint64(m.StatusLevel))
	}
	if m.Status != 0 {
		n += 1 + sovEvent(uint64(m.Status))
	}
	if m.JobCount != 0 {
		n += 1 + sovEvent(uint64(m.JobCount))
	}
	return n
}

func (m *Work) Size() (n int) {
	var l int
	_ = l
	l = len(m.BatchId)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.JobId != 0 {
		n += 1 + sovEvent(uint64(m.JobId))
	}
	if m.Duration != 0 {
		n += 1 + sovEvent(uint64(m.Duration))
	}
	return n
}

func sovEvent(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Event) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Event: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Event: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BatchId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BatchId = append(m.BatchId[:0], dAtA[iNdEx:postIndex]...)
			if m.BatchId == nil {
				m.BatchId = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field JobId", wireType)
			}
			m.JobId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.JobId |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatusLevel", wireType)
			}
			m.StatusLevel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StatusLevel |= (Event_Level(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= (Event_Status(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field JobCount", wireType)
			}
			m.JobCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.JobCount |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEvent
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
func (m *Work) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Work: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Work: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BatchId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BatchId = append(m.BatchId[:0], dAtA[iNdEx:postIndex]...)
			if m.BatchId == nil {
				m.BatchId = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field JobId", wireType)
			}
			m.JobId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.JobId |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			m.Duration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Duration |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEvent
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
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthEvent
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowEvent
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipEvent(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthEvent = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("event.proto", fileDescriptorEvent) }

var fileDescriptorEvent = []byte{
	// 315 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0xd1, 0xc1, 0x4e, 0xf2, 0x40,
	0x10, 0x07, 0x70, 0x96, 0xb2, 0x05, 0x06, 0xf2, 0x65, 0xb3, 0xc9, 0x97, 0x14, 0x48, 0x1a, 0xd2,
	0x13, 0xa7, 0x1e, 0xf0, 0xea, 0x05, 0x6a, 0xd5, 0x35, 0x58, 0xcd, 0x2e, 0xc4, 0x23, 0x69, 0x69,
	0x13, 0x40, 0xc2, 0x9a, 0x65, 0xe1, 0x39, 0x7c, 0x14, 0x1f, 0xc3, 0xa3, 0x8f, 0x60, 0xea, 0x8b,
	0x98, 0xdd, 0x1a, 0xe3, 0xc1, 0x93, 0xc7, 0x5f, 0xe6, 0x3f, 0xb3, 0x93, 0x59, 0xe8, 0x14, 0xa7,
	0x62, 0xaf, 0xc3, 0x27, 0x25, 0xb5, 0xa4, 0x5d, 0xa9, 0x56, 0xeb, 0xe2, 0xa0, 0x55, 0xaa, 0xa5,
	0x0a, 0x5e, 0xea, 0x80, 0x63, 0x53, 0xa5, 0x3d, 0x68, 0x65, 0xa9, 0x5e, 0xad, 0x97, 0x9b, 0xdc,
	0x43, 0x43, 0x34, 0xea, 0xf2, 0xa6, 0x35, 0xcb, 0xe9, 0x7f, 0x70, 0xb7, 0x32, 0x33, 0x85, 0xfa,
	0x10, 0x8d, 0x30, 0xc7, 0x5b, 0x99, 0xb1, 0x9c, 0x9e, 0x43, 0xf7, 0xa0, 0x53, 0x7d, 0x3c, 0x2c,
	0x77, 0xc5, 0xa9, 0xd8, 0x79, 0xce, 0x10, 0x8d, 0xfe, 0x8d, 0x7b, 0xe1, 0xcf, 0x07, 0x42, 0x3b,
	0x3c, 0x9c, 0x99, 0x00, 0xef, 0x54, 0x71, 0x0b, 0x3a, 0x06, 0xb7, 0xa2, 0xd7, 0xb0, 0x7d, 0xfd,
	0xdf, 0xfa, 0x84, 0x4d, 0xf0, 0xaf, 0x24, 0x1d, 0x40, 0xdb, 0x2c, 0xb2, 0x92, 0xc7, 0xbd, 0xf6,
	0xb0, 0xdd, 0xa5, 0xb5, 0x95, 0x59, 0x64, 0x1c, 0x0c, 0x00, 0x57, 0x93, 0xdb, 0x80, 0xa7, 0x93,
	0x79, 0x74, 0x4d, 0x6a, 0xb4, 0x09, 0xce, 0xcd, 0xdd, 0x94, 0xa0, 0x80, 0x81, 0x5b, 0xcd, 0xa2,
	0x1d, 0x68, 0xde, 0x32, 0x21, 0x58, 0x72, 0x45, 0x6a, 0x06, 0xf7, 0x71, 0x72, 0x61, 0x80, 0x0c,
	0xf8, 0x22, 0x49, 0x0c, 0xea, 0x06, 0x62, 0x11, 0x45, 0xb1, 0x10, 0xc4, 0x31, 0xb8, 0x9c, 0xb0,
	0xd9, 0x82, 0xc7, 0xa4, 0x11, 0xcc, 0xa1, 0xf1, 0x20, 0xd5, 0xe3, 0x1f, 0x0e, 0xd6, 0x87, 0x56,
	0x7e, 0x54, 0xa9, 0xde, 0xc8, 0xbd, 0x3d, 0x96, 0xc3, 0xbf, 0x3d, 0x25, 0xaf, 0xa5, 0x8f, 0xde,
	0x4a, 0x1f, 0xbd, 0x97, 0x3e, 0x7a, 0xfe, 0xf0, 0x6b, 0x99, 0x6b, 0xff, 0xeb, 0xec, 0x33, 0x00,
	0x00, 0xff, 0xff, 0x40, 0x8c, 0xdb, 0x02, 0xbe, 0x01, 0x00, 0x00,
}