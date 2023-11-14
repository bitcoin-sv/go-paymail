package paymail

import (
	"context"
	"errors"
	"testing"

	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of a service provider
type mockServiceProvider struct{}

// VerifyMerkleRoots is a mock implementation of this interface
func (m *mockServiceProvider) VerifyMerkleRoots(_ context.Context, _ []MerkleRootConfirmationRequestItem) error {
	// Verify the merkle roots
	return nil
}

func TestDecodeBEEF_DecodeBEEF_HappyPaths(t *testing.T) {
	testCases := []struct {
		name                       string
		hexStream                  string
		expectedDecodedBEEF        *DecodedBEEF
		pathIndexForTheOldestInput *bt.VarInt
		expectedError              error
	}{
		{
			name:      "valid BEEF with 1 CMP and 1 input transaction",
			hexStream: "0100beef01fe636d0c0007021400fe507c0c7aa754cef1f7889d5fd395cf1f785dd7de98eed895dbedfe4e5bc70d1502ac4e164f5bc16746bb0868404292ac8318bbac3800e4aad13a014da427adce3e010b00bc4ff395efd11719b277694cface5aa50d085a0bb81f613f70313acd28cf4557010400574b2d9142b8d28b61d88e3b2c3f44d858411356b49a28a4643b6d1a6a092a5201030051a05fc84d531b5d250c23f4f886f6812f9fe3f402d61607f977b4ecd2701c19010000fd781529d58fc2523cf396a7f25440b409857e7e221766c57214b1d38c7b481f01010062f542f45ea3660f86c013ced80534cb5fd4c19d66c56e7e8c5d4bf2d40acc5e010100b121e91836fd7cd5102b654e9f72f3cf6fdbfd0b161c53a9c54b12c841126331020100000001cd4e4cac3c7b56920d1e7655e7e260d31f29d9a388d04910f1bbd72304a79029010000006b483045022100e75279a205a547c445719420aa3138bf14743e3f42618e5f86a19bde14bb95f7022064777d34776b05d816daf1699493fcdf2ef5a5ab1ad710d9c97bfb5b8f7cef3641210263e2dee22b1ddc5e11f6fab8bcd2378bdd19580d640501ea956ec0e786f93e76ffffffff013e660000000000001976a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac0000000001000100000001ac4e164f5bc16746bb0868404292ac8318bbac3800e4aad13a014da427adce3e000000006a47304402203a61a2e931612b4bda08d541cfb980885173b8dcf64a3471238ae7abcd368d6402204cbf24f04b9aa2256d8901f0ed97866603d2be8324c2bfb7a37bf8fc90edd5b441210263e2dee22b1ddc5e11f6fab8bcd2378bdd19580d640501ea956ec0e786f93e76ffffffff013c660000000000001976a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac0000000000",
			expectedDecodedBEEF: &DecodedBEEF{
				BUMPs: BUMPs{
					BUMP{
						blockHeight: 814435,
						path: [][]BUMPLeaf{
							{
								BUMPLeaf{hash: "0dc75b4efeeddb95d8ee98ded75d781fcf95d35f9d88f7f1ce54a77a0c7c50fe", offset: 20},
								BUMPLeaf{hash: "3ecead27a44d013ad1aae40038acbb1883ac9242406808bb4667c15b4f164eac", txId: true, offset: 21},
							},
							{
								BUMPLeaf{hash: "5745cf28cd3a31703f611fb80b5a080da55acefa4c6977b21917d1ef95f34fbc", offset: 11},
							},
							{
								BUMPLeaf{hash: "522a096a1a6d3b64a4289ab456134158d8443f2c3b8ed8618bd2b842912d4b57", offset: 4},
							},
							{
								BUMPLeaf{hash: "191c70d2ecb477f90716d602f4e39f2f81f686f8f4230c255d1b534dc85fa051", offset: 3},
							},
							{
								BUMPLeaf{hash: "1f487b8cd3b11472c56617227e7e8509b44054f2a796f33c52c28fd5291578fd", offset: 0},
							},
							{
								BUMPLeaf{hash: "5ecc0ad4f24b5d8c7e6ec5669dc1d45fcb3405d8ce13c0860f66a35ef442f562", offset: 1},
							},
							{
								BUMPLeaf{hash: "31631241c8124bc5a9531c160bfddb6fcff3729f4e652b10d57cfd3618e921b1", offset: 1},
							},
						},
					},
				},
				InputsTxData: []*TxData{
					{
						Transaction: &bt.Tx{
							Version:  1,
							LockTime: 0,
							Inputs: []*bt.Input{
								{
									PreviousTxSatoshis: 0,
									PreviousTxOutIndex: 1,
									SequenceNumber:     4294967295,
									PreviousTxScript:   nil,
								},
							},
							Outputs: []*bt.Output{
								{
									Satoshis:      26174,
									LockingScript: bscript.NewFromBytes([]byte("76a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac")),
								},
							},
						},
						PathIndex: func(v bt.VarInt) *bt.VarInt { return &v }(0x0),
					},
				},
				ProcessedTxData: &bt.Tx{
					Version:  1,
					LockTime: 0,
					Inputs: []*bt.Input{
						{
							PreviousTxSatoshis: 0,
							PreviousTxOutIndex: 0,
							SequenceNumber:     4294967295,
							PreviousTxScript:   nil,
						},
					},
					Outputs: []*bt.Output{
						{
							Satoshis:      26172,
							LockingScript: bscript.NewFromBytes([]byte("76a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac")),
						},
					},
				},
			},
			pathIndexForTheOldestInput: func(v bt.VarInt) *bt.VarInt { return &v }(0x0),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			beef := tc.hexStream

			// when
			decodedBEEF, err := DecodeBEEF(beef)

			// then
			assert.Equal(t, tc.expectedError, err, "expected error %v, but got %v", tc.expectedError, err)

			assert.Equal(t, len(tc.expectedDecodedBEEF.InputsTxData), len(decodedBEEF.InputsTxData), "expected %v inputs, but got %v", len(tc.expectedDecodedBEEF.InputsTxData), len(decodedBEEF.InputsTxData))

			assert.Equal(t, len(tc.expectedDecodedBEEF.BUMPs), len(decodedBEEF.BUMPs), "expected %v BUMPs, but got %v", len(tc.expectedDecodedBEEF.BUMPs), len(decodedBEEF.BUMPs))

			for i, bump := range tc.expectedDecodedBEEF.BUMPs {
				assert.Equal(t, len(bump.path), len(decodedBEEF.BUMPs[i].path), "expected %v BUMPPaths for %v BUMP, but got %v", len(bump.path), i, len(decodedBEEF.BUMPs[i].path))
				assert.Equal(t, bump.path, decodedBEEF.BUMPs[i].path, "expected equal BUMPPaths for %v BUMP, expected: %v but got %v", i, bump, len(decodedBEEF.BUMPs[i].path))
			}

			assert.NotNil(t, decodedBEEF.ProcessedTxData, "expected original transaction to be not nil")

			assert.Equal(t, tc.expectedDecodedBEEF.InputsTxData[0].PathIndex, decodedBEEF.InputsTxData[0].PathIndex, "expected path index for the oldest input to be %v, but got %v", tc.expectedDecodedBEEF.InputsTxData[0].PathIndex, decodedBEEF.InputsTxData[0].PathIndex)
		})
	}
}

