package paymail

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bitcoin-sv/go-sdk/script"
)

/*
Example:
{
  "satoshis": 1000100
}
*/

// PaymentRequest is the request body for the P2P payment request
type PaymentRequest struct {
	Satoshis uint64 `json:"satoshis"` // The amount, in Satoshis, that the sender intends to transfer to the receiver
}

// PaymentDestinationResponse is the response from the GetP2PPaymentDestination() request
//
// The reference is unique for the payment destination request
type PaymentDestinationResponse struct {
	StandardResponse
	PaymentDestinationPayload
}

// PaymentDestinationPayload is the payload from the response
//
// The reference is unique for the payment destination request
type PaymentDestinationPayload struct {
	Outputs   []*PaymentOutput `json:"outputs"`   // A list of outputs
	Reference string           `json:"reference"` // A reference for the payment, created by the receiver of the transaction
}

// PaymentOutput is returned inside the payment destination response
//
// There can be several outputs in one response based on the amount of satoshis being transferred and
// the rules in place by the Paymail provider
type PaymentOutput struct {
	Address  string `json:"address,omitempty"`  // Hex encoded locking script
	Satoshis uint64 `json:"satoshis,omitempty"` // Number of satoshis for that output
	Script   string `json:"script"`             // Hex encoded locking script
}

// GetP2PPaymentDestination will return list of outputs for the P2P transactions to use
//
// Specs: https://docs.moneybutton.com/docs/paymail-07-p2p-payment-destination.html
func (c *Client) GetP2PPaymentDestination(p2pURL, alias, domain string,
	paymentRequest *PaymentRequest) (response *PaymentDestinationResponse, err error) {

	// Require a valid url
	if len(p2pURL) == 0 || !strings.Contains(p2pURL, "https://") {
		err = fmt.Errorf("invalid url: %s", p2pURL)
		return
	}

	// Basic requirements for request
	if paymentRequest == nil {
		err = errors.New("paymentRequest cannot be nil")
		return
	} else if paymentRequest.Satoshis == 0 {
		err = errors.New("satoshis is required")
		return
	} else if len(alias) == 0 {
		err = errors.New("missing alias")
		return
	} else if len(domain) == 0 {
		err = errors.New("missing domain")
		return
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/api/rawtx/{alias}@{domain.tld}
	// https://<host-discovery-target>/api/p2p-payment-destination/{alias}@{domain.tld}
	reqURL := replaceAliasDomain(p2pURL, alias, domain)

	// Fire the POST request
	var resp StandardResponse
	if resp, err = c.postRequest(reqURL, paymentRequest); err != nil {
		return
	}

	// Start the response
	response = &PaymentDestinationResponse{StandardResponse: resp}

	// Test the status code
	if response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusNotModified {

		// Paymail address not found?
		if response.StatusCode == http.StatusNotFound {
			err = errors.New("paymail address not found")
		} else {
			serverError := &ServerError{}
			if err = json.Unmarshal(resp.Body, serverError); err != nil {
				return
			}
			err = fmt.Errorf(
				"bad response from paymail provider: code %d, message: %s",
				response.StatusCode, serverError.Message,
			)
		}

		return
	}

	// Decode the body of the response
	if err = json.Unmarshal(resp.Body, &response); err != nil {
		return
	}

	// Check for a reference number
	if len(response.Reference) == 0 {
		err = errors.New("missing a returned reference value")
		return
	}

	// No outputs?
	if len(response.Outputs) == 0 {
		err = errors.New("missing a returned output")
		return
	}

	// Loop all outputs
	for index, out := range response.Outputs {
		// No script returned
		if len(out.Script) == 0 {
			err = fmt.Errorf("script was missing from output: %d", index)
			return
		}

		var sc *script.Script
		sc, err = script.NewFromHex(out.Script)
		if err != nil {
			return
		}

		var addresses []string
		addresses, err = sc.Addresses()
		if err != nil || len(addresses) == 0 {
			err = errors.New("invalid output script, missing an address")
			return
		}

		response.Outputs[index].Address = addresses[0]
	}

	return
}
