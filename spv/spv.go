package spv

import (
	"context"
	"errors"

	"github.com/libsv/go-bt/v2"

	"github.com/bitcoin-sv/go-paymail"
)

type MerkleRootVerifier interface {
	VerifyMerkleRoots(
		ctx context.Context,
		merkleRoots []paymail.MerkleRootConfirmationRequestItem,
	) error
}

// ExecuteSimplifiedPaymentVerification executes the SPV for decoded BEEF tx
func ExecuteSimplifiedPaymentVerification(dBeef *paymail.DecodedBEEF, provider MerkleRootVerifier) error {

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

func validateSatoshisSum(dBeef *paymail.DecodedBEEF) error {
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

func validateLockTime(dBeef *paymail.DecodedBEEF) error {
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

func findParentForInput(input *bt.Input, parentTxs []*paymail.TxData) *paymail.TxData {
	parentID := input.PreviousTxIDStr()

	for _, ptx := range parentTxs {
		if ptx.GetTxID() == parentID {
			return ptx
		}
	}

	return nil
}
