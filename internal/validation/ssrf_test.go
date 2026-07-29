package validation

import (
	"net"
	"testing"
)

func TestIsPrivateOrLoopbackIP(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"127.0.0.1", true},
		{"10.0.0.5", true},
		{"172.16.0.1", true},
		{"192.168.1.1", true},
		{"169.254.169.254", true},
		{"::1", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
	}

	for _, tt := range tests {
		ip := net.ParseIP(tt.ip)
		if got := IsPrivateOrLoopbackIP(ip); got != tt.expected {
			t.Errorf("IsPrivateOrLoopbackIP(%s) = %v, expected %v", tt.ip, got, tt.expected)
		}
	}
}

func TestValidateSafeTarget(t *testing.T) {
	if err := ValidateSafeTarget("127.0.0.1"); err == nil {
		t.Error("expected error for 127.0.0.1, got nil")
	}
	if err := ValidateSafeTarget("169.254.169.254"); err == nil {
		t.Error("expected error for metadata IP, got nil")
	}
}
