// Port Management Capability elements
package pmic

const (
	TSNTimeDomainNumber uint16 = 0x00D4
)

const (
	SupportedPTPInstanceTypes     uint16 = 0x00E2
	SupportedTransportTypes       uint16 = 0x00E3
	SupportedDelayMechanisms      uint16 = 0x00E4
	PTPGrandmasterCapable         uint16 = 0x00E5
	gPTPGrandmasterCapable        uint16 = 0x00E6
	SupportedPTPProfiles          uint16 = 0x00E7
	NumberOfSupportedPTPInstances uint16 = 0x00E8
	PTPInstanceList               uint16 = 0x00E9
)
