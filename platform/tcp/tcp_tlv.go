package tcp

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	TypeCmd     uint8 = 0
	TypePayload uint8 = 10
)

type Tlv struct {
	Type   uint8
	Length uint32
	Value  []byte
}

func (tlv Tlv) Pack() []byte {
	b := []byte{}
	b = append(b, tlv.Type&0xFF)
	b = binary.BigEndian.AppendUint32(b, tlv.Length)
	b = append(b, tlv.Value...)
	return b
}

func NewTlv(t uint8, data []byte) Tlv {
	return Tlv{
		Type:   t,
		Length: uint32(len(data)),
		Value:  data,
	}
}

func TlvInt8(t uint8, v int8) Tlv {
	return Tlv{
		Type:   t,
		Length: 1,
		Value:  []byte{byte(v)},
	}
}

func TlvInt16(t uint8, v int16) Tlv {
	b := binary.BigEndian.AppendUint16([]byte{}, uint16(v))
	return Tlv{
		Type:   t,
		Length: 2,
		Value:  b,
	}
}

func TlvInt32(t uint8, v int32) Tlv {
	b := binary.BigEndian.AppendUint32([]byte{}, uint32(v))
	return Tlv{
		Type:   t,
		Length: 4,
		Value:  b,
	}
}

func TlvInt64(t uint8, v int64) Tlv {
	b := binary.BigEndian.AppendUint64([]byte{}, uint64(v))
	return Tlv{
		Type:   t,
		Length: 8,
		Value:  b,
	}
}

func TlvUInt8(t uint8, v uint8) Tlv {
	return Tlv{
		Type:   t,
		Length: 4,
		Value:  []byte{byte(v)},
	}
}

func TlvUInt16(t uint8, v uint16) Tlv {
	b := binary.BigEndian.AppendUint16([]byte{}, v)
	return Tlv{
		Type:   t,
		Length: 4,
		Value:  b,
	}
}

func TlvUInt32(t uint8, v uint32) Tlv {
	b := binary.BigEndian.AppendUint32([]byte{}, v)
	return Tlv{
		Type:   t,
		Length: 4,
		Value:  b,
	}
}

func TlvUInt64(t uint8, v uint64) Tlv {
	b := binary.BigEndian.AppendUint64([]byte{}, v)
	return Tlv{
		Type:   t,
		Length: 4,
		Value:  b,
	}
}

func TlvFloat32(t uint8, v float32) Tlv {
	u := math.Float32bits(v)
	b := binary.BigEndian.AppendUint32([]byte{}, uint32(u))
	return Tlv{
		Type:   t,
		Length: 4,
		Value:  b,
	}
}

func TlvFloat64(t uint8, v float64) Tlv {
	u := math.Float64bits(v)
	b := binary.BigEndian.AppendUint64([]byte{}, uint64(u))
	return Tlv{
		Type:   t,
		Length: 8,
		Value:  b,
	}
}

func TlvString(t uint8, v string) Tlv {
	return NewTlv(t, []byte(v))
}

func Join(tlvs ...Tlv) []byte {
	rs := make([]byte, 0)
	for _, tlv := range tlvs {
		rs = append(rs, tlv.Pack()...)
	}
	return rs
}

func GetAll(b []byte) []Tlv {
	rs := make([]Tlv, 0)
	ml := uint32(len(b))
	for i := uint32(0); i < ml; {
		rt := b[i]
		i += 1
		s := binary.BigEndian.Uint32(b[i:])
		i += 4
		rs = append(rs, Tlv{
			Type:   rt,
			Length: s,
			Value:  b[i : i+s],
		})
		i += s
	}
	return rs
}

func GetTlv(t uint8, b []byte) (Tlv, error) {
	ml := uint32(len(b))
	for i := uint32(0); i < ml; {
		rt := b[i]
		i += 1
		s := binary.BigEndian.Uint32(b[i:])
		i += 4
		if rt == t {
			return Tlv{
				Type:   t,
				Length: s,
				Value:  b[i : i+s],
			}, nil
		}
		i += s
	}
	return Tlv{}, fmt.Errorf("not found given type")
}

func (t Tlv) IsNullOrEmpty() bool {
	return t.Value == nil || len(t.Value) == 0
}

func (t Tlv) GetString() string {
	if t.IsNullOrEmpty() {
		return ""
	}
	return string(t.Value)
}

func (t Tlv) GetBool() bool {
	if t.IsNullOrEmpty() {
		return false
	}
	return t.Value[0] == 0x01
}

func (t Tlv) GetInt8() int8 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return int8(t.Value[0])
}

func (t Tlv) GetInt16() int16 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return int16(binary.BigEndian.Uint16(t.Value))
}

func (t Tlv) GetInt32() int16 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return int16(binary.BigEndian.Uint32(t.Value))
}

func (t Tlv) GetInt64() int16 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return int16(binary.BigEndian.Uint64(t.Value))
}

func (t Tlv) GetUInt8() uint8 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return t.Value[0]
}

func (t Tlv) GetUInt16() uint16 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return binary.BigEndian.Uint16(t.Value)
}

func (t Tlv) GetUInt32() uint32 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return binary.BigEndian.Uint32(t.Value)
}

func (t Tlv) GetUInt64() uint64 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return binary.BigEndian.Uint64(t.Value)
}

func (t Tlv) GetFloat32() float32 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return math.Float32frombits(binary.BigEndian.Uint32(t.Value))
}

func (t Tlv) GetFloat64() float64 {
	if t.IsNullOrEmpty() {
		return 0
	}
	return math.Float64frombits(binary.BigEndian.Uint64(t.Value))
}
