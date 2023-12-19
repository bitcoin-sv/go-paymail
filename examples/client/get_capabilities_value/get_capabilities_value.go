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

	// Get the capabilities
	var capabilities *paymail.CapabilitiesResponse
	capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort)
	if err != nil {
		logger.Fatal().Msgf("error getting capabilities: %s", err.Error())
	}
	logger.Info().Msgf("found capabilities: %d", len(capabilities.Capabilities))

	// Get the URL for a capability
	endpoint := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)
	logger.Info().Msgf("capability endpoint found: %v", endpoint)
}
