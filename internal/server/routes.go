package server

import (
	"net/http"

	"github.com/DpkRn/devtunnel/internal/platform/mongo"
)

func SetupRoutes(reg *Registry, mongoClient mongo.Client) {
	http.HandleFunc("/", Handler(reg, mongoClient))

	http.HandleFunc("/logs", GetLogsHandler(mongoClient))
	http.HandleFunc("/logs/", GetLogByIDHandler(mongoClient))
	http.HandleFunc("/replay/", ReplayHandler(reg, mongoClient))
}
