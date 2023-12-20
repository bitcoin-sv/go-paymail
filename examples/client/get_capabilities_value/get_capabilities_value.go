package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"log"
)

func main() {
	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		log.Fatalf("error loading client: %s", err.Error())
	}

	// Get the capabilities
	var capabilities *paymail.CapabilitiesResponse
	capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort)
	if err != nil {
		log.Fatalf("error getting capabilities: %s", err.Error())
	}
	log.Printf("found capabilities: %d", len(capabilities.Capabilities))

	// Get the URL for a capability
	endpoint := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)
	log.Printf("capability endpoint found: %v", endpoint)
}
