package paymail

import (
	"errors"
	"github.com/libsv/go-bc"
	"github.com/libsv/go-bt/v2"
)

// BUMPs represents a slice of BUMPs - BSV Unified Merkle Paths
type BUMPs []BUMP

// BUMP is a struct that represents a whole BUMP format
type BUMP struct {
	blockHeight uint64
	path        [][]BUMPLeaf
}

// BUMPLeaf represents each BUMP path element
type BUMPLeaf struct {
	hash      string
	txId      bool
	duplicate bool
	offset    uint64
}

// Flags which are used to determine the type of BUMPLeaf
const (
	dataFlag bt.VarInt = iota
	duplicateFlag
	txIDFlag
)

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

func findLeafByOffset(offset uint64, bumpLeaves []BUMPLeaf) *BUMPLeaf {
	for _, bumpTx := range bumpLeaves {
		if bumpTx.offset == offset {
			return &bumpTx
		}
	}
	return nil
}

// calculateMerkleRoots will calculate one merkle root for tx in the BUMPPath
func calculateMerkleRoot(baseLeaf BUMPLeaf, bump BUMP) (string, error) {
	calculatedHash := baseLeaf.hash
	offset := baseLeaf.offset

	for _, bLevel := range bump.path {
		newOffset := getOffsetPair(offset)
		leafInPair := findLeafByOffset(newOffset, bLevel)
		if leafInPair == nil {
			return "", errors.New("could not find pair")
		}

		leftNode, rightNode := prepareNodes(baseLeaf, offset, *leafInPair, newOffset)

		str, err := bc.MerkleTreeParentStr(leftNode, rightNode)
		if err != nil {
			return "", err
		}
		calculatedHash = str

		offset = offset / 2

		baseLeaf = BUMPLeaf{
			hash:   calculatedHash,
			offset: offset,
		}
	}

	return calculatedHash, nil
}

func getOffsetPair(offset uint64) uint64 {
	if offset%2 == 0 {
		return offset + 1
	}
	return offset - 1
}

func prepareNodes(baseLeaf BUMPLeaf, offset uint64, leafInPair BUMPLeaf, newOffset uint64) (string, string) {
	var txHash, tx2Hash string

	if baseLeaf.duplicate {
		txHash = leafInPair.hash
	} else {
		txHash = baseLeaf.hash
	}

	if leafInPair.duplicate {
		tx2Hash = baseLeaf.hash
	} else {
		tx2Hash = leafInPair.hash
	}

	if newOffset > offset {
		return txHash, tx2Hash
	}
	return tx2Hash, txHash
}
