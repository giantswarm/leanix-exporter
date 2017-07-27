package exporter

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"k8s.io/api/core/v1"

	"github.com/giantswarm/leanix-exporter/server/middleware"
	"github.com/giantswarm/leanix-exporter/service"
	"github.com/giantswarm/microerror"
	micrologger "github.com/giantswarm/microkit/logger"
)

const (
	// Method is the HTTP method this endpoint is registered for.
	Method = "GET"
	// Name identifies the endpoint. It is aligned to the package path.
	Name = "exporter"
	// Path is the HTTP request path this endpoint is registered for.
	Path = "/exporter/"
)

// Config represents the configuration used to create a version endpoint.
type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
}

type pod struct {
	Name              string
	Status            string
	ContainerStatuses []v1.ContainerStatus
}
type namespace struct {
	Name string
	Pods []pod
}
type Response struct {
	Namespaces []namespace
	LastUpdate time.Time
}

// DefaultConfig provides a default configuration to create a new version
// endpoint by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:     nil,
		Middleware: nil,
		Service:    nil,
	}
}

// New creates a new configured version endpoint.
func New(config Config) (*Endpoint, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}
	if config.Middleware == nil {
		return nil, microerror.Maskf(invalidConfigError, "middleware must not be empty")
	}
	if config.Service == nil {
		return nil, microerror.Maskf(invalidConfigError, "service must not be empty")
	}

	newEndpoint := &Endpoint{
		Config: config,
	}

	return newEndpoint, nil
}

type Endpoint struct {
	Config
}

func (e *Endpoint) Decoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}
}

func (e *Endpoint) Encoder() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		return json.NewEncoder(w).Encode(response)
	}
}

func (e *Endpoint) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		serviceResponse, err := e.Service.Exporter.Get(ctx)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		nss := []namespace{}

		for _, ns := range serviceResponse.Namespaces {
			n := namespace{Name: ns.Name, Pods: []pod{}}
			for _, p := range ns.Pods {
				n.Pods = append(n.Pods, pod(p))
			}
			nss = append(nss, n)
		}
		r := Response{
			Namespaces: nss,
			LastUpdate: serviceResponse.LastUpdate,
		}

		return r, nil
	}
}

func (e *Endpoint) Method() string {
	return Method
}

func (e *Endpoint) Middlewares() []kitendpoint.Middleware {
	return []kitendpoint.Middleware{}
}

func (e *Endpoint) Name() string {
	return Name
}

func (e *Endpoint) Path() string {
	return Path
}
