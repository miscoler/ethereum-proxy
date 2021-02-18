// +build testing

package blocktsx

import (
	"context"
	"github.com/miscoler/ethereum-proxy/internal/application"
	"github.com/miscoler/ethereum-proxy/internal/blocktsx/config"
	"github.com/miscoler/ethereum-proxy/pkg/logger"
	"github.com/miscoler/ethereum-proxy/testutil/gomock/mockethclient"
	"github.com/miscoler/ethereum-proxy/testutil/teststats"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func test(
	t *testing.T,
	mock func(*mockethclient.MockEthClient),
	test func(*application.EContext, BlockProvider, *testing.T),
) {
	lg, err := logger.New()
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prv, err := New(config.Config{
		UncertainBlockLimit: 2,
		CacheSize:           2,
	}, teststats.New())
	require.NoError(t, err)

	ethclient := mockethclient.NewMockEthClient(ctrl)
	app := &application.Application{
		Clock:     clock.NewMock(),
		EthClient: ethclient,
		Stats:     teststats.New(),
		Logger:    lg,
	}

	mock(ethclient)
	ctx := application.UpgradeContext(context.Background(), lg, app)
	test(ctx, prv, t)
}

func TestLatest(t *testing.T) {
	t.Parallel()

	test(
		t,
		func(ethclient *mockethclient.MockEthClient) {
			ethclient.EXPECT().LatestBlockNumber(gomock.Any()).Return(int64(3), nil)
			ethclient.EXPECT().GetBlock(gomock.Any(), int64(3)).Return([]byte(
				"{\"transactions\":[{\"hash\":\"1\"},{\"hash\":\"2\"}]}",
			), nil)
		},
		func(ctx *application.EContext, prv BlockProvider, t *testing.T) {
			blockStored, err := prv.GetBlock(ctx, 0, true)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))
		},
	)
}

func TestCacheNew(t *testing.T) {
	t.Parallel()

	test(
		t,
		func(ethclient *mockethclient.MockEthClient) {
			ethclient.EXPECT().GetBlock(gomock.Any(), int64(1)).Return([]byte(
				"{\"transactions\":[{\"hash\":\"1\"},{\"hash\":\"2\"}]}",
			), nil)
			ethclient.EXPECT().LatestBlockNumber(gomock.Any()).Return(int64(1), nil)
			ethclient.EXPECT().GetBlock(gomock.Any(), int64(1)).Return([]byte(
				"{\"transactions\":[{\"hash\":\"1\"},{\"hash\":\"2\"}]}",
			), nil)
			ethclient.EXPECT().LatestBlockNumber(gomock.Any()).Return(int64(3), nil)
		},
		func(ctx *application.EContext, prv BlockProvider, t *testing.T) {
			blockStored, err := prv.GetBlock(ctx, 1, false)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))

			blockStored, err = prv.GetBlock(ctx, 1, false)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))

			blockStored, err = prv.GetBlock(ctx, 1, false)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))
		},
	)
}

func TestCacheOld(t *testing.T) {
	t.Parallel()

	test(
		t,
		func(ethclient *mockethclient.MockEthClient) {
			ethclient.EXPECT().LatestBlockNumber(gomock.Any()).Return(int64(3), nil)
			ethclient.EXPECT().GetBlock(gomock.Any(), int64(3)).Return([]byte(
				"{\"transactions\":[{\"hash\":\"1\"},{\"hash\":\"2\"}]}",
			), nil)
			ethclient.EXPECT().GetBlock(gomock.Any(), int64(0)).Return([]byte(
				"{\"transactions\":[{\"hash\":\"1\"},{\"hash\":\"2\"}]}",
			), nil)
		},
		func(ctx *application.EContext, prv BlockProvider, t *testing.T) {
			blockStored, err := prv.GetBlock(ctx, 0, true)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))

			blockStored, err = prv.GetBlock(ctx, 0, false)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))

			blockStored, err = prv.GetBlock(ctx, 0, false)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))
		},
	)
}

func TestGetTSXbyHash(t *testing.T) {
	t.Parallel()

	test(
		t,
		func(ethclient *mockethclient.MockEthClient) {
			ethclient.EXPECT().LatestBlockNumber(gomock.Any()).Return(int64(3), nil)
			ethclient.EXPECT().GetBlock(gomock.Any(), int64(3)).Return([]byte(
				"{\"transactions\":[{\"hash\":\"1\"},{\"hash\":\"2\"}]}",
			), nil)
		},
		func(ctx *application.EContext, prv BlockProvider, t *testing.T) {
			blockStored, err := prv.GetBlock(ctx, 0, true)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))

			_, err = prv.GetTSXbyHash(blockStored, "2")
			require.NoError(t, err)

			_, err = prv.GetTSXbyHash(blockStored, "3")
			require.Error(t, err)
		},
	)
}

func TestCacheGetTSXbyIdx(t *testing.T) {
	t.Parallel()

	test(
		t,
		func(ethclient *mockethclient.MockEthClient) {
			ethclient.EXPECT().LatestBlockNumber(gomock.Any()).Return(int64(3), nil)
			ethclient.EXPECT().GetBlock(gomock.Any(), int64(3)).Return([]byte(
				"{\"transactions\":[{\"hash\":\"1\"},{\"hash\":\"2\"}]}",
			), nil)
		},
		func(ctx *application.EContext, prv BlockProvider, t *testing.T) {
			blockStored, err := prv.GetBlock(ctx, 0, true)
			require.NoError(t, err)
			require.Equal(t, 2, len(blockStored.Transactions))
			require.Equal(t, 2, len(blockStored.TransactionsByHash))

			_, err = prv.GetTSXbyIndex(blockStored, 1)
			require.NoError(t, err)

			_, err = prv.GetTSXbyIndex(blockStored, 3)
			require.Error(t, err)
		},
	)
}
