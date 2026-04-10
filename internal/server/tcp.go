package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/DpkRn/devtunnel/internal/pkg"
	"github.com/DpkRn/devtunnel/internal/platform/config"
	"github.com/DpkRn/devtunnel/internal/platform/mongo"
	"github.com/hashicorp/yamux"
)

func StartTCP(reg *Registry, tcpConfig config.TCPCfg, mongoClient mongo.Client) {
	// config := config.NewConfig()
	tcpPort := tcpConfig.ListenAddrFunc()
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", tcpPort, err)
	}
	fmt.Println("✅TCP Connection Listening on port: ", tcpPort)

	for {
		conn, _ := listener.Accept()
		fmt.Println("Client connected:", conn.RemoteAddr())

		// connBytes, err := json.Marshal(conn)
		// if err != nil {
		// 	log.Fatalf("Failed to marshal connection: %v", err)
		// }
		mongoClient.InsertTunnelLog(context.Background(), map[string]any{
			"client_ip":       conn.RemoteAddr().String(),
			"client_port":     conn.RemoteAddr().String(),
			"connection_time": time.Now().Format(time.RFC3339),
			"connection_type": "tcp",
			// "conn":            string(connBytes),
		})
		go func() {
			handleClient(conn, reg, tcpConfig)
		}()
	}
}

func handleClient(conn net.Conn, reg *Registry, tcpConfig config.TCPCfg) {
	subdomain := pkg.GenerateID()
	scheme := strings.TrimSuffix(strings.ToLower(tcpConfig.PublicURLSchemeFunc()), "://")
	if scheme == "" {
		scheme = "https"
	}
	publicURL := fmt.Sprintf("%s://%s.%s\n", scheme, subdomain, tcpConfig.PublicHostSuffixFunc())

	// ✅ send BEFORE yamux
	conn.Write([]byte(publicURL))

	// now start yamux
	session, err := yamux.Server(conn, nil)
	if err != nil {
		conn.Close()
		return
	}

	reg.Add(subdomain, session)

	fmt.Println("Client connected:", subdomain)
}
