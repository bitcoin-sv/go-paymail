package paymail

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/libsv/go-bc"
	"github.com/libsv/go-bt/v2"
)

type CompoundMerklePath []map[string]uint64

type CMPSlice []CompoundMerklePath

const (
	BEEFMarkerPart1 = 0xBE
	BEEFMarkerPart2 = 0xEF
)

const (
	HasNoCMP = 0x00
	HasCMP   = 0x01
)

const (
	hashBytesCount    = 32
	markerBytesCount  = 2
	versionBytesCount = 2
)

type TxData struct {
	Transaction *bt.Tx
	PathIndex   *bt.VarInt
}

type DecodedBEEF struct {
	CMPSlice        CMPSlice
	InputsTxData    []TxData
	ProcessedTxData TxData
}

func (dBeef *DecodedBEEF) GetMerkleRoots() ([]string, error) {
	var merkleRoots []string
	for _, cmp := range dBeef.CMPSlice {
		fmt.Println("<----   CMP   ---->")
		fmt.Println(cmp)
		partialMerkleRoots, err := cmp.CalculateMerkleRoots()
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, partialMerkleRoots...)
	}
	return merkleRoots, nil
}

func (cmp *CompoundMerklePath) CalculateMerkleRoots() ([]string, error) {
	merkleRoots := make([]string, 0)
	cmpCopy := *cmp

	// Get first layer
	for tx, offset := range cmpCopy[len(cmpCopy)-1] {
		fmt.Println("<--- Offset:  ", offset, "<--- TX:  ", tx)
		// Get leaf
		// Calculate merkle root for one tx
		fmt.Println("Calculate merkle root for one tx")
		merkleRoot, err := calculateMerkleRoot(tx, offset, cmpCopy)
		merkleRoot2, err := calculateMerkleRoot2(tx, offset, cmpCopy)
		merkleRoot3, err := calculateMerkleRoot3(tx, offset, cmpCopy)
		fmt.Println("merkleRoot: ", merkleRoot)
		fmt.Println("merkleRoot2: ", merkleRoot2)
		fmt.Println("merkleRoot3: ", merkleRoot3)
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, merkleRoot)
		merkleRoots = append(merkleRoots, merkleRoot2)
		merkleRoots = append(merkleRoots, merkleRoot3)

	}
	return merkleRoots, nil
}

func calculateMerkleRoot(baseTx string, offset uint64, cmp []map[string]uint64) (string, error) {
	fmt.Println("<------------------")
	fmt.Println("<------------------")
	fmt.Println("<------------------ CALCULATE MERKLE ROOT")
	// Iterate through layers
	for i := len(cmp) - 1; i >= 0; i-- {
		var leftNode, rightNode string
		// Get pair tx for given tx
		fmt.Println("Get pair tx for given tx")
		fmt.Println("tx: ", baseTx)

		newOffset := offset - 1
		if offset%2 == 0 {
			newOffset = offset + 1
		}
		tx2 := keyByValue(cmp[i], newOffset)
		if tx2 == nil {
			fmt.Println("could not find pair")
			return "", errors.New("could not find pair")
		}
		fmt.Println("tx2: ", *tx2)

		if newOffset > offset {
			leftNode = baseTx
			rightNode = *tx2
		} else {
			leftNode = *tx2
			rightNode = baseTx
		}

		fmt.Println("leftNode: ", leftNode)
		fmt.Println("rightNode: ", rightNode)

		str, err := bc.MerkleTreeParentStr(leftNode, rightNode)
		if err != nil {
			fmt.Println("error: ", err)
			return "", err
		}
		baseTx = str
		fmt.Println("New hash: ", baseTx)

		// Reduce offset
		offset = offset / 2
		fmt.Println("New offset: ", offset)
	}

	fmt.Println("Final hash: ", baseTx)

	return baseTx, nil
}

