package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tonicpow/go-paymail"
)

/*
Incoming Data Object Example:
{
  "satoshis": 1000100,
}
*/

// p2pDestination will return an output script(s) for a destination (used with SendP2PTransaction)
//
// Specs: https://docs.moneybutton.com/docs/paymail-07-p2p-payment-destination.html
func (c *Configuration) p2pDestination(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params & paymail address submitted via URL request
	params := req.URL.Query()
	incomingPaymail := params.Get("paymailAddress")

	// Parse, sanitize and basic validation
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		ErrorResponse(w, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	satoshis, err := strconv.ParseUint(params.Get("satoshis"), 10, 64)
	if err != nil {
		ErrorResponse(w, ErrorInvalidParameter, "invalid satoshis: "+params.Get("satoshis"), http.StatusBadRequest)
		return
	}

	// Start the PaymentRequest
	paymentRequest := &paymail.PaymentRequest{
		Satoshis: satoshis,
	}

	// Did we get some satoshis?
	if paymentRequest.Satoshis == 0 {
		ErrorResponse(w, ErrorMissingSatoshis, "missing parameter: satoshis", http.StatusBadRequest)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(req, alias, domain, "")
	md.PaymentDestination = paymentRequest

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(req.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(w, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed)
		return
	} else if foundPaymail == nil {
		ErrorResponse(w, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	// Create the response
	var response *paymail.PaymentDestinationPayload
	if response, err = c.actions.CreateP2PDestinationResponse(
		req.Context(), alias, domain, paymentRequest.Satoshis, md,
	); err != nil {
		ErrorResponse(w, ErrorScript, "error creating output script(s): "+err.Error(), http.StatusExpectationFailed)
		return
	}

	// Set the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		ErrorResponse(w, ErrorEncodingResponse, err.Error(), http.StatusInternalServerError)
		return
	}
}
