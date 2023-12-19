package main

import (
	"github.com/bitcoin-sv/go-paymail/logging"
	"net"

	"github.com/bitcoin-sv/go-paymail"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		logger.Fatal().Msgf("error loading client: %s", err.Error())
	}

	// Get the SRV record
	var srv *net.SRV
	if srv, err = client.GetSRVRecord(paymail.DefaultServiceName, paymail.DefaultProtocol, "moneybutton.com"); err != nil {
		logger.Fatal().Msgf("error getting SRV record: %s", err.Error())
	}
	logger.Info().Msgf("found SRV record: %v", srv)
}
