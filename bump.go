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
	path        []BUMPPathMap
}

// BUMPPathMap is a map of BUMPTx where offset is the key
type BUMPPathMap map[uint64]BUMPPathElement

// BUMPPathElement is a struct that represents a single BUMP transaction
type BUMPPathElement struct {
	hash      string
	txId      bool
	duplicate bool
}

func (b BUMP) calculateMerkleRoots() ([]string, error) {
	merkleRoots := make([]string, 0)

	for offset, pathElement := range b.path[0] {
		merkleRoot, err := calculateMerkleRoot(pathElement, offset, b.path)
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, merkleRoot)
	}
	//}
	return merkleRoots, nil
}

func calculateMerkleRoot(baseElement BUMPPathElement, offset uint64, bumpPathMaps []BUMPPathMap) (string, error) {
	calculatedHash := baseElement.hash

	for _, bPathMap := range bumpPathMaps {
		newOffset := offset - 1
		if offset%2 == 0 {
			newOffset = offset + 1
		}
		pairElement := bPathMap[newOffset]
		if &pairElement == nil {
			return "", errors.New("could not find pair")
		}

		leftNode, rightNode := prepareNodes(baseElement, offset, pairElement, newOffset)

		str, err := bc.MerkleTreeParentStr(leftNode, rightNode)
		if err != nil {
			return "", err
		}
		calculatedHash = str

		offset = offset / 2

		baseElement = BUMPPathElement{
			hash: calculatedHash,
		}
	}

	return calculatedHash, nil
}

func prepareNodes(baseTx BUMPPathElement, offset uint64, pairElement BUMPPathElement, newOffset uint64) (string, string) {
	var txHash, tx2Hash string

	if baseTx.duplicate {
		txHash = pairElement.hash
	} else {
		txHash = baseTx.hash
	}

	if pairElement.duplicate {
		tx2Hash = baseTx.hash
	} else {
		tx2Hash = pairElement.hash
	}

	if newOffset > offset {
		return txHash, tx2Hash
	}
	return tx2Hash, txHash
}
