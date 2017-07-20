package version

import (
	"context"
	"runtime"

	microerror "github.com/giantswarm/microkit/error"
)

// Response is the return value of the service action.
type Response struct {
	Description string `json:"description"`
	GitCommit   string `json:"git_commit"`
	GoVersion   string `json:"go_version"`
	Name        string `json:"name"`
	OSArch      string `json:"os_arch"`
	Source      string `json:"source"`
}

type Config struct {
	// Settings.
	Description string
	GitCommit   string
	Name        string
	Source      string
}

// New creates a new configured version service.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Description == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "description commit must not be empty")
	}
	if config.GitCommit == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "git commit must not be empty")
	}
	if config.Name == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "name must not be empty")
	}
	if config.Source == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "name must not be empty")
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
	return &Response{
		Description: s.Description,
		GitCommit:   s.GitCommit,
		GoVersion:   runtime.Version(),
		Name:        s.Name,
		OSArch:      runtime.GOOS + "/" + runtime.GOARCH,
		Source:      s.Source,
	}, nil
}
