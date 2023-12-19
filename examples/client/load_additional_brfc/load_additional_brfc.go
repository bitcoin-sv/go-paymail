package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/logging"
)

func main() {
	logger := logging.GetDefaultLogger()
	// Create a client with options

	// Load additional specification(s)
	additionalSpec := `[{"author": "andy (nChain)","id": "57dd1f54fc67","title": "BRFC Specifications","url": "http://bsvalias.org/01-02-brfc-id-assignment.html","version": "1"}]`
	_, err := paymail.NewClient(paymail.WithUserAgent(additionalSpec))
	if err != nil {
		logger.Fatal().Msgf("error loading client: %s", err.Error())
	}

}
