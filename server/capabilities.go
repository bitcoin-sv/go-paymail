package server

import (
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/julienschmidt/httprouter"
)

type CallableCapability struct {
	Path    string
	Method  string
	Handler httprouter.Handle
}

type CallableCapabilitiesMap map[string]CallableCapability
type StaticCapabilitiesMap map[string]any

func (c *Configuration) SetGenericCapabilities() {
	_addCapabilities(c.callableCapabilities,
		CallableCapabilitiesMap{
			paymail.BRFCPaymentDestination: CallableCapability{
				Path:    fmt.Sprintf("/address/%s", PaymailAddressTemplate),
				Method:  http.MethodPost,
				Handler: c.resolveAddress,
			},
			paymail.BRFCPki: CallableCapability{
				Path:    fmt.Sprintf("/id/%s", PaymailAddressTemplate),
				Method:  http.MethodGet,
				Handler: c.showPKI,
			},
			paymail.BRFCPublicProfile: CallableCapability{
				Path:    fmt.Sprintf("/public-profile/%s", PaymailAddressTemplate),
				Method:  http.MethodGet,
				Handler: c.publicProfile,
			},
			paymail.BRFCVerifyPublicKeyOwner: CallableCapability{
				Path:    fmt.Sprintf("/verify-pubkey/%s/%s", PaymailAddressTemplate, PubKeyTemplate),
				Method:  http.MethodGet,
				Handler: c.verifyPubKey,
			},
		},
	)
	_addCapabilities(c.staticCapabilities,
		StaticCapabilitiesMap{
			paymail.BRFCSenderValidation: c.SenderValidationEnabled,
		},
	)
}

func (c *Configuration) SetP2PCapabilities() {
	_addCapabilities(c.callableCapabilities,
		CallableCapabilitiesMap{
			paymail.BRFCP2PTransactions: CallableCapability{
				Path:    fmt.Sprintf("/receive-transaction/%s", PaymailAddressTemplate),
				Method:  http.MethodPost,
				Handler: c.p2pReceiveTx,
			},
			paymail.BRFCP2PPaymentDestination: CallableCapability{
				Path:    fmt.Sprintf("/p2p-payment-destination/%s", PaymailAddressTemplate),
				Method:  http.MethodPost,
				Handler: c.p2pDestination,
			},
		},
	)
}

func (c *Configuration) SetBeefCapabilities() {
	_addCapabilities(c.callableCapabilities,
		CallableCapabilitiesMap{
			paymail.BRFCBeefTransaction: CallableCapability{
				Path:    fmt.Sprintf("/beef/%s", PaymailAddressTemplate),
				Method:  http.MethodPost,
				Handler: c.p2pReceiveBeefTx,
			},
		},
	)
}

func _addCapabilities[T any](base map[string]T, newCaps map[string]T) {
	for key, val := range newCaps {
		base[key] = val
	}
}

// showCapabilities will return the service discovery results for the server
// and list all active capabilities of the Paymail server
//
// Specs: http://bsvalias.org/02-02-capability-discovery.html
func (c *Configuration) showCapabilities(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Check the host (allowed, and used for capabilities response)
	// todo: bake this into middleware? This is protecting the "req" host name (like CORs)
	host := ""
	if req.URL.IsAbs() || len(req.URL.Host) == 0 {
		host = req.Host
	} else {
		host = req.URL.Host
	}

	if !c.IsAllowedDomain(host) {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+host, http.StatusBadRequest, c.Logger)
		return
	}

	capabilities, err := c.EnrichCapabilities(host)
	if err != nil {
		ErrorResponse(w, req, ErrorEncodingResponse, err.Error(), http.StatusBadRequest, c.Logger)
		return
	}

	writeJsonResponse(w, req, c.Logger, capabilities)
}

// EnrichCapabilities will update the capabilities with the appropriate service url
func (c *Configuration) EnrichCapabilities(host string) (*paymail.CapabilitiesPayload, error) {
	serviceUrl, err := GenerateServiceURL(c.Prefix, host, c.APIVersion, c.ServiceName)
	if err != nil {
		return nil, err
	}
	payload := &paymail.CapabilitiesPayload{
		BsvAlias:     c.BSVAliasVersion,
		Capabilities: make(map[string]interface{}),
	}
	for key, cap := range c.staticCapabilities {
		payload.Capabilities[key] = cap
	}
	for key, cap := range c.callableCapabilities {
		payload.Capabilities[key] = serviceUrl + string(cap.Path)
	}
	return payload, nil
}
