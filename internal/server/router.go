package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/srv/internal/config"
	"github.com/qdm12/srv/internal/metrics"
	logmware "github.com/qdm12/srv/internal/server/middlewares/log"
	metricsmware "github.com/qdm12/srv/internal/server/middlewares/metrics"
	fsroute "github.com/qdm12/srv/internal/server/routes/fs"
)

func newRouter(config config.HTTP, logger logging.Logger,
	metrics metrics.Metrics, srvFS http.FileSystem) http.Handler {
	router := chi.NewRouter()

	// Middlewares
	logMiddleware := logmware.New(logger, config.LogRequests)
	metricsMiddleware := metricsmware.New(metrics)
	router.Use(metricsMiddleware, logMiddleware)

	if config.RootURL != "" {
		router.Handle(config.RootURL, http.RedirectHandler(config.RootURL+"/", http.StatusMovedPermanently))
	}
	router.Mount(config.RootURL+"/", fsroute.NewHandler(config.RootURL, srvFS))

	return router
}
