package server

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/julienschmidt/httprouter"
)

type CapabilityEndpoint struct {
	Path    string
	Method  string
	Handler httprouter.Handle
}

type CapabilityInterface interface {
	Key() string
	Value() interface{}
}

type CallableCapability struct {
	key      string
	endpoint CapabilityEndpoint
}

func (c *CallableCapability) Key() string {
	return c.key
}

func (c *CallableCapability) Value() interface{} {
	return c.endpoint
}

type BooleanCapability struct {
	key   string
	value bool
}

func (c *BooleanCapability) Key() string {
	return c.key
}

func (c *BooleanCapability) Value() interface{} {
	return c.value
}

func MakeGenericCapabilities(c *Configuration) []CapabilityInterface {
	return []CapabilityInterface{
		&CallableCapability{
			key: paymail.BRFCPaymentDestination,
			endpoint: CapabilityEndpoint{
				Path:    "/address/{alias}@{domain.tld}",
				Method:  http.MethodPost,
				Handler: c.resolveAddress,
			},
		},
		&CallableCapability{
			key: paymail.BRFCPki,
			endpoint: CapabilityEndpoint{
				Path:    "/id/{alias}@{domain.tld}",
				Method:  http.MethodGet,
				Handler: c.showPKI,
			},
		},
		&CallableCapability{
			key: paymail.BRFCPublicProfile,
			endpoint: CapabilityEndpoint{
				Path:    "/public-profile/{alias}@{domain.tld}",
				Method:  http.MethodGet,
				Handler: c.publicProfile,
			},
		},
		&BooleanCapability{
			key:   paymail.BRFCSenderValidation,
			value: c.SenderValidationEnabled,
		},
		&CallableCapability{
			key: paymail.BRFCVerifyPublicKeyOwner,
			endpoint: CapabilityEndpoint{
				Path:    "/verify-pubkey/{alias}@{domain.tld}/{pubkey}",
				Method:  http.MethodGet,
				Handler: c.verifyPubKey,
			},
		},
	}
}

func MakeP2PCapabilities(c *Configuration) []CapabilityInterface {
	return []CapabilityInterface{
		&CallableCapability{
			key: paymail.BRFCP2PTransactions,
			endpoint: CapabilityEndpoint{
				Path:    "/receive-transaction/{alias}@{domain.tld}",
				Method:  http.MethodPost,
				Handler: c.p2pReceiveTx,
			},
		},
		&CallableCapability{
			key: paymail.BRFCP2PPaymentDestination,
			endpoint: CapabilityEndpoint{
				Path:    "/p2p-payment-destination/{alias}@{domain.tld}",
				Method:  http.MethodPost,
				Handler: c.p2pDestination,
			},
		},
	}
}

func MakeBeefCapabilities(c *Configuration) []CapabilityInterface {
	return []CapabilityInterface{
		&CallableCapability{
			key: paymail.BRFCBeefTransaction,
			endpoint: CapabilityEndpoint{
				Path:    "/beef/{alias}@{domain.tld}",
				Method:  http.MethodPost,
				Handler: c.p2pReceiveBeefTx,
			},
		},
	}
}
func generateCapabilitiesMap(array []CapabilityInterface) map[string]CapabilityInterface {
	dictionary := make(map[string]CapabilityInterface)
	for _, capability := range array {
		dictionary[capability.Key()] = capability
	}
	return dictionary
}

func extendCapabilitiesMap(base map[string]CapabilityInterface, array []CapabilityInterface) map[string]CapabilityInterface {
	newElements := generateCapabilitiesMap(array)
	for key, val := range newElements {
		base[key] = val
	}
	return base
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

// EnrichCapabilities will update the capabilities with the appropriate service url
func (c *Configuration) EnrichCapabilities(domain string) *paymail.CapabilitiesPayload {
	payload := &paymail.CapabilitiesPayload{
		BsvAlias:     c.BSVAliasVersion,
		Capabilities: make(map[string]interface{}),
	}
	for key, cap := range c.capabilities {
		switch capValue := cap.Value().(type) {
		case CapabilityEndpoint:
			payload.Capabilities[key] = GenerateServiceURL(c.Prefix, domain, c.APIVersion, c.ServiceName) + string(capValue.Path)
		default:
			payload.Capabilities[key] = capValue
		}
	}
	return payload
}
