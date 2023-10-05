package paymail

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/libsv/bitcoin-hc/transports/http/endpoints/api/merkleroots"
	"github.com/libsv/go-bc"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
)

type CompoundMerklePath []map[string]uint64

type CMPSlice []CompoundMerklePath

type MerkleRootVerifier interface {
	VerifyMerkleRoots(
		ctx context.Context,
		merkleRoots []string,
	) (*merkleroots.MerkleRootsConfirmationsResponse, error)
}

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
		partialMerkleRoots, err := cmp.calculateMerkleRoots()
		if err != nil {
			return nil, err
		}
		merkleRoots = append(merkleRoots, partialMerkleRoots...)
	}
	return merkleRoots, nil
}

// ExecuteSimplifiedPaymentVerification executes the SPV for decoded BEEF tx
func (dBeef *DecodedBEEF) ExecuteSimplifiedPaymentVerification(provider MerkleRootVerifier) error {
	err := dBeef.satoshisInInputsGreaterThanZero()
	if err != nil {
		return err
	}

	err = dBeef.satoshisInOutputsGreaterThanZero()
	if err != nil {
		return err
	}

	err = dBeef.validateSatoshisSum()
	if err != nil {
		return err
	}

	err = dBeef.validateLockTime()
	if err != nil {
		return err
	}

	err = dBeef.validateScripts()
	if err != nil {
		return err
	}

	err = dBeef.verifyMerkleRoots(provider)
	if err != nil {
		return err
	}

	return nil
}

func (dBeef *DecodedBEEF) satoshisInOutputsGreaterThanZero() error {
	if len(dBeef.ProcessedTxData.Transaction.Outputs) == 0 {
		return errors.New("invalid output, no outputs")
	}
	return nil
}

func (dBeef *DecodedBEEF) satoshisInInputsGreaterThanZero() error {
	if len(dBeef.ProcessedTxData.Transaction.Inputs) == 0 {
		return errors.New("invalid input, no inputs")
	}
	return nil
}

func (dBeef *DecodedBEEF) verifyMerkleRoots(provider MerkleRootVerifier) error {
	merkleRoots, err := dBeef.GetMerkleRoots()
	if err != nil {
		return err
	}

	res, err := provider.VerifyMerkleRoots(context.Background(), merkleRoots)
	if err != nil {
		return err
	}

	if !res.AllConfirmed {
		return errors.New("not all merkle roots were confirmed")
	}
	return nil
}

func (dBeef *DecodedBEEF) validateScripts() error {
	for _, input := range dBeef.ProcessedTxData.Transaction.Inputs {
		txId := input.PreviousTxID()
		for j, input2 := range dBeef.InputsTxData {
			if input2.Transaction.TxID() == string(txId) {
				result := verifyScripts(dBeef.ProcessedTxData.Transaction, input2.Transaction, j)
				if !result {
					return errors.New("invalid script")
				}
				break
			}
		}
	}
	return nil
}

func (dBeef *DecodedBEEF) validateSatoshisSum() error {
	inputSum, outputSum := uint64(0), uint64(0)
	for i, input := range dBeef.ProcessedTxData.Transaction.Inputs {
		input2 := dBeef.InputsTxData[i]
		inputSum += input2.Transaction.Outputs[input.PreviousTxOutIndex].Satoshis
	}
	for _, output := range dBeef.ProcessedTxData.Transaction.Outputs {
		outputSum += output.Satoshis
	}

	if inputSum <= outputSum {
		return errors.New("invalid input and output sum, outputs can not be larger than inputs")
	}
	return nil
}

func (dBeef *DecodedBEEF) validateLockTime() error {
	if dBeef.ProcessedTxData.Transaction.LockTime == 0 {
		for _, input := range dBeef.ProcessedTxData.Transaction.Inputs {
			if input.SequenceNumber != 0xffffffff {
				return errors.New("invalid sequence")
			}
		}
	} else {
		return errors.New("invalid locktime")
	}
	return nil
}

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

// Verify locking and unlocking scripts pair
func verifyScripts(tx, prevTx *bt.Tx, inputIdx int) bool {
	input := tx.InputIdx(inputIdx)
	prevOutput := prevTx.OutputIdx(int(input.PreviousTxOutIndex))

	if err := interpreter.NewEngine().Execute(
		interpreter.WithTx(tx, inputIdx, prevOutput),
		interpreter.WithForkID(),
		interpreter.WithAfterGenesis(),
	); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func calculateMerkleRoot(baseTx string, offset uint64, cmp []map[string]uint64) (string, error) {
	for i := len(cmp) - 1; i >= 0; i-- {
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
