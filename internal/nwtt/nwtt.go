package nwtt

import (
	"context"
	"fmt"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/pkg/factory"
	logger_util "github.com/free5gc/util/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PortIdentity struct {
	ClockIdentity uint64
	PortNumber    uint16
}

const (
	SYNC                  = 0x0
	DELAY_REQ             = 0x1
	PDELAY_REQ            = 0x2
	PDELAY_RESP           = 0x3
	FOLLOW_UP             = 0x8
	DELAY_RESP            = 0x9
	PDELAY_RESP_FOLLOW_UP = 0xA
	ANNOUNCE              = 0xB
	SIGNALING             = 0xC
	MANAGEMENT            = 0xD
)

type PTPHeader struct {
	SdoIDAndMsgType     uint8 // first 4 bits are SdoId, next 4 bits are msgtype
	Version             uint8 // first 4 bits are minorVersionPTP, next 4 bits are versionPTP
	MessageLength       uint16
	DomainNumber        uint8
	MinorSdoID          uint8
	FlagField           uint16
	CorrectionField     int64 // IntFloat is a float64 stored in int64
	MessageTypeSpecific uint32
	ClockIdentity       uint64
	SourcePortID        uint16
	SequenceID          uint16
	ControlField        uint8 // the use of this field is obsolete according to IEEE, unless it's ipv4
	LogMessageInterval  int8  // specified as a power of two in seconds. The default is 0 (1 second)
}

type TimeStamp struct {
	Reserved    uint16
	Seconds     uint32
	Nanoseconds uint32
}

type ClockQuality struct {
	ClockClass              uint8
	ClockAccuracy           uint8
	OffsetScaledLogVariance uint16
}

type TimeSource uint8

const (
	AtomicClock        TimeSource = 0x10
	GNSS               TimeSource = 0x20
	TerrestrialRadio   TimeSource = 0x30
	SerialTimeCode     TimeSource = 0x39
	PTP                TimeSource = 0x40
	NTP                TimeSource = 0x50
	HandSet            TimeSource = 0x60
	Other              TimeSource = 0x90
	InternalOscillator TimeSource = 0xa0
)

type Announce struct {
	header                  PTPHeader
	OriginTimestamp         TimeStamp
	CurrentUTCOffset        int16
	Reserved                uint8
	GrandmasterPriority1    uint8
	GrandmasterClockQuality ClockQuality
	GrandmasterPriority2    uint8
	GrandmasterIdentity     uint64
	StepsRemoved            uint16
	TimeSource              TimeSource
}

type Sync struct {
	header          PTPHeader
	OriginTimestamp TimeStamp
}

type FollowUp struct {
	header                 PTPHeader
	presiseOriginTimestamp TimeStamp
}

type UDPHeader struct {
	Src      uint16
	Dst      uint16
	Len      uint16
	Checksum uint16
}

func (hdr *UDPHeader) Bytes() []byte {
	var b []byte
	b = append(b, byte(hdr.Src>>8), byte(hdr.Src))
	b = append(b, byte(hdr.Dst>>8), byte(hdr.Dst))
	b = append(b, byte(hdr.Len>>8), byte(hdr.Len))
	b = append(b, byte(hdr.Checksum>>8), byte(hdr.Checksum))
	return b
}

type NWTTServer struct {
	cfg   *factory.Config //undone
	laddr string
	// BridgeMacAddress net.HardwareAddr
	// NumOfNWTTPorts   uint32
	// NumOfDSTTPorts   uint32
	// ListOfNWTTPorts  []uint32
	// ListOfDSTTPorts  []uint32
	/* NW-TT -> DS-TT */
	// PortPair map[uint32]uint32
	/* Traffic class -> IngressPort -> EgressPort -> txBridgeDelaly*/
	// BD5GS map[uint8]map[uint32]map[uint32]*BridgeDelay5GS
	/* EgressPort -> DelayValue*/
	// PDelayPerEgressPort map[uint32]uint64
	/* EgressPort -> Traffic class -> Priority*/
	// TrafficPriorityPerPort map[uint32]map[uint8]uint8
	/* PortNum -> Capability -> boolean*/
	// PortCapabilityList   map[uint32]map[uint16]bool
	// BridgeCapabilityList map[uint16]bool
	// LLDPAdminStatus      uint8
	// PortslldpAdminStatus map[uint32]uint8
	// pfcp                 report.NWTThandler
	// ThresholdVal         ThresholdMesurement
	ctx    context.Context
	cancel context.CancelFunc
	log    *logrus.Entry
}

func NewNWTTServer(cfg *factory.Config, driver forwarder.Driver) (*NWTTServer, error) {
	cfgGtpu := cfg.Gtpu
	if cfgGtpu == nil {
		return nil, errors.Errorf("no Gtpu config")
	}

	var gtpuAddr string
	if cfgGtpu.Forwarder == "gtp5g" {
		for _, ifInfo := range cfgGtpu.IfList {
			gtpuAddr = fmt.Sprintf("%s", ifInfo.Addr)
			break
		}
		if gtpuAddr == "" {
			return nil, errors.Errorf("not found GTP address")
		}
	}
	return &NWTTServer{
		cfg:   cfg,
		laddr: gtpuAddr,
		// BridgeMacAddress:       MacAddress,
		// NumOfNWTTPorts:         uint32(len(nwttports)),
		// NumOfDSTTPorts:         uint32(len(dsttports)),
		// ListOfNWTTPorts:        nwttports,
		// ListOfDSTTPorts:        dsttports,
		// PortPair:               make(map[uint32]uint32),
		// BD5GS:                  make(map[uint8]map[uint32]map[uint32]*BridgeDelay5GS),
		// PDelayPerEgressPort:    make(map[uint32]uint64),
		// TrafficPriorityPerPort: make(map[uint32]map[uint8]uint8),
		// PortCapabilityList:     make(map[uint32]map[uint16]bool),
		// BridgeCapabilityList:   make(map[uint16]bool),
		// LLDPAdminStatus:        LLDPrxtx,
		// PortslldpAdminStatus:   make(map[uint32]uint8),
		// pfcp:                   nil,
		// ThresholdVal:           ThresholdMesurement{0, 0, "", 0},
		ctx:    nil,
		cancel: nil,
		log:    logger.NwttLog.WithField(logger_util.FieldCategory, "NWTT"),
	}, nil
}

func (n *NWTTServer) Init() error {

	return nil
}

func (n *NWTTServer) CreatePortCapability() error {
	return nil
}
