package tcp

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

type Msg struct {
	Version   uint32
	Size      uint32
	TotalPage uint32
	Page      uint32
	Payload   []byte
}

func (m Msg) Pack() []byte {
	bs := []byte{
		Stx, 0x4D, 0x46, 0x45,
	}
	bs = binary.BigEndian.AppendUint32(bs, m.Version)
	bs = binary.BigEndian.AppendUint32(bs, m.Size)
	bs = binary.BigEndian.AppendUint32(bs, m.TotalPage)
	bs = binary.BigEndian.AppendUint32(bs, m.Page)
	bs = append(bs, m.Payload...)
	bs = append(bs, Etx)
	return bs
}

const (
	Version = uint32(1)
	Stx     = 0x02
	Etx     = 0x03
)

func ValidateHeader(bytes []byte) bool {
	return bytes[0] == 0x4D && bytes[1] == 0x46 && bytes[2] == 0x45
}

func ReadFull(r *bufio.Reader, length uint32) ([]byte, error) {
	bs := make([]byte, length)
	n := uint32(0)
	for n < length {
		m, err := r.Read(bs[n:])
		if err != nil {
			return nil, err
		}
		n += uint32(m)
	}
	return bs, nil
}

func Read(r *bufio.Reader) (Msg, error) {
	var msg Msg
	b, err := r.ReadByte()
	if err != nil {
		return msg, err
	}
	if b != Stx {
		return msg, fmt.Errorf("beginning of msg must be STX")
	}
	headers, err := ReadFull(r, 3)
	if err != nil {
		return msg, err
	}
	if !ValidateHeader(headers) {
		return msg, fmt.Errorf("header is not correct")
	}
	//read version
	arr, err := ReadFull(r, 4)
	if err != nil {
		return msg, err
	}
	msg.Version = binary.BigEndian.Uint32(arr)
	//read size
	arr, err = ReadFull(r, 4)
	if err != nil {
		return msg, err
	}
	msg.Size = binary.BigEndian.Uint32(arr)
	// read total page
	arr, err = ReadFull(r, 4)
	if err != nil {
		return msg, err
	}
	msg.TotalPage = binary.BigEndian.Uint32(arr)
	// read page
	arr, err = ReadFull(r, 4)
	if err != nil {
		return msg, err
	}
	msg.Page = binary.BigEndian.Uint32(arr)
	// read payload
	arr, err = ReadFull(r, msg.Size)
	if err != nil {
		return msg, err
	}
	msg.Payload = arr
	b, err = r.ReadByte()
	if err != nil {
		return msg, err
	}
	if b != Etx {
		return msg, fmt.Errorf("end of msg must be ETX")
	}
	return msg, nil
}

func Pack(payload []byte, maxSize uint32) ([][]byte, error) {
	totalPage := uint32(1)
	if uint32(len(payload)) > maxSize {
		totalPage = uint32(len(payload)) / maxSize
		if uint32(len(payload))%maxSize > 0 {
			totalPage += 1
		}
	}
	rs := make([][]byte, 0)
	for page := uint32(0); page < totalPage; page++ {
		size := uint32(len(payload))
		if size > maxSize {
			size = maxSize
		}
		msg := Msg{
			Version:   Version,
			Size:      size,
			TotalPage: totalPage,
			Page:      page + 1,
			Payload:   payload[:size],
		}
		payload = payload[size:]
		rs = append(rs, msg.Pack())
	}
	return rs, nil
}
