package server

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
)

// verifyPubKey will return a response if the pubkey matches the paymail given
//
// Specs: https://bsvalias.org/05-verify-public-key-owner.html
func (c *Configuration) verifyPubKey(context *gin.Context) {
	incomingPaymail := context.Param(PaymailAddressParamName)
	incomingPubKey := context.Param(PubKeyParamName)

	// Parse, sanitize and basic validation
	alias, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		context.JSON(http.StatusBadRequest, "invalid paymail: "+incomingPaymail)
		return
	} else if !c.IsAllowedDomain(domain) {
		context.JSON(http.StatusBadRequest, "domain unknown: "+domain)
		return
	}

	// Basic validation on pubkey
	if len(incomingPubKey) != paymail.PubKeyLength {
		context.JSON(http.StatusBadRequest, "invalid pubkey: "+incomingPubKey)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(context.Request, alias, domain, "")

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		context.JSON(http.StatusExpectationFailed, err.Error())
		return
	} else if foundPaymail == nil {
		context.JSON(http.StatusNotFound, "paymail not found: "+incomingPaymail)
		return
	}

	verPayload := paymail.VerificationPayload{
		BsvAlias: c.BSVAliasVersion,
		Handle:   address,
		PubKey:   foundPaymail.PubKey,
		Match:    foundPaymail.PubKey == incomingPubKey,
	}

	context.JSON(http.StatusOK, verPayload)
}
