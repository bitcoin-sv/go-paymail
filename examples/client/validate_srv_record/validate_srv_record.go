package main

import (
	"context"
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
	srv, err = client.GetSRVRecord(paymail.DefaultServiceName, paymail.DefaultProtocol, "moneybutton.com")
	if err != nil {
		logger.Fatal().Msgf("error getting SRV record: %s", err.Error())
	}

	// Found record!
	logger.Info().Msgf("found SRV record: %v", srv)

	// Validate the record (1 instead of 10, moneybutton deviated from the defaults)
	err = client.ValidateSRVRecord(
		context.Background(), srv, paymail.DefaultPort, 1, paymail.DefaultWeight,
	)
	if err != nil {
		logger.Fatal().Msgf("failed validating SRV record: " + err.Error())
	}
}
