package paymail

import (
	"encoding/hex"
	"errors"
	"fmt"

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

type TxData struct {
	Transaction *bt.Tx
	PathIndex   *bt.VarInt
}

type DecodedBEEF struct {
	CMPSlice        CMPSlice
	InputsTxData    []TxData
	ProcessedTxData TxData
}

func DecodeBEEF(beefHex string) (*DecodedBEEF, error) {
	beefBytes, err := extractBytesWithoutVersionAndMarker(beefHex)
	if err != nil {
		return nil, err
	}

	cmpSlice, remainingBytes, err := DecodeCMPSliceFromStream(beefBytes)
	if err != nil {
		return nil, err
	}

	transactions, err := DecodeTransactionsWithPathIndexes(remainingBytes)
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions found")
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

func DecodeCMPSliceFromStream(hexBytes []byte) (CMPSlice, []byte, error) {
	if len(hexBytes) == 0 {
		return nil, nil, errors.New("provided hexStream is empty")
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
	previousHeight := height
	bytesUsedToDecodeCMP := bytesUsed

	for previousHeight >= 0 {
		var pathMap map[string]uint64

		pathMap, bytesUsed, err = extractPathMap(hexBytes, previousHeight)
		if err != nil {
			return nil, 0, err
		}

		cmp = append(cmp, pathMap)
		hexBytes = hexBytes[bytesUsed:]

		previousHeight--
		bytesUsedToDecodeCMP += bytesUsed
	}

	return cmp, bytesUsedToDecodeCMP, nil
}

func DecodeTransactionsWithPathIndexes(bytes []byte) ([]TxData, error) {
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
			return nil, errors.New("invalid HasCMP flag")
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

	if height < 0 {
		return nil, 0, fmt.Errorf("unexpected negative value of height %d", height)
	}

	nLeaves, nLeavesBytesUsed := bt.NewVarIntFromBytes(hexBytes)
	bytesUsed := nLeavesBytesUsed
	var pathMap = make(map[string]uint64)

	for i := 0; i < int(nLeaves); i++ {
		if len(hexBytes[bytesUsed:]) < 1 {
			return nil, 0, fmt.Errorf("insufficient bytes to extract %d path of %d paths at %d height", i, int(nLeaves), height)
		}

		offsetValue, offsetBytesUsed := bt.NewVarIntFromBytes(hexBytes[bytesUsed:])
		bytesUsed += offsetBytesUsed

		if len(hexBytes[bytesUsed:]) < 32 {
			return nil, 0, fmt.Errorf("insufficient bytes to extract hash of path with offset %d at height %d", offsetValue, height)
		}

		hash := hex.EncodeToString(hexBytes[bytesUsed : bytesUsed+32])
		bytesUsed += 32

		pathMap[hash] = uint64(offsetValue)
	}

	return pathMap, bytesUsed, nil
}

func extractBytesWithoutVersionAndMarker(hexStream string) ([]byte, error) {
	versionBytesCount := 2
	markerBytesCount := 2

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
		return errors.New("invalid beef marker")
	}

	return nil
}
