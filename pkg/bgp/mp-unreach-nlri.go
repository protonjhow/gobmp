package bgp

import (
	"encoding/binary"
	"fmt"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/base"
	"github.com/sbezverk/gobmp/pkg/evpn"
	"github.com/sbezverk/gobmp/pkg/flowspec"
	"github.com/sbezverk/gobmp/pkg/l3vpn"
	"github.com/sbezverk/gobmp/pkg/ls"
	"github.com/sbezverk/gobmp/pkg/srpolicy"
	"github.com/sbezverk/gobmp/pkg/tools"
	"github.com/sbezverk/gobmp/pkg/unicast"
)

// MPUnReachNLRI defines an MP UnReach NLRI object
type MPUnReachNLRI struct {
	AddressFamilyID    uint16
	SubAddressFamilyID uint8
	WithdrawnRoutes    []byte
	EndOfRIB           bool
}

// GetAFISAFIType returns underlaying NLRI's type based on AFI/SAFI
func (mp *MPUnReachNLRI) GetAFISAFIType() int {
	return getNLRIMessageType(mp.AddressFamilyID, mp.SubAddressFamilyID)
}

// IsIPv6NLRI return true if NLRI is for IPv6 address family
func (mp *MPUnReachNLRI) IsIPv6NLRI() bool {
	return mp.AddressFamilyID == 2
}

// GetNextHop return a string representation of the next hop ip address.
func (mp *MPUnReachNLRI) GetNextHop() string {
	return ""
}

// IsNextHopIPv6 return true if the next hop is IPv6 address, otherwise it returns flase.
// in case of MP_UNREACH_NLRI there is no Next Hope field and this func should not be used.
func (mp *MPUnReachNLRI) IsNextHopIPv6() bool {
	return false
}

// GetNLRI71 check for presense of NLRI 71 in the NLRI 14 NLRI data and if exists, instantiate NLRI71 object
func (mp *MPUnReachNLRI) GetNLRI71() (*ls.NLRI71, error) {
	if mp.SubAddressFamilyID == 71 {
		if mp.EndOfRIB {
			return &ls.NLRI71{}, nil
		}
		nlri71, err := ls.UnmarshalLSNLRI71(mp.WithdrawnRoutes)
		if err != nil {
			return nil, err
		}
		return nlri71, nil
	}

	// TODO return new type of errors to be able to check for the code
	return nil, fmt.Errorf("not found")
}

// GetNLRI73 check for presense of NLRI 73 in the NLRI 14 NLRI data and if exists, instantiate NLRI73 object
func (mp *MPUnReachNLRI) GetNLRI73() (*srpolicy.NLRI73, error) {
	if mp.SubAddressFamilyID == 73 {
		if mp.EndOfRIB {
			return &srpolicy.NLRI73{}, nil
		}
		nlri73, err := srpolicy.UnmarshalLSNLRI73(mp.WithdrawnRoutes)
		if err != nil {
			return nil, err
		}
		return nlri73, nil
	}

	// TODO return new type of errors to be able to check for the code
	return nil, fmt.Errorf("not found")
}

// GetNLRIL3VPN check for presense of NLRI L3VPN AFI 1 and SAFI 128 in the NLRI 14 NLRI data and if exists, instantiate L3VPN object
func (mp *MPUnReachNLRI) GetNLRIL3VPN() (*base.MPNLRI, error) {
	if mp.AddressFamilyID == 1 && mp.SubAddressFamilyID == 128 {
		if mp.EndOfRIB {
			return &base.MPNLRI{}, nil
		}
		nlri, err := l3vpn.UnmarshalL3VPNNLRI(mp.WithdrawnRoutes)
		if err != nil {
			return nil, err
		}
		return nlri, nil
	}

	// TODO return new type of errors to be able to check for the code
	return nil, fmt.Errorf("not found")
}

// GetNLRIEVPN check for presense of NLRI EVPN AFI 25 and SAFI 70 in the NLRI 14 NLRI data and if exists, instantiate EVPN object
func (mp *MPUnReachNLRI) GetNLRIEVPN() (*evpn.Route, error) {
	if mp.AddressFamilyID == 25 && mp.SubAddressFamilyID == 70 {
		if mp.EndOfRIB {
			return &evpn.Route{}, nil
		}
		route, err := evpn.UnmarshalEVPNNLRI(mp.WithdrawnRoutes)
		if err != nil {
			return nil, err
		}
		return route, nil
	}

	// TODO return new type of errors to be able to check for the code
	return nil, fmt.Errorf("not found")
}

// GetNLRIUnicast check for presense of NLRI EVPN AFI 1 or 2  and SAFI 1 in the NLRI 14 NLRI data and if exists, instantiate Unicast object
func (mp *MPUnReachNLRI) GetNLRIUnicast() (*base.MPNLRI, error) {
	if (mp.AddressFamilyID == 1 || mp.AddressFamilyID == 2) && mp.SubAddressFamilyID == 1 {
		if mp.EndOfRIB {
			return &base.MPNLRI{}, nil
		}
		nlri, err := unicast.UnmarshalUnicastNLRI(mp.WithdrawnRoutes)
		if err != nil {
			return nil, err
		}
		return nlri, nil
	}

	// TODO return new type of errors to be able to check for the code
	return nil, fmt.Errorf("not found")
}

// GetNLRILU check for presense of NLRI EVPN AFI 1 or 2  and SAFI 4 in the NLRI 14 NLRI data and if exists, instantiate Unicast object
func (mp *MPUnReachNLRI) GetNLRILU() (*base.MPNLRI, error) {
	if (mp.AddressFamilyID == 1 || mp.AddressFamilyID == 2) && mp.SubAddressFamilyID == 4 {
		if mp.EndOfRIB {
			return &base.MPNLRI{}, nil
		}
		nlri, err := unicast.UnmarshalLUNLRI(mp.WithdrawnRoutes)
		if err != nil {
			return nil, err
		}
		return nlri, nil
	}

	// TODO return new type of errors to be able to check for the code
	return nil, fmt.Errorf("not found")
}

// GetFlowspecNLRI checks for presense of NLRI 133 IPv4 Flowspec in the NLRI 15 NLRI data and if exists, instantiate NLRI object
func (mp *MPUnReachNLRI) GetFlowspecNLRI() (*flowspec.NLRI, error) {
	if mp.SubAddressFamilyID == 133 {
		if mp.EndOfRIB {
			return &flowspec.NLRI{}, nil
		}
		return flowspec.UnmarshalFlowspecNLRI(mp.WithdrawnRoutes)
	}

	// TODO return new type of errors to be able to check for the code
	return nil, fmt.Errorf("not found")
}

// UnmarshalMPUnReachNLRI builds MP Reach NLRI attributes
func UnmarshalMPUnReachNLRI(b []byte) (MPNLRI, error) {
	if glog.V(6) {
		glog.Infof("MPUnReachNLRI Raw: %s", tools.MessageHex(b))
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("NLRI length is 0")
	}
	mp := MPUnReachNLRI{}
	p := 0
	mp.AddressFamilyID = binary.BigEndian.Uint16(b[p : p+2])
	p += 2
	mp.SubAddressFamilyID = uint8(b[p])
	p++
	if p < len(b) {
		mp.WithdrawnRoutes = make([]byte, len(b[p:]))
		copy(mp.WithdrawnRoutes, b[p:])
	} else {
		// https://www.rfc-editor.org/rfc/rfc4724.html#section-2
		// Marker for End-of-RIB
		mp.EndOfRIB = true
	}

	return &mp, nil
}
