package blocktsx

import (
	"encoding/json"

	"github.com/miscoler/ethereum-proxy/internal/application"
	"github.com/miscoler/ethereum-proxy/internal/blocktsx/config"
	"github.com/miscoler/ethereum-proxy/pkg/stats"

	cache "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
)

//nolint: lll
//go:generate mockgen -destination=../../testutil/gomock/mockblocktsx/mockblocktsx.go -package=mockblocktsx github.com/miscoler/ethereum-proxy/internal/blocktsx BlockProvider
type BlockProvider interface {
	GetBlock(ctx *application.EContext, blockNumber int64, isLatest bool) (*BlockStored, error)
	GetTSXbyHash(*BlockStored, string) ([]byte, error)
	GetTSXbyIndex(*BlockStored, int) ([]byte, error)
}

type blockProviderImpl struct {
	config.Config
	latestBlockNumber int64
	cache             *cache.ARCCache
	cacheHit          stats.Metric
	cacheMiss         stats.Metric
}

func New(conf config.Config, stats stats.Stats) (BlockProvider, error) {
	cache, err := cache.NewARC(conf.CacheSize)
	if err != nil {
		return nil, errors.Wrap(err, "initializing cache")
	}

	return &blockProviderImpl{
		Config: conf,
		cache:  cache,
		cacheHit: stats.NewMetric(
			"cache_hit",
			"cache hits counter",
		),
		cacheMiss: stats.NewMetric(
			"cache_miss",
			"cache miss counter",
		),
	}, nil
}

func parseBlock(data []byte) (*BlockStored, error) {
	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, errors.Wrap(err, "unmarshaling block")
	}

	b := BlockStored{
		Transactions:       block.Transactions,
		TransactionsByHash: map[string]*json.RawMessage{},
	}
	for i := range block.Transactions {
		var tsx Transaction
		err := json.Unmarshal(block.Transactions[i], &tsx)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshaling transaction")
		}
		b.TransactionsByHash[tsx.Hash] = &b.Transactions[i]
	}

	return &b, nil
}

func (b *blockProviderImpl) GetTSXbyHash(block *BlockStored, hash string) ([]byte, error) {
	tsx, ok := block.TransactionsByHash[hash]
	if !ok {
		return nil, errors.Errorf("no transaction with hash: %s in block", hash)
	}

	return *tsx, nil
}

func (b *blockProviderImpl) GetTSXbyIndex(block *BlockStored, idx int) ([]byte, error) {
	if idx >= len(block.Transactions) {
		return nil, errors.Errorf("no transaction with index: %d", idx)
	}

	return block.Transactions[idx], nil
}

func (b *blockProviderImpl) GetBlock(
	ctx *application.EContext,
	blockNumber int64,
	isLatest bool,
) (*BlockStored, error) {
	if isLatest {
		var err error
		blockNumber, err = ctx.EthClient.LatestBlockNumber(ctx)
		b.latestBlockNumber = blockNumber
		if err != nil {
			return nil, errors.Wrap(err, "getting latest block number")
		}
	}

	result, ok := b.cache.Get(blockNumber)
	if ok {
		b.cacheHit.Observe(1)
		ctx.Logger.Info("cache hit")

		return result.(*BlockStored), nil
	}
	ctx.Logger.Info("cache miss")
	b.cacheMiss.Observe(1)

	data, err := ctx.EthClient.GetBlock(ctx, blockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting block")
	}

	block, err := parseBlock(data)
	if err != nil {
		return nil, errors.Wrap(err, "pasing block")
	}

	if b.latestBlockNumber-blockNumber >= b.UncertainBlockLimit {
		ctx.Logger.Info("added block to cache")
		b.cache.Add(blockNumber, block)

		return block, nil
	}

	if isLatest {
		return block, nil
	}

	b.latestBlockNumber, err = ctx.EthClient.LatestBlockNumber(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting latest block number")
	}
	if b.latestBlockNumber-blockNumber >= b.UncertainBlockLimit {
		ctx.Logger.Info("added block to cache")
		b.cache.Add(blockNumber, block)
	}

	return block, nil
}
