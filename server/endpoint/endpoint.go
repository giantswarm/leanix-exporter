package endpoint

import (
	"github.com/giantswarm/leanix-exporter/server/endpoint/exporter"
	"github.com/giantswarm/leanix-exporter/server/endpoint/version"
	"github.com/giantswarm/leanix-exporter/server/middleware"
	"github.com/giantswarm/leanix-exporter/service"
	"github.com/giantswarm/microerror"
	micrologger "github.com/giantswarm/microkit/logger"
)

// Config represents the configuration used to create a endpoint.
type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
}

// DefaultConfig provides a default configuration to create a new endpoint by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:     nil,
		Middleware: nil,
		Service:    nil,
	}
}

// New creates a new configured endpoint.
func New(config Config) (*Endpoint, error) {
	var err error

	var exporterEndpoint *exporter.Endpoint
	{
		exporterEndpoint, err = exporter.New(exporter.Config{
			Logger:     config.Logger,
			Middleware: config.Middleware,
			Service:    config.Service,
		})

		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionEndpoint *version.Endpoint
	{
		versionConfig := version.DefaultConfig()
		versionConfig.Logger = config.Logger
		versionConfig.Middleware = config.Middleware
		versionConfig.Service = config.Service
		versionEndpoint, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newEndpoint := &Endpoint{
		Version:  versionEndpoint,
		Exporter: exporterEndpoint,
	}

	return newEndpoint, nil
}

// Endpoint is the endpoint collection.
type Endpoint struct {
	Version  *version.Endpoint
	Exporter *exporter.Endpoint
}
