package spv

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/libsv/go-bt/v2"
)

func verifyMerkleRoots(ctx context.Context, dBeef *paymail.DecodedBEEF, provider MerkleRootVerifier) error {
	if err := ensureInputTransactionsArePresentInBump(dBeef.GetLatestTx(), dBeef); err != nil {
		return err
	}

	merkleRoots, err := dBeef.GetMerkleRootsRequest()
	if err != nil {
		return err
	}

	err = provider.VerifyMerkleRoots(ctx, merkleRoots)
	if err != nil {
		return err
	}

	return nil
}

func ensureInputTransactionsArePresentInBump(tx *bt.Tx, dBeef *paymail.DecodedBEEF) error {

	for _, input := range tx.Inputs {

		minedAcestors := findMinedAncestors(input, dBeef.Transactions, 0)
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

func findMinedAncestors(input *bt.Input, parentTxs []*paymail.TxData, depth uint) []*paymail.TxData {
	if depth > 64 {
		return []*paymail.TxData{}
	}
	depth++

	parent := findParentForInput(input, parentTxs)

	if parent == nil { // oh oh- end of hierarchy
		return []*paymail.TxData{}
	}

	if !parent.Unmined() {
		return []*paymail.TxData{parent}
	}

	ancestors := make([]*paymail.TxData, 0)

	for _, in := range parent.Transaction.Inputs {
		ancestors = append(ancestors, findMinedAncestors(in, parentTxs, depth)...)
	}

	return ancestors
}

func existsInBumps(ancestorTx *paymail.TxData, bumps paymail.BUMPs) bool {
	bumpIdx := int(*ancestorTx.BumpIndex)
	parentTxID := ancestorTx.GetTxID()

	if len(bumps) > bumpIdx {
		leafs := bumps[bumpIdx].Path[0]

		for _, lf := range leafs {
			if parentTxID == lf.Hash {
				return true
			}
		}
	}

	return false
}
