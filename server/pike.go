package server

import (
	"encoding/json"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/gin-gonic/gin"
)

func (c *Configuration) pikeNewContact(rc *gin.Context) {
	receiverPaymail := rc.Param(PaymailAddressParamName)

	var requesterContact paymail.PikeContactRequestPayload
	err := json.NewDecoder(rc.Request.Body).Decode(&requesterContact)
	if err != nil {
		ErrorResponse(rc, ErrorInvalidParameter, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.pikeContactActions.AddContact(rc.Request.Context(), receiverPaymail, &requesterContact); err != nil {
		ErrorResponse(rc, ErrorAddContactRequest, err.Error(), http.StatusExpectationFailed)
		return
	}

	rc.Status(http.StatusCreated)
}

func (c *Configuration) pikeGetPaymentDestinations(rc *gin.Context) {
	var paymentDestinationRequest paymail.PikePaymentOutputsPayload
	err := json.NewDecoder(rc.Request.Body).Decode(&paymentDestinationRequest)
	defer func() {
		_ = rc.Request.Body.Close()
	}()
	if err != nil {
		ErrorResponse(rc, ErrorInvalidParameter, err.Error(), http.StatusBadRequest)
		return
	}

	alias, domain, md, ok := c.GetPaymailAndCreateMetadata(rc, paymentDestinationRequest.Amount)
	if !ok {
		return
	}

	var response *paymail.PikePaymentOutputsResponse
	if response, err = c.pikePaymentActions.CreatePikeDestinationResponse(
		rc.Request.Context(), alias, domain, paymentDestinationRequest.Amount, md,
	); err != nil {
		ErrorResponse(rc, ErrorScript, "error creating output script(s): "+err.Error(), http.StatusExpectationFailed)
		return
	}

	rc.JSON(http.StatusOK, response)
}
