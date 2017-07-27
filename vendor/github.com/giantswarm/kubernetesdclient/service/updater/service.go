package updater

import (
	"fmt"
	"net/http"
	"net/url"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/go-resty/resty"
	"golang.org/x/net/context"
)

const (
	// Endpoint is the API endpoint of the service this client action interacts
	// with.
	Endpoint = "/v1/clusters/%s/"
	// Name is the service name being implemented. This can be used for e.g.
	// logging.
	Name = "cluster/updater"
)

// Config represents the configuration used to create a creator service.
type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	RestClient *resty.Client

	// Settings.
	URL *url.URL
}

// DefaultConfig provides a default configuration to create a new creator
// service by best effort.
func DefaultConfig() Config {
	var err error

	var newLogger micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		newLogger, err = micrologger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	newConfig := Config{
		// Dependencies.
		Logger:     newLogger,
		RestClient: resty.New(),

		// Settings.
		URL: nil,
	}

	return newConfig
}

// New creates a new configured updater service.
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

func (s *Service) Update(ctx context.Context, request Request) (*Response, error) {
	u, err := s.URL.Parse(fmt.Sprintf(Endpoint, request.Cluster.ID))
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	s.Logger.Log("debug", fmt.Sprintf("sending PATCH request to %s", u.String()), "service", Name)
	r, err := s.RestClient.R().SetBody(request.Cluster.Patch).SetResult(DefaultResponse()).Patch(u.String())
	if err != nil {
		return nil, microerror.MaskAny(err)
	}
	s.Logger.Log("debug", fmt.Sprintf("received status code %d", r.StatusCode()), "service", Name)

	if r.StatusCode() == http.StatusNotFound {
		return nil, microerror.MaskAny(notFoundError)
	} else if r.StatusCode() != http.StatusOK {
		return nil, microerror.MaskAny(fmt.Errorf(string(r.Body())))
	}

	response := r.Result().(*Response)

	return response, nil
}
