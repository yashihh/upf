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

// User Plane Node management service message type definitions.
const (
	ManageUserPlaneNodeCommand       uint8 = 0b00000001
	ManageUserPlaneNodeComplete      uint8 = 0b00000010
	UserPlaneNodeManagementNotify    uint8 = 0b00000011
	UserPlaneNodeManagementNotifyAck uint8 = 0b00000100
)

// IEI of MANAGE PORT COMPLETE message content
const (
	PortManagementCapabilityIEI uint8 = 70
	PortStatusIEI               uint8 = 71
	PortUpdateResultIEI         uint8 = 72
)

// IEI of MANAGE User Plane Node COMPLETE message content
const (
	UserPlaneNodeManagementCapability uint8 = 70
	UserPlaneNodeStatus               uint8 = 71
	UserPlaneNodeUpdateResult         uint8 = 72
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

// Supported PTP Instance Types
const (
	OrdinaryClock       uint8 = 0x00
	BoundaryClock       uint8 = 0x01
	P2PTransparentClock uint8 = 0x02
	E2ETransparentClock uint8 = 0x03
	P2PRelayInstance    uint8 = 0x04
)

var ptpInstanceTypesMap = map[uint8]string{
	OrdinaryClock:       "OrdinaryClock",
	BoundaryClock:       "BOUNDARY_CLOCK",
	P2PTransparentClock: "P2P_TRANS_CLOCK",
	E2ETransparentClock: "E2E_TRANS_CLOCK",
	P2PRelayInstance:    "P2P_RELAY_INSTANCE",
}

// Supported transport types
const (
	IPv4     uint8 = 0b00000000
	IPv6     uint8 = 0b00000001
	Ethernet uint8 = 0b00000010
)

var transportTypesMap = map[uint8]string{
	IPv4:     "IPv4",
	IPv6:     "IPv6",
	Ethernet: "Ethernet",
}

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
	UserPlaneNodeAddress uint16 = 0x0001
	UserPlaneNodeID      uint16 = 0x0003
	NWTTPortNumbers      uint16 = 0x0004
)

/* Time synchronization information(Read only) */
const (
	SupportedPTPInstanceTypes           uint16 = 0x0074
	SupportedTransportTypes             uint16 = 0x0075
	SupportedDelayMechanisms            uint16 = 0x0076
	PTPGrandmasterCapable               uint16 = 0x0077
	gPTPGrandmasterCapable              uint16 = 0x0078
	SupportedPTPProfiles                uint16 = 0x0079
	NumberOfSupportedPTPInstances       uint16 = 0x007A
	DSTTPortTimeSynchronizationInfoList uint16 = 0x007B
	PTPInstanceSpecification            uint16 = 0x007C
)

// PTP instance specification
const (
	PTP_profile                                    uint16 = 0x0001
	Transport_type                                 uint16 = 0x0002
	Grandmaster_enabled                            uint16 = 0x0003
	Grandmaster_on_behalf_of_DSTT_enabled          uint16 = 0x0004
	Grandmaster_candidate_enabled                  uint16 = 0x0005
	DefaultDS_clockIdentity                        uint16 = 0x0006
	DefaultDS_clockQuality_clockClass              uint16 = 0x0007
	DefaultDS_clockQuality_clockAccuracy           uint16 = 0x0008
	DefaultDS_clockQuality_offsetScaledLogVariance uint16 = 0x0009
	DefaultDS_priority1                            uint16 = 0x000A
	DefaultDS_priority2                            uint16 = 0x000B
	DefaultDS_domainNumber                         uint16 = 0x000C
	DefaultDS_sdoId                                uint16 = 0x000D
	DefaultDS_instanceEnable                       uint16 = 0x000E
	DefaultDS_externalPortConfigurationEnabled     uint16 = 0x000F
	DefaultDS_instanceType                         uint16 = 0x0010
	PortDS_PortIdentity                            uint16 = 0x0011
	PortDS_PortState                               uint16 = 0x0012
	PortDS_LogMinDelayReqInterval                  uint16 = 0x0013
	PortDS_LogAnnounceInterval                     uint16 = 0x0014
	PortDS_AnnounceReceiptTimeout                  uint16 = 0x0015
	PortDS_LogSyncInterval                         uint16 = 0x0016
	PortDS_DelayMechanism                          uint16 = 0x0017
	PortDS_LogMinPdelayReqInterval                 uint16 = 0x0018
	PortDS_VersionNumber                           uint16 = 0x0019
	PortDS_MinorVersionNumber                      uint16 = 0x001A
	PortDS_DelayAsymmetry                          uint16 = 0x001B
	PortDS_PortEnable                              uint16 = 0x001C
)

type ConfigurationForPTP struct {
	// DefaultDS_clockIdentity
	// DefaultDS_clockQuality_clockClass
	// DefaultDS_clockQuality_clockAccuracy
	// DefaultDS_clockQuality_offsetScaledLogVariance
	// DefaultDS_priority1
	// DefaultDS_priority2
	// DefaultDS_domainNumber
	// DefaultDS_sdoId
	// DefaultDS_instanceEnable
	// DefaultDS_externalPortConfigurationEnabled
	DefaultDS_instanceType uint8
}

var ch = make(chan ConfigurationForPTP)
