package nwtt

import (
	"encoding/binary"
	"sort"
	"time"
)

type ClockIdentity [8]byte

type PortIdentity struct {
	ClockIdentity ClockIdentity
	PortNumber    uint16
}

type PortState byte

const (
	INITIALIZING PortState = iota + 1
	FAULTY
	DISABLED
	LISTENING
	PRE_MASTER
	MASTER
	PASSIVE
	UNCALIBRATED
	SLAVE
)

type DelayMechanism byte

const (
	E2E DelayMechanism = iota + 1
	P2P
	COMMON_P2P
	SPECIAL
	NO_MECHANISM DelayMechanism = 0xFE
)

/*
	PORT DATA SET (IEEE 1588)
*/
type PortDS struct {
	PortIdentity            PortIdentity
	PortState               PortState
	LogMinDelayReqInterval  int8
	PeerMeanPathDelay       time.Duration
	LogAnnounceInterval     int8
	AnnounceReceiptTimeout  uint8
	LogSyncInterval         int8
	DelayMechanism          DelayMechanism
	LogMinPdelayReqInterval int8
	VersionNumber           uint8
	DelayAsymmetry          int64
	PortEnable              bool // optional
}

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
