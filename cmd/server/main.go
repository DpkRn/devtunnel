package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	appconfig "github.com/DpkRn/devtunnel/internal/platform/config"
	"github.com/DpkRn/devtunnel/internal/platform/mongo"
	"github.com/DpkRn/devtunnel/internal/server"
)

func main() {
	reg := server.NewRegistry()
	cfg := appconfig.NewConfig()

	defer func() {
		if r := recover(); r != nil {
			log.Println("panic:", r)
		}
	}()

	mongoClient, err := mongo.NewMongoClient(cfg.MongoDB())
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	tcpCfg := cfg.TCPServer()
	go server.StartTCP(
		reg,
		tcpCfg,
		mongoClient,
	)

	server.SetupRoutes(reg, mongoClient, tcpCfg)
	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatalf("Failed to listen on port 3000: %v", err)
		}
	}()
	fmt.Println("✅HTTP Server Listening on port 3000")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
