package pprof

import (
	"github.com/miscoler/ethereum-proxy/pkg/graceful"
	_ "net/http/pprof" //nolint:gosec
)

func New(conf Config) {
	if conf.Enabled {
		go func() {
			graceful.ListenAndServe(conf.Addr, nil) //nolint:errcheck
		}()
	}
}
