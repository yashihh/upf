package nwtt

import (
	"bytes"
	"testing"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/pkg/factory"
)

func TestUMIC(t *testing.T) {
	cfg := &factory.Config{
		NWTT: &factory.NWTT{
			UpNodeID:  "02:00:5e:10:00:00:00:01",
			NwttPorts: "2,3",
			DsttPorts: "4",
		},
	}
	n, err := NewNWTTServer(
		cfg,
		forwarder.Empty{},
	)
	n.Init()
	t.Run("Create NWTT Server", func(t *testing.T) {

		if n == nil {
			t.Errorf("NewNWTTServer Create failed; got %v\n", err)
		}
	})
	t.Run("Encode User Plane Node Management Capability", func(t *testing.T) {
		truth := []byte{0, 0x01, 0, 0x03, 0, 116, 0, 117, 0, 118, 0, 119, 0, 121, 0, 122, 0, 123, 0, 124}
		tmp, err := n.EncodeUserPlaneNodeManagementCapability()
		if err != nil {
			t.Errorf("Encode User Plane Node Management Capability failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode User Plane Node Management Capability wrong, want %v but got %v", truth, tmp)
		}
	})
	t.Run("Encode UserPlaneNode Status with Supported PTP instance types", func(t *testing.T) {
		truth := []byte{0, 0x74, 0, 2, BoundaryClock, E2ETransparentClock}
		tmp, err := n.EncodeUserPlaneNodeStatus(SupportedPTPInstanceTypes)
		if err != nil {
			t.Errorf("Encode Supported PTP instance types failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported PTP instance types wrong, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Encode UserPlaneNode Status with Supported transport types", func(t *testing.T) {
		truth := []byte{0, 0x75, 0, 1, IPv4}
		tmp, err := n.EncodeUserPlaneNodeStatus(SupportedTransportTypes)
		if err != nil {
			t.Errorf("Encode Supported transport types failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported transport types, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Encode UserPlaneNode Status with Supported Delay Mechanisms", func(t *testing.T) {
		truth := []byte{0, 0x76, 0, 1, E2E}
		tmp, err := n.EncodeUserPlaneNodeStatus(SupportedDelayMechanisms)
		if err != nil {
			t.Errorf("Encode Supported Delay Mechanisms failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported Delay Mechanisms, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Encode UserPlaneNode Status with PTP Grandmaster Capable", func(t *testing.T) {
		truth := []byte{0, 0x77, 0, 1, TRUE}
		tmp, err := n.EncodeUserPlaneNodeStatus(PTPGrandmasterCapable)
		if err != nil {
			t.Errorf("Encode Supported PTP Grandmaster Capable failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported PTP Grandmaster Capable, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Encode UserPlaneNode Status with Supported PTP Profiles", func(t *testing.T) {
		truth := []byte{0, 0x79, 0, 1, E2EDefault}
		tmp, err := n.EncodeUserPlaneNodeStatus(SupportedPTPProfiles)
		if err != nil {
			t.Errorf("Encode Supported PTP Profiles failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported PTP Profiles, want %v but got %v", truth, tmp)
		}
	})

	// t.Run("Handle Manage Port Command GetCapabilities and ReadParameter ", func(t *testing.T) {
	// 	truth := []byte{70, 0, 10, 0, 0xE2, 0, 0xE3, 0, 0xE4, 0, 0xE5, 0, 0xE7, 71, 0, 11, 2, 0, 227, 0, 1, 0, 0, 228, 0, 1, 1}
	// 	tmp, err := n.HandleManagePortCommand([]byte{0, 0, 7, 1, 2, 0, 0xE3, 2, 0, 0xE4}, 2)
	// 	if err != nil {
	// 		t.Errorf("HandleManagePortCommand failed; got %v\n", err)
	// 	}
	// 	if !bytes.Equal(truth, tmp) {
	// 		t.Errorf("Build PortStatus wrong, want %v but got %v", truth, tmp)
	// 	}
	// })

	t.Run("Handle Manage UserPlaneNode Command SetParameter1 ", func(t *testing.T) {
		truth := []byte{72, 0, 13, 1, 0, 124, 9, 0, 5, 48, 57, 0, 2, 0, 1, 0}
		tmp, err := n.HandleManageUserPlaneNodeCommand([]byte{01, 00, 0x0e, 03, 00, 0x7c, 00, 0x09, 00, 05, 0x30, 0x39, 00, 02, 00, 01, 00})
		if err != nil {
			t.Errorf("HandleManagePortCommand failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Build PortStatus wrong, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Handle Manage UserPlaneNode Command SetParameter2 ", func(t *testing.T) {
		truth := []byte{72, 0, 22, 1, 0, 123, 18, 0, 14, 0, 4, 0, 10, 48, 57, 0, 4, 0, 1, 1, 0, 5, 0, 1, 1}
		tmp, err := n.HandleManageUserPlaneNodeCommand([]byte{01, 00, 0x17, 03, 00, 0x7B, 00, 0x12, 00, 0x0E, 00, 0x04, 00, 0x0A, 0x30, 0x39, 00, 0x04, 00, 0x01, 01, 00, 0x05, 00, 0x01, 01})
		if err != nil {
			t.Errorf("HandleManagePortCommand failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Build PortStatus wrong, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Handle Manage UserPlaneNode Command SetParameter of two parameter ", func(t *testing.T) {
		truth := []byte{72, 0, 34, 2, 0, 124, 9, 0, 5, 48, 57, 0, 2, 0, 1, 0,
			0, 123, 18, 0, 14, 0, 4, 0, 10, 48, 57, 0, 4, 0, 1, 1, 0, 5, 0, 1, 1}
		tmp, err := n.HandleManageUserPlaneNodeCommand([]byte{01, 00, 0x25, 03, 00, 0x7c, 00, 0x09, 00, 0x05, 0x30, 0x39, 00, 02, 00, 01, 00,
			03, 00, 0x7B, 00, 0x12, 00, 0x0E, 00, 0x04, 00, 0x0A, 0x30, 0x39, 00, 0x04, 00, 0x01, 01, 00, 0x05, 00, 0x01, 01})
		if err != nil {
			t.Errorf("HandleManagePortCommand failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Build PortStatus wrong, want %v but got %v", truth, tmp)
		}
	})
}
