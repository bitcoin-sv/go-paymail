package server

import (
	"encoding/json"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/julienschmidt/httprouter"
)

type parseError struct {
	code, msg string
}

func parseP2pReceiveTxRequest(c *Configuration, req *http.Request, params httprouter.Params, format p2pPayloadFormat) (*p2pReceiveTxReqPayload, *parseError) {
	incomingPaymail := params.ByName("paymailAddress")

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

	var p2pTransaction paymail.P2PTransaction
	err := json.NewDecoder(req.Body).Decode(&p2pTransaction)
	if err != nil {
		return nil, &parseError{ErrorInvalidParameter, "invalid request"}
	}
	if len(p2pTransaction.Reference) == 0 {
		return nil, &parseError{ErrorMissingField, "missing parameter: reference"}
	}
	if format == basicP2pPayload {
		if len(p2pTransaction.Hex) == 0 {
			return nil, &parseError{ErrorMissingField, "missing parameter: hex"}
		}
	} else if format == beefP2pPayload {
		if len(p2pTransaction.Beef) == 0 {
			return nil, &parseError{ErrorMissingField, "missing parameter: beef"}
		}
	}
	vErr := validateMetadata(c, p2pTransaction.MetaData)

	if vErr != nil {
		return nil, vErr
	}

	requestData.P2PTransaction = &p2pTransaction
	return &requestData, nil
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
