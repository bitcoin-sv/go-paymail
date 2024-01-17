package server

import (
	"net/http"

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
func (c *Configuration) RegisterBasicRoutes(r *nrhttprouter.Router) {
	c.registerBasicRoutes(r)
}

// RegisterRoutes register all the available paymail routes to the http router
func (c *Configuration) RegisterRoutes(r *nrhttprouter.Router) {
	c.registerPaymailRoutes(r)
}

// registerBasicRoutes will register basic server related routes
func (c *Configuration) registerBasicRoutes(router *nrhttprouter.Router) {

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

// registerPaymailRoutes will register all paymail related routes
func (c *Configuration) registerPaymailRoutes(router *nrhttprouter.Router) {

	// Capabilities (service discovery)
	router.GET(
		"/.well-known/"+c.ServiceName,
		c.showCapabilities,
	)

	for key, cap := range c.callableCapabilities {
		c.Logger.Info().Msgf("Registering endpoint for capability: %s", key)
		router.Handle(
			cap.Method,
			cap.Path,
			cap.Handler,
		)
	}
}
