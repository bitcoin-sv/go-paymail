package server

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/julienschmidt/httprouter"
)

// verifyPubKey will return a response if the pubkey matches the paymail given
//
// Specs: https://bsvalias.org/05-verify-public-key-owner.html
func (c *Configuration) verifyPubKey(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	incomingPaymail := p.ByName("paymailAddress")
	incomingPubKey := p.ByName("pubKey")

	// Parse, sanitize and basic validation
	alias, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		ErrorResponse(w, req, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest, c.Logger)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest, c.Logger)
		return
	}

	// Basic validation on pubkey
	if len(incomingPubKey) != paymail.PubKeyLength {
		ErrorResponse(w, req, ErrorInvalidPubKey, "invalid pubkey: "+incomingPubKey, http.StatusBadRequest, c.Logger)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(req, alias, domain, "")

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(req.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(w, req, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed, c.Logger)
		return
	} else if foundPaymail == nil {
		ErrorResponse(w, req, ErrorPaymailNotFound, "paymail not found: "+incomingPaymail, http.StatusBadRequest, c.Logger)
		return
	}

	verPayload := paymail.VerificationPayload{
		BsvAlias: c.BSVAliasVersion,
		Handle:   address,
		PubKey:   foundPaymail.PubKey,
		Match:    foundPaymail.PubKey == incomingPubKey,
	}

	// Set the response
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//err = json.NewEncoder(w).Encode(verPayload)
	//
	//if err != nil {
	//	ErrorResponse(w, req, ErrorEncodingResponse, err.Error(), http.StatusBadRequest, c.Logger)
	//	return
	//}
	writeJsonResponse(w, req, c.Logger, verPayload)
}
