package service

// Config represents the configuration used to create a new service.
import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/giantswarm/leanix-exporter/flag"
	"github.com/giantswarm/leanix-exporter/service/exporter"
	"github.com/giantswarm/leanix-exporter/service/version"
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
)

type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Viper: nil,

		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",
	}
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}
	config.Logger.Log("debug", fmt.Sprintf("creating cluster service with config: %#v", config))

	// Settings.
	if config.Flag == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Viper must not be empty")
	}

	var err error

	var versionService *version.Service
	{

		versionService, err = version.New(version.Config{
			Description: config.Description,
			GitCommit:   config.GitCommit,
			Name:        config.Name,
			Source:      config.Source,
		})
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var exporterService *exporter.Service
	{
		exporterService, err = exporter.New(exporter.Config{Excludes: config.Viper.GetStringSlice(config.Flag.Excludes)})
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newService := &Service{
		Version:  versionService,
		Exporter: exporterService,
	}

	return newService, nil
}

type Service struct {
	Version  *version.Service
	Exporter *exporter.Service
}
