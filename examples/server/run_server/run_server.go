package main

import (
	"time"

	"github.com/bitcoin-sv/go-paymail/logging"

	"github.com/bitcoin-sv/go-paymail/server"
)

func main() {
	logger := logging.GetDefaultLogger()

	// initialize the demo database
	if err := InitDemoDatabase(); err != nil {
		logger.Fatal().Msg(err.Error())
	}

	// Custom server with lots of customizable goodies
	config, err := server.NewConfig(
		new(demoServiceProvider),
		server.WithBasicRoutes(),
		server.WithDomain("localhost"), // todo: make this work locally?
		server.WithDomain("another.com"),
		server.WithDomain("test.com"),
		server.WithGenericCapabilities(),
		server.WithPort(3000),
		server.WithServiceName("BsvAliasCustom"),
		server.WithTimeout(15*time.Second),
	)
	config.Prefix = "http://"
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	// Create & start the server
	server.StartServer(server.CreateServer(config), config.Logger)
}
