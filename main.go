package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/netip"
	"slices"
	"strings"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
)

//go:embed .version
var version string

const loopDelay = 5 * time.Second

type DeviceInfo struct {
	IP        netip.Addr
	Hostname  string
	BytesSent int64
	BytesRecv int64
	SeenAt    time.Duration
	IsOnline  bool
}

type RenderConfig struct {
	InterfaceName  string
	StatusFilter   string
	SortBy         string
	MaxHostnameLen int
}

func main() {
	validSortOptions := []string{"ip", "name", "time", "sent", "recv"}
	validStatusOptions := []string{"online", "offline", "all"}
	sortOptionsStr := strings.Join(validSortOptions, ", ")
	statusOptionsStr := strings.Join(validStatusOptions, ", ")
	sortHelp := "Sort devices by: " + sortOptionsStr
	statusHelp := "Filter devices by status: " + statusOptionsStr

	sortBy := flag.String("sort", "ip", sortHelp)
	statusFilter := flag.String("status", "online", statusHelp)
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("wgtop version %s\n", version)
		return
	}

	if !slices.Contains(validSortOptions, *sortBy) {
		log.Fatalf(
			"Invalid sort option '%s'. Allowed: %s",
			*sortBy,
			sortOptionsStr,
		)
	}
	if !slices.Contains(validStatusOptions, *statusFilter) {
		log.Fatalf(
			"Invalid status option '%s'. Allowed: %s",
			*statusFilter,
			statusOptionsStr,
		)
	}

	client, err := wgctrl.New()
	if err != nil {
		log.Fatalf("Failed to initialize WireGuard client: %v", err)
	}
	defer client.Close()

	fmt.Print("\033[2J")

	for {
		wgDevices, err := client.Devices()
		if err != nil {
			log.Fatalf("Failed to list devices: %v", err)
		}

		fmt.Print("\033[H")

		for _, dev := range wgDevices {
			processed, maxLen := fetchDevices(dev.Peers, *statusFilter)
			sortDevices(processed, *sortBy)

			cfg := RenderConfig{
				InterfaceName:  dev.Name,
				StatusFilter:   *statusFilter,
				SortBy:         *sortBy,
				MaxHostnameLen: maxLen,
			}

			renderOutput(cfg, processed)
		}

		time.Sleep(loopDelay)
	}
}

func renderOutput(cfg RenderConfig, devices []DeviceInfo) {
	fmt.Printf(
		"--- Devices on Interface [%s] (Filter: %s | Sort: %s) --- %s\033[K\n",
		cfg.InterfaceName,
		cfg.StatusFilter,
		cfg.SortBy,
		version,
	)

	if len(devices) == 0 {
		fmt.Println("No devices found matching the criteria.\033[K")
		fmt.Print("\033[K\n")
		return
	}

	for _, d := range devices {
		statusIcon := "💻"
		if !d.IsOnline {
			statusIcon = "❌"
		}

		timeStr := formatDuration(d.SeenAt)
		sentStr := formatTraffic(d.BytesSent)
		recvStr := formatTraffic(d.BytesRecv)

		formatStr := fmt.Sprintf(
			"%s %-13s | 🏷️ %%-%ds | 📤 %%-8s | 📥 %%-8s | 🕒 %%s\033[K\n",
			statusIcon,
			d.IP.String(),
			cfg.MaxHostnameLen,
		)
		fmt.Printf(formatStr, d.Hostname, sentStr, recvStr, timeStr)
	}

	fmt.Print("\033[K\n")
}
