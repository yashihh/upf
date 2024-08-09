package nwtt

import (
	"encoding/binary"
	"sort"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-pfcp/ie"
)

func (n *NWTTServer) EncodeUserPlaneNodeManagementCapability() ([]byte, error) {
	/* TODO: Get capabilities from structure*/
	userPlaneNodeCapability := []byte{}
	b := make([]byte, 2)
	capabilities := []uint16{}
	for idx, c := range n.UserPlaneNodeCapabilityList {
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
		userPlaneNodeCapability = append(userPlaneNodeCapability, b...)
	}

	return userPlaneNodeCapability, nil
}
func (n *NWTTServer) HandleManageUserPlaneNodeCommand(managementList []byte) ([]byte, error) {
	buffer := []byte{}
	StatusContents := []byte{}
	var readNum uint8 = 0
	iEI := uint8(managementList[0])
	length := binary.BigEndian.Uint16(managementList[1:3])
	if int(length) != len(managementList[3:]) {
		return nil, errors.New("BMIC Length IE mismatch with length of Port management list contents")
	}
	n.log.Infof("HandleManageUserPlaneNodeCommand IEI: [%d], length:[%d], len(list):[%d]", iEI, length, len(managementList[3:]))
	for idx := 3; idx < int(length)+3; {
		switch uint8(managementList[idx]) {
		case GetCapabilities:
			n.log.Debugln("Handle Port GetCapabilities Operation")
			contents, err := n.EncodeUserPlaneNodeManagementCapability()
			if err != nil {
				n.log.Errorln(err)
				return nil, err
			}
			/* Insert length before contents */
			oplength := make([]byte, 2)
			binary.BigEndian.PutUint16(oplength, uint16(len(contents)))
			contents = append(oplength, contents...)
			/* Insert port management capability IEI */
			contents = append([]byte{PortManagementCapabilityIEI}, contents...)
			/* Append buffer */
			buffer = append(buffer, contents...)
			idx += 1
		case ReadParameter:
			n.log.Infof("Handle UserPlaneNode ReadParameter Operation")
			parameter := binary.BigEndian.Uint16(managementList[idx+1 : idx+3])
			contents, err := n.EncodeUserPlaneNodeStatus(parameter)
			if err != nil {
				n.log.Errorln(err)
				return nil, err
			}
			readNum++
			StatusContents = append(StatusContents, contents...)
			idx += 3
		case SetParameter:
			n.log.Infof("Handle UserPlaneNode SetParameter Operation")
			/* TODO：Hangle correct idx*/
			idx += 3
		case SubscribeNotifyForParameter:
			n.log.Infof("Handle UserPlaneNode SubscribeNotifyForParameter Operation")
			idx += 3
		case UnsubscribeForParameter:
			n.log.Infof("Handle UserPlaneNode UnsubscribeForParameter Operation")
			idx += 3
		default:
			return nil, errors.Errorf("Unsupport operation code in Manage UserPlaneNode Command")
		}

	}
	if readNum != 0 {
		/* Put Number of port parameters successfully read into Status */
		StatusContents = append([]byte{readNum}, StatusContents...)
		/* Put Length of port status and error contents */
		statusLength := make([]byte, 2)
		binary.BigEndian.PutUint16(statusLength, uint16(len(StatusContents)))
		StatusContents = append(statusLength, StatusContents...)
		/* Put Port status IEI  */
		StatusContents = append([]byte{PortStatusIEI}, StatusContents...)
		/* Merge to buffer*/
		buffer = append(buffer, StatusContents...)
	}
	return buffer, nil
}

func (n *NWTTServer) DecodeUserPlaneNodeManagementInformation(BMIC []byte) (*ie.IE, error) {
	switch uint8(BMIC[0]) {
	case ManageUserPlaneNodeCommand:
		buffer, err := n.EncodeUserPlaneNodeManagementCapability()
		n.log.Infof("buffer length = %v", len(buffer))
		if err != nil {
			n.log.Errorln(err)
			return nil, err
		}
		/* insert Manage Bridge Complete id*/
		buffer = append([]byte{ManageUserPlaneNodeComplete}, buffer...)
		/* Return buffer wrap in BMIC */
		return ie.NewBridgeManagementInformationContainer(string(buffer)), nil
	case ManageUserPlaneNodeComplete:
	case UserPlaneNodeManagementNotify:
	case UserPlaneNodeManagementNotifyAck:
	default:
	}
	return nil, nil
}

func (n *NWTTServer) EncodeUserPlaneNodeStatus(parameter uint16) ([]byte, error) {
	/* TODO: move default value to config */
	n.log.Infof("Build User Plane Node Status = [0x%x]", parameter)

	portStatus := []byte{}
	parameterName := make([]byte, 2)
	length := make([]byte, 2)
	buffer := []byte{}
	binary.BigEndian.PutUint16(parameterName, parameter)

	switch parameter {
	case SupportedPTPInstanceTypes:
		buffer = append(buffer, BoundaryClock)
		buffer = append(buffer, E2ETransparentClock)

	case SupportedTransportTypes: // only support ipv4 currently
		buffer = append(buffer, IPv4)

	case SupportedDelayMechanisms: // only support E2E currently
		buffer = append(buffer, E2E)

	case PTPGrandmasterCapable:
		buffer = append(buffer, TRUE)

	case gPTPGrandmasterCapable:
		buffer = append(buffer, FALSE)

	case SupportedPTPProfiles:
		buffer = append(buffer, E2EDefault)

	case NumberOfSupportedPTPInstances: // currently support only one instance
		value := make([]byte, 2)
		binary.BigEndian.PutUint16(value, 1)
		buffer = append(buffer, value...)

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
