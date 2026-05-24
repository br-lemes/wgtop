package main

import (
	"net"
	"net/netip"
	"testing"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func TestFetchDevices(t *testing.T) {
	_, ipValid1, _ := net.ParseCIDR("10.0.0.2/32")
	_, ipValid2, _ := net.ParseCIDR("10.0.0.3/32")
	_, ipValid3, _ := net.ParseCIDR("10.0.0.4/32")
	_, ipIgnored, _ := net.ParseCIDR("192.168.1.0/24")

	peers := []wgtypes.Peer{
		{
			LastHandshakeTime: time.Now().Add(-1 * time.Minute),
			AllowedIPs:        []net.IPNet{*ipValid1},
		},
		{
			LastHandshakeTime: time.Now().Add(-5 * time.Minute),
			AllowedIPs:        []net.IPNet{*ipValid2},
		},
		{
			LastHandshakeTime: time.Time{},
			AllowedIPs:        []net.IPNet{*ipValid3},
		},
		{
			LastHandshakeTime: time.Now().Add(-1 * time.Minute),
			AllowedIPs:        []net.IPNet{*ipIgnored},
		},
	}

	t.Run("Filter All", func(t *testing.T) {
		devices, maxLen := fetchDevices(peers, "all")

		if len(devices) != 3 {
			t.Fatalf("Expected 3 devices, got %d", len(devices))
		}

		for _, d := range devices {
			if d.Hostname != "unknown" {
				t.Errorf(
					"Expected hostname to be 'unknown' for IP %s, got %s",
					d.IP.String(),
					d.Hostname,
				)
			}
		}

		if maxLen != 7 {
			t.Errorf("Expected maxLen to be 7, got %d", maxLen)
		}
	})

	t.Run("Filter Online Only", func(t *testing.T) {
		devices, _ := fetchDevices(peers, "online")

		if len(devices) != 1 {
			t.Fatalf("Expected 1 online device, got %d", len(devices))
		}
		if !devices[0].IsOnline {
			t.Error("Device should be marked as online")
		}
	})

	t.Run("Filter Offline Only", func(t *testing.T) {
		devices, _ := fetchDevices(peers, "offline")

		if len(devices) != 2 {
			t.Fatalf("Expected 2 offline devices, got %d", len(devices))
		}

		for _, d := range devices {
			if d.IsOnline {
				t.Errorf(
					"Device %s should be offline, but was marked online",
					d.IP.String(),
				)
			}
		}
	})
}

func TestSortDevices(t *testing.T) {
	devA := DeviceInfo{
		IP:        netip.MustParseAddr("10.0.0.2"),
		Hostname:  "alpha",
		BytesSent: 100,
		BytesRecv: 500,
		SeenAt:    5 * time.Second,
	}
	devB := DeviceInfo{
		IP:        netip.MustParseAddr("10.0.0.5"),
		Hostname:  "bravo",
		BytesSent: 500,
		BytesRecv: 100,
		SeenAt:    10 * time.Second,
	}
	devC := DeviceInfo{
		IP:        netip.MustParseAddr("10.0.0.100"),
		Hostname:  "charlie",
		BytesSent: 100,
		BytesRecv: 100,
		SeenAt:    0,
	}
	devD := DeviceInfo{
		IP:        netip.MustParseAddr("10.0.0.10"),
		Hostname:  "delta",
		BytesSent: 100,
		BytesRecv: 100,
		SeenAt:    0,
	}
	devE := DeviceInfo{
		IP:        netip.MustParseAddr("10.0.0.3"),
		Hostname:  "echo",
		BytesSent: 100,
		BytesRecv: 100,
		SeenAt:    0,
	}

	t.Run("Sort By Name", func(t *testing.T) {
		devices := []DeviceInfo{devB, devA}
		sortDevices(devices, "name")
		if devices[0].Hostname != "alpha" {
			t.Errorf("Expected 'alpha' first, got %s", devices[0].Hostname)
		}
	})

	t.Run("Sort By IP Numeric Order", func(t *testing.T) {
		devices := []DeviceInfo{devC, devD, devE}
		sortDevices(devices, "ip")

		if devices[0].IP.String() != "10.0.0.3" {
			t.Errorf(
				"Expected '10.0.0.3' first, got %s",
				devices[0].IP.String(),
			)
		}
		if devices[1].IP.String() != "10.0.0.10" {
			t.Errorf(
				"Expected '10.0.0.10' second, got %s",
				devices[1].IP.String(),
			)
		}
		if devices[2].IP.String() != "10.0.0.100" {
			t.Errorf(
				"Expected '10.0.0.100' third, got %s",
				devices[2].IP.String(),
			)
		}
	})

	t.Run("Sort By Time with IP Fallback", func(t *testing.T) {
		devices := []DeviceInfo{devC, devD}
		sortDevices(devices, "time")
		if devices[0].IP.String() != "10.0.0.10" {
			t.Errorf(
				"Expected IP '10.0.0.10' first via fallback, got %s",
				devices[0].IP.String(),
			)
		}

		devicesNormal := []DeviceInfo{devB, devA}
		sortDevices(devicesNormal, "time")
		if devicesNormal[0].Hostname != "alpha" {
			t.Errorf(
				"Expected 'alpha' (5s) before 'bravo' (10s)",
			)
		}
	})

	t.Run("Sort By Sent (Traffic and Fallback)", func(t *testing.T) {
		devicesNormal := []DeviceInfo{devA, devB}
		sortDevices(devicesNormal, "sent")
		if devicesNormal[0].BytesSent != 500 {
			t.Errorf(
				"Expected highest BytesSent (500) first, got %d",
				devicesNormal[0].BytesSent,
			)
		}

		devicesFallback := []DeviceInfo{devC, devD}
		sortDevices(devicesFallback, "sent")
		if devicesFallback[0].IP.String() != "10.0.0.10" {
			t.Errorf(
				"Expected IP '10.0.0.10' first on sent fallback, got %s",
				devicesFallback[0].IP.String(),
			)
		}
	})

	t.Run("Sort By Recv (Traffic and Fallback)", func(t *testing.T) {
		devicesNormal := []DeviceInfo{devB, devA}
		sortDevices(devicesNormal, "recv")
		if devicesNormal[0].BytesRecv != 500 {
			t.Errorf(
				"Expected highest BytesRecv (500) first, got %d",
				devicesNormal[0].BytesRecv,
			)
		}

		devicesFallback := []DeviceInfo{devC, devD}
		sortDevices(devicesFallback, "recv")
		if devicesFallback[0].IP.String() != "10.0.0.10" {
			t.Errorf(
				"Expected IP '10.0.0.10' first on recv fallback, got %s",
				devicesFallback[0].IP.String(),
			)
		}
	})
}
