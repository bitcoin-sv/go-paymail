package server

import (
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/julienschmidt/httprouter"
)

type p2pPayloadFormat uint

const (
	basicP2pPayload p2pPayloadFormat = iota
	beefP2pPayload
)

/*
Incoming Data Object Example:
{
  "hex": "01000000012adda020db81f2155ebba69e7.........154888ac00000000",
  "metadata": {
	"sender": "someone@example.tld",
	"pubkey": "<sender-pubkey>",
	"signature": "signature(txid)",
	"note": "Human readable information related to the tx."
  },
  "reference": "someRefId"
}
*/

// p2pReceiveTx will receive a P2P transaction (from previous request: P2P Payment Destination)
//
// Specs: https://docs.moneybutton.com/docs/paymail-06-p2p-transactions.html
func (c *Configuration) p2pReceiveTx(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	p2pFormat := basicP2pPayload

	requestPayload, _, md, vErr := processP2pReceiveTxRequest(c, req, p2pFormat)
	if vErr != nil {
		ErrorResponse(w, vErr.code, vErr.msg, vErr.httpResponseCode)
		return
	}

	if len(requestPayload.Hex) == 0 {
		panic("empty hex after parsing!")
	}

	var response *paymail.P2PTransactionPayload
	var err error
	if response, err = c.actions.RecordTransaction(
		req.Context(), requestPayload.P2PTransaction, md,
	); err != nil {
		ErrorResponse(w, ErrorRecordingTx, err.Error(), http.StatusExpectationFailed)
		return
	}

	writeJsonResponse(w, http.StatusOK, response)
}

/*
Incoming Data Object Example:
{
  "beef": "01000000012adda020db81f2155ebba69e7.........154888ac00000000",
  "metadata": {
	"sender": "someone@example.tld",
	"pubkey": "<sender-pubkey>",
	"signature": "signature(txid)",
	"note": "Human readable information related to the tx."
  },
  "reference": "someRefId"
}
*/
// p2pReceiveBeefTx will receive a P2P transaction in BEEF format
func (c *Configuration) p2pReceiveBeefTx(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	p2pFormat := beefP2pPayload

	requestPayload, beefData, md, vErr := processP2pReceiveTxRequest(c, req, p2pFormat)
	if vErr != nil {
		ErrorResponse(w, vErr.code, vErr.msg, vErr.httpResponseCode)
		return
	}

	if len(requestPayload.Hex) == 0 {
		panic("empty hex after parsing!")
	}

	if beefData == nil {
		panic("empty beef after parsing!")
	}

	var err error
	if err = c.actions.ExecuteSimplifiedPaymentVerification(req.Context(), beefData); err != nil {
		ErrorResponse(w, ErrorSimplifiedPaymentVerification, err.Error(), http.StatusExpectationFailed)
		return
	}

	var response *paymail.P2PTransactionPayload
	if response, err = c.actions.RecordTransaction(
		req.Context(), requestPayload.P2PTransaction, md,
	); err != nil {
		ErrorResponse(w, ErrorRecordingTx, err.Error(), http.StatusExpectationFailed)
		return
	}

	writeJsonResponse(w, http.StatusOK, response)
}
