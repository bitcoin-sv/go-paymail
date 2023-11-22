package spv

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail/beef"
	"github.com/libsv/go-bt/v2"
)

func verifyMerkleRoots(ctx context.Context, dBeef *beef.DecodedBEEF, provider MerkleRootVerifier) error {
	if err := ensureInputParentsArePresentInBump(dBeef.GetLatestTx(), dBeef); err != nil {
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

func ensureInputParentsArePresentInBump(tx *bt.Tx, dBeef *beef.DecodedBEEF) error {

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

func findMinedAncestors(input *bt.Input, parentTxs []*beef.TxData, depth uint) []*beef.TxData {
	if depth > 64 {
		return []*beef.TxData{}
	}
	depth++

	parent := findParentForInput(input, parentTxs)

	if parent == nil { // oh oh- end of hierarchy
		return []*beef.TxData{}
	}

	if !parent.Unmined() {
		return []*beef.TxData{parent}
	}

	ancestors := make([]*beef.TxData, 0)

	for _, in := range parent.Transaction.Inputs {
		ancestors = append(ancestors, findMinedAncestors(in, parentTxs, depth)...)
	}

	return ancestors
}

func existsInBumps(ancestorTx *beef.TxData, bumps beef.BUMPs) bool {
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
