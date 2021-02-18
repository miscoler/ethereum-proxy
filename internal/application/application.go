package application

import (
	"context"
	"github.com/miscoler/ethereum-proxy/pkg/config"
	"github.com/miscoler/ethereum-proxy/pkg/ethclient"
	ethcfg "github.com/miscoler/ethereum-proxy/pkg/ethclient/config"
	"github.com/miscoler/ethereum-proxy/pkg/logger"
	"github.com/miscoler/ethereum-proxy/pkg/stats"
	statscfg "github.com/miscoler/ethereum-proxy/pkg/stats/config"

	"github.com/benbjohnson/clock"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Application struct {
	Clock     clock.Clock
	EthClient ethclient.EthClient
	Stats     stats.Stats
	Logger    logger.Logger
}

func NewApplication(name string) (*Application, error) {
	lg, err := logger.New()
	if err != nil {
		return nil, errors.Wrap(err, "creating logger")
	}
	if err := config.ParseConfig(name); err != nil {
		return nil, errors.Wrap(err, "parsing config")
	}

	clock := clock.New()

	stats, err := stats.New(statscfg.NewConfig("stats"), clock, lg)
	if err != nil {
		lg.Fatal("can not start metrics server", zap.Error(err))
	}
	defer stats.NewTiming("test", "test help").Start().Stop()

	return &Application{
		Clock:     clock,
		EthClient: ethclient.New(ethcfg.NewConfig("eth_client"), stats),
		Logger:    lg,
		Stats:     stats,
	}, nil
}

type EContext struct {
	context.Context
	*Application
	Logger logger.Logger
}

func UpgradeContext(
	ctx context.Context,
	lg logger.Logger,
	app *Application,
) *EContext {
	return &EContext{
		Context:     ctx,
		Logger:      lg,
		Application: app,
	}
}
