package server

import (
	"net/http"
	"strings"

	"github.com/DpkRn/devtunnel/internal/platform/config"
	"github.com/DpkRn/devtunnel/internal/platform/mongo"
)

// SetupRoutes registers a single handler that routes by Host:
//   - Subdomain of PUBLIC_HOST_SUFFIX → tunnel (HTTP through yamux)
//   - Apex / www / localhost → control plane (/, /logs, /replay/…)
func SetupRoutes(reg *Registry, mongoClient mongo.Client, tcpCfg config.TCPCfg) {
	http.Handle("/", EdgeHandler(reg, mongoClient, tcpCfg))
}

// EdgeHandler returns the root HTTP handler (host-based dispatch).
func EdgeHandler(reg *Registry, mongoClient mongo.Client, tcpCfg config.TCPCfg) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suffix := tcpCfg.PublicHostSuffixFunc()
		if tunnelID, ok := TunnelIDFromHost(r.Host, suffix); ok {
			HandleTunnelRequest(w, r, reg, mongoClient, tunnelID)
			return
		}
		serveControlPlane(w, r, reg, mongoClient)
	})
}

func serveControlPlane(w http.ResponseWriter, r *http.Request, reg *Registry, mongoClient mongo.Client) {
	path := r.URL.Path
	switch {
	case path == "/health":
		HealthHandler(w, r)
	case path == "/logs":
		GetLogsHandler(mongoClient).ServeHTTP(w, r)
	case strings.HasPrefix(path, "/logs/"):
		GetLogByIDHandler(mongoClient).ServeHTTP(w, r)
	case strings.HasPrefix(path, "/replay/"):
		ReplayHandler(reg, mongoClient).ServeHTTP(w, r)
	case path == "/":
		ServerHomeHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}
