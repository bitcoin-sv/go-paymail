package spv

import (
	"errors"
	"fmt"

	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript/interpreter"

	"github.com/bitcoin-sv/go-paymail"
)

func validateScripts(dBeef *paymail.DecodedBEEF) error {
	for i, input := range dBeef.ProcessedTxData.Inputs {
		inputParentTx := findParentForInput(input, dBeef.InputsTxData)
		if inputParentTx == nil {
			return errors.New("invalid parent transactions, no matching trasactions for input")
		}

		result := verifyScripts(dBeef.ProcessedTxData, inputParentTx.Transaction, i)
		if !result {
			return errors.New("invalid script")
		}
	}

	return nil
}

// Verify locking and unlocking scripts pair
func verifyScripts(tx, prevTx *bt.Tx, inputIdx int) bool {
	input := tx.InputIdx(inputIdx)
	prevOutput := prevTx.OutputIdx(int(input.PreviousTxOutIndex))

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
