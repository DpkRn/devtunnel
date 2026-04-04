package client

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/hashicorp/yamux"
)

// tunnelControlAddr returns host:port for the tunnel control plane.
// Override for local dev: DEVTUNNEL_SERVER=localhost:9000
func tunnelControlAddr() string {
	if s := strings.TrimSpace(os.Getenv("DEVTUNNEL_SERVER")); s != "" {
		return s
	}
	return "clickly.cv:9000"
}

func Start(port string) {
	conn, _ := net.Dial("tcp", tunnelControlAddr())

	session, _ := yamux.Client(conn, nil)

	go acceptStreams(session, port)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	line := strings.TrimSpace(string(buf[:n]))
	publicURL := line
	if !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
		publicURL = "http://" + line
	}
	localURL := "http://localhost:" + port

	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════════════╗")
	fmt.Println("  ║   🚇  mytunnel — tunnel is live                  ║")
	fmt.Println("  ╠══════════════════════════════════════════════════╣")
	fmt.Printf("  ║  🌍  Public   →  %-32s║\n", publicURL)
	fmt.Printf("  ║  💻  Local    →  %-32s║\n", localURL)
	fmt.Println("  ╠══════════════════════════════════════════════════╣")
	fmt.Println("  ║  ⚡  Forwarding requests...                      ║")
	fmt.Println("  ║  🛑  Press Ctrl+C to stop                        ║")
	fmt.Println("  ╚══════════════════════════════════════════════════╝")
	fmt.Println()
}