func calculateMerkleRoot2(baseTx string, baseOffset uint64, cmp []map[string]uint64) (string, error) {
	fmt.Println("<------------------")
	fmt.Println("<------------------")
	fmt.Println("<------------------ CALCULATE MERKLE ROOT 2")
	// Iterate through layers
	branches := make([]string, 0)
	offset := baseOffset
	for i := len(cmp) - 1; i >= 0; i-- {
		newOffset := offset - 1
		if offset%2 == 0 {
			newOffset = offset + 1
		}
		tx2 := keyByValue(cmp[i], newOffset)
		if tx2 == nil {
			fmt.Println("could not find pair")
			return "", errors.New("could not find pair")
		}
		fmt.Println("tx2: ", *tx2)
		branches = append(branches, *tx2)

		// Reduce offset
		offset = offset / 2
		fmt.Println("New offset: ", offset)
	}

	merkleRoot, err := bc.MerkleRootFromBranches(baseTx, int(baseOffset), branches)
	if err != nil {
		fmt.Println("error: ", err)
		return "", err
	}

	fmt.Println("Final hash: ", baseTx)

	return merkleRoot, nil
}

func calculateMerkleRoot3(baseTx string, offset uint64, cmp []map[string]uint64) (string, error) {
	fmt.Println("<------------------")
	fmt.Println("<------------------")
	fmt.Println("<------------------ CALCULATE MERKLE ROOT 3")
	// Iterate through layers
	for i := len(cmp) - 1; i >= 0; i-- {
		var leftNode, rightNode string
		// Get pair tx for given tx
		fmt.Println("Get pair tx for given tx")
		fmt.Println("tx: ", baseTx)

		newOffset := offset - 1
		if offset%2 == 0 {
			newOffset = offset + 1
		}
		tx2 := keyByValue(cmp[i], newOffset)
		if tx2 == nil {
			fmt.Println("could not find pair")
			return "", errors.New("could not find pair")
		}
		fmt.Println("tx2: ", *tx2)

		if newOffset > offset {
			leftNode = baseTx
			rightNode = *tx2
		} else {
			leftNode = *tx2
			rightNode = baseTx
		}
		fmt.Println("leftNode: ", leftNode)
		fmt.Println("rightNode: ", rightNode)

		hashBytes := sha256.Sum256([]byte(leftNode + rightNode))
		baseTx = fmt.Sprintf("%x", hashBytes)
		fmt.Println("New hash: ", baseTx)

		// Reduce offset
		offset = offset / 2
		fmt.Println("New offset: ", offset)
	}

	fmt.Println("Final hash: ", baseTx)

	return baseTx, nil
}

func keyByValue(m map[string]uint64, value uint64) *string {
	for k, v := range m {
		if value == v {
			return &k
		}
	}
	return nil
}

func DecodeBEEF(beefHex string) (*DecodedBEEF, error) {
	beefBytes, err := extractBytesWithoutVersionAndMarker(beefHex)
	if err != nil {
		return nil, err
	}

	cmpSlice, remainingBytes, err := decodeCMPSliceFromStream(beefBytes)
	if err != nil {
		return nil, err
	}

	transactions, err := decodeTransactionsWithPathIndexes(remainingBytes)
	if err != nil {
		return nil, err
	}

	if len(transactions) < 2 {
		return nil, errors.New("not enough transactions provided to decode BEEF")
	}

	// get the last transaction as the processed transaction - it should be the last one because of khan's ordering
	processedTx := transactions[len(transactions)-1]

	transactions = transactions[:len(transactions)-1]

	return &DecodedBEEF{
		CMPSlice:        cmpSlice,
		InputsTxData:    transactions,
		ProcessedTxData: processedTx,
	}, nil
}

func decodeCMPSliceFromStream(hexBytes []byte) (CMPSlice, []byte, error) {
	if len(hexBytes) == 0 {
		return nil, nil, errors.New("cannot decode cmp slice from stream - no bytes provided")
	}

	nCMPs, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
	hexBytes = hexBytes[bytesUsed:]

	var cmpPaths []CompoundMerklePath
	for i := 0; i < int(nCMPs); i++ {
		cmp, bytesUsedToDecodeCMP, err := NewCMPFromStream(hexBytes)
		if err != nil {
			return nil, nil, err
		}

		cmpPaths = append(cmpPaths, cmp)
		hexBytes = hexBytes[bytesUsedToDecodeCMP:]
	}

	cmpSlice := CMPSlice(cmpPaths)

	return cmpSlice, hexBytes, nil
}

