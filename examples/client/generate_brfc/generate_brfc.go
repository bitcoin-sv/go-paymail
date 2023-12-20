package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"log"
)

func main() {
	// Start with a new BRFC specification
	newBRFC := &paymail.BRFCSpec{
		Author:  "MrZ",
		Title:   "New BRFC",
		Version: "1",
	}

	// Generate the BRFC ID
	if err := newBRFC.Generate(); err != nil {
		log.Fatalf("error generating BRFC id: %s", err.Error())
	}
	log.Fatalf("id generated: %s", newBRFC.ID)
}
