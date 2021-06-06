// Package server implements an HTTP server.
package server

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/srv/internal/config"
	"github.com/qdm12/srv/internal/metrics"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup, crashed chan<- error)
}

type server struct {
	address string
	logger  logging.Logger
	handler http.Handler
}

var ErrAccessSrvDirectory = errors.New("cannot access srv directory in embedded filesystem")

func New(c config.HTTP, logger logging.Logger, metrics metrics.Metrics, fs http.FileSystem) Server {
	handler := newRouter(c, logger, metrics, fs)
	return &server{
		address: c.Address,
		logger:  logger,
		handler: handler,
	}
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup, crashed chan<- error) {
	defer wg.Done()
	server := http.Server{Addr: s.address, Handler: s.handler}
	go func() {
		<-ctx.Done()
		s.logger.Warn("context canceled: shutting down")
		defer s.logger.Warn("shut down")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: %s", err)
		}
	}()

	s.logger.Info("listening on %s", s.address)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(ctx.Err(), context.Canceled) { // server crashed
		crashed <- err
	}
}
