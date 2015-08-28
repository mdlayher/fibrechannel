package fibrechannel

import (
	"bytes"
	"io"
	"testing"
)

func TestFrameMarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		f    *Frame
		b    []byte
		err  error
	}{
		{
			desc: "Frame with nil Header",
			f:    &Frame{},
			err:  ErrInvalidFrame,
		},
		{
			desc: "empty Frame, empty Header",
			f: &Frame{
				Header: &Header{},
			},
			b: []byte{
				0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				163, 193, 202, 32,
				0, 0, 0, 0,
			},
		},
		{
			desc: "payload less than one word Frame, empty Header",
			f: &Frame{
				Header:  &Header{},
				Payload: []byte{1, 2, 3},
			},
			b: []byte{
				0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				1, 2, 3, 0,
				16, 101, 151, 33,
				0, 0, 0, 0,
			},
		},
		{
			desc: "payload less than three words Frame, empty Header",
			f: &Frame{
				Header: &Header{},
				Payload: []byte{
					1, 2, 3, 4,
					1, 2, 3, 4,
					1,
				},
			},
			b: []byte{
				0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				1, 2, 3, 4,
				1, 2, 3, 4,
				1, 0, 0, 0,
				143, 57, 15, 200,
				0, 0, 0, 0,
			},
		},
		{
			desc: "full Frame with Header",
			f: &Frame{
				SOF: SOFi3,
				Header: &Header{
					RoutingControl:       0x01,
					DestinationID:        [3]byte{3, 3, 3},
					Priority:             5,
					SourceID:             [3]byte{6, 6, 6},
					Type:                 2,
					FrameControl:         [3]byte{1, 2, 3},
					SequenceID:           5,
					DataFieldControl:     6,
					SequenceCount:        256,
					OriginatorExchangeID: 1000,
					ResponderExchangeID:  1001,
					Parameter:            20000,
				},
				Payload: []byte{1, 2, 3},
				EOF:     EOFt,
			},
			b: []byte{
				0, 0, 0, 0x2e,
				0x01, 3, 3, 3, 5, 6, 6, 6, 2, 1, 2, 3,
				5, 6, 1, 0, 3, 232, 3, 233, 0, 0, 78, 32,
				1, 2, 3, 0,
				9, 76, 188, 232,
				0x42, 0, 0, 0,
			},
		},
	}

	for i, tt := range tests {
		b, err := tt.f.MarshalBinary()
		if err != nil || tt.err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.b, b; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Frame bytes:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}

func TestFrameUnmarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		b    []byte
		f    *Frame
		err  error
	}{
		{
			desc: "nil buffer",
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "short buffer",
			b:    bytes.Repeat([]byte{0}, frameLen-1),
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "SOF reserved byte 0 not empty",
			b: append([]byte{
				1, 0, 0, 0x2e,
			}, bytes.Repeat([]byte{0}, frameLen-4)...),
			err: ErrInvalidFrame,
		},
		{
			desc: "SOF reserved byte 1 not empty",
			b: append([]byte{
				0, 1, 0, 0x2e,
			}, bytes.Repeat([]byte{0}, frameLen-4)...),
			err: ErrInvalidFrame,
		},
		{
			desc: "SOF reserved byte 2 not empty",
			b: append([]byte{
				0, 0, 1, 0x2e,
			}, bytes.Repeat([]byte{0}, frameLen-4)...),
			err: ErrInvalidFrame,
		},
		{
			desc: "EOF reserved byte 2 not empty",
			b: append(bytes.Repeat([]byte{0}, frameLen-4),
				[]byte{0x42, 1, 0, 0}...),
			err: ErrInvalidFrame,
		},
		{
			desc: "EOF reserved byte 3 not empty",
			b: append(bytes.Repeat([]byte{0}, frameLen-4),
				[]byte{0x42, 0, 1, 0}...),
			err: ErrInvalidFrame,
		},
		{
			desc: "EOF reserved byte 4 not empty",
			b: append(bytes.Repeat([]byte{0}, frameLen-4),
				[]byte{0x42, 0, 0, 1}...),
			err: ErrInvalidFrame,
		},
		{
			desc: "payload not on word boundary",
			b: []byte{
				0, 0, 0, 0x2e,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				1, 2, 3,
				0, 0, 0, 0,
				0x42, 0, 0, 0,
			},
			err: ErrInvalidFrame,
		},
		{
			desc: "invalid CRC",
			b: []byte{
				0, 0, 0, 0x2e,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0xff, 0xff, 0xff, 0xff,
				0x42, 0, 0, 0,
			},
			err: ErrInvalidCRC,
		},
		{
			desc: "OK Frame, empty header and payload",
			b: []byte{
				0, 0, 0, 0x2e,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				163, 193, 202, 32,
				0x42, 0, 0, 0,
			},
			f: &Frame{
				SOF:    SOFi3,
				Header: &Header{},
				EOF:    EOFt,
			},
		},
		{
			desc: "OK Frame",
			b: []byte{
				0, 0, 0, 0x2e,
				0x01, 3, 3, 3, 5, 6, 6, 6, 2, 1, 2, 3,
				5, 6, 1, 0, 3, 232, 3, 233, 0, 0, 78, 32,
				1, 2, 3, 0,
				9, 76, 188, 232,
				0x42, 0, 0, 0,
			},
			f: &Frame{
				SOF: SOFi3,
				Header: &Header{
					RoutingControl:       0x01,
					DestinationID:        [3]byte{3, 3, 3},
					Priority:             5,
					SourceID:             [3]byte{6, 6, 6},
					Type:                 2,
					FrameControl:         [3]byte{1, 2, 3},
					SequenceID:           5,
					DataFieldControl:     6,
					SequenceCount:        256,
					OriginatorExchangeID: 1000,
					ResponderExchangeID:  1001,
					Parameter:            20000,
				},
				Payload: []byte{1, 2, 3, 0},
				EOF:     EOFt,
			},
		},
	}

	for i, tt := range tests {
		f := new(Frame)
		if err := f.UnmarshalBinary(tt.b); err != nil || tt.err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		fb, err := f.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		if want, got := tt.b, fb; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Frame bytes:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}
