package nwtt

import "time"

const (
	SUPPORT   bool = true
	UNSUPPORT bool = false
)

// PMIC
const (
	TSNTimeDomainNumber uint16 = 0x00D4
)
const (
	PMIC_SupportedPTPInstanceTypes     uint16 = 0x00E2
	PMIC_SupportedTransportTypes       uint16 = 0x00E3
	PMIC_SupportedDelayMechanisms      uint16 = 0x00E4
	PMIC_PTPGrandmasterCapable         uint16 = 0x00E5
	PMIC_gPTPGrandmasterCapable        uint16 = 0x00E6
	PMIC_SupportedPTPProfiles          uint16 = 0x00E7
	PMIC_NumberOfSupportedPTPInstances uint16 = 0x00E8
	PMIC_PTPInstanceList               uint16 = 0x00E9
)

// Supported PTP Instance Types
const (
	OrdinaryClock       uint8 = 0x00
	BoundaryClock       uint8 = 0x01
	P2PTransparentClock uint8 = 0x02
	E2ETransparentClock uint8 = 0x03
)

// Supported transport types
const (
	IPv4     uint8 = 0b00000000
	IPv6     uint8 = 0b00000001
	Ethernet uint8 = 0b00000010
)

// Supported PTP delay mechanisms
const (
	E2E          uint8 = 0x01
	P2P          uint8 = 0x02
	COMMON_P2P   uint8 = 0x03
	SPECIAL      uint8 = 0x04
	NO_MECHANISM uint8 = 0xFE
)

const (
	FALSE uint8 = 0x00
	TRUE  uint8 = 0x01
)

// Supported PTP profile
const (
	SMPTE               uint8 = 0b00000000
	IEEE8021AS          uint8 = 0b00000001
	E2EDefault          uint8 = 0b00000010 // Default delay request-response profile
	P2PDefault          uint8 = 0b00000011 // Default delay peer-to-peer delay profile
	HighAccuracyDefault uint8 = 0b00000100 // High Accuracy Delay Request-Response Default PTP profile
)

type PortIdentity struct {
	ClockIdentity [8]uint8
	PortNumber    uint16
}

const (
	INITIALIZING uint8 = iota + 1
	FAULTY
	DISABLED
	LISTENING
	PRE_MASTER
	MASTER
	PASSIVE
	UNCALIBRATED
	SLAVE
)

// Supported PTP Instance List
type PortDS struct {
	PortIdentity            PortIdentity
	PortState               uint8
	LogMinDelayReqInterval  int8
	PeerMeanPathDelay       time.Duration
	LogAnnounceInterval     int8
	AnnounceReceiptTimeout  uint8
	LogSyncInterval         int8
	DelayMechanism          uint8
	LogMinPdelayReqInterval int8
	VersionNumber           uint8
	DelayAsymmetry          int64
	PortEnable              bool
}

// UMIC
/* Information for 5GS Bridge(Read only) */
const (
	UMIC_UserPlaneNodeAddress uint16 = 0x0001
	UMIC_UserPlaneNodeID      uint16 = 0x0003
	UMIC_NWTTPortNumbers      uint16 = 0x0004
)

/* Time synchronization information(Read only) */
const (
	UMIC_SupportedPTPInstanceType            uint16 = 0x0074
	UMIC_SupportedTransportTypes             uint16 = 0x0075
	UMIC_SupportedDelayMechanisms            uint16 = 0x0076
	UMIC_PTPGrandmasterCapable               uint16 = 0x0077
	UMIC_gPTPGrandmasterCapable              uint16 = 0x0078
	UMIC_SupportedPTPProfiles                uint16 = 0x0079
	UMIC_NumberOfSupportedPTPInstances       uint16 = 0x007A
	UMIC_DSTTPortTimeSynchronizationInfoList uint16 = 0x007B
	UMIC_PTPInstanceSpecification            uint16 = 0x007C
)
