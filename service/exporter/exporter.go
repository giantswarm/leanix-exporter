package exporter

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/leanix-exporter/service/exporter/k8s"
	"github.com/giantswarm/microerror"
	micrologger "github.com/giantswarm/microkit/logger"
)

type Response struct {
	Namespaces []k8s.Namespace
	LastUpdate time.Time
}

type Config struct {
	Logger   micrologger.Logger
	Excludes []string
}

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
	newService := &Service{
		Config: config,
	}

	return newService, nil
}

// Service implements the version service interface.
type Service struct {
	Config
}

func (s *Service) Get(ctx context.Context) (*Response, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, microerror.Mask(err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
	}

	return &Response{
		Namespaces: k8s.GetNamespaces(clientset, s.Config.Excludes, s.Config.Logger),
		LastUpdate: time.Now(),
	}, nil
}
