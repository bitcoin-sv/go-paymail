package paymail

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock implementation of a service provider
type mockServiceProvider struct{}

// VerifyMerkleRoots is a mock implementation of this interface
func (m *mockServiceProvider) VerifyMerkleRoots(_ context.Context, _ []string) error {
	// Verify the merkle roots
	return nil
}

// func TestDecodeBEEF_DecodeBEEF_HappyPaths(t *testing.T) {
// 	testCases := []struct {
// 		name                       string
// 		hexStream                  string
// 		expectedDecodedBEEF        *DecodedBEEF
// 		pathIndexForTheOldestInput *bt.VarInt
// 		expectedError              error
// 	}{
// 		{
// 			name:      "valid BEEF with 1 CMP and 1 input transaction",
// 			hexStream: "0100beef01020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b02020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac000000000100020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000",
// 			expectedDecodedBEEF: &DecodedBEEF{
// 				CMPSlice: CMPSlice{
// 					{
// 						{
// 							"8c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee8": 0x0,
// 							"da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb": 0x1,
// 							"b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef": 0x2,
// 							"e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b": 0x3,
// 						},
// 						{
// 							"3470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e": 0x0, "c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c": 0x1,
// 						},
// 						{"cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b": 0x1},
// 					},
// 				},
// 				InputsTxData: []TxData{
// 					{
// 						Transaction: &bt.Tx{
// 							Version:  2,
// 							LockTime: 0,
// 							Inputs: []*bt.Input{
// 								{
// 									PreviousTxSatoshis: 0,
// 									PreviousTxOutIndex: 0,
// 									SequenceNumber:     4294967295,
// 									PreviousTxScript:   nil,
// 								},
// 							},
// 							Outputs: []*bt.Output{
// 								{
// 									Satoshis: 1,
// 								},
// 							}},
// 						PathIndex: func(v bt.VarInt) *bt.VarInt { return &v }(0x0),
// 					},
// 				},
// 				ProcessedTxData: TxData{
// 					Transaction: &bt.Tx{
// 						Version:  2,
// 						LockTime: 0,
// 						Inputs: []*bt.Input{
// 							{
// 								PreviousTxSatoshis: 0,
// 								PreviousTxOutIndex: 0,
// 								SequenceNumber:     4294967295,
// 								PreviousTxScript:   nil,
// 							},
// 						},
// 						Outputs: []*bt.Output{
// 							{
// 								Satoshis: 1,
// 							},
// 						}},
// 					PathIndex: nil,
// 				},
// 			},
// 			pathIndexForTheOldestInput: func(v bt.VarInt) *bt.VarInt { return &v }(0x0),
// 		},
// 		{
// 			name:      "valid BEEF with 2 CMP and 2 input transaction - all input transactions have no CMP flag set",
// 			hexStream: "0100beef02020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b03020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000",
// 			expectedDecodedBEEF: &DecodedBEEF{
// 				CMPSlice: CMPSlice{
// 					{
// 						{
// 							"8c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee8": 0x0,
// 							"da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb": 0x1,
// 							"b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef": 0x2,
// 							"e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b": 0x3,
// 						},
// 						{
// 							"3470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e": 0x0,
// 							"c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c": 0x1,
// 						},
// 						{"cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b": 0x1},
// 					},
// 					{
// 						{
// 							"8c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee8": 0x0,
// 							"da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb": 0x1,
// 							"b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef": 0x2,
// 							"e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b": 0x3,
// 						},
// 						{
// 							"3470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e": 0x0,
// 							"c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c": 0x1,
// 						},
// 						{"cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b": 0x1},
// 					},
// 				},
// 				InputsTxData: []TxData{
// 					{
// 						Transaction: &bt.Tx{
// 							Version:  2,
// 							LockTime: 0,
// 							Inputs: []*bt.Input{
// 								{
// 									PreviousTxSatoshis: 0,
// 									PreviousTxOutIndex: 0,
// 									SequenceNumber:     4294967295,
// 									PreviousTxScript:   nil,
// 								},
// 							},
// 							Outputs: []*bt.Output{
// 								{
// 									Satoshis: 1,
// 								},
// 							}},
// 						PathIndex: nil,
// 					},
// 					{
// 						Transaction: &bt.Tx{
// 							Version:  2,
// 							LockTime: 0,
// 							Inputs: []*bt.Input{
// 								{
// 									PreviousTxSatoshis: 0,
// 									PreviousTxOutIndex: 0,
// 									SequenceNumber:     4294967295,
// 									PreviousTxScript:   nil,
// 								},
// 							},
// 							Outputs: []*bt.Output{
// 								{
// 									Satoshis: 1,
// 								},
// 							}},
// 						PathIndex: nil,
// 					},
// 				},
// 				ProcessedTxData: TxData{
// 					Transaction: &bt.Tx{
// 						Version:  2,
// 						LockTime: 0,
// 						Inputs: []*bt.Input{
// 							{
// 								PreviousTxSatoshis: 0,
// 								PreviousTxOutIndex: 0,
// 								SequenceNumber:     4294967295,
// 								PreviousTxScript:   nil,
// 							},
// 						},
// 						Outputs: []*bt.Output{
// 							{
// 								Satoshis: 1,
// 							},
// 						}},
// 				},
// 			},
// 			pathIndexForTheOldestInput: nil,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// given
// 			beef := tc.hexStream

