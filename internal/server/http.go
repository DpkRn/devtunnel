package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/DpkRn/devtunnel/internal/platform/mongo"
	"github.com/DpkRn/devtunnel/internal/protocol"
)

func Handler(reg *Registry, mongoClient mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}
		// Do not put r.Body in the document: io.ReadCloser is not BSON-encodable and InsertOne would fail.
		if _, err := mongoClient.InsertRequestLog(context.Background(), map[string]any{
			"client_ip":       r.RemoteAddr,
			"request_time":    time.Now(),
			"request_type":    "http",
			"request_uri":     r.RequestURI,
			"request_method":  r.Method,
			"request_headers": r.Header.Clone(),
			"content_length":  r.ContentLength,
		}); err != nil {
			log.Printf("mongodb request log: %v", err)
		}

		fmt.Println("r.Host:", r.Host)

		parts := strings.Split(host, ".")
		if len(parts) < 2 {
			http.Error(w, "Invalid host", http.StatusBadRequest)
			return
		}

		subdomain := parts[0]

		session, ok := reg.Get(subdomain)
		if !ok {
			http.Error(w, "Tunnel not found", http.StatusNotFound)
			return
		}

		stream, err := session.Open()
		if err != nil {
			reg.Remove(subdomain)
			http.Error(w, "Tunnel session closed", http.StatusBadGateway)
			return
		}
		defer stream.Close()

		body, _ := io.ReadAll(r.Body)

		req := protocol.TunnelRequest{
			Method:  r.Method,
			Path:    r.URL.String(),
			Headers: r.Header,
			Body:    body,
		}

		fmt.Println("req:", req)

		data, err := json.Marshal(req)
		if err != nil {
			http.Error(w, "Bad request", http.StatusInternalServerError)
			return
		}
		if _, err := stream.Write(append(data, '\n')); err != nil {
			reg.Remove(subdomain)
			http.Error(w, "Tunnel write failed", http.StatusBadGateway)
			return
		}

		reader := bufio.NewReader(stream)
		respBytes, err := reader.ReadBytes('\n')
		if err != nil || len(respBytes) == 0 {
			reg.Remove(subdomain)
			http.Error(w, "Tunnel closed before response", http.StatusBadGateway)
			return
		}
		fmt.Println("respBytes:", string(respBytes))

		var resp protocol.TunnelResponse
		if err := json.Unmarshal(respBytes, &resp); err != nil {
			http.Error(w, "Invalid tunnel response", http.StatusBadGateway)
			return
		}

		if resp.Status < 100 || resp.Status > 599 {
			http.Error(w, "Bad status from tunnel", http.StatusBadGateway)
			return
		}

		for k, v := range resp.Headers {
			for _, val := range v {
				w.Header().Add(k, val)
			}
		}

		w.WriteHeader(resp.Status)
		_, _ = w.Write(resp.Body)
	}
}
