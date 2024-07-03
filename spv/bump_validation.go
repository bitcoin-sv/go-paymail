package spv

import (
	"github.com/bitcoin-sv/go-paymail/errors"

	"github.com/bitcoin-sv/go-paymail/beef"
	"github.com/libsv/go-bt/v2"
)

func ensureAncestorsArePresentInBump(tx *bt.Tx, dBeef *beef.DecodedBEEF) error {
	ancestors, err := findMinedAncestors(tx, dBeef.Transactions)
	if err != nil {
		return err
	}

	for _, tx := range ancestors {
		if !existsInBumps(tx, dBeef.BUMPs) {
			return errors.ErrBUMPAncestorNotPresent
		}
	}

	return nil
}

func findMinedAncestors(tx *bt.Tx, ancestors []*beef.TxData) (map[string]*beef.TxData, error) {
	am := make(map[string]*beef.TxData)

	for _, input := range tx.Inputs {

		if err := findMinedAncestorsForInput(input, ancestors, am); err != nil {
			return nil, err
		}
	}

	return am, nil
}

func findMinedAncestorsForInput(input *bt.Input, ancestors []*beef.TxData, ma map[string]*beef.TxData) error {
	parent := findParentForInput(input, ancestors)
	if parent == nil {
		return errors.ErrBUMPCouldNotFindMinedParent
	}

	if !parent.Unmined() {
		ma[parent.GetTxID()] = parent
		return nil
	}

	for _, in := range parent.Transaction.Inputs {
		err := findMinedAncestorsForInput(in, ancestors, ma) // we don't have to worry about infinite recursion - the graph will always be acyclic due to the nature of the transactions
		if err != nil {
			return err
		}
	}

	return nil
}

func existsInBumps(tx *beef.TxData, bumps beef.BUMPs) bool {
	bumpIdx := int(*tx.BumpIndex)
	txID := tx.GetTxID()

	if len(bumps) > bumpIdx {
		leafs := bumps[bumpIdx].Path[0]

		for _, lf := range leafs {
			if txID == lf.Hash {
				return true
			}
		}
	}

	return false
}
