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
		context.JSON(http.StatusBadRequest, "invalid paymail: "+incomingPaymail)
		return
	} else if !c.IsAllowedDomain(domain) {
		context.JSON(http.StatusBadRequest, "domain unknown: "+domain)
		return
	}
	var b p2pDestinationRequestBody
	err := context.Bind(&b)
	if err != nil {
		context.JSON(http.StatusBadRequest, "error decoding body: "+err.Error())
		return
	}

	// Start the PaymentRequest
	paymentRequest := &paymail.PaymentRequest{
		Satoshis: b.Satoshis,
	}

	// Did we get some satoshis?
	if paymentRequest.Satoshis == 0 {
		context.JSON(http.StatusBadRequest, "missing parameter: satoshis")
		return
	}

	// Create the metadata struct
	md := CreateMetadata(context.Request, alias, domain, "")
	md.PaymentDestination = paymentRequest

	// Get from the data layer
	foundPaymail, err := c.actions.GetPaymailByAlias(context.Request.Context(), alias, domain, md)
	if err != nil {
		context.JSON(http.StatusExpectationFailed, err.Error())
		return
	} else if foundPaymail == nil {
		context.JSON(http.StatusNotFound, "paymail not found: "+incomingPaymail)
		return
	}

	// Create the response
	var response *paymail.PaymentDestinationPayload
	if response, err = c.actions.CreateP2PDestinationResponse(
		context.Request.Context(), alias, domain, paymentRequest.Satoshis, md,
	); err != nil {
		context.JSON(http.StatusExpectationFailed, "error creating output script(s): "+err.Error())
		return
	}

	context.JSON(http.StatusOK, response)
}
