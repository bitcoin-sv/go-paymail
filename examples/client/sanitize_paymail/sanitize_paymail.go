package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"log"
)

func main() {
	// Start with a paymail address
	paymailAddress := "MrZ@MoneyButton.com"

	// Sanitize the address, extract the parts
	alias, domain, address := paymail.SanitizePaymail(paymailAddress)
	log.Printf("alias: %s domain: %s address: %s", alias, domain, address)
}
