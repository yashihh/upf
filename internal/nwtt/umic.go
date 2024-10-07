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
	UpdatedStatusContents := []byte{}

	var readNum uint8 = 0
	var setNum uint8 = 0

	var transportType uint8
	var grandmaster_candidate_enabled uint8
	var grandmaster_on_behalf_of_DSTT_enabled uint8
	var instanceType uint8
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
			capability := binary.BigEndian.Uint16(managementList[idx+1 : idx+3])
			switch capability {
			case PTPInstanceSpecification:
				UpdatedStatusContents = append(UpdatedStatusContents, byte(PTPInstanceSpecification>>8), byte(PTPInstanceSpecification&0xFF))
				// TODO : support more than one PTP intance
				listLength := binary.BigEndian.Uint16(managementList[idx+3 : idx+5])
				for ptpI := 0; ptpI < int(listLength); {
					ptpInstance := managementList[idx+5 : idx+5+int(listLength)]
					ptpILength := binary.BigEndian.Uint16(ptpInstance[0:2])
					ptpID := binary.BigEndian.Uint16(ptpInstance[2:4])
					if n.PTPInstanceID == 0 {
						n.PTPInstanceID = ptpID
						n.log.Infof("Set PTP Instance ID with [%d] successfully", n.PTPInstanceID)
					}
					for i := 4; i < int(ptpILength)+4; {
						parameter := binary.BigEndian.Uint16(ptpInstance[i : i+2])
						valLength := binary.BigEndian.Uint16(ptpInstance[i+2 : i+4])

						if parameter == PTP_profile {
							value := ptpInstance[i+4 : i+4+int(valLength)]
							n.log.Infof("PTP profile :[%x]", value)

						} else if parameter == Transport_type {
							transportType = ptpInstance[i+4 : i+4+int(valLength)][0]
							n.log.Infof("Transport type :[%s]", transportTypesMap[transportType])
						} else if parameter == Grandmaster_candidate_enabled {
							grandmaster_candidate_enabled = ptpInstance[i+4 : i+4+int(valLength)][0]
							n.log.Infof("Grandmaster candidate enabled :[%x]", grandmaster_candidate_enabled)
						} else if parameter == DefaultDS_instanceType {
							instanceType = ptpInstance[i+4 : i+4+int(valLength)][0]
							n.log.Infof("DefaultDS.instanceType :[%s]", ptpInstanceTypesMap[instanceType])
						} else {
							n.log.Infof("parameter [%d] not supported.", parameter)
						}
						i += 4 + int(valLength)
					}
					ptpI += 4 + int(ptpILength)
					// User plane node parameter update
					UpdatedStatusContents = append(UpdatedStatusContents, byte(listLength&0xFF))
					UpdatedStatusContents = append(UpdatedStatusContents, ptpInstance...)
				}
				idx += int(listLength) + 5

			case DSTTPortTimeSynchronizationInfoList:
				UpdatedStatusContents = append(UpdatedStatusContents, byte(DSTTPortTimeSynchronizationInfoList>>8), byte(DSTTPortTimeSynchronizationInfoList&0xFF))

				listLength := binary.BigEndian.Uint16(managementList[idx+3 : idx+5])

				for dsttI := 0; dsttI < int(listLength); {
					dsttInfo := managementList[idx+5 : idx+5+int(listLength)]
					dsttILength := binary.BigEndian.Uint16(managementList[idx+5 : idx+7])
					dsttPortNum := binary.BigEndian.Uint16(managementList[idx+7 : idx+9])

					ptpInstance := dsttInfo[4:int(dsttILength)]
					ptpILength := binary.BigEndian.Uint16(ptpInstance[0:2])
					ptpID := binary.BigEndian.Uint16(ptpInstance[2:4])
					n.log.Infof("Set parameter with PTP Instance ID[%d]", ptpID)

					for i := 4; i < int(ptpILength)+4; {
						parameter := binary.BigEndian.Uint16(ptpInstance[i : i+2])
						valLength := binary.BigEndian.Uint16(ptpInstance[i+2 : i+4])

						if parameter == Grandmaster_on_behalf_of_DSTT_enabled {
							grandmaster_on_behalf_of_DSTT_enabled = ptpInstance[i+4 : i+4+int(valLength)][0]
							n.log.Infof("set DSTT PortNum:%d with Grandmaster on behalf of DSTT enabled :[%x]", uint32(dsttPortNum), grandmaster_on_behalf_of_DSTT_enabled)

						} else if parameter == Grandmaster_candidate_enabled {
							grandmaster_candidate_enabled = ptpInstance[i+4 : i+4+int(valLength)][0]
							n.log.Infof("set DSTT PortNum:%d with Grandmaster candidate enabled :[%x]", uint32(dsttPortNum), grandmaster_candidate_enabled)

						} else {
							n.log.Infof("parameter [%d] not supported.", parameter)
						}
						i += 4 + int(valLength)
					}
					dsttI += 4 + int(dsttILength)

					UpdatedStatusContents = append(UpdatedStatusContents, byte(listLength&0xFF))
					UpdatedStatusContents = append(UpdatedStatusContents, dsttInfo...)
				}
				idx += int(listLength) + 5
			}
			setNum++
			/* TODO：Hangle correct idx*/

		case SubscribeNotifyForParameter:
			n.log.Infof("Handle UserPlaneNode SubscribeNotifyForParameter Operation")
			idx += 3
		case UnsubscribeForParameter:
			n.log.Infof("Handle UserPlaneNode UnsubscribeForParameter Operation")
			idx += 3
		default:
			return nil, errors.Errorf("Unsupport operation code in Manage UserPlaneNode Command %x with -%d", managementList[idx], idx)
		}

	}
	if readNum != 0 {
		/* Put Number of UpNode parameters successfully read into Status */
		StatusContents = append([]byte{readNum}, StatusContents...)
		/* Put Length of UpNode status and error contents */
		statusLength := make([]byte, 2)
		binary.BigEndian.PutUint16(statusLength, uint16(len(StatusContents)))
		StatusContents = append(statusLength, StatusContents...)
		/* Put UpNode status IEI  */
		StatusContents = append([]byte{UserPlaneNodeStatus}, StatusContents...)
		/* Merge to buffer*/
		buffer = append(buffer, StatusContents...)
	}
	if setNum != 0 {
		/* Put Number of UpNode parameters successfully set into UpdatedStatusContents */
		UpdatedStatusContents = append([]byte{setNum}, UpdatedStatusContents...)
		/* Put Length of UpNode Updated status and error contents */
		updatedstatusLength := make([]byte, 2)
		binary.BigEndian.PutUint16(updatedstatusLength, uint16(len(UpdatedStatusContents)))
		UpdatedStatusContents = append(updatedstatusLength, UpdatedStatusContents...)
		/* Put UpNode Updated status IEI  */
		UpdatedStatusContents = append([]byte{UserPlaneNodeUpdateResult}, UpdatedStatusContents...)

		/* Merge to buffer*/
		buffer = append(buffer, UpdatedStatusContents...)

	}
	if grandmaster_candidate_enabled == TRUE && grandmaster_on_behalf_of_DSTT_enabled == TRUE && transportType == IPv4 {
		config := ConfigurationForPTP{
			DefaultDS_instanceType: instanceType,
		}
		n.log.Infof("Send ConfigurationForPTP")

		ch <- config
	}
	return buffer, nil
}

func (n *NWTTServer) DecodeUserPlaneNodeManagementInformation(BMIC []byte) (*ie.IE, error) {
	switch uint8(BMIC[0]) {
	case ManageUserPlaneNodeCommand:
		buffer, err := n.HandleManageUserPlaneNodeCommand(BMIC)
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

	upNodeStatus := []byte{}
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
	upNodeStatus = append(upNodeStatus, parameterName...)

	/* Put Length of Port parameter value */
	upNodeStatus = append(upNodeStatus, length...)

	/* Put Port parameter value */
	upNodeStatus = append(upNodeStatus, buffer...)

	return upNodeStatus, nil
}
