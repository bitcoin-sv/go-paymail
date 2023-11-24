package spv

import (
	"errors"

	"github.com/bitcoin-sv/go-paymail/beef"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
)

func validateScripts(tx *bt.Tx, inputTxs []*beef.TxData) error {
	for i, input := range tx.Inputs {
		inputParentTx := findParentForInput(input, inputTxs)
		if inputParentTx == nil {
			return errors.New("invalid parent transactions, no matching trasactions for input")
		}

		err := verifyScripts(tx, inputParentTx.Transaction, i)
		if err != nil {
			return errors.New("invalid script")
		}
	}

	return nil
}

// Verify locking and unlocking scripts pair
func verifyScripts(tx, prevTx *bt.Tx, inputIdx int) error {
	input := tx.InputIdx(inputIdx)
	prevOutput := prevTx.OutputIdx(int(input.PreviousTxOutIndex))

	err := interpreter.NewEngine().Execute(
		interpreter.WithTx(tx, inputIdx, prevOutput),
		interpreter.WithForkID(),
		interpreter.WithAfterGenesis(),
	)

	return err
}
