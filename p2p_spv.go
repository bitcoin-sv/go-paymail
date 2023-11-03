package paymail

import (
	"context"
	"errors"
	"fmt"

	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
)

type MerkleRootVerifier interface {
	VerifyMerkleRoots(
		ctx context.Context,
		merkleRoots []MerkleRootConfirmationRequestItem,
	) error
}

// ExecuteSimplifiedPaymentVerification executes the SPV for decoded BEEF tx
func ExecuteSimplifiedPaymentVerification(dBeef *DecodedBEEF, provider MerkleRootVerifier) error {
	err := validateSatoshisSum(dBeef)
	if err != nil {
		return err
	}

	err = validateLockTime(dBeef)
	if err != nil {
		return err
	}

	err = validateScripts(dBeef)
	if err != nil {
		return err
	}

	err = verifyMerkleRoots(dBeef, provider)
	if err != nil {
		return err
	}

	return nil
}

func verifyMerkleRoots(dBeef *DecodedBEEF, provider MerkleRootVerifier) error {
	merkleRoots, err := dBeef.GetMerkleRootsRequest()
	if err != nil {
		return err
	}

	err = provider.VerifyMerkleRoots(context.Background(), merkleRoots)
	if err != nil {
		return err
	}

	return nil
}

func validateScripts(dBeef *DecodedBEEF) error {
	for _, input := range dBeef.ProcessedTxData.Transaction.Inputs {
		txId := input.PreviousTxID()
		for j, input2 := range dBeef.InputsTxData {
			if input2.Transaction.TxID() == string(txId) {
				result := verifyScripts(dBeef.ProcessedTxData.Transaction, input2.Transaction, j)
				if !result {
					return errors.New("invalid script")
				}
				break
			}
		}
	}
	return nil
}

func validateSatoshisSum(dBeef *DecodedBEEF) error {
	if len(dBeef.ProcessedTxData.Transaction.Outputs) == 0 {
		return errors.New("invalid output, no outputs")
	}

	if len(dBeef.ProcessedTxData.Transaction.Inputs) == 0 {
		return errors.New("invalid input, no inputs")
	}

	inputSum, outputSum := uint64(0), uint64(0)
	for i, input := range dBeef.ProcessedTxData.Transaction.Inputs {
		inputParentTx := dBeef.InputsTxData[i]
		inputSum += inputParentTx.Transaction.Outputs[input.PreviousTxOutIndex].Satoshis
	}
	for _, output := range dBeef.ProcessedTxData.Transaction.Outputs {
		outputSum += output.Satoshis
	}

	if inputSum <= outputSum {
		return errors.New("invalid input and output sum, outputs can not be larger than inputs")
	}
	return nil
}

func validateLockTime(dBeef *DecodedBEEF) error {
	if dBeef.ProcessedTxData.Transaction.LockTime == 0 {
		for _, input := range dBeef.ProcessedTxData.Transaction.Inputs {
			if input.SequenceNumber != 0xffffffff {
				return errors.New("unexpected transaction with nSequence")
			}
		}
	} else {
		return errors.New("nexpected transaction with nLockTime")
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
