package paymail

import (
	"errors"

	"github.com/libsv/go-bc"
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
	dataFlag byte = iota
	duplicateFlag
	txIDFlag
)

func (b BUMP) calculateMerkleRoot() (string, error) {
	merkleRoot := ""

	for _, bumpPathElement := range b.path[0] {
		if bumpPathElement.txId {
			calcMerkleRoot, err := calculateMerkleRoot(bumpPathElement, b)
			if err != nil {
				return "", err
			}

			if merkleRoot == "" {
				merkleRoot = calcMerkleRoot
				continue
			}

			if calcMerkleRoot != merkleRoot {
				return "", errors.New("different merkle roots for the same block")
			}
		}
	}
	return merkleRoot, nil
}

func findLeafByOffset(offset uint64, bumpLeaves []BUMPLeaf) *BUMPLeaf {
	for _, bumpTx := range bumpLeaves {
		if bumpTx.offset == offset {
			return &bumpTx
		}
	}
	return nil
}

// calculateMerkleRoots will calculate one merkle root for tx in the BUMPLeaf
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
	var baseLeafHash, pairLeafHash string

	if baseLeaf.duplicate {
		baseLeafHash = leafInPair.hash
	} else {
		baseLeafHash = baseLeaf.hash
	}

	if leafInPair.duplicate {
		pairLeafHash = baseLeaf.hash
	} else {
		pairLeafHash = leafInPair.hash
	}

	if newOffset > offset {
		return baseLeafHash, pairLeafHash
	}
	return pairLeafHash, baseLeafHash
}
