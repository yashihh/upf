package nwtt

import (
	"bytes"
	"testing"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/pkg/factory"
)

func TestPMIC(t *testing.T) {
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
	t.Run("Build PortCapability", func(t *testing.T) {
		truth := []byte{0, 0xE2, 0, 0xE3, 0, 0xE4, 0, 0xE5, 0, 0xE7}
		tmp, err := n.EncodePortManagementCapability(2)
		if err != nil {
			t.Errorf("Build PortCapability failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Build PortCapability wrong, want %v but got %v", truth, tmp)
		}
	})
	t.Run("Encode Port Status with Supported PTP instance types", func(t *testing.T) {
		truth := []byte{0, 0xE2, 0, 2, BoundaryClock, E2ETransparentClock}
		tmp, err := n.EncodePortStatus(2, PMIC_SupportedPTPInstanceTypes)
		if err != nil {
			t.Errorf("Encode Supported PTP instance types failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported PTP instance types wrong, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Encode Port Status with Supported transport types", func(t *testing.T) {
		truth := []byte{0, 0xE3, 0, 1, IPv4}
		tmp, err := n.EncodePortStatus(2, PMIC_SupportedTransportTypes)
		if err != nil {
			t.Errorf("Encode Supported transport types failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported transport types, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Encode Port Status with Supported Delay Mechanisms", func(t *testing.T) {
		truth := []byte{0, 0xE4, 0, 1, E2E}
		tmp, err := n.EncodePortStatus(2, PMIC_SupportedDelayMechanisms)
		if err != nil {
			t.Errorf("Encode Supported Delay Mechanisms failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported Delay Mechanisms, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Encode Port Status with PTP Grandmaster Capable", func(t *testing.T) {
		truth := []byte{0, 0xE5, 0, 1, TRUE}
		tmp, err := n.EncodePortStatus(2, PMIC_PTPGrandmasterCapable)
		if err != nil {
			t.Errorf("Encode Supported PTP Grandmaster Capable failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported PTP Grandmaster Capable, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Encode Port Status with Supported PTP Profiles", func(t *testing.T) {
		truth := []byte{0, 0xE7, 0, 1, E2EDefault}
		tmp, err := n.EncodePortStatus(2, PMIC_SupportedPTPProfiles)
		if err != nil {
			t.Errorf("Encode Supported PTP Profiles failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Encode Supported PTP Profiles, want %v but got %v", truth, tmp)
		}
	})

	t.Run("Handle Manage Port Command GetCapabilities and ReadParameter ", func(t *testing.T) {
		truth := []byte{70, 0, 10, 0, 0xE2, 0, 0xE3, 0, 0xE4, 0, 0xE5, 0, 0xE7, 71, 0, 11, 2, 0, 227, 0, 1, 0, 0, 228, 0, 1, 1}
		tmp, err := n.HandleManagePortCommand([]byte{0, 0, 7, 1, 2, 0, 0xE3, 2, 0, 0xE4}, 2)
		if err != nil {
			t.Errorf("HandleManagePortCommand failed; got %v\n", err)
		}
		if !bytes.Equal(truth, tmp) {
			t.Errorf("Build PortStatus wrong, want %v but got %v", truth, tmp)
		}
	})
}
