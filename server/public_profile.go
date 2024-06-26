package server

import (
	"github.com/bitcoin-sv/go-paymail/errors"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
)

// publicProfile will return the public profile for the corresponding paymail address
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-paymail/pull/7/files
func (c *Configuration) publicProfile(context *gin.Context) {
	incomingPaymail := context.Param(PaymailAddressParamName)

	// Parse, sanitize and basic validation
	alias, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		errors.ErrorResponse(context, errors.ErrInvalidPaymail)
		return
	} else if !c.IsAllowedDomain(domain) {
		errors.ErrorResponse(context, errors.ErrDomainUnknown)
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

	payload := paymail.PublicProfilePayload{
		Avatar: foundPaymail.Avatar,
		Name:   foundPaymail.Name,
	}

	// Set the response
	context.JSON(http.StatusOK, payload)
}
