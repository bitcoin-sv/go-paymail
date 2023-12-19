package server

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/julienschmidt/httprouter"
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

// BeefCapabilities will add beef capabilities to given ones
func BeefCapabilities(c *paymail.CapabilitiesPayload) *paymail.CapabilitiesPayload {
	c.Capabilities[paymail.BRFCBeefTransaction] = "/beef/{alias}@{domain.tld}"
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
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest, c.Logger)
		return
	}

	// Set the service URL
	capabilities := c.EnrichCapabilities(domain)
	writeJsonResponse(w, req, c.Logger, capabilities)
}
