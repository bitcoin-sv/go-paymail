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
	// This is required first to get the corresponding VerifyPubKey endpoint url
	var capabilities *paymail.CapabilitiesResponse
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		logger.Fatal().Msgf("error getting capabilities: %s", err.Error())
	}
	logger.Info().Msgf("found capabilities: %d", len(capabilities.Capabilities))

	// Extract the verify URL from the capabilities response
	verifyURL := capabilities.GetString(paymail.BRFCVerifyPublicKeyOwner, "")

	// Verify the pubkey
	var verification *paymail.VerificationResponse
	verification, err = client.VerifyPubKey(verifyURL, "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err != nil {
		logger.Fatal().Msgf("error getting verification: %s", err.Error())
	}
	if verification.Match {
		logger.Info().Msgf("pubkey: %s matched handle: %s", verification.PubKey[:12]+"...", verification.Handle)
	} else {
		logger.Info().Msgf("pubkey: %s DID NOT MATCH handle: %s", verification.PubKey[:12]+"...", verification.Handle)
	}
}
