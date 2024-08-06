package nwtt

import (
	"encoding/binary"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-pfcp/ie"
)

func (n *NWTTServer) DecodePortManagementInformation(PMIC []byte, PortNumber uint32) (*ie.IE, error) {
	switch uint8(PMIC[0]) {
	case ManagePortCommand:
		Buffer, err := n.HandleManagePortCommand(PMIC[1:], PortNumber)
		if err != nil {
			n.log.Errorln(err)
			return nil, err
		}
		/* append PORT MANAGEMENT COMPLETE ID*/
		Buffer = append([]byte{ManagePortComplete}, Buffer...)
		/* Return buffer wrap in PMIC */
		return ie.NewPortManagementInformationContainer(string(Buffer)), nil
	case ManagePortComplete:
	case PortManagementNotify:
	case PortManagementNotifyAck:
	case PortManagementNotifyComplete:
	case PortManagementCapability:
	default:
	}
	return nil, nil
}

func (n *NWTTServer) HandleManagePortCommand(managementList []byte, portNumber uint32) ([]byte, error) {
	buffer := []byte{}
	StatusContents := []byte{}
	var readNum uint8 = 0
	iEI := uint8(managementList[0])
	length := binary.BigEndian.Uint16(managementList[1:3])
	if int(length) != len(managementList[3:]) {
		return nil, errors.New("PMIC Length IE mismatch with length of Port management list contents")
	}
	n.log.Infof("HandleManagePortCommand IEI: [%d], length:[%d], len(list):[%d]", iEI, length, len(managementList[3:]))
	for idx := 3; idx < int(length)+3; {
		switch uint8(managementList[idx]) {
		case GetCapabilities:
			n.log.Debugln("Handle Port GetCapabilities Operation")
			contents, err := n.EncodePortManagementCapability(portNumber)
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
			n.log.Infof("Handle Port ReadParameter Operation")
			parameter := binary.BigEndian.Uint16(managementList[idx+1 : idx+3])
			contents, err := n.EncodePortStatus(portNumber, parameter)
			if err != nil {
				n.log.Errorln(err)
				return nil, err
			}
			readNum++
			StatusContents = append(StatusContents, contents...)
			idx += 3
		case SetParameter:
			n.log.Infof("Handle Port SetParameter Operation")
			/* TODO：Hangle correct idx*/
			idx += 3
		case SubscribeNotifyForParameter:
			n.log.Infof("Handle Port SubscribeNotifyForParameter Operation")
			idx += 3
		case UnsubscribeForParameter:
			n.log.Infof("Handle Port UnsubscribeForParameter Operation")
			idx += 3
		default:
			return nil, errors.Errorf("Unsupport operation code in Manage Ethernet Port Command")
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

func (n *NWTTServer) EncodePortStatus(portNumber uint32, parameter uint16) ([]byte, error) {
	/* TODO: move default value to config */
	n.log.Infof("Build PortStatus = [0x%x]", parameter)

	portStatus := []byte{}
	parameterName := make([]byte, 2)
	length := make([]byte, 2)
	buffer := []byte{}
	binary.BigEndian.PutUint16(parameterName, parameter)

	switch parameter {
	case PMIC_SupportedPTPInstanceTypes:
		buffer = append(buffer, BoundaryClock)
		buffer = append(buffer, E2ETransparentClock)

	case PMIC_SupportedTransportTypes: // only support ipv4 currently
		buffer = append(buffer, IPv4)

	case PMIC_SupportedDelayMechanisms: // only support E2E currently
		buffer = append(buffer, E2E)

	case PMIC_PTPGrandmasterCapable:
		buffer = append(buffer, TRUE)

	case PMIC_gPTPGrandmasterCapable:
		buffer = append(buffer, FALSE)

	case PMIC_SupportedPTPProfiles:
		buffer = append(buffer, E2EDefault)

	case PMIC_NumberOfSupportedPTPInstances: // currently support only one instance
		value := make([]byte, 2)
		binary.BigEndian.PutUint16(value, 1)
		buffer = append(buffer, value...)

	case PMIC_PTPInstanceList: // currently support only one instance
		ptpInstance := getCurrentPortDataSet()

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

// func (n *NWTTServer) EncodePortUpdateResult(portNumber uint32, parameter uint16) ([]byte, error) {
// 	portUpdateResult := []byte{}
// 	parameterName := make([]byte, 2)
// 	length := make([]byte, 2)
// 	return portUpdateResult, nil
// }

// TODO: make ptp instance IE
func getCurrentPortDataSet() PortDataSet {
	cmd := exec.Command("sudo", "pmc", "-u", "-b", "0", "GET PORT_DATA_SET")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("exec ptp4l fail. [%s]", err)
	}

	fmt.Println("CURRENT_DATA_SET:")
	fmt.Println(string(output))

	length := make([]byte, 2)
	// ptpInstanceID :=
	var portDS PortDataSet

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "portIdentity"):
			portIdentity := parsePortIdentity(line)
			portDS.PortIdentity = portIdentity

		case strings.HasPrefix(line, "portState"):
			portState := parsePortState(line)
			portDS.PortState = portState

		case strings.HasPrefix(line, "logMinDelayReqInterval"):
			fmt.Sscanf(line, "logMinDelayReqInterval %d", &portDS.LogMinDelayReqInterval)

		case strings.HasPrefix(line, "peerMeanPathDelay"):
			fmt.Sscanf(line, "peerMeanPathDelay %d", &portDS.PeerMeanPathDelay)

		case strings.HasPrefix(line, "logAnnounceInterval"):
			fmt.Sscanf(line, "logAnnounceInterval %d", &portDS.LogAnnounceInterval)

		case strings.HasPrefix(line, "announceReceiptTimeout"):
			fmt.Sscanf(line, "announceReceiptTimeout %d", &portDS.AnnounceReceiptTimeout)

		case strings.HasPrefix(line, "logSyncInterval"):
			fmt.Sscanf(line, "logSyncInterval %d", &portDS.LogSyncInterval)

		case strings.HasPrefix(line, "delayMechanism"):
			fmt.Sscanf(line, "delayMechanism %d", &portDS.DelayMechanism)

		case strings.HasPrefix(line, "logMinPdelayReqInterval"):
			fmt.Sscanf(line, "logMinPdelayReqInterval %d", &portDS.LogMinPdelayReqInterval)

		case strings.HasPrefix(line, "versionNumber"):
			fmt.Sscanf(line, "versionNumber %d", &portDS.VersionNumber)

		case strings.HasPrefix(line, "delayAsymmetry"):

		case strings.HasPrefix(line, "portEnable"):

		}
	}
	return portDS
}

func parsePortIdentity(line string) PortIdentity {
	var portIdentity PortIdentity
	var clockIdentityStr string
	var portNumberStr string

	fmt.Sscanf(line, "portIdentity %s-%s", &clockIdentityStr, &portNumberStr)

	for i := 0; i < 8; i++ {
		byteStr := clockIdentityStr[i*2 : i*2+2]
		byteVal, _ := strconv.ParseUint(byteStr, 16, 8)
		portIdentity.ClockIdentity[i] = uint8(byteVal)
	}

	portNumber, _ := strconv.ParseUint(portNumberStr, 10, 16)
	portIdentity.PortNumber = uint16(portNumber)

	return portIdentity
}

func parsePortState(line string) uint8 {
	var portStateStr string
	fmt.Sscanf(line, "portState %s", &portStateStr)

	switch portStateStr {
	case "INITIALIZING":
		return INITIALIZING
	case "FAULTY":
		return FAULTY
	case "DISABLED":
		return DISABLED
	case "LISTENING":
		return LISTENING
	case "PRE_MASTER":
		return PRE_MASTER
	case "MASTER":
		return MASTER
	case "PASSIVE":
		return PASSIVE
	case "UNCALIBRATED":
		return UNCALIBRATED
	case "SLAVE":
		return SLAVE
	default:
		return 0
	}
}

func parsePeerMeanPathDelay(line string) time.Duration {
	var peerMeanPathDelay int64
	fmt.Sscanf(line, "peerMeanPathDelay %d", &peerMeanPathDelay)
	return time.Duration(peerMeanPathDelay) * time.Nanosecond
}
