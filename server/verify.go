package server

import (
	"github.com/bitcoin-sv/go-paymail/errors"
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
		errors.ErrorResponse(context, errors.ErrInvalidPaymail)
		return
	} else if !c.IsAllowedDomain(domain) {
		errors.ErrorResponse(context, errors.ErrDomainUnknown)
		return
	}

	// Basic validation on pubkey
	if len(incomingPubKey) != paymail.PubKeyLength {
		errors.ErrorResponse(context, errors.ErrInvalidPubKey)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(context.Request, alias, domain, "")

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		errors.ErrorResponse(context, err)
		return
	} else if foundPaymail == nil {
		errors.ErrorResponse(context, errors.ErrCouldNotFindPaymail)
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
