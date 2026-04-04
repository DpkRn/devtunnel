package server

import (
	"fmt"
	"log"
	"net"

	"github.com/DpkRn/devtunnel/internal/config"
	"github.com/DpkRn/devtunnel/internal/pkg"
	"github.com/hashicorp/yamux"
)

func StartTCP(reg *Registry) {
	config := config.NewConfig()
	tcpPort := config.ControlTCPListenAddr
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", tcpPort, err)
	}
	fmt.Println("✅TCP Connection Listening on port: ", tcpPort)

	for {
		conn, _ := listener.Accept()
		fmt.Println("Client connected:", conn.RemoteAddr())
		go handleClient(conn, reg, config)
	}
}

func handleClient(conn net.Conn, reg *Registry, config *config.Config) {
	subdomain := pkg.GenerateID()
	publicUrl := subdomain + "." + config.PublicHostSuffix + "\n"

	// ✅ send BEFORE yamux
	conn.Write([]byte(publicUrl))

	// now start yamux
	session, err := yamux.Server(conn, nil)
	if err != nil {
		conn.Close()
		return
	}

	reg.Add(subdomain, session)

	fmt.Println("Client connected:", subdomain)
}
