package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/logging"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		logger.Fatal().Msgf("error loading client: %s", err.Error())
	}

	// Check the SSL certificate
	var valid bool
	if valid, err = client.CheckSSL("moneybutton.com"); err != nil {
		logger.Fatal().Msg("error getting SSL certificate: " + err.Error())
	} else if !valid {
		logger.Fatal().Msg("SSL certificate validation failed")
	}
	logger.Info().Msg("valid SSL certificate found for: moneybutton.com")
}
