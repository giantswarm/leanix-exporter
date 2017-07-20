package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/giantswarm/leanix-exporter/server/endpoint"
	"github.com/giantswarm/leanix-exporter/server/middleware"
	"github.com/giantswarm/leanix-exporter/service"
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Config represents the configuration used to create a new server object.
type Config struct {
	// Dependencies.
	Service *service.Service

	// Settings.
	MicroServerConfig microserver.Config
}

// New creates a new configured server object.
func New(config Config) (microserver.Server, error) {
	var err error

	var middlewareCollection *middleware.Middleware
	{
		middlewareConfig := middleware.DefaultConfig()
		middlewareConfig.Logger = config.MicroServerConfig.Logger
		middlewareConfig.Service = config.Service
		middlewareCollection, err = middleware.New(middlewareConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var endpointCollection *endpoint.Endpoint
	{

		endpointCollection, err = endpoint.New(endpoint.Config{
			Logger:     config.MicroServerConfig.Logger,
			Middleware: middlewareCollection,
			Service:    config.Service,
		})

		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newServer := &server{
		// Dependencies.
		logger: config.MicroServerConfig.Logger,

		// Internals.
		bootOnce:     sync.Once{},
		config:       config.MicroServerConfig,
		shutdownOnce: sync.Once{},
	}

	// Apply internals to the micro server config.
	newServer.config.Endpoints = []microserver.Endpoint{
		endpointCollection.Version,
		endpointCollection.Exporter,
	}
	newServer.config.ErrorEncoder = newServer.newErrorEncoder()

	return newServer, nil
}

type server struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	bootOnce     sync.Once
	config       microserver.Config
	shutdownOnce sync.Once
}

func (s *server) Boot() {
	s.bootOnce.Do(func() {
		// Here goes your custom boot logic for your server/endpoint/middleware, if
		// any.
	})
}

func (s *server) Config() microserver.Config {
	return s.config
}

func (s *server) Shutdown() {
	s.shutdownOnce.Do(func() {
		// Here goes your custom shutdown logic for your server/endpoint/middleware,
		// if any.
	})
}

func (s *server) newErrorEncoder() kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		rErr := err.(microserver.ResponseError)
		uErr := rErr.Underlying()
		if uErr != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
