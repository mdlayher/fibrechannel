package fip

import (
	"encoding/binary"
	"errors"
	"io"
)

type Frame struct {
	Version               uint8
	ProtocolCode          uint16
	Subcode               uint8
	DescriptorListLength  uint16
	FlagFPMA              bool
	FlagSPMA              bool
	FlagAvailableForLogin bool
	FlagSolicited         bool
	FlagFCF               bool
}

const (
	frameLen = 2 + 2 + 1 + 1 + 2 + 2
)

var (
	ErrInvalidFrame = errors.New("invalid frame")
)

func (f *Frame) MarshalBinary() ([]byte, error) {
	if f.Version != Version {
		return nil, ErrInvalidFrame
	}

	b := make([]byte, frameLen)

	b[0] = f.Version << 4

	binary.BigEndian.PutUint16(b[2:4], f.ProtocolCode)

	b[5] = f.Subcode

	binary.BigEndian.PutUint16(b[6:8], f.DescriptorListLength)

	if f.FlagFPMA {
		b[8] |= 1 << 7
	}
	if f.FlagSPMA {
		b[8] |= 1 << 6
	}
	if f.FlagAvailableForLogin {
		b[9] |= 1 << 2
	}
	if f.FlagSolicited {
		b[9] |= 1 << 1
	}
	if f.FlagFCF {
		b[9] |= 1
	}

	return b, nil
}

func (f *Frame) UnmarshalBinary(b []byte) error {
	if len(b) < frameLen {
		return io.ErrUnexpectedEOF
	}

	//log.Printf("%04b, %04b, %08b", ((b[0] & 0x10) >> 4), b[0]&0x01, b[1])

	v := (b[0] & 0xf0) >> 4
	if v != Version || b[0]&0x0f != 0 || b[1] != 0 {
		return ErrInvalidFrame
	}
	if b[4] != 0 {
		return ErrInvalidFrame
	}
	if b[8]&0x3f != 0 || b[9]&0xf8 != 0 {
		return ErrInvalidFrame
	}

	f.Version = v

	f.ProtocolCode = binary.BigEndian.Uint16(b[2:4])

	subcode := b[5]
	f.Subcode = subcode

	f.DescriptorListLength = binary.BigEndian.Uint16(b[6:8])

	f.FlagFPMA = b[8]&0x80 != 0
	f.FlagSPMA = b[8]&0x40 != 0
	f.FlagAvailableForLogin = b[9]&0x04 != 0
	f.FlagSolicited = b[9]&0x02 != 0
	f.FlagFCF = b[9]&0x01 != 0

	return nil
}
