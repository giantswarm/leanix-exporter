package deleter

import (
	"fmt"
	"net/url"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/go-resty/resty"
	"golang.org/x/net/context"
)

const (
	Endpoint = "/v1/clusters/%s/"
)

// Config represents the configuration used to create a deleter service.
type Config struct {
	// Dependencies.
	RestClient *resty.Client

	// Settings.
	URL *url.URL
}

// DefaultConfig provides a default configuration to create a new deleter
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		RestClient: resty.New(),

		// Settings.
		URL: nil,
	}
}

// New creates a new configured deleter service.
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

func (s *Service) Delete(ctx context.Context, request Request) (*Response, error) {
	u, err := s.URL.Parse(fmt.Sprintf(Endpoint, request.Cluster.ID))
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	r, err := s.RestClient.R().SetBody(request).SetResult(DefaultResponse()).Delete(u.String())
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	if r.StatusCode() != 202 {
		return nil, microerror.MaskAnyf(executionFailedError, string(r.Body()))
	}

	response := r.Result().(*Response)

	return response, nil
}
