package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host    string
	Timeout time.Duration
	RPS     int
}

func NewConfig(path string) Config {
	path += ".eth_client"
	const (
		host    = ".host"
		timeout = ".timeout"
		rps     = ".rps"
	)

	viper.SetDefault(path+host, "https://cloudflare-eth.com")
	viper.SetDefault(path+timeout, 5*time.Second)
	viper.SetDefault(path+rps, "500")

	return Config{
		Host:    viper.GetString(path + host),
		Timeout: viper.GetDuration(path + timeout),
		RPS:     viper.GetInt(path + rps),
	}
}
