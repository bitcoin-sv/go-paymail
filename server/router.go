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

	// PKI request (public key information)
	router.GET(
		"/"+c.APIVersion+"/"+c.ServiceName+"/id/:paymailAddress",
		c.showPKI,
	)

	// Verify PubKey request (public key verification to paymail address)
	router.GET(
		"/"+c.APIVersion+"/"+c.ServiceName+"/verify-pubkey/:paymailAddress/:pubKey",
		c.verifyPubKey,
	)

	// Payment Destination request (address resolution)
	router.POST(
		"/"+c.APIVersion+"/"+c.ServiceName+"/address/:paymailAddress",
		c.resolveAddress,
	)

	// Public Profile request (returns Name & Avatar)
	router.GET(
		"/"+c.APIVersion+"/"+c.ServiceName+"/public-profile/:paymailAddress",
		c.publicProfile,
	)

	// P2P Destination request (returns output & reference)
	router.POST(
		"/"+c.APIVersion+"/"+c.ServiceName+"/p2p-payment-destination/:paymailAddress",
		c.p2pDestination,
	)

	// P2P Receive Tx request (receives the P2P transaction, broadcasts, returns tx_id)
	router.POST(
		"/"+c.APIVersion+"/"+c.ServiceName+"/receive-transaction/:paymailAddress",
		c.p2pReceiveTx,
	)

	// P2P BEEF capability Receive Tx request
	router.POST(
		"/"+c.APIVersion+"/"+c.ServiceName+"/beef/:paymailAddress",
		c.p2pReceiveBeefTx,
	)
}
