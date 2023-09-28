package server

import (
	"encoding/json"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/libsv/go-bt/v2/bscript"
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
func (c *Configuration) p2pReceiveTx(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	incomingPaymail := p.ByName("paymailAddress")
	// Parse, sanitize and basic validation
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		ErrorResponse(w, ErrorInvalidParameter, "invalid paymail", http.StatusBadRequest)
		return
	} else if !c.IsAllowedDomain(domain) {
		ErrorResponse(w, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	var p2pTransaction paymail.P2PTransaction
	err := json.NewDecoder(req.Body).Decode(&p2pTransaction)
	if err != nil {
		ErrorResponse(w, ErrorInvalidParameter, "invalid request body", http.StatusBadRequest)
		return
	}

	// Check for required fields
	if len(p2pTransaction.Hex) == 0 {
		ErrorResponse(w, ErrorMissingHex, "missing parameter: hex", http.StatusBadRequest)
		return
	} else if len(p2pTransaction.Reference) == 0 {
		ErrorResponse(w, ErrorMissingReference, "missing parameter: reference", http.StatusBadRequest)
		return
	}

	// Convert the raw tx into a transaction
	transaction, err := bitcoin.TxFromHex(p2pTransaction.Hex)
	if err != nil {
		ErrorResponse(w, ErrorInvalidParameter, "invalid parameter: hex", http.StatusBadRequest)
		return
	}

	// Start the final response
	response := &paymail.P2PTransactionPayload{
		Note: p2pTransaction.MetaData.Note,
		TxID: transaction.TxID(),
	}

	// Check signature if: 1) sender validation enabled or 2) a signature was given (optional)
	if c.SenderValidationEnabled || len(p2pTransaction.MetaData.Signature) > 0 {

		// Check required fields for signature validation
		if len(p2pTransaction.MetaData.Signature) == 0 {
			ErrorResponse(w, ErrorInvalidSignature, "missing parameter: signature", http.StatusBadRequest)
			return
		} else if len(p2pTransaction.MetaData.PubKey) == 0 {
			ErrorResponse(w, ErrorInvalidPubKey, "missing parameter: pubkey", http.StatusBadRequest)
			return
		}

		// Get the address from pubKey
		var rawAddress *bscript.Address
		if rawAddress, err = bitcoin.GetAddressFromPubKeyString(p2pTransaction.MetaData.PubKey, true); err != nil {
			ErrorResponse(w, ErrorInvalidPubKey, "invalid pubkey: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the signature of the tx id
		if err = bitcoin.VerifyMessage(rawAddress.AddressString, p2pTransaction.MetaData.Signature, response.TxID); err != nil {
			ErrorResponse(w, ErrorInvalidSignature, "invalid signature: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Create the metadata struct
	md := CreateMetadata(req, alias, domain, "")

	// Get from the data layer
	var foundPaymail *paymail.AddressInformation
	foundPaymail, err = c.actions.GetPaymailByAlias(req.Context(), alias, domain, md)
	if err != nil {
		ErrorResponse(w, ErrorFindingPaymail, err.Error(), http.StatusExpectationFailed)
		return
	} else if foundPaymail == nil {
		ErrorResponse(w, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	// Record the transaction (verify, save, broadcast...)
	if response, err = c.actions.RecordTransaction(
		req.Context(), &p2pTransaction, md,
	); err != nil {
		ErrorResponse(w, ErrorRecordingTx, err.Error(), http.StatusExpectationFailed)
		return
	}

	// Set the response
	writeJsonResponse(w, http.StatusOK, response)
}

// p2pReceiveBeefTx will receive a P2P transaction in BEEF format
func (c *Configuration) p2pReceiveBeefTx(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	parms := req.URL.Query()
	beef := parms.Get("beef")

	// TODO: Use those values in future processing
	_, _ = paymail.DecodeBEEF(beef)

	ErrorResponse(w, ErrorNotImplmented, "Receive BEEF transactions not implemented", http.StatusNotImplemented)
}
