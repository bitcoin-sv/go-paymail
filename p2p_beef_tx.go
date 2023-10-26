package paymail

import (
	"encoding/hex"
	"errors"
	"fmt"
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
	BUMPs           BUMPPaths
	InputsTxData    []TxData
	ProcessedTxData TxData
}

// GetMerkleRoots will calculate the merkle roots for the BUMPs in the BEEF transaction
func (dBeef *DecodedBEEF) GetMerkleRoots() ([]string, error) {
	var merkleRoots []string
	for _, bump := range dBeef.BUMPs {
		partialMerkleRoots, err := bump.calculateMerkleRoots()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		merkleRoots = append(merkleRoots, partialMerkleRoots...)
	}
	return merkleRoots, nil
}

// calculateMerkleRoots will calculate one merkle root for tx in the BUMPPath
func calculateMerkleRoot(baseTx BUMPPathElement, bump BUMP) (string, error) {
	calculatedHash := baseTx.hash
	offset := baseTx.offset

	for _, bLevel := range bump.path {
		newOffset := offset - 1
		if offset%2 == 0 {
			newOffset = offset + 1
		}
		tx2 := bLevel.findTxByOffset(newOffset)
		if &tx2 == nil {
			return "", errors.New("could not find pair")
		}

		leftNode, rightNode := prepareNodes(baseTx, offset, *tx2, newOffset)

		str, err := bc.MerkleTreeParentStr(leftNode, rightNode)
		if err != nil {
			return "", err
		}
		calculatedHash = str

		offset = offset / 2

		baseTx = BUMPPathElement{
			hash:   calculatedHash,
			offset: offset,
		}
	}

	return calculatedHash, nil
}

func prepareNodes(baseTx BUMPPathElement, offset uint64, tx2 BUMPPathElement, newOffset uint64) (string, string) {
	var txHash, tx2Hash string

	if baseTx.duplicate {
		txHash = tx2.hash
	} else {
		txHash = baseTx.hash
	}

	if tx2.duplicate {
		tx2Hash = baseTx.hash
	} else {
		tx2Hash = tx2.hash
	}

	if newOffset > offset {
		return txHash, tx2Hash
	}
	return tx2Hash, txHash
}

func DecodeBEEF(beefHex string) (*DecodedBEEF, error) {
	beefBytes, err := extractBytesWithoutVersionAndMarker(beefHex)
	if err != nil {
		return nil, err
	}

	bumps, remainingBytes, err := decodeBUMPs(beefBytes)
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
	bumps := make([]BUMP, 0)
	nBump, bytesUsed := bt.NewVarIntFromBytes(beefBytes)
	beefBytes = beefBytes[bytesUsed:]

	for i := 0; i < int(nBump); i++ {
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

func decodeBUMPPathsFromStream(hexBytes []byte) ([]BUMPPath, []byte, error) {
	if len(hexBytes) == 0 {
		return nil, nil, errors.New("cannot decode cmp slice from stream - no bytes provided")
	}

	treeHeight, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
	hexBytes = hexBytes[bytesUsed:]

	fmt.Println("treeHeight", treeHeight)
	fmt.Println("bytesUsed", bytesUsed)
	fmt.Println()

	bumpPaths := make([]BUMPPath, 0)

	for i := 0; i < int(treeHeight); i++ {
		nLeaves, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]

		fmt.Println("nLeaves", nLeaves)
		fmt.Println("bytesUsed", bytesUsed)
		fmt.Println("<- decoding bump level")

		bumpPath, remainingBytes := decodeBUMPPath(nLeaves, hexBytes)
		hexBytes = remainingBytes
		bumpPaths = append(bumpPaths, bumpPath)
	}

	return bumpPaths, hexBytes, nil
}

func decodeBUMPPath(nLeaves bt.VarInt, hexBytes []byte) (BUMPPath, []byte) {
	var bumpPath BUMPPath
	for i := 0; i < int(nLeaves); i++ {
		fmt.Println("<-------------------------------------------------------------------------- decoding bump tx")
		if len(hexBytes) < 1 {
			panic(fmt.Errorf("insufficient bytes to extract offset for %d leaf of %d leaves", i, int(nLeaves)))
		}

		offset, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]

		fmt.Println("offset", offset)
		fmt.Println("bytesUsed", bytesUsed)
		fmt.Println("hexBytes", hexBytes)
		fmt.Println()

		if len(hexBytes[bytesUsed:]) < 1 {
			panic(fmt.Errorf("insufficient bytes to extract flag for %d leaf of %d leaves", i, int(nLeaves)))
		}

		flag, bytesUsed := bt.NewVarIntFromBytes(hexBytes)
		hexBytes = hexBytes[bytesUsed:]

		fmt.Println("flag", flag)
		fmt.Println("bytesUsed", bytesUsed)
		fmt.Println("hexBytes", hexBytes)
		fmt.Println()

		if flag == 1 {
			bTx := BUMPPathElement{
				offset:    uint64(offset),
				duplicate: true,
			}
			bumpPath = append(bumpPath, bTx)
			continue
		}

		if len(hexBytes[bytesUsed:]) < hashBytesCount-1 {
			panic("insufficient bytes to extract hash of path")
		}

		hash := hex.EncodeToString(hexBytes[:hashBytesCount])
		bytesUsed += hashBytesCount - 1
		hexBytes = hexBytes[bytesUsed:]
		hash = reverse(hash)

		fmt.Println("hash", hash)
		fmt.Println("bytesUsed", bytesUsed)
		fmt.Println("hexBytes", hexBytes)
		fmt.Println()

		if flag == 0 {
			bTx := BUMPPathElement{
				hash:   hash,
				offset: uint64(offset),
			}
			bumpPath = append(bumpPath, bTx)
		} else {
			bTx := BUMPPathElement{
				hash:   hash,
				txId:   true,
				offset: uint64(offset),
			}
			bumpPath = append(bumpPath, bTx)
		}
	}

	return bumpPath, hexBytes
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
