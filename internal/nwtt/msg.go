package nwtt

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
	UMIC_GPTPGrandmasterCapable              uint16 = 0x0078
	UMIC_SupportedPTPProfiles                uint16 = 0x0079
	UMIC_NumberOfSupportedPTPInstances       uint16 = 0x007A
	UMIC_DSTTPortTimeSynchronizationInfoList uint16 = 0x007B
	UMIC_PTPInstanceSpecification            uint16 = 0x007C
)
