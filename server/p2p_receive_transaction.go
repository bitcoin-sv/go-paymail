package server

import (
	"fmt"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
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
func (c *Configuration) p2pReceiveTx(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	p2pFormat := basicP2pPayload

	requestPayload, _, md, vErr := processP2pReceiveTxRequest(c, req, p, p2pFormat)
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
func (c *Configuration) p2pReceiveBeefTx(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	p2pFormat := beefP2pPayload

	requestPayload, beefData, md, vErr := processP2pReceiveTxRequest(c, req, p, p2pFormat)
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
	if err = ExecuteSimplifiedPaymentVerification(req.Context(), beefData); err != nil {
		ErrorResponse(w, ErrorSimplifiedPaymentVerification, err.Error(), http.StatusExpectationFailed)
		//var err error
		//if err = c.actions.ExecuteSimplifiedPaymentVerification(req.Context(), beefData); err != nil {
		//	ErrorResponse(w, ErrorSimplifiedPaymentVerification, err.Error(), http.StatusExpectationFailed)
		//	return
		//}

		// verify merkle proofs
		merkleRoots, err := beefData.GetMerkleRoots()
		fmt.Println("<------- Merkle Roots")
		fmt.Println(merkleRoots)

		err = c.actions.VerifyMerkleRoots(req.Context(), merkleRoots)
		if err != nil {
			ErrorResponse(w, ErrorInvalidParameter, "invalid parameter: merkle proofs", http.StatusBadRequest)
			return
		}

		// TODO: get values from BEEF decode
		//inputSatoshis := 100000
		//outputSatoshis := 100000

		//if inputSatoshis <= outputSatoshis {
		//	ErrorResponse(w, ErrorInvalidParameter, "invalid parameter: input satoshis has to be larger than output satoshis", http.StatusBadRequest)
		//	return
		//}

		// TODO: check scripts pair
		//verifyScripts(nil, nil)

		var response *paymail.P2PTransactionPayload
		if response, err = c.actions.RecordTransaction(
			req.Context(), requestPayload.P2PTransaction, md,
		); err != nil {
			ErrorResponse(w, ErrorRecordingTx, err.Error(), http.StatusExpectationFailed)
			return
		}

		writeJsonResponse(w, http.StatusOK, response)
	}
}

func verifyScripts(tx, prevTx *bt.Tx) bool {

	//bscript.NewFromHexString(firstTx)
	//bscript.NewFromHexString(secondTX)
	//
	//tx, err := bt.NewTxFromString(firstTx)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//prevTx, err := bt.NewTxFromString(secondTX)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	inputIdx := 0
	input := tx.InputIdx(inputIdx)
	prevOutput := prevTx.OutputIdx(int(input.PreviousTxOutIndex))

	inputASM, err := input.UnlockingScript.ToASM()
	if err != nil {
		fmt.Println(err)
		return false
	}

	outputASM, err := prevOutput.LockingScript.ToASM()
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(inputASM)
	fmt.Println(outputASM)

	if err := interpreter.NewEngine().Execute(
		interpreter.WithTx(tx, inputIdx, prevOutput),
		interpreter.WithForkID(),
		interpreter.WithAfterGenesis(),
	); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
