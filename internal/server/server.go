// Package server implements an HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/srv/internal/config"
	"github.com/qdm12/srv/internal/metrics"
)

type Server interface {
	Run(ctx context.Context) error
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

var (
	ErrCrashed  = errors.New("server crashed")
	ErrShutdown = errors.New("server could not be shutdown")
)

func (s *server) Run(ctx context.Context) error {
	server := http.Server{Addr: s.address, Handler: s.handler}
	shutdownErrCh := make(chan error)
	go func() {
		<-ctx.Done()
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		shutdownErrCh <- server.Shutdown(shutdownCtx)
	}()

	s.logger.Info("listening on " + s.address)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(ctx.Err(), context.Canceled) { // server crashed
		return fmt.Errorf("%w: %s", ErrCrashed, err)
	}

	if err := <-shutdownErrCh; err != nil {
		return fmt.Errorf("%w: %s", ErrShutdown, err)
	}

	return nil
}
