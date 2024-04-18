// Package server is the server package for Paymail
package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

// CreateServer will create a basic Paymail Server
func CreateServer(c *Configuration) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", c.Port), // Address to run the server on
		Handler:           Handlers(c),                // Load all the routes
		ReadHeaderTimeout: c.Timeout,                  // Basic default timeout for header read requests
		ReadTimeout:       c.Timeout,                  // Basic default timeout for read requests
		WriteTimeout:      c.Timeout,                  // Basic default timeout for write requests
	}
}

// StartServer will run the Paymail server
func StartServer(srv *http.Server, logger *zerolog.Logger) {
	logger.Info().Str("address", srv.Addr).Msg("starting go paymail server...")
	logger.Fatal().Msg(srv.ListenAndServe().Error())
}
