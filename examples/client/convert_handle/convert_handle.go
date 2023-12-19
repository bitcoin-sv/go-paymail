package main

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/logging"
)

func main() {
	logger := logging.GetDefaultLogger()

	// Start with a handle
	handle := "$Mr-Z"

	// Convert the handle to paymail address
	address := paymail.ConvertHandle(handle, false)
	logger.Info().Msgf("handle %s was converted to: %s", handle, address)

	// Try another handle
	handle = "1MrZ"

	// Convert the handle to paymail address
	address = paymail.ConvertHandle(handle, false)
	logger.Info().Msgf("handle %s was converted to: %s", handle, address)
}
