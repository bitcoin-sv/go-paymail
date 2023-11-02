package paymail

import (
	"encoding/hex"
	"errors"
	"fmt"
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
	BUMPs           BUMPs
	InputsTxData    []TxData
	ProcessedTxData TxData
}

// GetMerkleRoots will calculate the merkle roots for the BUMPs in the BEEF transaction
func (dBeef *DecodedBEEF) GetMerkleRoots() ([]string, error) {
	var merkleRoots []string
	for _, bump := range dBeef.BUMPs {
		partialMerkleRoots, err := bump.calculateMerkleRoots()
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, partialMerkleRoots...)
	}
	return merkleRoots, nil
}

func DecodeBEEF(beefHex string) (*DecodedBEEF, error) {
	beefBytes, err := extractBytesWithoutVersionAndMarker(beefHex)
	if err != nil {
		return nil, err
	}

	bumps, remainingBytes, err := decodeBUMPs(beefBytes)
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
		BUMPs:           bumps,
		InputsTxData:    transactions,
		ProcessedTxData: processedTx,
	}, nil
}

func decodeBUMPs(beefBytes []byte) ([]BUMP, []byte, error) {
	if len(beefBytes) == 0 {
		return nil, nil, errors.New("cannot decode BUMP - no bytes provided")
	}

	bumps := make([]BUMP, 0)
	nBump, bytesUsed := bt.NewVarIntFromBytes(beefBytes)
	beefBytes = beefBytes[bytesUsed:]

	for i := 0; i < int(nBump); i++ {
		if len(beefBytes) == 0 {
			return nil, nil, errors.New("insufficient bytes to extract BUMP blockHeight")
		}
		blockHeight, bytesUsed := bt.NewVarIntFromBytes(beefBytes)
		beefBytes = beefBytes[bytesUsed:]
		bumpPaths, remainingBytes, err := decodeBUMPPathsFromStream(beefBytes)
		if err != nil {
			return nil, nil, err
		}
		beefBytes = remainingBytes

		bump := BUMP{
			blockHeight: uint64(blockHeight),
			path:        bumpPaths,
		}

		bumps = append(bumps, bump)
	}

	return bumps, beefBytes, nil
}

func decodeBUMPPathsFromStream(hexBytes []byte) ([][]BUMPLeaf, []byte, error) {
	if len(hexBytes) == 0 {
		return nil, nil, errors.New("cannot decode BUMP paths from stream - no bytes provided")
	}

	treeHeight, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
	hexBytes = hexBytes[bytesUsed:]
	bumpPaths := make([][]BUMPLeaf, 0)

	for i := 0; i < int(treeHeight); i++ {
		if len(hexBytes) == 0 {
			return nil, nil, errors.New("cannot decode BUMP paths number of leaves from stream - no bytes provided")
		}
		nLeaves, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]
		bumpPath, remainingBytes, err := decodeBUMPLevel(nLeaves, hexBytes)
		if err != nil {
			return nil, nil, err
		}
		hexBytes = remainingBytes
		bumpPaths = append(bumpPaths, bumpPath)
	}

	return bumpPaths, hexBytes, nil
}

func decodeBUMPLevel(nLeaves bt.VarInt, hexBytes []byte) ([]BUMPLeaf, []byte, error) {
	bumpPath := make([]BUMPLeaf, 0)
	for i := 0; i < int(nLeaves); i++ {
		if len(hexBytes) == 0 {
			return nil, nil, fmt.Errorf("insufficient bytes to extract offset for %d leaf of %d leaves", i, int(nLeaves))
		}

		offset, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]

		if len(hexBytes) == 0 {
			return nil, nil, fmt.Errorf("insufficient bytes to extract flag for %d leaf of %d leaves", i, int(nLeaves))
		}

		flag, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]

		if flag != dataFlag && flag != duplicateFlag && flag != txIDFlag {
			return nil, nil, fmt.Errorf("invalid flag: %d for %d leaf of %d leaves", flag, i, int(nLeaves))
		}

		if flag == duplicateFlag {
			bumpLeaf := BUMPLeaf{
				offset:    uint64(offset),
				duplicate: true,
			}
			bumpPath = append(bumpPath, bumpLeaf)
			continue
		}

		if len(hexBytes) < hashBytesCount {
			return nil, nil, errors.New("insufficient bytes to extract hash of path")
		}

		hash := hex.EncodeToString(bt.ReverseBytes(hexBytes[:hashBytesCount]))
		bytesUsed += hashBytesCount - 1
		hexBytes = hexBytes[bytesUsed:]

		bumpLeaf := BUMPLeaf{
			hash:   hash,
			offset: uint64(offset),
		}
		if flag == txIDFlag {
			bumpLeaf.txId = true
		}
		bumpPath = append(bumpPath, bumpLeaf)
	}

	return bumpPath, hexBytes, nil
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
