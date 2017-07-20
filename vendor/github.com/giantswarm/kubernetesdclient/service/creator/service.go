package creator

import (
	"net/url"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/go-resty/resty"
	"golang.org/x/net/context"
)

const (
	// Endpoint is the API endpoint of the service this client action interacts
	// with.
	Endpoint = "/v1/clusters/"
)

// Config represents the configuration used to create a creator service.
type Config struct {
	// Dependencies.
	RestClient *resty.Client

	// Settings.
	URL *url.URL
}

// DefaultConfig provides a default configuration to create a new creator
// service by best effort.
func DefaultConfig() Config {
	newConfig := Config{
		// Dependencies.
		RestClient: resty.New(),

		// Settings.
		URL: nil,
	}

	return newConfig
}

// New creates a new configured creator service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.RestClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "rest client must not be empty")
	}

	// Settings.
	if config.URL == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "URL must not be empty")
	}

	newService := &Service{
		Config: config,
	}

	return newService, nil
}

type Service struct {
	Config
}

func (s *Service) Create(ctx context.Context, request Request) (*Response, error) {
	u, err := s.URL.Parse(Endpoint)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	r, err := s.RestClient.R().SetBody(request).SetResult(DefaultResponse()).Post(u.String())
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	if r.StatusCode() != 201 {
		return nil, microerror.MaskAnyf(executionFailedError, string(r.Body()))
	}

	response := r.Result().(*Response)

	return response, nil
}
