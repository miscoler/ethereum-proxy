package pprof

import (
	_ "net/http/pprof" //nolint:gosec

	"github.com/miscoler/ethereum-proxy/pkg/graceful"
)

func New(conf Config) {
	if conf.Enabled {
		go func() {
			graceful.ListenAndServe(conf.Addr, nil) //nolint:errcheck
		}()
	}
}
