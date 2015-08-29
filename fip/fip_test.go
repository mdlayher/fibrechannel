package fip

import (
	"testing"
)

func TestParseOperation(t *testing.T) {
	var tests = []struct {
		p uint16
		s uint8
		o Operation
	}{
		{
			p: 0x0000,
			s: 0x00,
			o: OperationReserved,
		},
		{
			p: 0x0001,
			s: 0x01,
			o: OperationDiscoverySolicitation,
		},
		{
			p: 0x0001,
			s: 0x02,
			o: OperationDiscoveryAdvertisement,
		},
		{
			p: 0x0002,
			s: 0x01,
			o: OperationVirtualLinkInstantiationRequest,
		},
		{
			p: 0x0002,
			s: 0x02,
			o: OperationVirtualLinkInstantiationReply,
		},
		{
			p: 0x0003,
			s: 0x01,
			o: OperationKeepAlive,
		},
		{
			p: 0x0003,
			s: 0x02,
			o: OperationClearVirtualLinks,
		},
		{
			p: 0x0004,
			s: 0x01,
			o: OperationVLANRequest,
		},
		{
			p: 0x0004,
			s: 0x02,
			o: OperationVLANNotification,
		},
		{
			p: 0xfff8,
			s: 0x01,
			o: OperationVendorSpecific,
		},
	}

	for i, tt := range tests {
		if want, got := tt.o, ParseOperation(tt.p, tt.s); want != got {
			t.Fatalf("[%02d] (0x%04x, 0x%02x), unexpected Operation:\n- want: %v\n-  got: %v",
				i, tt.p, tt.s, want, got)
		}
	}
}
