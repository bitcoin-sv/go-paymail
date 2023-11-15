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

func validateSatoshisSum(dBeef *DecodedBEEF) error {
	if len(dBeef.ProcessedTxData.Outputs) == 0 {
		return errors.New("invalid output, no outputs")
	}

	if len(dBeef.ProcessedTxData.Inputs) == 0 {
		return errors.New("invalid input, no inputs")
	}

	inputSum, outputSum := uint64(0), uint64(0)

	for _, input := range dBeef.ProcessedTxData.Inputs {
		inputParentTx := findParentForInput(input, dBeef.InputsTxData)

		if inputParentTx == nil {
			return errors.New("invalid parent transactions, no matching trasactions for input")
		}

		inputSum += inputParentTx.Transaction.Outputs[input.PreviousTxOutIndex].Satoshis
	}
	for _, output := range dBeef.ProcessedTxData.Outputs {
		outputSum += output.Satoshis
	}

	if inputSum <= outputSum {
		return errors.New("invalid input and output sum, outputs can not be larger than inputs")
	}
	return nil
}

func validateLockTime(dBeef *DecodedBEEF) error {
	if dBeef.ProcessedTxData.LockTime == 0 {
		for _, input := range dBeef.ProcessedTxData.Inputs {
			if input.SequenceNumber != 0xffffffff {
				return errors.New("unexpected transaction with nSequence")
			}
		}
	} else {
		return errors.New("unexpected transaction with nLockTime")
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

func findParentForInput(input *bt.Input, parentTxs []*TxData) *TxData {
	parentID := input.PreviousTxIDStr()

	for _, ptx := range parentTxs {
		if ptx.Transaction.TxID() == parentID {
			return ptx
		}
	}

	return nil
}
