package main

import (
	"cmp"
	"net/netip"
	"slices"
	"strings"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const onlineTimeout = 3 * time.Minute

func fetchDevices(peers []wgtypes.Peer, filter string) ([]DeviceInfo, int) {
	var processed []DeviceInfo
	maxHostnameLen := 0

	for _, peer := range peers {
		isOnline := !peer.LastHandshakeTime.IsZero() &&
			time.Since(peer.LastHandshakeTime) <= onlineTimeout

		if filter == "online" && !isOnline {
			continue
		}
		if filter == "offline" && isOnline {
			continue
		}

		for _, ipNet := range peer.AllowedIPs {
			ones, bits := ipNet.Mask.Size()
			if ones != 32 || bits != 32 {
				continue
			}

			ipAddr, err := netip.ParseAddr(ipNet.IP.String())
			if err != nil {
				continue
			}

			hostname := lookupHostname(ipAddr.String())

			if len(hostname) > maxHostnameLen {
				maxHostnameLen = len(hostname)
			}

			var seenAt time.Duration
			if !peer.LastHandshakeTime.IsZero() {
				seenAt = time.Since(peer.LastHandshakeTime)
			}

			processed = append(processed, DeviceInfo{
				IP:        ipAddr,
				Hostname:  hostname,
				SeenAt:    seenAt,
				IsOnline:  isOnline,
				BytesSent: peer.TransmitBytes,
				BytesRecv: peer.ReceiveBytes,
			})
		}
	}

	return processed, maxHostnameLen
}

func sortDevices(devices []DeviceInfo, sortBy string) {
	slices.SortFunc(devices, func(a, b DeviceInfo) int {
		switch sortBy {
		case "name":
			return strings.Compare(a.Hostname, b.Hostname)
		case "time":
			c := cmp.Compare(a.SeenAt, b.SeenAt)
			if c != 0 {
				return c
			}
			return a.IP.Compare(b.IP)
		case "sent":
			c := cmp.Compare(b.BytesSent, a.BytesSent)
			if c != 0 {
				return c
			}
			return a.IP.Compare(b.IP)
		case "recv":
			c := cmp.Compare(b.BytesRecv, a.BytesRecv)
			if c != 0 {
				return c
			}
			return a.IP.Compare(b.IP)
		default: // "ip"
			return a.IP.Compare(b.IP)
		}
	})
}
