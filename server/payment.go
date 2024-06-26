package server

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/errors"
	"github.com/gin-gonic/gin"
)

// GetPaymailAndCreateMetadata is a helper function to get the paymail from the request, check it in database and create the metadata based on that.
func (c *Configuration) GetPaymailAndCreateMetadata(context *gin.Context, satoshis uint64) (alias, domain string, md *RequestMetadata, ok bool) {
	incomingPaymail := context.Param(PaymailAddressParamName)

	// Parse, sanitize and basic validation
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		errors.ErrorResponse(context, errors.ErrInvalidPaymail)
		return
	}
	if !c.IsAllowedDomain(domain) {
		errors.ErrorResponse(context, errors.ErrDomainUnknown)
		return
	}

	// Start the PaymentRequest
	paymentRequest := &paymail.PaymentRequest{
		Satoshis: satoshis,
	}

	// Did we get some satoshis?
	if paymentRequest.Satoshis == 0 {
		errors.ErrorResponse(context, errors.ErrMissingFieldSatoshis)
		return
	}

	// Create the metadata struct
	md = CreateMetadata(context.Request, alias, domain, "")
	md.PaymentDestination = paymentRequest

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		errors.ErrorResponse(context, err)
		return
	}
	if foundPaymail == nil {
		errors.ErrorResponse(context, errors.ErrCouldNotFindPaymail)
		return
	}

	ok = true
	return
}
