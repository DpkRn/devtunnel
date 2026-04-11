package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DpkRn/devtunnel/internal/platform/mongo"
	"github.com/DpkRn/devtunnel/internal/protocol"
)

// ServerHomeHandler serves the control-plane root (no tunnel subdomain).
func ServerHomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"service":"devtunnel","role":"control"}` + "\n"))
}

func GetLogsHandler(mongoClient mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tunnelID := r.URL.Query().Get("tunnel_id")

		logs, err := mongoClient.GetLogs(context.Background(), tunnelID, 50)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(logs)
	}
}

func GetLogByIDHandler(mongoClient mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/logs/")

		logData, err := mongoClient.GetLogByID(context.Background(), id)
		if err != nil {
			http.Error(w, "Not found", 404)
			return
		}

		json.NewEncoder(w).Encode(logData)
	}
}

func ReplayHandler(reg *Registry, mongoClient mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := strings.TrimPrefix(r.URL.Path, "/replay/")

		logData, err := mongoClient.GetLogByID(context.Background(), id)
		if err != nil {
			http.Error(w, "Not found", 404)
			return
		}

		tunnelID := logData["tunnel_id"].(string)

		session, ok := reg.Get(tunnelID)
		if !ok {
			http.Error(w, "Tunnel not active", 400)
			return
		}

		reqMap := logData["request"].(map[string]any)

		bodyBase64 := reqMap["body"].(string)
		body, _ := base64.StdEncoding.DecodeString(bodyBase64)

		req := protocol.TunnelRequest{
			Method:  reqMap["method"].(string),
			Path:    reqMap["path"].(string),
			Headers: map[string][]string{}, // optional decode
			Body:    body,
		}

		stream, err := session.OpenStream()
		if err != nil {
			http.Error(w, "Tunnel stream error", http.StatusBadGateway)
			return
		}
		defer stream.Close()

		data, _ := json.Marshal(req)
		stream.Write(append(data, '\n'))

		reader := bufio.NewReader(stream)
		respBytes, _ := reader.ReadBytes('\n')

		var resp protocol.TunnelResponse
		json.Unmarshal(respBytes, &resp)

		w.WriteHeader(resp.Status)
		w.Write(resp.Body)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"status":"ok"}` + "\n"))
}
