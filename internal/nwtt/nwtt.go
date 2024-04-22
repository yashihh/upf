package nwtt

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/pkg/factory"
	logger_util "github.com/free5gc/util/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

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
	go n.BC()
	return nil
}

func (n *NWTTServer) BC() {
	n.log.Info("5GS as BC over UDP/IPv4")

	cmd := exec.Command("sudo", "ptp4l", "-i", "enp0s9", "-SmE4")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		n.log.Error(err)
	}

}

func (n *NWTTServer) CreatePortCapability() error {
	return nil
}
