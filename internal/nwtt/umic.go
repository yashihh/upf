package nwtt

import "time"

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
