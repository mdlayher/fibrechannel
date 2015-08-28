package fibrechannel

import (
	"encoding/binary"
	"io"
)

const (
	// headerLen specifies the exact required length for a Header.
	headerLen = 24
)

// A Header is a Fibre Channel header.  A Header contains metadata regarding
// a Frame.
type Header struct {
	// RoutingControl (R_CTL) identifies a frame's type and information category.
	RoutingControl RoutingControl

	// DestinationID (D_ID) identifies a frame's destination.
	DestinationID [3]byte

	// Priority (CS_CTL/PRI or Class Specific Control) identifies the priority
	// of traffic contained in a frame.
	Priority byte

	// SourceID (S_ID) identifies a frame's source.
	SourceID [3]byte

	// Type identifies the type of traffic contained in a frame.
	Type byte

	// FrameControl (F_CTL) contains a variety of frame options.
	FrameControl [3]byte

	// SequenceID (SEQ_ID) identifies which sequence a frame belongs to.
	SequenceID uint8

	// DataFieldControl (DF_CTL) indicates the presence of optional headers
	// and their size.
	DataFieldControl byte

	// SequenceCount (SEQ_CNT) identifies a frame's number within a sequence.
	SequenceCount uint16

	// OriginatorExchangeID (OX_ID) is an ID assigned by an initiator, used to
	// group related sequences.
	OriginatorExchangeID uint16

	// ResponderExchangeID (RX_ID) is an ID assigned by a target, also used to
	// group related sequences.
	ResponderExchangeID uint16

	// Parameter is used as a relative offset in sequences.
	Parameter uint32
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
//
// MarshalBinary never returns an error.
func (h *Header) MarshalBinary() ([]byte, error) {
	b := make([]byte, headerLen)

	b[0] = byte(h.RoutingControl)
	copy(b[1:4], h.DestinationID[:])
	b[4] = h.Priority
	copy(b[5:8], h.SourceID[:])
	b[8] = h.Type
	copy(b[9:12], h.FrameControl[:])
	b[12] = h.SequenceID
	b[13] = h.DataFieldControl

	binary.BigEndian.PutUint16(b[14:16], h.SequenceCount)
	binary.BigEndian.PutUint16(b[16:18], h.OriginatorExchangeID)
	binary.BigEndian.PutUint16(b[18:20], h.ResponderExchangeID)
	binary.BigEndian.PutUint32(b[20:24], h.Parameter)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a Header.
//
// If the byte slice does not contain exactly enough data to unmarshal a valid
// Header, io.ErrUnexpectedEOF is returned.
func (h *Header) UnmarshalBinary(b []byte) error {
	if len(b) != headerLen {
		return io.ErrUnexpectedEOF
	}

	// TODO(mdlayher): optimize to avoid allocating a bunch of times.
	// Unfortunately, this may mean using slices instead of arrays,
	// which means having to do lots of length checks when marshaling.

	bb := make([]byte, headerLen)
	copy(bb, b)

	h.RoutingControl = RoutingControl(bb[0])

	var dstID [3]byte
	copy(dstID[:], bb[1:4])
	h.DestinationID = dstID

	h.Priority = bb[4]

	var srcID [3]byte
	copy(srcID[:], bb[5:8])
	h.SourceID = srcID

	h.Type = bb[8]

	var fctl [3]byte
	copy(fctl[:], bb[9:12])
	h.FrameControl = fctl

	h.SequenceID = bb[12]
	h.DataFieldControl = bb[13]

	h.SequenceCount = binary.BigEndian.Uint16(bb[14:16])
	h.OriginatorExchangeID = binary.BigEndian.Uint16(bb[16:18])
	h.ResponderExchangeID = binary.BigEndian.Uint16(bb[18:20])
	h.Parameter = binary.BigEndian.Uint32(bb[20:24])

	return nil
}
