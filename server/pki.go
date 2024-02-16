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
		context.JSON(http.StatusBadRequest, "invalid paymail: "+incomingPaymail)
		return
	} else if !c.IsAllowedDomain(domain) {
		context.JSON(http.StatusBadRequest, "domain unknown: "+domain)
		return
	}

	md := CreateMetadata(context.Request, alias, domain, "")

	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		context.JSON(http.StatusExpectationFailed, err.Error())
		return
	} else if foundPaymail == nil {
		context.JSON(http.StatusNotFound, "paymail not found: "+incomingPaymail)
		return
	}

	pkiPayload := paymail.PKIPayload{
		BsvAlias: c.BSVAliasVersion,
		Handle:   address,
		PubKey:   foundPaymail.PubKey,
	}

	context.JSON(http.StatusOK, pkiPayload)
}
