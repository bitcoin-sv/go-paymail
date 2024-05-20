package server

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Configuration) GetPaymailAndCreateMetadata(context *gin.Context, satoshis uint64) (alias, domain string, md *RequestMetadata, ok bool) {
	incomingPaymail := context.Param(PaymailAddressParamName)

	// Parse, sanitize and basic validation
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		ErrorResponse(context, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(context, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Start the PaymentRequest
	paymentRequest := &paymail.PaymentRequest{
		Satoshis: satoshis,
	}

	// Did we get some satoshis?
	if paymentRequest.Satoshis == 0 {
		ErrorResponse(context, ErrorMissingField, "missing parameter: satoshis", http.StatusBadRequest)
		return
	}

	// Create the metadata struct
	md = CreateMetadata(context.Request, alias, domain, "")
	md.PaymentDestination = paymentRequest

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(context, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed)
		return
	} else if foundPaymail == nil {
		ErrorResponse(context, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	ok = true
	return
}
