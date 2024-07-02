package nwtt

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/internal/report"
	"github.com/free5gc/go-upf/pkg/factory"
	logger_util "github.com/free5gc/util/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type NWTTServer struct {
	cfg *factory.Config //undone
	// BridgeMacAddress net.HardwareAddr
	NumOfNWTTPorts  uint32
	NumOfDSTTPorts  uint32
	ListOfNWTTPorts []uint32
	ListOfDSTTPorts []uint32
	/* NW-TT -> DS-TT */
	// PortPair map[uint32]uint32
	/* Traffic class -> IngressPort -> EgressPort -> txBridgeDelaly*/
	// BD5GS map[uint8]map[uint32]map[uint32]*BridgeDelay5GS
	/* EgressPort -> DelayValue*/
	// PDelayPerEgressPort map[uint32]uint64
	/* EgressPort -> Traffic class -> Priority*/
	// TrafficPriorityPerPort map[uint32]map[uint8]uint8
	/* PortNum -> Capability -> boolean*/
	PortCapabilityList          map[uint32]map[uint16]bool
	UserPlaneNodeCapabilityList map[uint16]bool
	handler                     report.Handler
	// ThresholdVal         ThresholdMesurement
	ctx    context.Context
	cancel context.CancelFunc
	log    *logrus.Entry
}

func ParsePorts(cfg string) ([]uint32, error) {
	tmp := strings.Split(cfg, ",")
	res := []uint32{}
	for _, v := range tmp {
		r, err := strconv.ParseUint(v, 10, 32)
		res = append(res, uint32(r))
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
func NewNWTTServer(cfg *factory.Config, driver forwarder.Driver) (*NWTTServer, error) {
	nwttports, err := ParsePorts(cfg.NWTT.NwttPorts)
	if err != nil {
		fmt.Printf("Nwttports is not valid and NWTT Server crushed with %s\n", err)
		return nil, errors.Errorf("Parse UPF Nwttports error")
	}
	dsttports, err := ParsePorts(cfg.NWTT.DsttPorts)
	if err != nil {
		fmt.Printf("Dsttports is not valid and NWTT Server crushed with %s\n", err)
		return nil, errors.Errorf("Parse UPF Dsttports error")
	}
	return &NWTTServer{
		cfg: cfg,
		// BridgeMacAddress:       MacAddress,
		NumOfNWTTPorts:  uint32(len(nwttports)),
		NumOfDSTTPorts:  uint32(len(dsttports)),
		ListOfNWTTPorts: nwttports,
		ListOfDSTTPorts: dsttports,
		// PortPair:               make(map[uint32]uint32),
		// BD5GS:                  make(map[uint8]map[uint32]map[uint32]*BridgeDelay5GS),
		// PDelayPerEgressPort:    make(map[uint32]uint64),
		// TrafficPriorityPerPort: make(map[uint32]map[uint8]uint8),
		PortCapabilityList:          make(map[uint32]map[uint16]bool),
		UserPlaneNodeCapabilityList: make(map[uint16]bool),
		// LLDPAdminStatus:        LLDPrxtx,
		// PortslldpAdminStatus:   make(map[uint32]uint8),
		handler: nil,
		// ThresholdVal:           ThresholdMesurement{0, 0, "", 0},
		ctx:    nil,
		cancel: nil,
		log:    logger.NwttLog.WithField(logger_util.FieldCategory, "NWTT"),
	}, nil
}

func (n *NWTTServer) Init() error {
	n.CreatePortCapability()
	n.CreateUserPlaneNodeCapability()

	// Init
	n.ReportTSCmanagemantInformation()
	return nil
}

/* 24.539 9.3 Port management capability */
func (n *NWTTServer) CreatePortCapability() error {
	for _, i := range n.ListOfNWTTPorts {
		n.PortCapabilityList[i] = map[uint16]bool{
			PMIC_SupportedPTPInstanceTypes:     SUPPORT,
			PMIC_SupportedTransportTypes:       SUPPORT,
			PMIC_SupportedDelayMechanisms:      SUPPORT,
			PMIC_PTPGrandmasterCapable:         UNSUPPORT,
			PMIC_gPTPGrandmasterCapable:        UNSUPPORT,
			PMIC_SupportedPTPProfiles:          SUPPORT,
			PMIC_NumberOfSupportedPTPInstances: UNSUPPORT,
			PMIC_PTPInstanceList:               UNSUPPORT,
		}
	}
	return nil
}

/* 24.539 9.5C User Plane Node management capability */
func (n *NWTTServer) CreateUserPlaneNodeCapability() error {
	n.UserPlaneNodeCapabilityList = map[uint16]bool{
		/* Information for 5GS Bridge(Read only) */
		UMIC_UserPlaneNodeAddress: UNSUPPORT,
		UMIC_UserPlaneNodeID:      SUPPORT,
		UMIC_NWTTPortNumbers:      UNSUPPORT,
		/* Time synchronization information(Read only) */
		UMIC_SupportedPTPInstanceType:            UNSUPPORT,
		UMIC_SupportedTransportTypes:             UNSUPPORT,
		UMIC_SupportedDelayMechanisms:            SUPPORT,
		UMIC_PTPGrandmasterCapable:               UNSUPPORT,
		UMIC_GPTPGrandmasterCapable:              UNSUPPORT,
		UMIC_SupportedPTPProfiles:                SUPPORT,
		UMIC_NumberOfSupportedPTPInstances:       UNSUPPORT,
		UMIC_DSTTPortTimeSynchronizationInfoList: UNSUPPORT,
		UMIC_PTPInstanceSpecification:            UNSUPPORT,
	}
	return nil
}

func (n *NWTTServer) ReportTSCmanagemantInformation() error {
	// var tmirs []report.TMIReport
	umic, err := n.EncodeUserPlaneNodeManagementCapability()
	if err != nil {
		n.log.Errorln("EncodeUserPlaneNodeManagementCapability", err)
		return err
	}
	if n.handler != nil {
		// TODO: get SEID
		n.handler.NotifySessReport(report.SessReport{
			SEID: 1, // SEID(Session Endpoint Identifier)
			Reports: []report.Report{
				report.TMIReport{
					PMIC: umic,
				}},
		})
	}
	// TODO: different pmic map to different port
	for _, i := range n.ListOfNWTTPorts {
		pmic, err := n.EncodePortManagementCapability(i)
		if err != nil {
			n.log.Errorln("EncodePortManagementCapability", err)
			return err
		}
		if n.handler != nil {
			// TODO: get SEID
			n.handler.NotifySessReport(report.SessReport{
				SEID: 1, // SEID(Session Endpoint Identifier)
				Reports: []report.Report{
					report.TMIReport{
						PMIC:    pmic,
						PortNum: i,
					}},
			})
		}
	}

	return nil
}
