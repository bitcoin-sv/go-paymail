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

	if err = c.pikeActions.AddContact(rc.Request.Context(), receiverPaymail, &requesterContact); err != nil {
		ErrorResponse(rc, ErrorAddContactRequest, err.Error(), http.StatusExpectationFailed)
		return
	}

	rc.Status(http.StatusCreated)
}