func TestDecodeBEEF_DecodeBEEF_HandlingErrors(t *testing.T) {
	testCases := []struct {
		name                         string
		hexStream                    string
		expectedDecodedBEEF          *DecodedBEEF
		expectedCMPForTheOldestInput bool
		expectedError                error
	}{
		{
			name:                         "too short hex stream",
			hexStream:                    "001",
			expectedDecodedBEEF:          nil,
			expectedError:                errors.New("invalid beef hex stream"),
			expectedCMPForTheOldestInput: false,
		},
		{
			name:                         "unable to decode BEEF - only marker and version has been provided",
			hexStream:                    "0100beef",
			expectedDecodedBEEF:          nil,
			expectedError:                errors.New("cannot decode BUMP - no bytes provided"),
			expectedCMPForTheOldestInput: false,
		},
		{
			name:                         "unable to decode BEEF - wrong marker",
			hexStream:                    "0100efbe",
			expectedDecodedBEEF:          nil,
			expectedError:                errors.New("invalid format of transaction, BEEF marker not found"),
			expectedCMPForTheOldestInput: false,
		},
		{
			name:                "unable to decode BUMP block height - proper BEEF marker and number of bumps",
			hexStream:           "0100beef01",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract BUMP blockHeight"),
		},
		{
			name:                "unable to decode BUMP tree height - proper BEEF marker, number of bumps and block height",
			hexStream:           "0100beef01fe8a6a0c00",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("cannot decode BUMP paths from stream - no bytes provided"),
		},
		{
			name:                "unable to decode BUMP number of leaves - proper BEEF marker, number of bumps, block height and tree height but end of stream at this point",
			hexStream:           "0100beef01fe8a6a0c000c",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("cannot decode BUMP paths number of leaves from stream - no bytes provided"),
		},
		{
			name:                "unable to decode BUMP leaf - no offset - proper BEEF marker, number of bumps, block height and tree height and nLeaves but end of stream at this point",
			hexStream:           "0100beef01fe8a6a0c000c04",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract offset for 0 leaf of 4 leaves"),
		},
		{
			name:                "unable to decode BUMP leaf - no flag - proper BEEF marker, number of bumps, block height and tree height, nLeaves and offset but end of stream at this point",
			hexStream:           "0100beef01fe8a6a0c000c04fde80b",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract flag for 0 leaf of 4 leaves"),
		},
		{
			name:                "unable to decode BUMP leaf - wrong flag - proper BEEF marker, number of bumps, block height and tree height, nLeaves and offset",
			hexStream:           "0100beef01fe8a6a0c000c04fde80b03",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("invalid flag: 3 for 0 leaf of 4 leaves"),
		},
		{
			name:                "unable to decode BUMP leaf - no hash with flag 0 - proper BEEF marker, number of bumps, block height and tree height, nLeaves, offset and flag",
			hexStream:           "0100beef01fe8a6a0c000c04fde80b00",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract hash of path"),
		},
		{
			name:                "unable to decode BUMP leaf - no hash with flag 2 - proper BEEF marker, number of bumps, block height and tree height, nLeaves, offset and flag",
			hexStream:           "0100beef01fe8a6a0c000c04fde80b00",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract hash of path"),
		},
		{
			name:                "unable to decode BUMP leaf - flag 1 - proper BEEF marker, number of bumps, block height and tree height, nLeaves, offset and flag but end of stream at this point - flag 1 means that there is no hash",
			hexStream:           "0100beef01fe8a6a0c000c04fde80b01",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract offset for 1 leaf of 4 leaves"),
		},
		{
			name:                "unable to decode BUMP leaf - not enough bytes for hash - proper BEEF marker, number of bumps, block height and tree height, nLeaves, offset and flag but with not enough bytes for hash",
			hexStream:           "0100beef01fe8a6a0c000c04fde80b0011774f01d26412f0d16ea3f0447be0b5ebec67b0782e321a7a01cbdf7f734e",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract hash of path"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			beef := tc.hexStream

			// when
			result, err := DecodeBEEF(beef)

			// then
			assert.Equal(t, tc.expectedError, err, "expected error %v, but got %v", tc.expectedError, err)
			assert.Nil(t, result, "expected nil result, but got %v", result)
		})
	}
}
