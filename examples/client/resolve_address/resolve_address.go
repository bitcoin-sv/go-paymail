package main

import (
	"github.com/bitcoin-sv/go-paymail/logging"
	"time"

	"github.com/bitcoin-sv/go-paymail"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		logger.Fatal().Msgf("error loading client: %s", err.Error())
	}

	// Get the capabilities
	// This is required first to get the corresponding AddressResolution endpoint url
	var capabilities *paymail.CapabilitiesResponse
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		logger.Fatal().Msgf("error getting capabilities: %s", err.Error())
	}
	logger.Info().Msgf("found capabilities: %d", len(capabilities.Capabilities))

	// Extract the resolution URL from the capabilities response
	resolveURL := capabilities.GetString(paymail.BRFCBasicAddressResolution, paymail.BRFCPaymentDestination)

	// Create the basic senderRequest to achieve an address resolution request
	senderRequest := &paymail.SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Get the address resolution results
	var resolution *paymail.ResolutionResponse
	if resolution, err = client.ResolveAddress(resolveURL, "mrz", "moneybutton.com", senderRequest); err != nil {
		logger.Fatal().Msgf("error getting resolution: %s", err.Error())
	}
	logger.Info().Msgf("resolved address: %v", resolution.Address)
}
