package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/logging"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Start with a BRFC specification
	existingBRFC := &paymail.BRFCSpec{
		Author:  "MrZ",
		ID:      "e898079d7d1a",
		Title:   "New BRFC",
		Version: "1",
	}

	// Validate the BRFC ID
	if valid, id, err := existingBRFC.Validate(); err != nil {
		logger.Fatal().Msgf("error validating BRFC id: %s", err.Error())
	} else if !valid {
		logger.Fatal().Msgf("brfc is invalid: %s", id)
	} else if valid {
		logger.Info().Msgf("brfc: %s is valid", existingBRFC.ID)
	}
}
