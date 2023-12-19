package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/logging"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Start with a domain name
	domainName := "MoneyButton.com"

	// Validate the domain name
	if err := paymail.ValidateDomain(domainName); err != nil {
		logger.Info().Msgf("error validating domain: %s", err.Error())
	} else {
		logger.Info().Msg("domain format is valid!")
	}
}
