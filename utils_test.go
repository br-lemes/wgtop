package main

import (
	"testing"
	"time"
)

func TestLookupHostname(t *testing.T) {
	t.Run("Valid Localhost Resolution", func(t *testing.T) {
		result := lookupHostname("127.0.0.1")

		if result == "unknown" {
			t.Skip("Skipping test: environment does not resolve 127.0.0.1 reversely")
		}

		if len(result) > 0 && result[len(result)-1] == '.' {
			t.Errorf("Expected trailing dot to be removed, but got %q", result)
		}
	})

	t.Run("Unknown IP Resolution", func(t *testing.T) {
		result := lookupHostname("192.0.2.1")
		if result != "unknown" {
			t.Errorf("Expected 'unknown' for unresolvable IP, got %q", result)
		}
	})
}

func TestFormatTraffic(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{5368709120, "5.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatTraffic(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatTraffic(%d) = %q; want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"Never", 0, "never"},
		{"Seconds", 45 * time.Second, "45s ago"},
		{"Minutes", 5 * time.Minute, "5m0s ago"},
		{"Hours and Minutes", 2*time.Hour + 15*time.Minute, "2h 15m ago"},
		{"Days", 48 * time.Hour, "2d ago"},
		{"Days and Hours", 51 * time.Hour, "2d 3h ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf(
					"formatDuration(%v) = %q; want %q",
					tt.duration,
					result,
					tt.expected,
				)
			}
		})
	}
}
