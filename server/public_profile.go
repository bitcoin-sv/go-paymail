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
		context.JSON(http.StatusBadRequest, "invalid paymail: "+incomingPaymail)
		return
	} else if !c.IsAllowedDomain(domain) {
		context.JSON(http.StatusBadRequest, "domain unknown: "+domain)
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
		context.JSON(http.StatusNotFound, "paymail not found")
		return
	}

	payload := paymail.PublicProfilePayload{
		Avatar: foundPaymail.Avatar,
		Name:   foundPaymail.Name,
	}

	// Set the response
	context.JSON(http.StatusOK, payload)
}
