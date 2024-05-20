// User Plane Node Capability elements.
package umic

/* Information for 5GS Bridge(Read only) */
const (
	UserPlaneNodeAddress uint16 = 0x0001
	UserPlaneNodeID      uint16 = 0x0003
	NWTTPortNumbers      uint16 = 0x0004
)

/* Time synchronization information(Read only) */
const (
	SupportedPTPInstanceType            uint16 = 0x0074
	SupportedTransportTypes             uint16 = 0x0075
	SupportedDelayMechanisms            uint16 = 0x0076
	PTPGrandmasterCapable               uint16 = 0x0077
	GPTPGrandmasterCapable              uint16 = 0x0078
	SupportedPTPProfiles                uint16 = 0x0079
	NumberOfSupportedPTPInstances       uint16 = 0x007A
	DSTTPortTimeSynchronizationInfoList uint16 = 0x007B
	PTPInstanceSpecification            uint16 = 0x007C
)
