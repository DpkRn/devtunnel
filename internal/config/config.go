package config

type Config struct {
	ControlTCPListenAddr string
	HTTPListenAddr       string
	PublicHostSuffix     string
}

func NewConfig() *Config {
	return &Config{
		ControlTCPListenAddr: ":9000",
		HTTPListenAddr:       ":3000",
		PublicHostSuffix:     "13.233.127.241:3000",
	}
}
