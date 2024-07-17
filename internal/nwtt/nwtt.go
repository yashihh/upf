package nwtt

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/internal/report"
	"github.com/free5gc/go-upf/pkg/factory"
	logger_util "github.com/free5gc/util/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wmnsk/go-pfcp/ie"
)

type NWTTServer struct {
	cfg             *factory.Config
	BridgeAddr      net.HardwareAddr
	NumOfNWTTPorts  uint32
	NumOfDSTTPorts  uint32
	ListOfNWTTPorts []uint32
	ListOfDSTTPorts []uint32
	/* NW-TT -> DS-TT */
	PortPair map[uint32]uint32
	/* Traffic class -> IngressPort -> EgressPort -> txBridgeDelaly*/
	// BD5GS map[uint8]map[uint32]map[uint32]*BridgeDelay5GS
	/* EgressPort -> DelayValue*/
	// PDelayPerEgressPort map[uint32]uint64
	/* EgressPort -> Traffic class -> Priority*/
	// TrafficPriorityPerPort map[uint32]map[uint8]uint8
	/* PortNum -> Capability -> boolean*/
	PortCapabilityList          map[uint32]map[uint16]bool
	UserPlaneNodeCapabilityList map[uint16]bool
	pfcpHandler                 report.Handler
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
	upNodeID, err := net.ParseMAC(cfg.NWTT.UpNodeID)
	if err != nil {
		fmt.Printf("MacAddress is not valid and NWTT Server crushed with %s\n", err)
		return nil, errors.Errorf("Parse UPF MacAddress error")
	}
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
		cfg:             cfg,
		BridgeAddr:      upNodeID,
		NumOfNWTTPorts:  uint32(len(nwttports)),
		NumOfDSTTPorts:  uint32(len(dsttports)),
		ListOfNWTTPorts: nwttports,
		ListOfDSTTPorts: dsttports,
		PortPair:        make(map[uint32]uint32),
		// BD5GS:                  make(map[uint8]map[uint32]map[uint32]*BridgeDelay5GS),
		// PDelayPerEgressPort:    make(map[uint32]uint64),
		// TrafficPriorityPerPort: make(map[uint32]map[uint8]uint8),
		PortCapabilityList:          make(map[uint32]map[uint16]bool),
		UserPlaneNodeCapabilityList: make(map[uint16]bool),
		// LLDPAdminStatus:        LLDPrxtx,
		// PortslldpAdminStatus:   make(map[uint32]uint8),
		pfcpHandler: nil,
		// ThresholdVal:           ThresholdMesurement{0, 0, "", 0},
		ctx:    nil,
		cancel: nil,
		log:    logger.NwttLog.WithField(logger_util.FieldCategory, "NWTT"),
	}, nil
}

func (n *NWTTServer) HandlePfcp(handler report.Handler) error {
	if handler != nil {
		n.pfcpHandler = handler
		return nil
	}
	return errors.Errorf("PFCP server assign nil in NWTT.")
}

func (n *NWTTServer) Init() error {
	n.CreatePortCapability()
	n.CreateUserPlaneNodeCapability()
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

func (n *NWTTServer) NewCreatedBridgeInfo() *ie.IE {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	NwttAllocatePort := n.ListOfNWTTPorts[rand.Intn(len(n.ListOfNWTTPorts))]
	DsttAllocatePort := n.ListOfDSTTPorts[rand.Intn(len(n.ListOfDSTTPorts))]
	n.PortPair[NwttAllocatePort] = DsttAllocatePort
	return ie.NewCreatedBridgeInfoForTSC(
		ie.NewDSTTPortNumber(DsttAllocatePort),
		ie.NewFGUserPlaneNode(n.BridgeAddr),
	)
}

func (n *NWTTServer) ReportTSCmanagemantInformation(seid uint64) error {
	n.log.Infoln("Report TSC managemant Information")

	umic, err := n.EncodeUserPlaneNodeManagementCapability()
	if err != nil {
		n.log.Errorln("EncodeUserPlaneNodeManagementCapability", err)
		return err
	}

	tmiReport := report.TMIReport{
		UMIC: umic,
	}

	for _, i := range n.ListOfNWTTPorts {
		pmic, err := n.EncodePortManagementCapability(i)
		if err != nil {
			n.log.Errorln("EncodePortManagementCapability", err)
			return err
		}
		tmiReport.PMIC = append(tmiReport.PMIC, pmic)
		tmiReport.PortNum = append(tmiReport.PortNum, i)
	}

	if n.pfcpHandler != nil {
		n.pfcpHandler.NotifySessReport(report.SessReport{
			SEID:    seid, // SEID(Session Endpoint Identifier)
			Reports: []report.Report{tmiReport},
		})
	}
	return nil
}
