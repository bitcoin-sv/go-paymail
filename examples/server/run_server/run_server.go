package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
		server.WithDomain("localhost"),
		server.WithDomain("another.com"),
		server.WithDomain("test.com"),
		server.WithGenericCapabilities(),
		server.WithPort(3000),
		server.WithServiceName("BsvAliasCustom"),
		server.WithTimeout(15*time.Second),
		server.WithCapabilities(customCapabilities()),
	)
	config.Prefix = "http://" //normally paymail requires https, but for demo purposes we'll use http
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	// Create & start the server
	server.StartServer(server.CreateServer(config), config.Logger)
}

func customCapabilities() map[string]any {
	exampleBrfcKey := "406cef0ae2d6"
	return map[string]any{
		"custom_static_boolean": false,
		"custom_static_int":     10,
		exampleBrfcKey:          true,
		"custom_callable_cap": server.CallableCapability{
			Path:   fmt.Sprintf("/display_paymail/%s", server.PaymailAddressTemplate),
			Method: http.MethodGet,
			Handler: func(c *gin.Context) {
				incomingPaymail := c.Param(server.PaymailAddressParamName)

				response := map[string]string{
					"paymail": incomingPaymail,
				}

				c.Header("Content-Type", "application/json")
				c.JSON(http.StatusOK, response)
			},
		},
	}
}
