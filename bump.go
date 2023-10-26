package paymail

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

func (bLevel BUMPPath) findTxByOffset(offset uint64) *BUMPPathElement {
	for _, bumpTx := range bLevel {
		if bumpTx.offset == offset {
			return &bumpTx
		}
	}
	return nil
}
