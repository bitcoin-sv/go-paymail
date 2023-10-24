package paymail

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/libsv/go-bc"
	"github.com/libsv/go-bt/v2"
)

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
	BUMP            BUMP
	InputsTxData    []TxData
	ProcessedTxData TxData
}

func (dBeef *DecodedBEEF) GetMerkleRoots() ([]string, error) {
	var merkleRoots []string
	for _, cmp := range dBeef.CMPSlice {
		partialMerkleRoots, err := cmp.calculateMerkleRoots()
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, partialMerkleRoots...)
	}
	return merkleRoots, nil
}

func calculateMerkleRoot(baseTx string, offset uint64, cmp []map[string]uint64) (string, error) {
	for i := 0; i < len(cmp); i++ {
		var leftNode, rightNode string
		newOffset := offset - 1
		if offset%2 == 0 {
			newOffset = offset + 1
		}
		tx2 := keyByValue(cmp[i], newOffset)
		if tx2 == nil {
			fmt.Println("could not find pair")
			return "", errors.New("could not find pair")
		}

		if newOffset > offset {
			leftNode = baseTx
			rightNode = *tx2
		} else {
			leftNode = *tx2
			rightNode = baseTx
		}

		// Calculate new merkle tree parent
		str, err := bc.MerkleTreeParentStr(leftNode, rightNode)
		if err != nil {
			return "", err
		}
		baseTx = str

		// Reduce offset
		offset = offset / 2
	}

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

	blockHeight, bytesUsed := bt.NewVarIntFromBytes(beefBytes)
	beefBytes = beefBytes[bytesUsed:]

	bumpSlice, remainingBytes, err := decodeBUMPSliceFromStream(beefBytes)
	if err != nil {
		return nil, err
	}

	bump := BUMP{
		blockHeight: uint64(blockHeight),
		path:        bumpSlice,
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
		BUMP:            bump,
		InputsTxData:    transactions,
		ProcessedTxData: processedTx,
	}, nil
}

func DecodeBUMP(beefHex string) (*DecodedBEEF, error) {
	beefBytes, err := hex.DecodeString(beefHex)
	if err != nil {
		return nil, err
	}

	blockHeight, bytesUsed := bt.NewVarIntFromBytes(beefBytes)
	beefBytes = beefBytes[bytesUsed:]

	bumpSlice, _, err := decodeBUMPSliceFromStream(beefBytes)
	if err != nil {
		return nil, err
	}

	bump := BUMP{
		blockHeight: uint64(blockHeight),
		path:        bumpSlice,
	}

	fmt.Println(bump)

	return nil, nil
}

func decodeBUMPSliceFromStream(hexBytes []byte) (BUMPSlice, []byte, error) {
	if len(hexBytes) == 0 {
		return nil, nil, errors.New("cannot decode cmp slice from stream - no bytes provided")
	}

	treeHeight, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
	hexBytes = hexBytes[bytesUsed:]

	var bumpSlice BUMPSlice
	for i := 0; i < int(treeHeight); i++ {
		nLeaves, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]

		bumpMap, remainingBytes := decodeBUMPLeaves(nLeaves, hexBytes)
		hexBytes = remainingBytes
		bumpSlice = append(bumpSlice, bumpMap)
	}

	return bumpSlice, hexBytes, nil
}

func decodeBUMPLeaves(nLeaves bt.VarInt, hexBytes []byte) (BUMPMap, []byte) {
	bumpMap := make(map[string]BUMPTx)
	for i := 0; i < int(nLeaves); i++ {
		if len(hexBytes) < 1 {
			panic(fmt.Errorf("insufficient bytes to extract offset for %d leaf of %d leaves", i, int(nLeaves)))
		}

		offset, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]

		if len(hexBytes[bytesUsed:]) < 1 {
			panic(fmt.Errorf("insufficient bytes to extract flag for %d leaf of %d leaves", i, int(nLeaves)))
		}

		flag, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]

		if flag == 1 {
			key := strconv.FormatUint(uint64(offset), 10)
			bumpMap[key] = BUMPTx{
				duplicate: true,
			}
			continue
		}

		if len(hexBytes[bytesUsed:]) < hashBytesCount-1 {
			panic("insufficient bytes to extract hash of path")
		}

		hash := hex.EncodeToString(hexBytes[:hashBytesCount])
		bytesUsed += hashBytesCount - 1
		hexBytes = hexBytes[bytesUsed:]
		hash = reverse(hash)

		if flag == 0 {
			bumpMap[strconv.FormatUint(uint64(offset), 10)] = BUMPTx{
				hash: hash,
			}
		} else {
			bumpMap[strconv.FormatUint(uint64(offset), 10)] = BUMPTx{
				hash: hash,
				txId: true,
			}
		}
	}

	return BUMPMap{bumpMap}, hexBytes
}

// reverse will reverse a hex string but it takes a pair of character at a time
func reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+2, j-2 {
		rns[i], rns[j], rns[i+1], rns[j-1] = rns[j-1], rns[i+1], rns[j], rns[i]
	}
	return string(rns)
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