func NewCMPFromStream(hexBytes []byte) (CompoundMerklePath, int, error) {
	height, bytesUsed, err := extractHeight(hexBytes)
	if err != nil {
		return nil, 0, err
	}
	hexBytes = hexBytes[bytesUsed:]

	var cmp CompoundMerklePath
	currentHeight := height
	bytesUsedToDecodeCMP := bytesUsed

	for currentHeight >= 0 {
		var pathMap map[string]uint64

		pathMap, bytesUsed, err = extractPathMap(hexBytes, currentHeight)
		if err != nil {
			return nil, 0, err
		}

		cmp = append(cmp, pathMap)
		hexBytes = hexBytes[bytesUsed:]

		currentHeight--
		bytesUsedToDecodeCMP += bytesUsed
	}

	return cmp, bytesUsedToDecodeCMP, nil
}

func decodeTransactionsWithPathIndexes(bytes []byte) ([]TxData, error) {
	nTransactions, offset := bt.NewVarIntFromBytes(bytes)
	bytes = bytes[offset:]

	var transactions []TxData

	for i := 0; i < int(nTransactions); i++ {
		tx, offset, err := bt.NewTxFromStream(bytes)
		if err != nil {
			return nil, err
		}
		bytes = bytes[offset:]

		var pathIndex *bt.VarInt
		if bytes[0] == HasCMP {
			value, offset := bt.NewVarIntFromBytes(bytes[1:])
			pathIndex = &value
			bytes = bytes[1+offset:]
		} else if bytes[0] == HasNoCMP {
			bytes = bytes[1:]
		} else {
			return nil, fmt.Errorf("invalid HasCMP flag for transaction at index %d", i)
		}

		transactions = append(transactions, TxData{
			Transaction: tx,
			PathIndex:   pathIndex,
		})
	}

	return transactions, nil
}

func extractHeight(hexBytes []byte) (int, int, error) {
	if len(hexBytes) < 1 {
		return 0, 0, errors.New("insufficient bytes to extract height of compount merkle path")
	}
	height := int(hexBytes[0])
	if height > 64 {
		return 0, 0, errors.New("height exceeds maximum allowed value of 64")
	}
	return height, 1, nil
}

func extractPathMap(hexBytes []byte, height int) (map[string]uint64, int, error) {
	if len(hexBytes) < 1 {
		return nil, 0, fmt.Errorf("insufficient bytes to extract Compound Merkle Path at height %d", height)
	}

	nLeaves, nLeavesBytesUsed := bt.NewVarIntFromBytes(hexBytes)
	bytesUsed := nLeavesBytesUsed
	var pathMap = make(map[string]uint64)

	for i := 0; i < int(nLeaves); i++ {
		if len(hexBytes[bytesUsed:]) < 1 {
			return nil, 0, fmt.Errorf("insufficient bytes to extract index %d leaf of %d leaves at %d height", i, int(nLeaves), height)
		}

		offsetValue, offsetBytesUsed := bt.NewVarIntFromBytes(hexBytes[bytesUsed:])
		bytesUsed += offsetBytesUsed

		if len(hexBytes[bytesUsed:]) < hashBytesCount {
			return nil, 0, fmt.Errorf("insufficient bytes to extract hash of path with offset %d at height %d", offsetValue, height)
		}

		hash := hex.EncodeToString(hexBytes[bytesUsed : bytesUsed+hashBytesCount])
		bytesUsed += hashBytesCount

		pathMap[hash] = uint64(offsetValue)
	}

	return pathMap, bytesUsed, nil
}

func extractBytesWithoutVersionAndMarker(hexStream string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexStream)
	if err != nil {
		return nil, errors.New("invalid beef hex stream")
	}
	if len(bytes) < 4 {
		return nil, errors.New("invalid beef hex stream")
	}

	// removes version bytes
	bytes = bytes[versionBytesCount:]
	err = validateMarker(bytes)
	if err != nil {
		return nil, err
	}

	// removes marker bytes
	bytes = bytes[markerBytesCount:]

	return bytes, nil
}

func validateMarker(bytes []byte) error {
	if bytes[0] != BEEFMarkerPart1 || bytes[1] != BEEFMarkerPart2 {
		return errors.New("invalid format of transaction, BEEF marker not found")
	}

	return nil
}
