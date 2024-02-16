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
		return
	} else if !c.IsAllowedDomain(domain) {
		context.JSON(http.StatusBadRequest, "domain unknown: "+domain)
		return
	}

	var senderRequest paymail.SenderRequest
	err := context.Bind(&senderRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// Check for required fields
	if len(senderRequest.SenderHandle) == 0 {
		context.JSON(http.StatusBadRequest, "senderHandle is empty")
		return
	} else if len(senderRequest.Dt) == 0 {
		context.JSON(http.StatusBadRequest, "dt is empty")
		return
	}

	// Validate the timestamp
	if err = paymail.ValidateTimestamp(senderRequest.Dt); err != nil {
		context.JSON(http.StatusBadRequest, "invalid dt: "+err.Error())
		return
	}

	// Basic validation on sender handle
	if err = paymail.ValidatePaymail(senderRequest.SenderHandle); err != nil {
		context.JSON(http.StatusBadRequest, "invalid senderHandle: "+err.Error())
		return
	}

	// Only validate signatures if sender validation is enabled (skip if disabled)
	if c.SenderValidationEnabled {
		if len(senderRequest.Signature) > 0 {

			// Get the pubKey from the corresponding sender paymail address
			var senderPubKey *bec.PublicKey
			senderPubKey, err = getSenderPubKey(senderRequest.SenderHandle)
			if err != nil {
				context.JSON(http.StatusBadRequest, "invalid senderHandle: "+err.Error())
				return
			}

			// Derive address from pubKey
			var rawAddress *bscript.Address
			if rawAddress, err = bitcoin.GetAddressFromPubKey(senderPubKey, true); err != nil {
				context.JSON(http.StatusBadRequest, "invalid senderHandle: "+err.Error())
				return
			}

			// Verify the signature
			if err = senderRequest.Verify(rawAddress.AddressString, senderRequest.Signature); err != nil {
				context.JSON(http.StatusBadRequest, "invalid signature: "+err.Error())
				return
			}
		} else {
			context.JSON(http.StatusBadRequest, "missing required signature")
			return
		}
	}

	// Create the metadata struct
	md := CreateMetadata(context.Request, alias, domain, "")
	md.ResolveAddress = &senderRequest

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		context.JSON(http.StatusExpectationFailed, err.Error())
		return
	} else if foundPaymail == nil {
		context.JSON(http.StatusNotFound, "paymail not found: "+incomingPaymail)
		return
	}

	// Get the resolution information
	var response *paymail.ResolutionPayload
	if response, err = c.actions.CreateAddressResolutionResponse(
		context.Request.Context(), alias, domain, c.SenderValidationEnabled, md,
	); err != nil {
		context.JSON(http.StatusExpectationFailed, err.Error())
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
