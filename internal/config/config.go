package config

type Config struct {
	ControlTCPListenAddr string
	HTTPListenAddr       string
	// PublicURLScheme is sent to clients for the printed public URL (https when nginx terminates TLS).
	PublicURLScheme string
	// PublicHostSuffix is everything after the first dot (e.g. tunnel.example.com → Host abc.tunnel.example.com).
	PublicHostSuffix string
}

func NewConfig() *Config {
	return &Config{
		ControlTCPListenAddr: ":9000",
		HTTPListenAddr:       ":3000",
		PublicURLScheme:      "https",
		PublicHostSuffix:     "13.233.127.241:3000",
	}
}
