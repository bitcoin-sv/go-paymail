package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tonicpow/go-paymail"
)

// verifyPubKey will return a response if the pubkey matches the paymail given
//
// Specs: https://bsvalias.org/05-verify-public-key-owner.html
func (c *Configuration) verifyPubKey(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params submitted via URL request
	params := req.URL.Query()
	incomingPaymail := params.Get("paymailAddress")
	incomingPubKey := params.Get("pubKey")

	// Parse, sanitize and basic validation
	alias, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		ErrorResponse(w, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Basic validation on pubkey
	if len(incomingPubKey) != paymail.PubKeyLength {
		ErrorResponse(w, ErrorInvalidPubKey, "invalid pubkey: "+incomingPubKey, http.StatusBadRequest)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(req, alias, domain, "")

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(req.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(w, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed)
		return
	} else if foundPaymail == nil {
		ErrorResponse(w, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	verPayload := paymail.VerificationPayload{
		BsvAlias: c.BSVAliasVersion,
		Handle:   address,
		PubKey:   foundPaymail.PubKey,
		Match:    foundPaymail.PubKey == incomingPubKey,
	}

	// Set the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(verPayload)

	if err != nil {
		ErrorResponse(w, ErrorEncodingResponse, err.Error(), http.StatusInternalServerError)
		return
	}
}
