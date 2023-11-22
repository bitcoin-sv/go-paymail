package spv

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail/beef"
	"github.com/libsv/go-bt/v2"
)

type MerkleRootVerifier interface {
	VerifyMerkleRoots(
		ctx context.Context,
		merkleRoots []*MerkleRootConfirmationRequestItem,
	) error
}

// MerkleRootConfirmationRequestItem is a request type for verification
// of Merkle Roots inclusion in the longest chain.
type MerkleRootConfirmationRequestItem struct {
	MerkleRoot  string `json:"merkleRoot"`
	BlockHeight uint64 `json:"blockHeight"`
}

// ExecuteSimplifiedPaymentVerification executes the SPV for decoded BEEF tx
func ExecuteSimplifiedPaymentVerification(ctx context.Context, dBeef *beef.DecodedBEEF, provider MerkleRootVerifier) error {

	for _, txDt := range dBeef.Transactions {
		tx := txDt.Transaction

		if len(tx.Outputs) == 0 {
			return errors.New("invalid output, no outputs")
		}

		if len(tx.Inputs) == 0 {
			return errors.New("invalid input, no inputs")
		}

		err := validateLockTime(tx)
		if err != nil {
			return err
		}

		if txDt.Unmined() {
			err = validateSatoshisSum(tx, dBeef.Transactions)
			if err != nil {
				return err
			}

			err = validateScripts(tx, dBeef.Transactions)
			if err != nil {
				return err
			}
		}
	}

	err := verifyMerkleRoots(ctx, dBeef, provider)
	if err != nil {
		return err
	}

	return nil
}

func validateLockTime(tx *bt.Tx) error {
	if tx.LockTime == 0 {
		for _, input := range tx.Inputs {
			if input.SequenceNumber != 0xffffffff {
				return errors.New("unexpected transaction with nSequence")
			}
		}
	} else {
		return errors.New("unexpected transaction with nLockTime")
	}
	return nil
}

func validateSatoshisSum(tx *bt.Tx, inputTxs []*beef.TxData) error {
	inputSum, outputSum := uint64(0), uint64(0)

	for _, input := range tx.Inputs {
		inputParentTx := findParentForInput(input, inputTxs)

		if inputParentTx == nil {
			return errors.New("invalid parent transactions, no matching trasactions for input")
		}

		inputSum += inputParentTx.Transaction.Outputs[input.PreviousTxOutIndex].Satoshis
	}
	for _, output := range tx.Outputs {
		outputSum += output.Satoshis
	}

	if inputSum <= outputSum {
		return errors.New("invalid input and output sum, outputs can not be larger than inputs")
	}

	return nil
}

func findParentForInput(input *bt.Input, parentTxs []*beef.TxData) *beef.TxData {
	parentID := input.PreviousTxIDStr()

	for _, ptx := range parentTxs {
		if ptx.GetTxID() == parentID {
			return ptx
		}
	}

	return nil
}
