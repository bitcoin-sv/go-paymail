package paymail

import (
	"errors"
	"github.com/libsv/go-bc"
)

// BUMPPaths represents BUMP format for all inputs
type BUMPPaths []BUMP

// BUMP is a struct that represents a whole BUMP format
type BUMP struct {
	blockHeight uint64
	path        []BUMPPath
}

// BUMPPath is a slice of BUMPLevel objects which represents a path
type BUMPPath []BUMPPathElement

// BUMPPathElement is a struct that represents a single BUMP transaction
type BUMPPathElement struct {
	hash      string
	txId      bool
	duplicate bool
	offset    uint64
}

func (b BUMP) calculateMerkleRoots() ([]string, error) {
	merkleRoots := make([]string, 0)

	for _, bumpPathElement := range b.path[0] {
		merkleRoot, err := calculateMerkleRoot(bumpPathElement, b)
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, merkleRoot)
	}
	return merkleRoots, nil
}

func (bPath BUMPPath) findTxByOffset(offset uint64) *BUMPPathElement {
	for _, bumpTx := range bPath {
		if bumpTx.offset == offset {
			return &bumpTx
		}
	}
	return nil
}

// calculateMerkleRoots will calculate one merkle root for tx in the BUMPPath
func calculateMerkleRoot(baseTx BUMPPathElement, bump BUMP) (string, error) {
	calculatedHash := baseTx.hash
	offset := baseTx.offset

	for _, bLevel := range bump.path {
		newOffset := offset - 1
		if offset%2 == 0 {
			newOffset = offset + 1
		}
		tx2 := bLevel.findTxByOffset(newOffset)
		if &tx2 == nil {
			return "", errors.New("could not find pair")
		}

		leftNode, rightNode := prepareNodes(baseTx, offset, *tx2, newOffset)

		str, err := bc.MerkleTreeParentStr(leftNode, rightNode)
		if err != nil {
			return "", err
		}
		calculatedHash = str

		offset = offset / 2

		baseTx = BUMPPathElement{
			hash:   calculatedHash,
			offset: offset,
		}
	}

	return calculatedHash, nil
}

func prepareNodes(baseTx BUMPPathElement, offset uint64, tx2 BUMPPathElement, newOffset uint64) (string, string) {
	var txHash, tx2Hash string

	if baseTx.duplicate {
		txHash = tx2.hash
	} else {
		txHash = baseTx.hash
	}

	if tx2.duplicate {
		tx2Hash = baseTx.hash
	} else {
		tx2Hash = tx2.hash
	}

	if newOffset > offset {
		return txHash, tx2Hash
	}
	return tx2Hash, txHash
}
