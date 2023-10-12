package paymail

type CompoundMerklePath []map[string]uint64

type CMPSlice []CompoundMerklePath

func (cmp CompoundMerklePath) calculateMerkleRoots() ([]string, error) {
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
