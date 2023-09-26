package server

import (
	"encoding/json"
	"net/url"

	"github.com/bitcoin-sv/go-paymail"
)

type parseError struct {
	code, msg string
}

func parseP2pReceiveTxRequest(c *Configuration, parms url.Values, format p2pPayloadFormat) (*p2pReceiveTxReqPayload, *parseError) {
	incomingPaymail := parms.Get("paymailAddress")

	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		return nil, &parseError{ErrorInvalidParameter, "invalid paymail: " + incomingPaymail}

	} else if !c.IsAllowedDomain(domain) {
		return nil, &parseError{ErrorUnknownDomain, "domain unknown: " + domain}
	}

	requestData := p2pReceiveTxReqPayload{
		incomingPaymailAlias:  alias,
		incomingPaymailDomain: domain,
	}

	requestTx := paymail.P2PTransaction{}

	requestTx.Reference = parms.Get("reference")
	if len(requestTx.Reference) == 0 {
		return nil, &parseError{ErrorMissingField, "missing parameter: reference"}
	}

	if format == basicP2pPayload {
		requestTx.Hex = parms.Get("hex")
		if len(requestTx.Hex) == 0 {
			return nil, &parseError{ErrorMissingField, "missing parameter: hex"}
		}
	} else if format == beefP2pPayload {
		requestTx.Beef = parms.Get("beef")
		if len(requestTx.Beef) == 0 {
			return nil, &parseError{ErrorMissingField, "missing parameter: beef"}
		}
	}

	requestTx.MetaData = parseMetadata(parms.Get("metadata"))
	vErr := validateMetadata(c, requestTx.MetaData)

	if vErr != nil {
		return nil, vErr
	}

	requestData.P2PTransaction = &requestTx
	return &requestData, nil
}

func parseMetadata(metadata string) *paymail.P2PMetaData {
	result := paymail.P2PMetaData{}

	if len(metadata) > 0 {
		_ = json.Unmarshal([]byte(metadata), &result) // ignore metadata deserialization errors
	}

	return &result
}

func validateMetadata(c *Configuration, metadata *paymail.P2PMetaData) *parseError {
	// Check signature if: 1) sender validation enabled or 2) a signature was given (optional)
	if c.SenderValidationEnabled || len(metadata.Signature) > 0 {

		// Check required fields for signature validation
		if len(metadata.Signature) == 0 {
			return &parseError{ErrorMissingField, "missing parameter: signature"}
		}

		if len(metadata.PubKey) == 0 {
			return &parseError{ErrorMissingField, "missing parameter: pubkey"}
		}
	}

	return nil
}
