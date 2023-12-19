package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/logging"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Start with a paymail address
	paymailAddress := "MrZ@MoneyButton.com"

	// Validate the paymail address format
	if err := paymail.ValidatePaymail(paymailAddress); err != nil {
		logger.Info().Msgf("error validating paymail: %s", err.Error())
	} else {
		logger.Info().Msg("paymail format is valid!")
	}
}
