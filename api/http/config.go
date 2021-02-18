package http

import (
	"github.com/spf13/viper"
)

type Config struct {
	Addr string
}

func NewConfig(path string) Config {
	path += ".http"
	const (
		addr = ".addr"
	)

	viper.SetDefault(path+addr, ":8080")

	return Config{
		Addr: viper.GetString(path + addr),
	}
}
