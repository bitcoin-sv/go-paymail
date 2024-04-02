package server

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
)

// showPKI will return the public key information for the corresponding paymail address
//
// Specs: http://bsvalias.org/03-public-key-infrastructure.html
func (c *Configuration) showPKI(context *gin.Context) {
	incomingPaymail := context.Param(PaymailAddressParamName)

	alias, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		ErrorResponse(context, ErrorInvalidParameter, "domain unknown: "+domain, http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(context, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	md := CreateMetadata(context.Request, alias, domain, "")

	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(context, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed)
		return
	} else if foundPaymail == nil {
		ErrorResponse(context, ErrorPaymailNotFound, "paymail not found: "+incomingPaymail, http.StatusBadRequest)
		return
	}

	pkiPayload := paymail.PKIPayload{
		BsvAlias: c.BSVAliasVersion,
		Handle:   address,
		PubKey:   foundPaymail.PubKey,
	}

	context.JSON(http.StatusOK, pkiPayload)
}
