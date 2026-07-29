package validation

import (
	"fmt"
	"net"
	"slices"
	"strings"
)

var privateIPBlocks []*net.IPNet

func initPrivateBlocks() {
	if privateIPBlocks != nil {
		return
	}
	cidrs := []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"::1/128",
		"fe80::/10",
		"fc00::/7",
	}
	for _, cidr := range cidrs {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPBlocks = append(privateIPBlocks, block)
		}
	}
}

func IsPrivateOrLoopbackIP(ip net.IP) bool {
	initPrivateBlocks()
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func ValidateSafeTarget(targetHost string) error {
	host := strings.TrimSpace(targetHost)
	if host == "" {
		return fmt.Errorf("target host is empty")
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		parsed := net.ParseIP(host)
		if parsed != nil {
			ips = []net.IP{parsed}
		} else {
			return fmt.Errorf("resolve target host %s: %w", host, err)
		}
	}
	if slices.ContainsFunc(ips, IsPrivateOrLoopbackIP) {
		return fmt.Errorf("target host %s resolves to private/loopback IP address", host)
	}
	return nil
}
