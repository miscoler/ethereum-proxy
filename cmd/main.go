package main

import (
	"github.com/miscoler/ethereum-proxy/api/http"
	"github.com/miscoler/ethereum-proxy/internal/application"
	"github.com/miscoler/ethereum-proxy/internal/blocktsx"
	blocktsxconf "github.com/miscoler/ethereum-proxy/internal/blocktsx/config"

	"github.com/miscoler/ethereum-proxy/pkg/pprof"
	"log"

	"go.uber.org/zap"
)

func main() {
	app, err := application.NewApplication("eth_proxy")
	if err != nil {
		log.Fatalf("can not initialize application: %v", err)
	}

	pprof.New(pprof.NewConfig())

	btsxConf := blocktsxconf.NewConfig()
	btsx, err := blocktsx.New(btsxConf, app.Stats)
	if err != nil {
		app.Logger.Fatal("can not create block provider", zap.Error(err))
	}

	conf := http.NewConfig("api")
	api, err := http.New(conf, app, btsx)
	if err != nil {
		app.Logger.Fatal("can not start http api", zap.Error(err))
	}

	app.Logger.Info("stating to serve", zap.String("address", conf.Addr))
	if err := api.ListenAndServe(); err != nil {
		app.Logger.Warn("stoped serving API", zap.Error(err))
	}
}
