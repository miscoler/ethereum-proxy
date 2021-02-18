package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	UncertainBlockLimit int64
	CacheSize           int
}

func NewConfig() Config {
	path := "block_tsx"
	const (
		uncertainBlockLimit = ".uncertain_block_limit"
		cacheSize           = ".cache_size"
	)

	viper.SetDefault(path+uncertainBlockLimit, "20")
	viper.SetDefault(path+cacheSize, "256")

	return Config{
		UncertainBlockLimit: viper.GetInt64(path + uncertainBlockLimit),
		CacheSize:           viper.GetInt(path + cacheSize),
	}
}
