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
	// This is required first to get the corresponding PKI endpoint url
	var capabilities *paymail.CapabilitiesResponse
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		logger.Fatal().Msgf("error getting capabilities: %s", err.Error())
	}
	logger.Info().Msgf("found capabilities: %d", len(capabilities.Capabilities))

	// Extract the PKI URL from the capabilities response
	pkiURL := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)

	// Get the actual PKI
	var pki *paymail.PKIResponse
	if pki, err = client.GetPKI(pkiURL, "mrz", "moneybutton.com"); err != nil {
		logger.Fatal().Msgf("error getting pki: %s", err.Error())
	}
	logger.Info().Msgf("found pki: %v", pki)
}
