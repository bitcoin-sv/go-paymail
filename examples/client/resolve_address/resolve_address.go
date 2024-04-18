package main

import (
	"log"
	"time"

	"github.com/bitcoin-sv/go-paymail"
)

func main() {
	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		log.Fatalf("error loading client: %s", err.Error())
	}

	// Get the capabilities
	// This is required first to get the corresponding AddressResolution endpoint url
	var capabilities *paymail.CapabilitiesResponse
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		log.Fatalf("error getting capabilities: %s", err.Error())
	}
	log.Printf("found capabilities: %d", len(capabilities.Capabilities))

	// Extract the resolution URL from the capabilities response
	resolveURL := capabilities.GetString(paymail.BRFCBasicAddressResolution, paymail.BRFCPaymentDestination)

	// Create the basic senderRequest to achieve an address resolution request
	senderRequest := &paymail.SenderRequest{
		Dt:           time.Now().UTC().Format(time.RFC3339),
		SenderHandle: "mrz@moneybutton.com",
		SenderName:   "MrZ",
	}

	// Get the address resolution results
	var resolution *paymail.ResolutionResponse
	if resolution, err = client.ResolveAddress(resolveURL, "mrz", "moneybutton.com", senderRequest); err != nil {
		log.Fatalf("error getting resolution: %s", err.Error())
	}
	log.Printf("resolved address: %v", resolution.Address)
}
