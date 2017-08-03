package exporter

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/leanix-exporter/service/exporter/k8s"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

// Response exporter is the exporter service response
type Response struct {
	Namespaces []k8s.Namespace
	LastUpdate time.Time
}

// Config is the Exporter service configuration
type Config struct {
	Logger   micrologger.Logger
	Excludes []string
}

// DefaultConfig is a configuration is default value
func DefaultConfig() Config {
	return Config{
		Excludes: []string{},
		Logger:   nil,
	}
}

// New creates a new configured version service.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Excludes == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Excludes must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	newService := &Service{
		Config: config,
	}

	return newService, nil
}

// Service implements the exporter service interface.
type Service struct {
	Config
}

// Get return the exporter Response
func (s *Service) Get(ctx context.Context) (*Response, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, microerror.Mask(err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	namespaces, err := k8s.GetNamespaces(clientset, s.Config.Excludes, s.Config.Logger)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return &Response{
		Namespaces: namespaces,
		LastUpdate: time.Now(),
	}, nil
}
