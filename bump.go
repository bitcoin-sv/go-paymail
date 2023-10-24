package paymail

// BUMPTx is a struct that represents a single BUMP transaction
type BUMPTx struct {
	hash      string
	txId      bool
	duplicate bool
}

// BUMPMap is a map of BUMPTx where offset is the key
type BUMPMap []map[uint64]BUMPTx

// BUMPSlice is a slice of BUMPMap which contain transactions required to calculate merkle roots
type BUMPSlice []BUMPMap

// BUMP is a struct that represents a whole BUMP format
type BUMP struct {
	blockHeight uint64
	path        BUMPSlice
}

func (b BUMPMap) calculateMerkleRoots() ([]string, error) {
	merkleRoots := make([]string, 0)

	for offset, bumpTx := range b[0] {
		merkleRoot, err := calculateMerkleRoot(bumpTx, offset, b)
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, merkleRoot)
	}
	return merkleRoots, nil
}
