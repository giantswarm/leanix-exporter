package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/giantswarm/leanix-exporter/flag"
	"github.com/giantswarm/leanix-exporter/server"
	"github.com/giantswarm/leanix-exporter/service"
	"github.com/giantswarm/microkit/command"
	"github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"
)

var (
	description = "Leanix exporter microservice using the microkit framework."
	gitCommit   = "n/a"
	name        = "leanix-exporter"
	source      = "https://github.com/giantswarm/leanix-exporter"
	f           = flag.New()
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger logger.Logger
	{
		loggerConfig := logger.DefaultConfig()
		loggerConfig.IOWriter = os.Stdout
		newLogger, err = logger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	newServerFactory := func(v *viper.Viper) microserver.Server {

		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			newService, err = service.New(service.Config{
				Logger:      newLogger,
				Viper:       v,
				Flag:        f,
				Description: description,
				GitCommit:   gitCommit,
				Name:        name,
				Source:      source,
			})
			if err != nil {
				panic(err)
			}
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			serverConfig := server.Config{MicroServerConfig: microserver.DefaultConfig()}

			serverConfig.MicroServerConfig.Logger = newLogger
			serverConfig.MicroServerConfig.ServiceName = name
			serverConfig.Service = newService

			newServer, err = server.New(serverConfig)
			if err != nil {
				panic(err)
			}
		}

		return newServer
	}

	var newCommand command.Command
	{
		commandConfig := command.DefaultConfig()

		commandConfig.Logger = newLogger
		commandConfig.ServerFactory = newServerFactory

		commandConfig.Description = description
		commandConfig.GitCommit = gitCommit
		commandConfig.Name = name
		commandConfig.Source = source

		newCommand, err = command.New(commandConfig)
		if err != nil {
			panic(err)
		}
	}
	daemonCommand := newCommand.DaemonCommand().CobraCommand()
	daemonCommand.PersistentFlags().String(f.Excludes, "", "Namespace to exclude")

	newCommand.CobraCommand().Execute()

}
