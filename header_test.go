package fibrechannel

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestHeaderMarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		h    *Header
		b    []byte
	}{
		{
			desc: "empty Header",
			h:    &Header{},
			b: []byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
		{
			desc: "full Header",
			h: &Header{
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
			b: []byte{
				0x01, 3, 3, 3, 5, 6, 6, 6, 2, 1, 2, 3,
				5, 6, 1, 0, 3, 232, 3, 233, 0, 0, 78, 32,
			},
		},
	}

	for i, tt := range tests {
		b, err := tt.h.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		if want, got := tt.b, b; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Header bytes:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}

func TestHeaderUnmarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		b    []byte
		h    *Header
		err  error
	}{
		{
			desc: "nil buffer",
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "short buffer",
			b:    bytes.Repeat([]byte{0}, headerLen-1),
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "long buffer",
			b:    bytes.Repeat([]byte{0}, headerLen+1),
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "empty Header",
			b: []byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			h: &Header{},
		},
		{
			desc: "full Header",
			b: []byte{
				0x01, 3, 3, 3, 5, 6, 6, 6, 2, 1, 2, 3,
				5, 6, 1, 0, 3, 232, 3, 233, 0, 0, 78, 32,
			},
			h: &Header{
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
		},
	}

	for i, tt := range tests {
		h := new(Header)
		if err := h.UnmarshalBinary(tt.b); err != nil || tt.err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.h, h; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected Header:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}
