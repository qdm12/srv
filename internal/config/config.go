// Package config takes care of reading and checking the program configuration
// from environment variables.
package config

import (
	"errors"
	"fmt"

	"github.com/qdm12/golibs/params"
)

type Config struct {
	HTTP      HTTP
	Filepaths Filepaths
	Metrics   Metrics
	Log       Log
	Health    Health
}

var (
	ErrHTTPConfig     = errors.New("cannot obtain HTTP server config")
	ErrFilepathConfig = errors.New("cannot obtain file paths config")
	ErrMetricsConfig  = errors.New("cannot obtain metrics config")
	ErrLogConfig      = errors.New("cannot obtain log config")
	ErrHealthConfig   = errors.New("cannot obtain health config")
)

func (c *Config) get(env params.Env) (warnings []string, err error) {
	warning, err := c.HTTP.get(env)
	if len(warning) > 0 {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrHTTPConfig, err)
	}

	err = c.Filepaths.get(env)
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrLogConfig, err)
	}

	warning, err = c.Metrics.get(env)
	if len(warning) > 0 {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrMetricsConfig, err)
	}

	err = c.Log.get(env)
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrLogConfig, err)
	}

	warning, err = c.Health.get(env)
	if len(warning) > 0 {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrHealthConfig, err)
	}

	return warnings, nil
}
