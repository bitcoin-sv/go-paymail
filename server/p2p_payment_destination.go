package server

import (
	"encoding/json"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/julienschmidt/httprouter"
)

/*
Incoming Data Object Example:

	{
	  "satoshis": 1000100,
	}
*/
type p2pDestinationRequestBody struct {
	Satoshis uint64 `json:"satoshis,omitempty"`
}

// p2pDestination will return an output script(s) for a destination (used with SendP2PTransaction)
//
// Specs: https://docs.moneybutton.com/docs/paymail-07-p2p-payment-destination.html
func (c *Configuration) p2pDestination(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	incomingPaymail := p.ByName("paymailAddress")

	// Parse, sanitize and basic validation
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		ErrorResponse(w, req, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest, c.Logger)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest, c.Logger)
		return
	}
	var b p2pDestinationRequestBody
	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil {
		ErrorResponse(w, req, ErrorInvalidParameter, "error decoding body: "+err.Error(), http.StatusBadRequest, c.Logger)
		return
	}

	// Start the PaymentRequest
	paymentRequest := &paymail.PaymentRequest{
		Satoshis: b.Satoshis,
	}

	// Did we get some satoshis?
	if paymentRequest.Satoshis == 0 {
		ErrorResponse(w, req, ErrorMissingField, "missing parameter: satoshis", http.StatusBadRequest, c.Logger)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(req, alias, domain, "")
	md.PaymentDestination = paymentRequest

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(req.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(w, req, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed, c.Logger)
		return
	} else if foundPaymail == nil {
		ErrorResponse(w, req, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound, c.Logger)
		return
	}

	// Create the response
	var response *paymail.PaymentDestinationPayload
	if response, err = c.actions.CreateP2PDestinationResponse(
		req.Context(), alias, domain, paymentRequest.Satoshis, md,
	); err != nil {
		ErrorResponse(w, req, ErrorScript, "error creating output script(s): "+err.Error(), http.StatusExpectationFailed, c.Logger)
		return
	}

	// Set the response
	writeJsonResponse(w, req, c.Logger, response)
}
