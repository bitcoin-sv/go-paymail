package server

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bt/v2/bscript"
)

/*
Incoming Data Object Example:
{
    "senderName": "UserName",
    "senderHandle": "alias@domain.com",
    "dt": "2020-04-09T16:08:06.419Z",
    "amount": 551,
    "purpose": "message to receiver",
	"signature": "SIGNATURE-IF-REQUIRED-IN-CONFIG"
}
*/

// resolveAddress will return the payment destination (bitcoin address) for the corresponding paymail address
//
// Specs: http://bsvalias.org/04-01-basic-address-resolution.html
func (c *Configuration) resolveAddress(context *gin.Context) {
	incomingPaymail := context.Param(PaymailAddressParamName)

	// Parse, sanitize and basic validation
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		context.JSON(http.StatusBadRequest, "invalid paymail: "+incomingPaymail)
		ErrorResponse(context, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(context, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	var senderRequest paymail.SenderRequest
	err := context.Bind(&senderRequest)
	if err != nil {
		ErrorResponse(context, ErrorInvalidParameter, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check for required fields
	if len(senderRequest.SenderHandle) == 0 {
		ErrorResponse(context, ErrorInvalidSenderHandle, "senderHandle is empty", http.StatusBadRequest)
		return
	} else if len(senderRequest.Dt) == 0 {
		ErrorResponse(context, ErrorInvalidDt, "dt is empty", http.StatusBadRequest)
		return
	}

	// Validate the timestamp
	if err = paymail.ValidateTimestamp(senderRequest.Dt); err != nil {
		ErrorResponse(context, ErrorInvalidDt, "invalid dt: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation on sender handle
	if err = paymail.ValidatePaymail(senderRequest.SenderHandle); err != nil {
		ErrorResponse(context, ErrorInvalidSenderHandle, "invalid senderHandle: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Only validate signatures if sender validation is enabled (skip if disabled)
	if c.SenderValidationEnabled {
		if len(senderRequest.Signature) > 0 {

			// Get the pubKey from the corresponding sender paymail address
			var senderPubKey *bec.PublicKey
			senderPubKey, err = getSenderPubKey(senderRequest.SenderHandle)
			if err != nil {
				ErrorResponse(context, ErrorInvalidSenderHandle, "invalid senderHandle: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Derive address from pubKey
			var rawAddress *bscript.Address
			if rawAddress, err = bitcoin.GetAddressFromPubKey(senderPubKey, true); err != nil {
				ErrorResponse(context, ErrorInvalidSenderHandle, "invalid senderHandle: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Verify the signature
			if err = senderRequest.Verify(rawAddress.AddressString, senderRequest.Signature); err != nil {
				ErrorResponse(context, ErrorInvalidSignature, "invalid signature: "+err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			ErrorResponse(context, ErrorInvalidSignature, "missing required signature", http.StatusBadRequest)
			return
		}
	}

	// Create the metadata struct
	md := CreateMetadata(context.Request, alias, domain, "")
	md.ResolveAddress = &senderRequest

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(context, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed)
		return
	} else if foundPaymail == nil {
		ErrorResponse(context, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	// Get the resolution information
	var response *paymail.ResolutionPayload
	if response, err = c.actions.CreateAddressResolutionResponse(
		context.Request.Context(), alias, domain, c.SenderValidationEnabled, md,
	); err != nil {
		ErrorResponse(context, ErrorScript, "error creating output script: "+err.Error(), http.StatusExpectationFailed)
		return
	}

	// Set the response
	context.JSON(http.StatusOK, response)
}

// getSenderPubKey will fetch the pubKey from a PKI request for the sender handle
func getSenderPubKey(senderPaymailAddress string) (*bec.PublicKey, error) {

	// Sanitize and break apart
	alias, domain, _ := paymail.SanitizePaymail(senderPaymailAddress)

	// Load the client
	client, err := paymail.NewClient(paymail.WithHTTPTimeout(15 * time.Second))
	if err != nil {
		return nil, err
	}

	// Get the SRV record
	var srv *net.SRV
	if srv, err = client.GetSRVRecord(
		paymail.DefaultServiceName, paymail.DefaultProtocol, domain,
	); err != nil {
		return nil, err
	}

	// Get the capabilities
	// This is required first to get the corresponding PKI endpoint url
	var capabilities *paymail.CapabilitiesResponse
	if capabilities, err = client.GetCapabilities(
		srv.Target, paymail.DefaultPort,
	); err != nil {
		return nil, err
	}

	// Extract the PKI URL from the capabilities response
	pkiURL := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)

	// Get the actual PKI
	var pki *paymail.PKIResponse
	if pki, err = client.GetPKI(
		pkiURL, alias, domain,
	); err != nil {
		return nil, err
	}

	// Convert the string pubKey to a bec.PubKey
	return bitcoin.PubKeyFromString(pki.PubKey)
}
