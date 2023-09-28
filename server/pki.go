package server

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/julienschmidt/httprouter"
)

// showPKI will return the public key information for the corresponding paymail address
//
// Specs: http://bsvalias.org/03-public-key-infrastructure.html
func (c *Configuration) showPKI(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	incomingPaymail := p.ByName("paymailAddress")

	// Parse, sanitize and basic validation
	alias, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		ErrorResponse(w, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
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

	pkiPayload := paymail.PKIPayload{
		BsvAlias: c.BSVAliasVersion,
		Handle:   address,
		PubKey:   foundPaymail.PubKey,
	}

	// Set the response
	writeJsonResponse(w, http.StatusOK, pkiPayload)
}
