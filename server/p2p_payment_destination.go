package server

import (
	"github.com/bitcoin-sv/go-paymail/errors"
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
	var b p2pDestinationRequestBody
	err := context.Bind(&b)
	if err != nil {
		errors.ErrorResponse(context, errors.ErrCannotBindRequest)
		return
	}

	alias, domain, md, ok := c.GetPaymailAndCreateMetadata(context, b.Satoshis)
	if !ok {
		// ErrorResponse already set up in GetPaymailAndCreateMetadata
		return
	}

	var response *paymail.PaymentDestinationPayload
	if response, err = c.actions.CreateP2PDestinationResponse(
		context.Request.Context(), alias, domain, b.Satoshis, md,
	); err != nil {
		errors.ErrorResponse(context, err)
		return
	}

	context.JSON(http.StatusOK, response)
}
