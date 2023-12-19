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
	// This is required first to get the corresponding PublicProfile endpoint url
	var capabilities *paymail.CapabilitiesResponse
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		logger.Fatal().Msgf("error getting capabilities: %s", err.Error())
	}
	logger.Info().Msgf("found capabilities: %d", len(capabilities.Capabilities))

	// Extract the PublicProfile URL from the capabilities response
	publicProfileURL := capabilities.GetString(paymail.BRFCPublicProfile, "")

	// Get the public profile
	var profile *paymail.PublicProfileResponse
	if profile, err = client.GetPublicProfile(publicProfileURL, "mrz", "moneybutton.com"); err != nil {
		logger.Fatal().Msgf("error getting profile: %s", err.Error())
	}
	logger.Info().Msgf("found profile: %s : %s", profile.Name, profile.Avatar)
}
