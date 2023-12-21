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
		ErrorResponse(w, req, ErrorInvalidParameter, "domain unknown: "+domain, http.StatusBadRequest, c.Logger)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest, c.Logger)
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

	pkiPayload := paymail.PKIPayload{
		BsvAlias: c.BSVAliasVersion,
		Handle:   address,
		PubKey:   foundPaymail.PubKey,
	}

	// Set the response
	writeJsonResponse(w, req, c.Logger, pkiPayload)
}
