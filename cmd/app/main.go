package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/srv/internal/config"
	"github.com/qdm12/srv/internal/filesystem"
	"github.com/qdm12/srv/internal/health"
	"github.com/qdm12/srv/internal/metrics"
	"github.com/qdm12/srv/internal/models"
	"github.com/qdm12/srv/internal/server"
	"github.com/qdm12/srv/internal/shutdown"
	"github.com/qdm12/srv/internal/splash"
)

var (
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

func main() {
	buildInfo := models.BuildInformation{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	configReader := config.NewReader()

	logger := logging.NewParent(logging.Settings{})

	args := os.Args

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, configReader)
	}()

	select {
	case <-ctx.Done():
		logger.Warn("Caught OS signal, shutting down\n")
		stop()
	case err := <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck
			os.Exit(0)
		}
		logger.Error(err)
	}

	const shutdownGracePeriod = 5 * time.Second
	timer := time.NewTimer(shutdownGracePeriod)
	select {
	case <-errorCh:
		if !timer.Stop() {
			<-timer.C
		}
		logger.Info("Shutdown successful")
	case <-timer.C:
		logger.Warn("Shutdown timed out")
	}

	os.Exit(1)
}

func _main(ctx context.Context, buildInfo models.BuildInformation,
	args []string, logger logging.ParentLogger, configReader config.Reader) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()

		config, _, err := configReader.ReadConfig()
		if err != nil {
			return err
		}

		return client.Query(ctx, config.Health.Address)
	}

	fmt.Println(splash.Splash(buildInfo))

	config, warnings, err := configReader.ReadConfig()
	for _, warning := range warnings {
		logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	files, directories, err := filesystem.WalkSrv(config.HTTP.SrvFilepath)
	if err != nil {
		return err
	}
	for _, file := range files {
		logger.Debug("Found file: " + file)
	}
	logger.Info("Found " + strconv.Itoa(len(files)) + " files and " +
		strconv.Itoa(len(directories)) + " directories in " + config.HTTP.SrvFilepath)

	shutdownServersGroup := shutdown.NewGroup("servers: ")

	logger = logger.NewChild(logging.Settings{Level: config.Log.Level})

	metricsLogger := logger.NewChild(logging.Settings{Prefix: "metrics server: "})
	metricsServer := metrics.NewServer(config.Metrics.Address, metricsLogger)
	const registerMetrics = true
	metrics, err := metrics.New(registerMetrics)
	if err != nil {
		return err
	}
	metricsServerCtx, metricsServerDone := shutdownServersGroup.Add("metrics", time.Second)
	go func() {
		defer close(metricsServerDone)
		if err := metricsServer.Run(metricsServerCtx); err != nil {
			logger.Error(err.Error())
		}
	}()

	serverLogger := logger.NewChild(logging.Settings{Prefix: "http server: "})
	srvFS := http.Dir(config.HTTP.SrvFilepath)
	mainServer := server.New(config.HTTP, serverLogger, metrics, srvFS)
	serverCtx, serverDone := shutdownServersGroup.Add("server", time.Second)
	go func() {
		defer close(serverDone)
		if err := mainServer.Run(serverCtx); err != nil {
			logger.Error(err.Error())
			if errors.Is(err, server.ErrCrashed) {
				cancel() // stop other routines
			}
		}
	}()

	healthcheck := func() error { return nil }
	heathcheckLogger := logger.NewChild(logging.Settings{Prefix: "healthcheck: "})
	healthServer := health.NewServer(config.Health.Address, heathcheckLogger, healthcheck)
	healthServerCtx, healthServerDone := shutdownServersGroup.Add("health", time.Second)
	go func() {
		defer close(healthServerDone)
		if err := healthServer.Run(healthServerCtx); err != nil {
			logger.Error(err.Error())
		}
	}()

	shutdownOrder := shutdown.NewOrder()
	shutdownOrder.Append(shutdownServersGroup)

	<-ctx.Done()
	return shutdownOrder.Shutdown(time.Second, logger)
}
