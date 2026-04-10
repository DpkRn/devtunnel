package config

import "os"

type HTTPServerCfg struct {
	ListenAddr string
}

func (c config) HTTPServer() HTTPServerCfg {
	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	return HTTPServerCfg{
		ListenAddr: listenAddr,
	}
}

func (c HTTPServerCfg) ListenAddrFunc() string {
	return c.ListenAddr
}
