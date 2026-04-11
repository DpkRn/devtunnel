package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/DpkRn/devtunnel/internal/platform/config"
	"github.com/hashicorp/yamux"
)

// tunnelControlAddr returns host:port for the tunnel control plane.
// Override for local dev: DEVTUNNEL_SERVER=localhost:9000
func tunnelControlAddr() string {

	config := config.NewConfig()
	switch config.TierFunc() {
	case "prod":
		return "clickly.cv:9000"
	default:
		return "localhost:9000"
	}
}

func Start(port string) string {
	conn, err := net.Dial("tcp", tunnelControlAddr())
	if err != nil {
		fmt.Println("Error connecting to tunnel control plane:", err)
		return ""
	}
	fmt.Println("Connected to tunnel control plane:", conn.RemoteAddr())
	tunnelReq := TunnelRequest{
		TunnelType:   "gotunnel",
		Version:      "1.0.8",
		TunnelID:     "random-tunnel-id", //todo: generate a fixed tunnel for user
		ConnectionID: GenerateConnectionID(),
	}
	tunnelReqBytes, err := json.Marshal(tunnelReq)
	if err != nil {
		fmt.Println("Error marshalling tunnel request:", err)
		return ""
	}
	conn.Write(append(tunnelReqBytes, '\n'))

	session, err := yamux.Client(conn, nil)
	if err != nil {
		fmt.Println("Error creating yamux session:", err)
		return ""
	}

	go acceptStreams(session, port)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return ""
	}
	line := strings.TrimSpace(string(buf[:n]))
	publicURL := line
	if !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
		publicURL = "http://" + line
	}
	localURL := "http://localhost:" + port

	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════════════╗")
	fmt.Println("  ║   🫦  🤢 mytunnel — tunnel is live 🤢 🫦                 ║")
	fmt.Println("  ╠══════════════════════════════════════════════════╣")
	fmt.Printf("  ║  🌍  Public   →  %-32s║\n", publicURL)
	fmt.Printf("  ║  💻  Local    →  %-32s║\n", localURL)
	fmt.Println("  ╠══════════════════════════════════════════════════╣")
	fmt.Println("  ║  ⚡  Forwarding requests...                      ║")
	fmt.Println("  ║  🛑  Press Ctrl+C to stop                        ║")
	fmt.Println("  ╚══════════════════════════════════════════════════╝")
	fmt.Println()
	return publicURL
}
