package beef

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
	HasNoBump = 0x00
	HasBump   = 0x01
)

const (
	hashBytesCount    = 32
	markerBytesCount  = 2
	versionBytesCount = 2
	maxTreeHeight     = 64
)

type TxData struct {
	Transaction *bt.Tx
	BumpIndex   *bt.VarInt

	txID string
}

func (td *TxData) Unmined() bool {
	return td.BumpIndex == nil
}

func (td *TxData) GetTxID() string {
	if len(td.txID) == 0 {
		td.txID = td.Transaction.TxID()
	}

	return td.txID
}

type DecodedBEEF struct {
	BUMPs        BUMPs
	Transactions []*TxData
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

	return &DecodedBEEF{
		BUMPs:        bumps,
		Transactions: transactions,
	}, nil
}

func (d *DecodedBEEF) GetLatestTx() *bt.Tx {
	return d.Transactions[len(d.Transactions)-1].Transaction // get the last transaction as the processed transaction - it should be the last one because of khan's ordering
}

func decodeBUMPs(beefBytes []byte) ([]*BUMP, []byte, error) {
	if len(beefBytes) == 0 {
		return nil, nil, errors.New("cannot decode BUMP - no bytes provided")
	}

	nBump, bytesUsed := bt.NewVarIntFromBytes(beefBytes)

	if nBump == 0 {
		return nil, nil, errors.New("invalid BEEF- lack of BUMPs")
	}

	beefBytes = beefBytes[bytesUsed:]

	bumps := make([]*BUMP, 0, uint64(nBump))
	for i := uint64(0); i < uint64(nBump); i++ {
		if len(beefBytes) == 0 {
			return nil, nil, errors.New("insufficient bytes to extract BUMP blockHeight")
		}
		blockHeight, bytesUsed := bt.NewVarIntFromBytes(beefBytes)
		beefBytes = beefBytes[bytesUsed:]

		treeHeight := beefBytes[0]
		if int(treeHeight) > maxTreeHeight {
			return nil, nil, fmt.Errorf("invalid BEEF - treeHeight cannot be grater than %d", maxTreeHeight)
		}
		beefBytes = beefBytes[1:]

		bumpPaths, remainingBytes, err := decodeBUMPPathsFromStream(int(treeHeight), beefBytes)
		if err != nil {
			return nil, nil, err
		}
		beefBytes = remainingBytes

		bump := &BUMP{
			BlockHeight: uint64(blockHeight),
			Path:        bumpPaths,
		}

		bumps = append(bumps, bump)
	}

	return bumps, beefBytes, nil
}

func decodeBUMPPathsFromStream(treeHeight int, hexBytes []byte) ([][]BUMPLeaf, []byte, error) {
	bumpPaths := make([][]BUMPLeaf, 0)

	for i := 0; i < treeHeight; i++ {
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

		flag := hexBytes[0]
		hexBytes = hexBytes[1:]

		if flag != dataFlag && flag != duplicateFlag && flag != txIDFlag {
			return nil, nil, fmt.Errorf("invalid flag: %d for %d leaf of %d leaves", flag, i, int(nLeaves))
		}

		if flag == duplicateFlag {
			bumpLeaf := BUMPLeaf{
				Offset:    uint64(offset),
				Duplicate: true,
			}
			bumpPath = append(bumpPath, bumpLeaf)
			continue
		}

		if len(hexBytes) < hashBytesCount {
			return nil, nil, errors.New("insufficient bytes to extract hash of path")
		}

		hash := hex.EncodeToString(bt.ReverseBytes(hexBytes[:hashBytesCount]))
		hexBytes = hexBytes[hashBytesCount:]

		bumpLeaf := BUMPLeaf{
			Hash:   hash,
			Offset: uint64(offset),
		}
		if flag == txIDFlag {
			bumpLeaf.TxId = true
		}
		bumpPath = append(bumpPath, bumpLeaf)
	}

	return bumpPath, hexBytes, nil
}

func decodeTransactionsWithPathIndexes(bytes []byte) ([]*TxData, error) {
	nTransactions, offset := bt.NewVarIntFromBytes(bytes)

	if nTransactions < 2 {
		return nil, errors.New("invalid BEEF- not enough transactions provided to decode BEEF")
	}

	bytes = bytes[offset:]

	transactions := make([]*TxData, 0, int(nTransactions))

	for i := 0; i < int(nTransactions); i++ {
		tx, offset, err := bt.NewTxFromStream(bytes)
		if err != nil {
			return nil, err
		}
		bytes = bytes[offset:]

		var pathIndex *bt.VarInt

		if bytes[0] == HasBump {
			value, offset := bt.NewVarIntFromBytes(bytes[1:])
			pathIndex = &value
			bytes = bytes[1+offset:]
		} else if bytes[0] == HasNoBump {
			bytes = bytes[1:]
		} else {
			return nil, fmt.Errorf("invalid HasCMP flag for transaction at index %d", i)
		}

		transactions = append(transactions, &TxData{
			Transaction: tx,
			BumpIndex:   pathIndex,
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