// 			// when
// 			decodedBEEF, err := DecodeBEEF(beef)

// 			// then
// 			assert.Equal(t, tc.expectedError, err, "expected error %v, but got %v", tc.expectedError, err)

// 			assert.Equal(t, len(tc.expectedDecodedBEEF.InputsTxData), len(decodedBEEF.InputsTxData), "expected %v inputs, but got %v", len(tc.expectedDecodedBEEF.InputsTxData), len(decodedBEEF.InputsTxData))

// 			assert.Equal(t, len(tc.expectedDecodedBEEF.CMPSlice), len(decodedBEEF.CMPSlice), "expected %v CMPs, but got %v", len(tc.expectedDecodedBEEF.CMPSlice), len(decodedBEEF.CMPSlice))

// 			assert.Equal(t, tc.expectedDecodedBEEF.CMPSlice, decodedBEEF.CMPSlice, "expected decoded CMP to be %v, but got %v", tc.expectedDecodedBEEF.CMPSlice, decodedBEEF.CMPSlice)

// 			assert.NotNil(t, decodedBEEF.ProcessedTxData.Transaction, "expected original transaction to be not nil")

// 			assert.Equal(t, tc.expectedDecodedBEEF.InputsTxData[0].PathIndex, decodedBEEF.InputsTxData[0].PathIndex, "expected path index for the oldest input to be %v, but got %v", tc.expectedDecodedBEEF.InputsTxData[0].PathIndex, decodedBEEF.InputsTxData[0].PathIndex)
// 		})
// 	}
// }

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
			expectedError:                errors.New("cannot decode cmp slice from stream - no bytes provided"),
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
			name:                         "unable to decode height - exceeded maximum allowed value",
			hexStream:                    "0100beef01660101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b02020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac000000000100020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000",
			expectedDecodedBEEF:          nil,
			expectedError:                errors.New("height exceeds maximum allowed value of 64"),
			expectedCMPForTheOldestInput: false,
		},
		{
			name:                "unable to decode nOfLeaves - proper BEEF marker, number of CMPs and starting height, but end of stream at this point",
			hexStream:           "0100beef0101",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract Compound Merkle Path at height 1"),
		},
		{
			name:                "unable to decode CMP offset - proper BEEF marker, number of CMPs and starting height, and number of leaves, but end of stream at this point",
			hexStream:           "0100beef010101",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract index 0 leaf of 1 leaves at 1 height"),
		},
		{
			name:                "unable to decode CMP leaf - proper BEEF marker, number of CMPs and starting height, and number of leaves, offset, but hash length is less than 32 bytes",
			hexStream:           "0100beef01010201cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("insufficient bytes to extract hash of path with offset 1 at height 1"),
		},
		{
			name:                "not enough transactions provided to decode BEEF properly (0 transactions)",
			hexStream:           "0100beef01020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b00",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("not enough transactions provided to decode BEEF"),
		},
		{
			name:                "invalid HasCMP flag provided for decoded transaction",
			hexStream:           "0100beef01020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b02020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac000000000100020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000007",
			expectedDecodedBEEF: nil,
			expectedError:       errors.New("invalid HasCMP flag for transaction at index 1"),
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

// TestDecodingBeef will test methods on the DecodedBEEF struct
func TestDecodedBeef(t *testing.T) {
	t.Parallel()

	const validBeefHex = "0100beef01020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b020100000001d0924efc6eb21c88ec91538edfb1fa8ae73e1e2417d6fdec0119998d6042778b0a0000006a47304402205d31e8777edd5d609d3ad9b3090c37016eacf9ab3b150d8badc6d9088ed1ba99022032af2a0b7b8d9cd6a92da5972dfd9d84722e86c213497bbe5a09d30acf9893ee412102d395073f0b4866d64d10015cb016924b1f79cad522911e0b884cd362304f6fd5ffffffff09f4010000000000001976a9147568534fbfc766d05a85c0a18adf71b736c9ad6888acf4010000000000001976a914005d343495af9904df7058ca255dfc7a6271b80f88acf4010000000000001976a914bdd0a2081a29b10c66b76534de0b3c4742fbe35688acf4010000000000001976a91479e158f460cedabf2ed37793e2c2b8f39c79909688acf4010000000000001976a9141a23b7405448ddb2fc687b8479fe9ba16d83edd888ac88130000000000001976a914d107abe806862ac2afa80e77ae5bc4c38eb93a7f88ac10270000000000001976a9149138e8bc3fad2076a9335b0a1f7ea29502b13ce588ac0000000000000000fd2401006a22314d6a4a7251744a735959647a753254487872654e69524c53586548637a417778550a5361746f73656e6465720a746578742f706c61696e057574662d38017c223150755161374b36324d694b43747373534c4b79316b683536575755374d7455523503534554036170700a7361746f73656e64657204747970650c7265676973746572557365720c7061796d61696c204861736840376538303531633662306330633339373231303238656134326361333836383733323236656263396264353732323336353935376637333461353365633339350970686f6e65486173684037363339633564383239646138373935373030353133613865626338623762656361643532333336346365356335626135666436636263356239626532333261c3090000000000001976a91479dcbb510e68557c8a791e439cd9f8b0d8d3429b88ac00000000010001000000010202f24b9ae7399cc6b218053b3b0800cd48c93131fc71442921eb46e9b2ea5a060000006a4730440220338e92e521529e433a2c6b9afbe02e30602c9a553570855692b03e8cfab5b65802204196d7bf136f9768d094808d6bfd6cade3030bf812affea1d71bd51b1c2b104b412102eb33b0cbffb1e3490033348e9d47bcffbeb2e917210958c013f3260864b86c4bffffffff08f4010000000000001976a914c0b3640ed2d59b31d90f1eca2b87db733fb303db88ac88130000000000001976a91417386cd7256887615d214d3e0f70fede265b52cc88acf4010000000000001976a91402b6128d583aa1588f617e5980f4727891e71a9b88acf4010000000000001976a91449412664a4231edb2dfc03cbabff3b404ea4776588acf4010000000000001976a914cb9610a9d2bf1805022779d6f97e2cfdd7c2c8c488acf4010000000000001976a91416b6a5e45d5b2fb7d64e255504734f7c8f7762fa88ac0000000000000000fd2401006a22314d6a4a7251744a735959647a753254487872654e69524c53586548637a417778550a5361746f73656e6465720a746578742f706c61696e057574662d38017c223150755161374b36324d694b43747373534c4b79316b683536575755374d7455523503534554036170700a7361746f73656e64657204747970650c7265676973746572557365720c7061796d61696c204861736840336135343561343561306535313837666262383264663538376438656266623061336665336130363665333338373034643539396234623132343335333362650970686f6e65486173684066376661396133303131616462356364346535336135363631646564363337633564666363663836323864643130363666663737393764613130383334303261c3090000000000001976a914caaf40bc699eb34363f25e961d72f7045dbd4d2688ac0000000000"
	validDecodedBeef, err := DecodeBEEF(validBeefHex)
	require.Nil(t, err)

	t.Run("SPV on valid beef", func(t *testing.T) {
		require.Nil(t, ExecuteSimplifiedPaymentVerification(validDecodedBeef, new(mockServiceProvider)))
	})
}
