package spv

import (
	"errors"
	"fmt"

	"github.com/bitcoin-sv/go-paymail/beef"
	"github.com/libsv/go-bt/v2"
)

const recursiveMaxDepth = 128 // arbitrarily chosen value

func ensureAncestorsArePresentInBump(tx *bt.Tx, dBeef *beef.DecodedBEEF) error {
	ancestors, err := findMinedAncestors(tx, dBeef.Transactions)
	if err != nil {
		return err
	}

	for _, tx := range ancestors {
		if !existsInBumps(tx, dBeef.BUMPs) {
			return errors.New("invalid BUMP - input mined ancestor is not present in BUMPs")
		}
	}

	return nil
}

func findMinedAncestors(tx *bt.Tx, ancestors []*beef.TxData) (map[string]*beef.TxData, error) {
	am := make(map[string]*beef.TxData)

	for _, input := range tx.Inputs {
		err := findMinedAncestorsForInput(input, ancestors, am, 0)

		if err != nil {
			return nil, err
		}
	}

	return am, nil
}

func findMinedAncestorsForInput(input *bt.Input, ancestors []*beef.TxData, ma map[string]*beef.TxData, depth uint) error {
	if depth > recursiveMaxDepth { //primitive protection against Cyclic Graph (and therefore infinite loop)
		return fmt.Errorf("invalid BUMP - cannot find mined parent for input %s on %d depth", input.String(), depth)
	}
	depth++

	parent := findParentForInput(input, ancestors)
	if parent == nil {
		return fmt.Errorf("invalid BUMP - cannot find mined parent for input %s", input.String())
	}

	if !parent.Unmined() {
		ma[parent.GetTxID()] = parent
		return nil
	}

	for _, in := range parent.Transaction.Inputs {
		err := findMinedAncestorsForInput(in, ancestors, ma, depth)
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
