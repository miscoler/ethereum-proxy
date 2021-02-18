package config

import "github.com/spf13/viper"

type Config struct {
	Addr     string
	Endpoint string
}

func NewConfig(path string) Config {
	path += ".prometheus"
	const (
		addr     = ".addr"
		endpoint = ".endpoint"
	)

	viper.SetDefault(path+addr, ":2112")
	viper.SetDefault(path+endpoint, "/metrics")

	return Config{
		Addr:     viper.GetString(path + addr),
		Endpoint: viper.GetString(path + endpoint),
	}
}
