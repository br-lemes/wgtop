package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

const lookupTimeout = 500 * time.Millisecond

func lookupHostname(ipStr string) string {
	resolver := &net.Resolver{PreferGo: true}

	ctx, cancel := context.WithTimeout(context.Background(), lookupTimeout)
	defer cancel()

	names, err := resolver.LookupAddr(ctx, ipStr)
	if err != nil || len(names) == 0 {
		return "unknown"
	}

	name := names[0]
	if len(name) > 0 && name[len(name)-1] == '.' {
		name = name[:len(name)-1]
	}

	return name
}

func formatTraffic(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	value := float64(bytes) / float64(div)
	suffix := []string{"KB", "MB", "GB", "TB"}[exp]

	return fmt.Sprintf("%.1f %s", value, suffix)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "never"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd %dh ago", days, hours)
		}
		return fmt.Sprintf("%dd ago", days)
	}

	if d.Hours() >= 1 {
		return fmt.Sprintf("%dh %dm ago", hours, minutes)
	}

	return fmt.Sprintf("%v ago", d.Round(time.Second))
}
