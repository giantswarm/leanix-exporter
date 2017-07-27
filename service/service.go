package service

// Config represents the configuration used to create a new service.
import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/giantswarm/leanix-exporter/flag"
	"github.com/giantswarm/leanix-exporter/service/exporter"
	"github.com/giantswarm/leanix-exporter/service/version"
	"github.com/giantswarm/microerror"
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
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	config.Logger.Log("debug", fmt.Sprintf("creating cluster service with config: %#v", config))

	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Viper must not be empty")
	}

	var err error

	var exporterService *exporter.Service
	{
		exporterConfig := exporter.DefaultConfig()
		exporterConfig.Excludes = config.Viper.GetStringSlice(config.Flag.Service.Excludes)
		exporterService, err = exporter.New(exporterConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
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
