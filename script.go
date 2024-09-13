package paymail

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/script"
)

// INFO: This function has been moved from Bitcoinschema/go-bitcoin repository
// Use the equivalent from go-sdk repository when available

func GetAddressFromScript(s string) (string, error) {
	// No script?
	if len(s) == 0 {
		return "", errors.New("missing script")
	}

	// Decode the hex string into bytes
	scriptBytes, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}

	// Extract the addresses from the script
	decodedScript := script.NewFromBytes(scriptBytes)
	var addresses []string
	addresses, err = decodedScript.Addresses()
	if err != nil {
		return "", err
	}

	// Missing an address?
	if len(addresses) == 0 {
		// This error case should not occur since the error above will occur when no address is found,
		// however we ensure that we have an address for the NewLegacyAddressPubKeyHash() below
		return "", fmt.Errorf("invalid output script, missing an address")
	}

	// Use the encoded version of the address
	return addresses[0], nil
}
