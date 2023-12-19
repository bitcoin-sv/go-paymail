package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/logging"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Start with a new BRFC specification
	newBRFC := &paymail.BRFCSpec{
		Author:  "MrZ",
		Title:   "New BRFC",
		Version: "1",
	}

	// Generate the BRFC ID
	if err := newBRFC.Generate(); err != nil {
		logger.Fatal().Msgf("error generating BRFC id: %s", err.Error())
	}
	logger.Info().Msgf("id generated: %s", newBRFC.ID)
}
