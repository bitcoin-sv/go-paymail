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
	// This is required first to get the corresponding P2P PaymentResolution endpoint url
	var capabilities *paymail.CapabilitiesResponse
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		logger.Fatal().Msgf("error getting capabilities: %s", err.Error())
	}
	logger.Info().Msgf("found capabilities: %d", len(capabilities.Capabilities))

	// Extract the URL from the capabilities response
	p2pURL := capabilities.GetString(paymail.BRFCP2PPaymentDestination, "")

	// Create the basic paymentRequest to achieve a payment destination (how many sats are you planning to send?)
	paymentRequest := &paymail.PaymentRequest{Satoshis: 1000}

	// Get the P2P destination
	var destination *paymail.PaymentDestinationResponse
	destination, err = client.GetP2PPaymentDestination(p2pURL, "mrz", "moneybutton.com", paymentRequest)
	if err != nil {
		logger.Fatal().Msgf("error getting destination: %s", err.Error())
	}
	logger.Info().Msgf("destination returned reference: %s and outputs: %d", destination.Reference, len(destination.Outputs))
}
