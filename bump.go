package paymail

type BUMPTx struct {
	hash      string
	txId      bool
	duplicate bool
}

type BUMPMap []map[string]BUMPTx

type BUMPSlice []BUMPMap

type BUMP struct {
	blockHeight uint64
	path        BUMPSlice
}

func (b BUMPMap) calculateMerkleRoots() ([]string, error) {
	merkleRoots := make([]string, 0)

	for tx, offset := range cmp[len(cmp)-1] {
		merkleRoot, err := calculateMerkleRoot(tx, offset, cmp)
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, merkleRoot)
	}
	return merkleRoots, nil
}
