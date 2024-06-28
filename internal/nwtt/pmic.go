package nwtt

import (
	"encoding/binary"
	"sort"
)

func (n *NWTTServer) EncodePortManagementCapability(portNumber uint32) ([]byte, error) {
	/* TODO: Get capabilities from structure*/
	portCapability := []byte{}
	b := make([]byte, 2)
	capabilities := []uint16{}
	for idx, c := range n.PortCapabilityList[portNumber] {
		if c {
			capabilities = append(capabilities, idx)
		}
	}
	//Sort capabilities for test file return
	sort.Slice(capabilities, func(i, j int) bool {
		return capabilities[i] < capabilities[j]
	})
	for _, c := range capabilities {
		binary.BigEndian.PutUint16(b, uint16(c))
		portCapability = append(portCapability, b...)
	}

	return portCapability, nil
}
