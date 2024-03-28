package server

import (
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
		ErrorResponse(context, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(context, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(context.Request, alias, domain, "")

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(context, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed)
		return
	} else if foundPaymail == nil {
		ErrorResponse(context, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	payload := paymail.PublicProfilePayload{
		Avatar: foundPaymail.Avatar,
		Name:   foundPaymail.Name,
	}

	// Set the response
	context.JSON(http.StatusOK, payload)
}
