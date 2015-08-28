// Package fibrechannel implements marshaling and unmarshaling of Fibre
// Channel frames.
package fibrechannel

//go:generate stringer -output=string.go -type=EOF,RoutingControl,SOF

// A SOF is a Start-of-Frame byte, which appears at the beginning of a
// Fibre Channel frame.
type SOF byte

// SOF constants which indicate different types of Frame traffic.
const (
	SOFf  SOF = 0x28 // Start fabric
	SOFi2 SOF = 0x2d // Start class 2
	SOFi3 SOF = 0x2e // Start class 3
	SOFi4 SOF = 0x29 // Start class 4
	SOFn2 SOF = 0x35 // Normal (continue) class 2
	SOFn3 SOF = 0x36 // Normal (continue) class 3
	SOFn4 SOF = 0x31 // Normal (continue) class 4
	SOFc4 SOF = 0x39
)

// A EOF is an End-of-Frame byte, which appears at the end of a Fibre Channel
// frame.
type EOF byte

// EOF constants which indicate continuation or termination.
const (
	EOFn   EOF = 0x41 // Normal (not last frame of sequence)
	EOFt   EOF = 0x42 // Terminate (last frame of sequence)
	EOFrt  EOF = 0x44
	EOFdt  EOF = 0x46 // Disconnect-terminate class-1
	EOFni  EOF = 0x49 // Normal-invalid
	EOFdti EOF = 0x4e // Disconnect-terminate-invalid
	EOFrti EOF = 0x4f
	EOFa   EOF = 0x50 // Abort
)

// A RoutingControl is a byte which appears in the R_CTL field of a Fibre
// Channel header.
type RoutingControl byte

// RoutingControl constants which indicate different frame types and
// information categories.
const (
	RoutingControlDeviceDataUncategorized       RoutingControl = 0x00
	RoutingControlDeviceDataSolicitedData       RoutingControl = 0x01
	RoutingControlDeviceDataUnsolicitiedControl RoutingControl = 0x02
	RoutingControlDeviceDataSolicitedControl    RoutingControl = 0x03
	RoutingControlDeviceDataUnsolicitedData     RoutingControl = 0x04
	RoutingControlDeviceDataDataDescriptor      RoutingControl = 0x05
	RoutingControlDeviceDataUnsolicitedCommand  RoutingControl = 0x06
	RoutingControlDeviceDataCommandStatus       RoutingControl = 0x07
)
