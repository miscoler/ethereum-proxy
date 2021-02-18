package ethclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/miscoler/ethereum-proxy/pkg/ethclient/config"
	"github.com/miscoler/ethereum-proxy/pkg/stats"

	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
)

//nolint: lll
//go:generate mockgen -destination=../../testutil/gomock/mockethclient/mockethclient.go -package=mockethclient github.com/miscoler/ethereum-proxy/pkg/ethclient EthClient
type EthClient interface {
	LatestBlockNumber(context.Context) (int64, error)
	GetBlock(context.Context, int64) ([]byte, error)
}

type ethClientImpl struct {
	config.Config
	*http.Client
	rl                      ratelimit.Limiter
	blockNumberTiming       stats.Timing
	blockNumberOk           stats.Metric
	blockNumberNetworkError stats.Metric
	blockTiming             stats.Timing
	blockOk                 stats.Metric
	blockNetworkError       stats.Metric
}

func New(conf config.Config, stats stats.Stats) EthClient {
	rl := ratelimit.NewUnlimited()
	if conf.RPS != 0 {
		rl = ratelimit.New(conf.RPS)
	}

	return &ethClientImpl{
		Config: conf,
		Client: &http.Client{ //nolint: exhaustivestruct
			Timeout: conf.Timeout,
		},
		rl: rl,
		blockNumberTiming: stats.NewTiming(
			"block_number_timing",
			"request time of eth_blockNumber",
		),
		blockNumberOk: stats.NewMetric(
			"block_number_ok",
			"request counter of successful eth_blockNumber requests",
		),
		blockNumberNetworkError: stats.NewMetric(
			"block_number_err",
			"request counter of errorness eth_blockNumber requests",
		),
		blockTiming: stats.NewTiming(
			"block_timing",
			"request time of eth_getBlockByNumber",
		),
		blockOk: stats.NewMetric(
			"block_ok",
			"request counter of successful eth_getBlockByNumber requests",
		),
		blockNetworkError: stats.NewMetric(
			"block_err",
			"request counter of errorness eth_getBlockByNumber requests",
		),
	}
}

type blockNumberRequest struct {
	JSONrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`
}

type blockNumberResponse struct {
	Result string `json:"result"`
}

func (e *ethClientImpl) LatestBlockNumber(ctx context.Context) (int64, error) {
	defer e.blockNumberTiming.Start().Stop()
	e.rl.Take()

	params, err := json.Marshal(&blockNumberRequest{
		JSONrpc: "2.0",
		Method:  "eth_blockNumber",
		ID:      1,
	})
	if err != nil {
		return 0, errors.Wrap(err, "marshaling params")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.Host, bytes.NewReader(params))
	if err != nil {
		return 0, errors.Wrap(err, "preparing request")
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := e.Do(req)
	if err != nil {
		e.blockNumberNetworkError.Observe(1)

		return 0, errors.Wrap(err, "making http request")
	}
	defer resp.Body.Close()

	var respData blockNumberResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respData); err != nil {
		return 0, errors.Wrap(err, "decoding http response")
	}

	result, err := strconv.ParseInt(respData.Result, 0, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parsing http response")
	}

	e.blockNumberOk.Observe(1)

	return result, nil
}

type blockRequest struct {
	JSONrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type blockResponse struct {
	Result json.RawMessage `json:"result"`
}

func (e *ethClientImpl) GetBlock(ctx context.Context, blockNumber int64) ([]byte, error) {
	defer e.blockTiming.Start().Stop()
	e.rl.Take()

	params, err := json.Marshal(&blockRequest{
		JSONrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{fmt.Sprintf("0x%x", blockNumber), true},
		ID:      1,
	})
	if err != nil {
		return nil, errors.Wrap(err, "marshaling params")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.Host, bytes.NewReader(params))
	if err != nil {
		return nil, errors.Wrap(err, "preparing request")
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := e.Do(req)
	if err != nil {
		e.blockNetworkError.Observe(1)

		return nil, errors.Wrap(err, "making http request")
	}
	defer resp.Body.Close()

	var respData blockResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respData); err != nil {
		return nil, errors.Wrap(err, "decoding http response")
	}

	e.blockOk.Observe(1)

	return respData.Result, nil
}
