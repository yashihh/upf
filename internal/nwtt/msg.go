package nwtt

const (
	SUPPORT   bool = true
	UNSUPPORT bool = false
)

// Ethernet port management service message type definitions.
const (
	ManagePortCommand            uint8 = 0b00000001
	ManagePortComplete           uint8 = 0b00000010
	PortManagementNotify         uint8 = 0b00000011
	PortManagementNotifyAck      uint8 = 0b00000100
	PortManagementNotifyComplete uint8 = 0b00000101
	PortManagementCapability     uint8 = 0b00000110
)

// IEI of MANAGE PORT COMPLETE message content
const (
	PortManagementCapabilityIEI uint8 = 70
	PortStatusIEI               uint8 = 71
	PortUpdateResultIEI         uint8 = 72
)

// Operation code
const (
	GetCapabilities                  uint8 = 0b00000001
	ReadParameter                    uint8 = 0b00000010
	SetParameter                     uint8 = 0b00000011
	SubscribeNotifyForParameter      uint8 = 0b00000100
	UnsubscribeForParameter          uint8 = 0b00000101
	SelevtiveReadParameter           uint8 = 0b00000110
	SelevtiveSubscribeForParameter   uint8 = 0b00000111
	SelevtiveUnsubscribeForParameter uint8 = 0b00001000
	DeleteParameterEntry             uint8 = 0b00001001
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
type PortDataSet struct {
	PortIdentity            PortIdentity
	PortState               uint8
	LogMinDelayReqInterval  int8
	PeerMeanPathDelay       int64
	LogAnnounceInterval     int8
	AnnounceReceiptTimeout  uint8
	LogSyncInterval         int8
	DelayMechanism          uint8
	LogMinPdelayReqInterval int8
	VersionNumber           uint8
	DelayAsymmetry          int64
	PortEnable              bool
}

const (
	PortDS_PortIdentity            uint16 = 0x0011
	PortDS_PortState               uint16 = 0x0012
	PortDS_LogMinDelayReqInterval  uint16 = 0x0013
	PortDS_LogAnnounceInterval     uint16 = 0x0014
	PortDS_AnnounceReceiptTimeout  uint16 = 0x0015
	PortDS_LogSyncInterval         uint16 = 0x0016
	PortDS_DelayMechanism          uint16 = 0x0017
	PortDS_LogMinPdelayReqInterval uint16 = 0x0018
	PortDS_VersionNumber           uint16 = 0x0019
	PortDS_MinorVersionNumber      uint16 = 0x001A
	PortDS_DelayAsymmetry          uint16 = 0x001B
	PortDS_PortEnable              uint16 = 0x001C
)

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
