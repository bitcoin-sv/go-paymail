package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/logging"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Start with a paymail address
	paymailAddress := "MrZ@MoneyButton.com"

	// Sanitize the address, extract the parts
	alias, domain, address := paymail.SanitizePaymail(paymailAddress)
	logger.Info().Msgf("alias: %s domain: %s address: %s", alias, domain, address)
}
