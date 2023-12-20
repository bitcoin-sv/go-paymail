package main

import (
	"log"
	"net"

	"github.com/bitcoin-sv/go-paymail"
)

func main() {
	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		log.Fatalf("error loading client: %s", err.Error())
	}

	// Get the SRV record
	var srv *net.SRV
	if srv, err = client.GetSRVRecord(paymail.DefaultServiceName, paymail.DefaultProtocol, "moneybutton.com"); err != nil {
		log.Fatalf("error getting SRV record: %s", err.Error())
	}
	log.Printf("found SRV record: %v", srv)
}
