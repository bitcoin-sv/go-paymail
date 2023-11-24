package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"

	"github.com/bitcoin-sv/go-paymail"
)

type p2pReceiveTxReqPayload struct {
	*paymail.P2PTransaction
	incomingPaymailAlias, incomingPaymailDomain string
}

type processingError struct {
	*parseError
	httpResponseCode int
}

func processP2pReceiveTxRequest(c *Configuration, req *http.Request, p httprouter.Params, format p2pPayloadFormat) (
	*p2pReceiveTxReqPayload, *paymail.DecodedBEEF, *RequestMetadata, *processingError,
) {
	payload, vErr := parseP2pReceiveTxRequest(c, req, p, format)
	if vErr != nil {
		return returnError(&processingError{vErr, http.StatusBadRequest})
	}

	md := CreateMetadata(req, payload.incomingPaymailAlias, payload.incomingPaymailDomain, "")
	vErr = verifyIncomingPaymail(req.Context(), c, md, payload.incomingPaymailAlias, payload.incomingPaymailDomain)

	if vErr != nil {
		if vErr.code == ErrorPaymailNotFound {
			return returnError(&processingError{vErr, http.StatusNotFound})
		}

		return returnError(&processingError{vErr, http.StatusExpectationFailed})
	}

	tx, beefData, pErr := getProcessedTxData(payload, format)
	if pErr != nil {
		return returnError(pErr)
	}

	if c.SenderValidationEnabled || len(payload.MetaData.Signature) > 0 {
		if vErr = verifySignature(payload.MetaData, tx.TxID()); vErr != nil {
			return returnError(&processingError{vErr, http.StatusBadRequest})
		}
	}

	if format == beefP2pPayload {
		payload.Hex = tx.String()
		payload.DecodedBeef = beefData
	}

	return payload, beefData, md, nil
}

func getProcessedTxData(payload *p2pReceiveTxReqPayload, format p2pPayloadFormat) (*bt.Tx, *paymail.DecodedBEEF, *processingError) {
	var processedTx *bt.Tx
	var beefData *paymail.DecodedBEEF
	var err error

	switch format {
	case basicP2pPayload:
		processedTx, err = bitcoin.TxFromHex(payload.Hex)
		if err != nil {
			errorMsg := fmt.Sprintf("error while parsing hex: %s", err.Error())
			log.Println(errorMsg)
			return nil, nil, &processingError{&parseError{ErrorInvalidParameter, errorMsg}, http.StatusBadRequest}
		}

	case beefP2pPayload:
		beefData, err = paymail.DecodeBEEF(payload.Beef)
		if err != nil {
			errorMsg := fmt.Sprintf("error while parsing beef: %s", err.Error())
			log.Println(errorMsg)
			return nil, nil, &processingError{&parseError{ErrorInvalidParameter, errorMsg}, http.StatusBadRequest}
		}

		processedTx = beefData.ProcessedTxData

	default:
		panic("Unexpected transaction format!")
	}

	return processedTx, beefData, nil
}

func verifyIncomingPaymail(ctx context.Context, c *Configuration, md *RequestMetadata, alias, domain string) *parseError {
	var foundPaymail *paymail.AddressInformation
	var err error

	foundPaymail, err = c.actions.GetPaymailByAlias(ctx, alias, domain, md)
	if err != nil {
		return &parseError{ErrorFindingPaymail, err.Error()}
	} else if foundPaymail == nil {
		return &parseError{ErrorPaymailNotFound, "paymail not found"}
	}

	return nil
}

func verifySignature(metadata *paymail.P2PMetaData, txID string) *parseError {
	// Get the address from pubKey
	var rawAddress *bscript.Address
	var err error

	if rawAddress, err = bitcoin.GetAddressFromPubKeyString(metadata.PubKey, true); err != nil {
		return &parseError{ErrorInvalidPubKey, "invalid pubkey: " + err.Error()}
	}

	// Validate the signature of the tx id
	if err = bitcoin.VerifyMessage(rawAddress.AddressString, metadata.Signature, txID); err != nil {
		return &parseError{ErrorInvalidSignature, "invalid signature: " + err.Error()}
	}

	return nil
}

func returnError(err *processingError) (
	*p2pReceiveTxReqPayload, *paymail.DecodedBEEF, *RequestMetadata, *processingError,
) {
	return nil, nil, nil, err
}
