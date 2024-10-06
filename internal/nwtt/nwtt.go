package nwtt

import (
	"context"
	"encoding/binary"
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
	PTPInstanceID   uint16
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
		cfg:        cfg,
		BridgeAddr: upNodeID,
		// TODO: extend for more ptp instance
		PTPInstanceID:   0,
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
			PortDS_PortIdentity:            SUPPORT,
			PortDS_PortState:               SUPPORT,
			PortDS_LogMinDelayReqInterval:  SUPPORT,
			PortDS_LogAnnounceInterval:     SUPPORT,
			PortDS_AnnounceReceiptTimeout:  SUPPORT,
			PortDS_LogSyncInterval:         SUPPORT,
			PortDS_DelayMechanism:          SUPPORT,
			PortDS_LogMinPdelayReqInterval: SUPPORT,
			PortDS_VersionNumber:           SUPPORT,
			PortDS_MinorVersionNumber:      SUPPORT,
			PortDS_DelayAsymmetry:          SUPPORT,
			PortDS_PortEnable:              SUPPORT,
		}
	}
	return nil
}

/* 24.539 9.5C User Plane Node management capability */
func (n *NWTTServer) CreateUserPlaneNodeCapability() error {
	n.UserPlaneNodeCapabilityList = map[uint16]bool{
		/* Information for 5GS Bridge(Read only) */
		UserPlaneNodeAddress: SUPPORT,
		UserPlaneNodeID:      SUPPORT,
		NWTTPortNumbers:      UNSUPPORT,
		/* Time synchronization information(Read only) */
		SupportedPTPInstanceTypes:           SUPPORT,
		SupportedTransportTypes:             SUPPORT,
		SupportedDelayMechanisms:            SUPPORT,
		PTPGrandmasterCapable:               SUPPORT,
		gPTPGrandmasterCapable:              UNSUPPORT,
		SupportedPTPProfiles:                SUPPORT,
		NumberOfSupportedPTPInstances:       SUPPORT,
		DSTTPortTimeSynchronizationInfoList: SUPPORT,
		PTPInstanceSpecification:            SUPPORT,
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

//TODO: handle NWTT port init capability notify
//TODO: handle UMIC init notify
// add ptpinstanceID, upnode id (for pfcp node id)

func (n *NWTTServer) ReportTSCmanagemantInformation(seid uint64) error {
	n.log.Infoln("Report TSC managemant Information")
	umic := []byte{}
	// umic = append(umic, UserPlaneNodeManagementNotify)
	// for parameter, _ := range n.UserPlaneNodeCapabilityList {

	// }
	// if err != nil {
	// 	n.log.Errorln("EncodeUserPlaneNodeManagementCapability", err)
	// 	return err
	// }

	tmiReport := report.TMIReport{
		UMIC: umic,
	}

	for _, i := range n.ListOfNWTTPorts {
		pmic := []byte{}
		pmic = append(pmic, PortManagementCapability)
		length := make([]byte, 2)
		capability, err := n.EncodePortManagementCapability(i)
		if err != nil {
			n.log.Errorln("EncodePortManagementCapability", err)
			return err
		}
		binary.BigEndian.PutUint16(length, uint16(len(capability)))
		pmic = append(pmic, length...)
		pmic = append(pmic, capability...)

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

func (n *NWTTServer) HandleTSCManagementInformation(TSCMInfoIEs []*ie.IE) (*ie.IE, error) {
	n.log.Infoln("HandleTSCManagementInformation")
	var PMIC, UMIC []byte
	var NWTTPortNumber uint32
	var err error
	for _, i := range TSCMInfoIEs {
		switch i.Type {
		case ie.PortManagementInformationContainer:
			PMIC = i.Payload
			if PMIC == nil {
				n.log.Errorln("pmic error")
				return nil, err
			}
		case ie.BridgeManagementInformationContainer:
			UMIC = i.Payload
			if UMIC == nil {
				n.log.Errorln("pmic error")
				return nil, err
			}
		case ie.NWTTPortNumber:
			NWTTPortNumber, err = i.NWTTPortNumber()
			if err != nil {
				n.log.Errorln(err)
				return nil, err
			}
		default:
			n.log.Errorln("Wrong IE type in SessionModificationResquest TSCManagementInformation")
		}
	}
	var PMICRsp *ie.IE
	var UMICRsp *ie.IE
	var NWTTPortNumberRsp *ie.IE
	var err2 error

	if PMIC != nil && NWTTPortNumber != 0 {
		PMICRsp, err2 = n.DecodePortManagementInformation(PMIC, NWTTPortNumber)
		if err2 != nil {
			n.log.Errorln(err2)
			return nil, err2
		}
		NWTTPortNumberRsp = ie.NewNWTTPortNumber(NWTTPortNumber)
	}
	if UMIC != nil {
		UMICRsp, err2 = n.DecodeUserPlaneNodeManagementInformation(UMIC)
		if err2 != nil {
			n.log.Errorln(err2)
			return nil, err2
		}
	}

	n.log.Infof("NWTT get PMIC:[%x] ,BMIC:[%x] ,NWTTPort:[%d]\n", PMIC, UMIC, NWTTPortNumber)
	return ie.NewTSCManagementInformationWithinSessionModificationResponse(PMICRsp, UMICRsp, NWTTPortNumberRsp), nil
}
