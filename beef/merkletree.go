package beef

import (
	"encoding/hex"

	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	"github.com/bitcoin-sv/go-sdk/util"
)

// INFO: This function is moved to go-paymail from go-bt
// try to use go-sdk implementation when available

// merkleTreeParentStr returns the Merkle Tree parent of two Merkle
func merkleTreeParentStr(leftNode, rightNode string) (string, error) {
	l, err := hex.DecodeString(leftNode)
	if err != nil {
		return "", err
	}
	r, err := hex.DecodeString(rightNode)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(merkleTreeParent(l, r)), nil
}

// merkleTreeParent returns the Merkle Tree parent of two Merkle
// Tree children.
func merkleTreeParent(leftNode, rightNode []byte) []byte {
	// swap endianness before concatenating
	l := util.ReverseBytes(leftNode)
	r := util.ReverseBytes(rightNode)

	// concatenate leaves
	concat := append(l, r...)

	// hash the concatenation
	hash := crypto.Sha256d(concat)

	// swap endianness at the end and convert to hex string
	return util.ReverseBytes(hash)
}
