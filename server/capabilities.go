package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tonicpow/go-paymail"
)

// GenericCapabilities will make generic capabilities
func GenericCapabilities(bsvAliasVersion string, senderValidation bool) *paymail.CapabilitiesPayload {
	return &paymail.CapabilitiesPayload{
		BsvAlias: bsvAliasVersion,
		Capabilities: map[string]interface{}{
			paymail.BRFCPaymentDestination:   "/address/{alias}@{domain.tld}",
			paymail.BRFCPki:                  "/id/{alias}@{domain.tld}",
			paymail.BRFCPublicProfile:        "/public-profile/{alias}@{domain.tld}",
			paymail.BRFCSenderValidation:     senderValidation,
			paymail.BRFCVerifyPublicKeyOwner: "/verify-pubkey/{alias}@{domain.tld}/{pubkey}",
		},
	}
}

// P2PCapabilities will make generic capabilities & add additional p2p capabilities
func P2PCapabilities(bsvAliasVersion string, senderValidation bool) *paymail.CapabilitiesPayload {
	c := GenericCapabilities(bsvAliasVersion, senderValidation)
	c.Capabilities[paymail.BRFCP2PTransactions] = "/receive-transaction/{alias}@{domain.tld}"
	c.Capabilities[paymail.BRFCP2PPaymentDestination] = "/p2p-payment-destination/{alias}@{domain.tld}"
	return c
}

// showCapabilities will return the service discovery results for the server
// and list all active capabilities of the Paymail server
//
// Specs: http://bsvalias.org/02-02-capability-discovery.html
func (c *Configuration) showCapabilities(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Check the domain (allowed, and used for capabilities response)
	// todo: bake this into middleware? This is protecting the "req" domain name (like CORs)
	domain := getHost(req)
	if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Set the service URL
	capabilities := c.EnrichCapabilities(domain)
	jsonData, err := json.Marshal(capabilities)
	if err != nil {
		ErrorResponse(w, ErrorFailedMarshalJSON, "failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}
