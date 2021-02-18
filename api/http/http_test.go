// +build testing

package http

import (
	"errors"
	"github.com/miscoler/ethereum-proxy/internal/application"
	"github.com/miscoler/ethereum-proxy/internal/blocktsx"
	"github.com/miscoler/ethereum-proxy/pkg/logger"
	"github.com/miscoler/ethereum-proxy/testutil/gomock/mockblocktsx"
	"github.com/miscoler/ethereum-proxy/testutil/teststats"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func test(
	t *testing.T,
	mock func(prv *mockblocktsx.MockBlockProvider),
	test func(t *testing.T, srv *httptest.Server),
) {
	lg, err := logger.New()
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prv := mockblocktsx.NewMockBlockProvider(ctrl)
	mock(prv)

	app := &application.Application{
		Clock:  clock.NewMock(),
		Stats:  teststats.New(),
		Logger: lg,
	}

	api, err := New(Config{}, app, prv)
	require.NoError(t, err)
	srv := httptest.NewServer(api.handler)
	defer srv.Close()

	test(t, srv)
}

func TestOk(t *testing.T) {
	t.Parallel()

	test(
		t,
		func(prv *mockblocktsx.MockBlockProvider) {
			prv.EXPECT().GetBlock(gomock.Any(), int64(0), true).Return(&blocktsx.BlockStored{}, nil)
			prv.EXPECT().GetTSXbyIndex(gomock.Any(), 1).Return([]byte("transaction data"), nil)
		},
		func(t *testing.T, srv *httptest.Server) {
			resp, err := http.Get(srv.URL + "/block/latest/txs/1")
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)
			data, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, []byte("transaction data"), data)
		},
	)
}

func TestBadRequest(t *testing.T) {
	t.Parallel()

	test(
		t,
		func(prv *mockblocktsx.MockBlockProvider) {},
		func(t *testing.T, srv *httptest.Server) {
			resp, err := http.Get(srv.URL + "/block/sdfs/txs/1")
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		},
	)
}

func TestInternalError(t *testing.T) {
	t.Parallel()

	test(
		t,
		func(prv *mockblocktsx.MockBlockProvider) {
			prv.EXPECT().GetBlock(gomock.Any(), int64(0), true).Return(nil, errors.New("some error"))
		},
		func(t *testing.T, srv *httptest.Server) {
			resp, err := http.Get(srv.URL + "/block/latest/txs/1")
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		},
	)
}
