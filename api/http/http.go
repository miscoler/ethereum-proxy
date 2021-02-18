package http

import (
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/miscoler/ethereum-proxy/internal/application"
	"github.com/miscoler/ethereum-proxy/internal/blocktsx"
	"github.com/miscoler/ethereum-proxy/pkg/graceful"
	"github.com/miscoler/ethereum-proxy/pkg/logger"
	"github.com/miscoler/ethereum-proxy/pkg/stats"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type errInternal struct {
	error
}

type errRequest struct {
	error
}

func (h *API) wrapLogsAndStats(
	fn func(
		lg logger.Logger,
		w http.ResponseWriter,
		r *http.Request,
	) ([]byte, error),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		defer h.requestTiming.Start().Stop()

		reqID := atomic.AddUint64(&h.reqID, 1)
		lg := h.app.Logger.With(
			zap.String("component", "http"),
			zap.Uint64("request_id", reqID),
			zap.String("block_number", vars["block_number"]),
			zap.String("transaction", vars["transaction"]),
		)

		result, err := fn(lg, w, r)

		if errors.As(err, &errInternal{}) {
			lg.Error("internal error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			h.requestCount500.Observe(1)

			return
		}

		if errors.As(err, &errRequest{}) || err != nil {
			lg.Warn("bad request error", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			h.requestCount400.Observe(1)

			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(result); err != nil {
			lg.Warn("can not write response", zap.Error(err))
		}

		lg.Info("served request")
		h.requestCount200.Observe(1)
	}
}

func (h *API) handleTransaction(
	lg logger.Logger,
	w http.ResponseWriter,
	r *http.Request,
) ([]byte, error) {
	vars := mux.Vars(r)

	var (
		blockNumber int64
		isLatest    bool
	)

	if vars["block_number"] == "latest" {
		isLatest = true
	} else {
		var err error
		blockNumber, err = strconv.ParseInt(vars["block_number"], 0, 64)
		if err != nil {
			return nil, errRequest{
				error: errors.Wrap(err, "parsing block number"),
			}
		}
	}

	block, err := h.blocktsx.GetBlock(
		application.UpgradeContext(r.Context(), lg, h.app),
		blockNumber,
		isLatest,
	)
	if err != nil {
		return nil, errInternal{
			error: errors.Wrap(err, "getting block"),
		}
	}

	const (
		hashLength = 66
	)
	if len(vars["transaction"]) == hashLength {
		return h.blocktsx.GetTSXbyHash(block, vars["transaction"])
	}

	tsxIdx, err := strconv.ParseInt(vars["transaction"], 0, 32)
	if err != nil {
		return nil, errRequest{
			error: errors.Wrap(err, "parsing transaction number"),
		}
	}

	return h.blocktsx.GetTSXbyIndex(block, int(tsxIdx))
}

type API struct {
	Config
	requestTiming   stats.Timing
	requestCount200 stats.Metric
	requestCount400 stats.Metric
	requestCount500 stats.Metric
	app             *application.Application
	blocktsx        blocktsx.BlockProvider

	handler http.Handler
	reqID   uint64
}

func New(
	conf Config,
	app *application.Application,
	btsx blocktsx.BlockProvider,
) (*API, error) {
	h := &API{
		Config: conf,
		app:    app,
		requestTiming: app.Stats.NewTiming(
			"http_request_timing",
			"full http request serving time",
		),
		requestCount200: app.Stats.NewMetric(
			"http_request_200_cnt",
			"served 20x request counter",
		),
		requestCount400: app.Stats.NewMetric(
			"http_request_400_cnt",
			"served 40x request counter",
		),
		requestCount500: app.Stats.NewMetric(
			"http_request_500_cnt",
			"served 50x request counter",
		),
		blocktsx: btsx,
	}

	r := mux.NewRouter()
	r.HandleFunc(
		"/block/{block_number}/txs/{transaction}",
		h.wrapLogsAndStats(h.handleTransaction),
	).Methods(http.MethodGet)
	h.handler = r

	return h, nil
}

func (h *API) ListenAndServe() error {
	return graceful.ListenAndServe(h.Addr, h.handler)
}
