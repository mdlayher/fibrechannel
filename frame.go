package fibrechannel

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

const (
	// frameLen specifies the minimum required length for a Frame.
	frameLen = 4 + headerLen + 4 + 4
)

var (
	// ErrInvalidFrame is returned when one of the following occur:
	//  - an invalid SOF or EOF sequence is detected
	//  - a payload does not end on a word (4 byte) boundary
	//  - no Header is present when attempting to binary marshal a Frame
	ErrInvalidFrame = errors.New("invalid frame")

	// ErrInvalidCRC is returned when Frame.UnmarshalBinary detects an incorrect
	// CRC checksum in a byte slice for a Frame.
	ErrInvalidCRC = errors.New("invalid frame checksum")
)

// A Frame is a Fibre Channel frame.  A Frame contains information such
// as SOF and EOF bytes, a header, and payload data.
type Frame struct {
	// SOF specifies the Start-of-Frame byte contained in this Frame.
	SOF SOF

	// Header specifies a Fibre Channel header, which contains metadata
	// regarding this Frame.
	Header *Header

	// Payload is a variable length data payload encapsulated by this Frame.
	// Payload is automatically padded to the next word (4 byte) boundary when
	// marshaled into binary form.
	Payload []byte

	// EOF specifies the End-of-Frame byte contained in this Frame.
	EOF EOF
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
//
// Frame.Header must not be nil, or ErrInvalidFrame will be returned.
func (f *Frame) MarshalBinary() ([]byte, error) {
	// Header must not be nil
	if f.Header == nil {
		return nil, ErrInvalidFrame
	}

	b := make([]byte, f.length())

	// Insert SOF at starting byte 4, EOF at byte 4 from end
	b[3] = byte(f.SOF)
	b[len(b)-4] = byte(f.EOF)

	// Marshal header and insert after SOF
	hb, err := f.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(b[4:4+headerLen], hb)

	// Copy payload up until CRC
	copy(b[4+headerLen:len(b)-8], f.Payload)

	// Calculate CRC and place at 8 bytes from end
	binary.BigEndian.PutUint32(b[len(b)-8:len(b)-4], crc32.ChecksumIEEE(b[4:len(b)-8]))

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a Frame.
//
// If the byte slice does not contain enough data to unmarshal a valid Frame,
// io.ErrUnexpectedEOF is returned.
//
// If an invalid SOF or EOF sequence is detected, or a payload does not end on
// a word (4 byte) boundary, ErrInvalidFrame is returned.
//
// If the CRC checksum present in the byte slice does not match the one
// computed by UnmarshalBinary, ErrInvalidCRC is returned.
func (f *Frame) UnmarshalBinary(b []byte) error {
	// Must have enough data to create a frame
	if len(b) < frameLen {
		return io.ErrUnexpectedEOF
	}

	// Some SOF and EOF bytes must be reserved
	if b[0] != 0 || b[1] != 0 || b[2] != 0 {
		return ErrInvalidFrame
	}
	if b[len(b)-3] != 0 || b[len(b)-2] != 0 || b[len(b)-1] != 0 {
		return ErrInvalidFrame
	}

	// Payload must end on word (4 byte) boundary
	if len(b[4+headerLen:len(b)-8])%4 != 0 {
		return ErrInvalidFrame
	}

	// Must have valid CRC checksum
	want := binary.BigEndian.Uint32(b[len(b)-8 : len(b)-4])
	got := crc32.ChecksumIEEE(b[4 : len(b)-8])
	if want != got {
		return ErrInvalidCRC
	}

	// Retrieve SOF and EOF bytes
	f.SOF = SOF(b[3])
	f.EOF = EOF(b[len(b)-4])

	// Unmarshal Header for Frame
	h := new(Header)
	if err := h.UnmarshalBinary(b[4:28]); err != nil {
		return err
	}
	f.Header = h

	// Copy payload up until bytes before CRC
	payload := make([]byte, len(b[28:len(b)-8]))
	copy(payload, b[28:len(b)-4])
	f.Payload = payload

	return nil
}

// length calculates the number of bytes required to store a Frame.
func (f *Frame) length() int {
	// Payload length must end on a word (4 byte) boundary.
	// If payload is not a multiple of 4 bytes, pad it up
	// to the next word boundary.
	pl := len(f.Payload)
	if r := pl % 4; r != 0 {
		pl += 4 - r
	}

	//  4 bytes: SOF
	// 24 bytes: header
	//  N bytes: payload
	//  4 bytes: CRC
	//  4 bytes: EOF
	return 4 + 24 + pl + 4 + 4
}
