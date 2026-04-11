package server

import "testing"

func TestTunnelIDFromHost(t *testing.T) {
	cases := []struct {
		host, suffix string
		wantID        string
		wantTunnel    bool
	}{
		{"abc.clickly.cv", "clickly.cv", "abc", true},
		{"clickly.cv", "clickly.cv", "", false},
		{"www.clickly.cv", "clickly.cv", "", false},
		{"localhost:3000", "localhost", "", false},
		{"foo.localhost", "localhost", "foo", true},
		{"abc.clickly.cv:443", "clickly.cv", "abc", true},
	}
	for _, tc := range cases {
		id, ok := TunnelIDFromHost(tc.host, tc.suffix)
		if ok != tc.wantTunnel || id != tc.wantID {
			t.Errorf("TunnelIDFromHost(%q, %q) = (%q, %v); want (%q, %v)",
				tc.host, tc.suffix, id, ok, tc.wantID, tc.wantTunnel)
		}
	}
}
