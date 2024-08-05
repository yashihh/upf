package nwtt

import (
	"encoding/binary"
	"sort"

	"github.com/pkg/errors"
)

func (n *NWTTServer) EncodePortManagementCapability(portNumber uint32) ([]byte, error) {
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

func (n *NWTTServer) EncodePortStatus(portNumber uint32, parameter uint16, value ...byte) ([]byte, error) {
	/* TODO: Get status from AF-request */
	n.log.Infof("Build PortStatus = [0x%x]", parameter)

	portStatus := []byte{}
	parameterName := make([]byte, 2)
	length := make([]byte, 2)
	buffer := value

	binary.BigEndian.PutUint16(parameterName, parameter)

	switch parameter {
	case PMIC_SupportedPTPInstanceTypes:
		if len(value) == 0 { // default
			buffer = append(buffer, BoundaryClock)
			buffer = append(buffer, E2ETransparentClock)
		}
	case PMIC_SupportedTransportTypes: // only support ipv4 currently
		if len(value) == 0 { // default
			buffer = append(buffer, IPv4)
		}
	case PMIC_SupportedDelayMechanisms: // only support E2E currently
		if len(value) == 0 { // default
			buffer = append(buffer, E2E)
		}
	case PMIC_PTPGrandmasterCapable:
		if len(value) == 0 { // default
			buffer = append(buffer, TRUE)
		}
	case PMIC_gPTPGrandmasterCapable:
		if len(value) == 0 { // default
			buffer = append(buffer, FALSE)
		}
	case PMIC_SupportedPTPProfiles:
		if len(value) == 0 { // default
			buffer = append(buffer, E2EDefault)
		}
	case PMIC_NumberOfSupportedPTPInstances:
	case PMIC_PTPInstanceList:
	default:
		return nil, errors.Errorf("Reading unknown parameter:[%v]", parameter)
	}

	binary.BigEndian.PutUint16(length, uint16(len(buffer)))

	/* Put Name of parameter */
	portStatus = append(portStatus, parameterName...)

	/* Put Length of Port parameter value */
	portStatus = append(portStatus, length...)

	/* Put Port parameter value */
	portStatus = append(portStatus, buffer...)

	return portStatus, nil
}
