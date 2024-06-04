package server

import (
	"encoding/json"
	"fmt"
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

func (c *Configuration) pikeGetOutputTemplates(rc *gin.Context) {
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
		// ErrorResponse already set up in GetPaymailAndCreateMetadata
		return
	}

	pki, err := getPKI(paymentDestinationRequest.SenderPaymail)
	if err != nil {
		ErrorResponse(rc, ErrorScript, "error getting pki: "+err.Error(), http.StatusExpectationFailed)
		return
	}

	var response *paymail.PikePaymentOutputsResponse
	if response, err = c.pikePaymentActions.CreatePikeOutputResponse(
		rc.Request.Context(), alias, domain, pki.PubKey, paymentDestinationRequest.Amount, md,
	); err != nil {
		ErrorResponse(rc, ErrorScript, "error creating output script(s): "+err.Error(), http.StatusExpectationFailed)
		return
	}

	rc.JSON(http.StatusOK, response)
}

func getPKI(paymailAddress string) (*paymail.PKIResponse, error) {
	alias, domain, paymailAddress := paymail.SanitizePaymail(paymailAddress)
	if len(paymailAddress) == 0 {
		return nil, fmt.Errorf("invalid paymail: %s", paymailAddress)
	}

	client, err := paymail.NewClient()
	if err != nil {
		return nil, err
	}

	var capabilities *paymail.CapabilitiesResponse
	if capabilities, err = client.GetCapabilities(domain, paymail.DefaultPort); err != nil {
		return nil, err
	}

	pkiURL := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)

	var pki *paymail.PKIResponse
	if pki, err = client.GetPKI(pkiURL, alias, domain); err != nil {
		return nil, err
	}
	return pki, nil
}
