package config

import (
	"flag"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ParseConfig(filename string) error {
	cpath := flag.String("c", "", "config path")
	flag.Parse()

	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(*cpath)
	viper.AutomaticEnv()

	return errors.Wrap(viper.ReadInConfig(), "parsing config")
}
