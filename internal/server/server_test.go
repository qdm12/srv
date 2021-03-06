package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/qdm12/srv/internal/config"
	"github.com/qdm12/srv/internal/metrics/mock_metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	config := config.HTTP{
		Address: "test",
	}
	logger := mock_logging.NewMockLogger(ctrl)
	metrics := mock_metrics.NewMockMetrics(ctrl)
	fs := http.Dir("path")

	serverInterface := New(config, logger, metrics, fs)
	serverImpl, ok := serverInterface.(*server)
	require.True(t, ok)
	assert.Equal(t, config.Address, serverImpl.address)
	assert.Equal(t, logger, serverImpl.logger)
	assert.NotNil(t, serverImpl.handler)
}

func Test_server_Run(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	const address = "127.0.0.1:9000"

	logger := mock_logging.NewMockLogger(ctrl)
	logger.EXPECT().Info("listening on " + address)

	server := &server{
		address: address,
		handler: nil,
		logger:  logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)

	go func() {
		errCh <- server.Run(ctx)
	}()

	cancel()
	err := <-errCh
	assert.NoError(t, err)
}
