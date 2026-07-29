package engine

import "testing"

func TestResolveTunnelID(t *testing.T) {
	tests := []struct {
		host     string
		expected string
		valid    bool
	}{
		{host: "abc.tunnel.localhost", expected: "abc", valid: true},
		{host: "abc.tunnel.localhost:443", expected: "abc", valid: true},
		{host: "tunnel.localhost", valid: false},
		{host: "abc.other.localhost", valid: false},
		{host: "www.tunnel.localhost", valid: false},
	}
	for _, test := range tests {
		actual, valid := resolveTunnelID(test.host, "tunnel.localhost")
		if valid != test.valid || actual != test.expected {
			t.Fatalf("resolveTunnelID(%q) = %q, %t", test.host, actual, valid)
		}
	}
}
