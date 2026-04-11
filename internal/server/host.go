package server

import (
	"net"
	"strings"
)

// stripHostPort returns the host part of r.Host (without port).
func stripHostPort(hostport string) string {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return strings.TrimSpace(hostport)
	}
	return strings.TrimSpace(host)
}

// TunnelIDFromHost reports whether the request should be handled as tunnel traffic.
// host is the full Host header value; publicSuffix is e.g. "clickly.cv" or "localhost" (dev).
// Returns (tunnelID, true) when the first label is the tunnel id (e.g. abc.clickly.cv → abc).
func TunnelIDFromHost(hostport, publicSuffix string) (tunnelID string, isTunnel bool) {
	host := strings.ToLower(stripHostPort(hostport))
	suffix := strings.ToLower(strings.TrimSpace(publicSuffix))
	if suffix == "" {
		// Misconfigured: cannot classify tunnel subdomains.
		return "", false
	}

	// Apex / dashboard on the public domain
	if host == suffix || host == "www."+suffix {
		return "", false
	}

	switch host {
	case "localhost", "127.0.0.1", "::1":
		return "", false
	}

	// Dev: <id>.localhost
	if suffix == "localhost" {
		if !strings.HasSuffix(host, ".localhost") {
			return "", false
		}
		sub := strings.TrimSuffix(host, ".localhost")
		if sub == "" || sub == "www" {
			return "", false
		}
		return sub, true
	}

	// Prod: <id>.<suffix> (single tunnel label before the configured suffix)
	dot := "." + suffix
	if !strings.HasSuffix(host, dot) {
		return "", false
	}
	sub := strings.TrimSuffix(host, dot)
	if sub == "" || sub == "www" {
		return "", false
	}
	return sub, true
}
