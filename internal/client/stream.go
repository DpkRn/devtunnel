package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"

	protocol "github.com/DpkRn/devtunnel/internal/protocol"
	"github.com/hashicorp/yamux"
)

func acceptStreams(session *yamux.Session, port string) {
	for {
		stream, err := session.Accept()
		if err != nil {
			return
		}

		go handle(stream, port)
	}
}

func handle(stream net.Conn, port string) {
	defer stream.Close()

	reader := bufio.NewReader(stream)
	data, _ := reader.ReadBytes('\n')

	var req protocol.TunnelRequest
	json.Unmarshal(data, &req)

	httpReq, _ := http.NewRequest(
		req.Method,
		"http://localhost:"+port+req.Path,
		bytes.NewReader(req.Body),
	)

	for k, v := range req.Headers {
		for _, val := range v {
			httpReq.Header.Add(k, val)
		}
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		writeErr(stream, http.StatusBadGateway, []byte("local request failed: "+err.Error()))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeErr(stream, http.StatusBadGateway, []byte("read local response: "+err.Error()))
		return
	}

	response := protocol.TunnelResponse{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Body:    body,
	}

	out, err := json.Marshal(response)
	if err != nil {
		writeErr(stream, http.StatusInternalServerError, []byte("marshal response"))
		return
	}
	_, _ = stream.Write(append(out, '\n'))
}

func writeErr(stream net.Conn, code int, msg []byte) {
	r := protocol.TunnelResponse{Status: code, Body: msg}
	out, _ := json.Marshal(r)
	_, _ = stream.Write(append(out, '\n'))
}
