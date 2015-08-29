package fip

import (
	"bytes"
	"io"
	"reflect"
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
			desc: "wrong Version",
			f: &Frame{
				Version: 2,
			},
			err: ErrInvalidFrame,
		},
		{
			desc: "OK",
			f: &Frame{
				Version:               1,
				ProtocolCode:          1,
				Subcode:               1,
				DescriptorListLength:  1,
				FlagFPMA:              true,
				FlagSPMA:              true,
				FlagAvailableForLogin: true,
				FlagSolicited:         true,
				FlagFCF:               true,
			},
			b: []byte{
				0x10, 0x00,
				0x00, 0x01,
				0x00, 0x01,
				0x00, 0x01,
				0xc0, 0x07,
			},
		},
	}

	for i, tt := range tests {
		fb, err := tt.f.MarshalBinary()
		if err != nil || tt.err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.b, fb; !bytes.Equal(want, got) {
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
			desc: "wrong Version",
			b: []byte{
				0x20, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
			},
			err: ErrInvalidFrame,
		},
		{
			desc: "reserved Version bits not empty",
			b: []byte{
				0x1f, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
			},
			err: ErrInvalidFrame,
		},
		{
			desc: "reserved byte after Version not empty",
			b: []byte{
				0x10, 0x01,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
			},
			err: ErrInvalidFrame,
		},
		{
			desc: "reserved byte after ProtocolCode not empty",
			b: []byte{
				0x10, 0x00,
				0x00, 0x00,
				0x01, 0x00,
				0x00, 0x00,
				0x00, 0x00,
			},
			err: ErrInvalidFrame,
		},
		{
			desc: "reserved bits after FlagSPMA not empty",
			b: []byte{
				0x10, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x3f, 0x00,
			},
			err: ErrInvalidFrame,
		},
		{
			desc: "reserved bits before FlagAvailableForLogin not empty",
			b: []byte{
				0x10, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x00, 0xf8,
			},
			err: ErrInvalidFrame,
		},
		{
			desc: "OK",
			b: []byte{
				0x10, 0x00,
				0x00, 0x01,
				0x00, 0x01,
				0x00, 0x01,
				0xc0, 0x07,
			},
			f: &Frame{
				Version:               1,
				ProtocolCode:          1,
				Subcode:               1,
				DescriptorListLength:  1,
				FlagFPMA:              true,
				FlagSPMA:              true,
				FlagAvailableForLogin: true,
				FlagSolicited:         true,
				FlagFCF:               true,
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

		if want, got := tt.f, f; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Frame:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}
