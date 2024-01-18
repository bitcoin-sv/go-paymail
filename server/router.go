package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
)

// Handlers are used to isolate loading the routes (used for testing)
func Handlers(configuration *Configuration) *nrhttprouter.Router {

	// Create a new router
	r := nrhttprouter.New(nil)

	// Register the routes
	configuration.RegisterBasicRoutes(r)
	configuration.RegisterRoutes(r)

	// Return the router
	return r
}

// RegisterBasicRoutes register the basic routes to the http router
func (c *Configuration) RegisterBasicRoutes(router *nrhttprouter.Router) {
	// Skip if not set
	if c.BasicRoutes == nil {
		return
	}

	// Set the main index page (navigating to slash)
	if c.BasicRoutes.AddIndexRoute {
		router.GET("/", index)
		// router.OPTIONS("/", router.SetCrossOriginHeaders) // Disabled for security
	}

	// Set the health request (used for load balancers)
	if c.BasicRoutes.AddHealthRoute {
		router.GET("/health", health)
		router.OPTIONS("/health", health)
		router.HEAD("/health", health)
	}

	// Set the 404 handler (any request not detected)
	if c.BasicRoutes.Add404Route {
		router.NotFound = http.HandlerFunc(notFound)
	}

	// Set the method not allowed
	if c.BasicRoutes.AddNotAllowed {
		router.MethodNotAllowed = http.HandlerFunc(methodNotAllowed)
	}
}

// RegisterRoutes register all the available paymail routes to the http router
func (c *Configuration) RegisterRoutes(router *nrhttprouter.Router) {
	router.GET("/.well-known/"+c.ServiceName, c.showCapabilities) // Capabilities (service discovery)

	for key, cap := range c.callableCapabilities {
		routerPath := c.templateToRouterPath(cap.Path)
		router.Handle(
			cap.Method,
			routerPath,
			cap.Handler,
		)

		c.Logger.Info().Msgf("Registering endpoint for capability: %s", key)
		c.Logger.Debug().Msgf("Endpoint[%s]: %s %s", key, cap.Method, routerPath)
	}
}

func (c *Configuration) templateToRouterPath(template string) string {
	urlParam := func(name string) string { return ":" + name }

	template = strings.ReplaceAll(template, PaymailAddressTemplate, urlParam(PaymailAddressParamName))
	template = strings.ReplaceAll(template, PubKeyTemplate, urlParam(PubKeyParamName))
	return fmt.Sprintf("/%s/%s/%s", c.APIVersion, c.ServiceName, strings.TrimPrefix(template, "/"))
}
