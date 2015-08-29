// Package fip implements marshaling and unmarshaling of Fibre Channel over
// Ethernet (FCoE) Initialization Protocol (FIP) frames, as described in
// FC-BB-5.
package fip

const (
	// Version is the current FIP version number, as specified by FC-BB-5.
	Version uint8 = 1
)

//go:generate stringer -output=string.go -type=Operation

// An Operation is a FIP operation.  Operations are specified by the
// ProtocolCode and Subcode fields of a Frame.
type Operation int

// List of known Operations.  Their values have no meaning, but these constants
// are provided as a convenience to avoid having to check both the ProtocolCode
// and Subcode fields in a Frame.
const (
	OperationReserved Operation = iota
	OperationDiscoverySolicitation
	OperationDiscoveryAdvertisement
	OperationVirtualLinkInstantiationRequest
	OperationVirtualLinkInstantiationReply
	OperationKeepAlive
	OperationClearVirtualLinks
	OperationVLANRequest
	OperationVLANNotification
	OperationVendorSpecific
)

// ParseOperation accepts an input protocol code p and subcode s, and parses
// an Operation value from them, according the values described in FC-BB-5,
// Table 26.
func ParseOperation(p uint16, s uint8) Operation {
	switch {
	case p == 0x0001 && s == 0x01:
		return OperationDiscoverySolicitation
	case p == 0x0001 && s == 0x02:
		return OperationDiscoveryAdvertisement
	case p == 0x0002 && s == 0x01:
		return OperationVirtualLinkInstantiationRequest
	case p == 0x0002 && s == 0x02:
		return OperationVirtualLinkInstantiationReply
	case p == 0x0003 && s == 0x01:
		return OperationKeepAlive
	case p == 0x0003 && s == 0x02:
		return OperationClearVirtualLinks
	case p == 0x0004 && s == 0x01:
		return OperationVLANRequest
	case p == 0x0004 && s == 0x02:
		return OperationVLANNotification
	case p >= 0xfff8 && p <= 0xfffe:
		return OperationVendorSpecific
	}

	return OperationReserved
}
