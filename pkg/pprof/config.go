package pprof

import "github.com/spf13/viper"

type Config struct {
	Enabled bool
	Addr    string
}

func NewConfig() Config {
	viper.SetDefault("pprof.enabled", true)
	viper.SetDefault("pprof.addr", ":6060")

	return Config{
		Enabled: viper.GetBool("pprof.enabled"),
		Addr:    viper.GetString("pprof.addr"),
	}
}
