package spv

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/libsv/go-bt/v2"
)

func verifyMerkleRoots(dBeef *paymail.DecodedBEEF, provider MerkleRootVerifier) error {
	if err := ensureInputTransactionArePresentInBump(dBeef); err != nil {
		return err
	}

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

func ensureInputTransactionArePresentInBump(dBeef *paymail.DecodedBEEF) error {

	for _, input := range dBeef.ProcessedTxData.Inputs {

		minedAcestors := findMinedAncestors(input, dBeef.InputsTxData)
		if len(minedAcestors) == 0 {
			return errors.New("invalid BUMP - input mined ancestor is not present in BUMPs")
		}

		for _, ancestorTx := range minedAcestors {

			if !existsInBumps(ancestorTx, dBeef.BUMPs) {
				return errors.New("invalid BUMP - input mined ancestor is not present in BUMPs")
			}
		}
	}

	return nil
}

func findMinedAncestors(input *bt.Input, parentTxs []*paymail.TxData) []*paymail.TxData {
	parent := findParentForInput(input, parentTxs)

	if parent == nil { // oh oh- end of hierarchy
		return []*paymail.TxData{}
	}

	if parent.PathIndex != nil { // mined parent
		return []*paymail.TxData{parent}
	}

	ancestors := make([]*paymail.TxData, 0)

	for _, in := range parent.Transaction.Inputs {
		ancestors = append(ancestors, findMinedAncestors(in, parentTxs)...)
	}

	return ancestors
}

func existsInBumps(ancestorTx *paymail.TxData, bumps paymail.BUMPs) bool {
	bumpIdx := int(*ancestorTx.PathIndex) // TODO: disscuss max value of index
	parentTxID := ancestorTx.GetTxID()

	if len(bumps) > bumpIdx {
		leafs := bumps[bumpIdx].Path[0]

		for _, l := range leafs {
			if parentTxID == l.Hash {
				return true
			}
		}
	}

	return false
}
