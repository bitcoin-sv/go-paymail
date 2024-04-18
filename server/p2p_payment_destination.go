package server

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
)

/*
Incoming Data Object Example:

	{
	  "satoshis": 1000100,
	}
*/
type p2pDestinationRequestBody struct {
	Satoshis uint64 `json:"satoshis,omitempty"`
}

// p2pDestination will return an output script(s) for a destination (used with SendP2PTransaction)
//
// Specs: https://docs.moneybutton.com/docs/paymail-07-p2p-payment-destination.html
func (c *Configuration) p2pDestination(context *gin.Context) {
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
	var b p2pDestinationRequestBody
	err := context.Bind(&b)
	if err != nil {
		ErrorResponse(context, ErrorInvalidParameter, "error decoding body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Start the PaymentRequest
	paymentRequest := &paymail.PaymentRequest{
		Satoshis: b.Satoshis,
	}

	// Did we get some satoshis?
	if paymentRequest.Satoshis == 0 {
		ErrorResponse(context, ErrorMissingField, "missing parameter: satoshis", http.StatusBadRequest)
		return
	}

	// Create the metadata struct
	md := CreateMetadata(context.Request, alias, domain, "")
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

	// Create the response
	var response *paymail.PaymentDestinationPayload
	if response, err = c.actions.CreateP2PDestinationResponse(
		context.Request.Context(), alias, domain, paymentRequest.Satoshis, md,
	); err != nil {
		ErrorResponse(context, ErrorScript, "error creating output script(s): "+err.Error(), http.StatusExpectationFailed)
		return
	}

	context.JSON(http.StatusOK, response)
}
