package paymail

import (
	"errors"
	"testing"

	"github.com/libsv/go-bt/v2"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBEEF_DecodeBEEF(t *testing.T) {
	testCases := []struct {
		name                         string
		hexStream                    string
		expectedDecodedBEEF          *DecodedBEEF
		expectedCMPForTheOldestInput bool
		expectedError                error
	}{
		{
			name:      "valid BEEF with 1 CMP and 1 input transaction",
			hexStream: "0100beef01020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b02020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac000000000100020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000",
			expectedDecodedBEEF: &DecodedBEEF{
				CMPSlice: CMPSlice{
					{
						{"cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b": 0x1},
						{"3470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e": 0x0, "c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c": 0x1},
						{"8c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee8": 0x0, "da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb": 0x1, "b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef": 0x2, "e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b": 0x3},
					},
				},
				InputsTxData: []TxData{
					{
						Transaction: &bt.Tx{},
						PathIndex:   func(v bt.VarInt) *bt.VarInt { return &v }(0x0),
					},
				},
				ProcessedTxData: TxData{
					Transaction: &bt.Tx{},
				},
			},
			expectedCMPForTheOldestInput: true,
		},
		{
			name:      "valid BEEF with 2 CMP and 2 input transaction - all input transactions have no CMP flag set",
			hexStream: "0100beef02020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b03020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000",
			expectedDecodedBEEF: &DecodedBEEF{
				CMPSlice: CMPSlice{
					{
						{"cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b": 0x1},
						{"3470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e": 0x0, "c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c": 0x1},
						{"8c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee8": 0x0, "da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb": 0x1, "b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef": 0x2, "e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b": 0x3},
					},
					{
						{"cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b": 0x1},
						{"3470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e": 0x0, "c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c": 0x1},
						{"8c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee8": 0x0, "da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb": 0x1, "b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef": 0x2, "e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b": 0x3},
					},
				},
				InputsTxData: []TxData{
					{
						Transaction: &bt.Tx{},
						PathIndex:   nil,
					},
					{
						Transaction: &bt.Tx{},
						PathIndex:   nil,
					},
				},
				ProcessedTxData: TxData{
					Transaction: &bt.Tx{},
				},
			},
			expectedCMPForTheOldestInput: false,
		},
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
			decodedBEEF, err := DecodeBEEF(beef)

			// then
			assert.Equal(t, tc.expectedError, err, "expected error %v, but got %v", tc.expectedError, err)

			// only if there is no error go through the rest of the assertions
			if tc.expectedError == nil {
				assert.Equal(t, len(tc.expectedDecodedBEEF.InputsTxData), len(decodedBEEF.InputsTxData), "expected %v inputs, but got %v", len(tc.expectedDecodedBEEF.InputsTxData), len(decodedBEEF.InputsTxData))
				assert.Equal(t, len(tc.expectedDecodedBEEF.CMPSlice), len(decodedBEEF.CMPSlice), "expected %v CMPs, but got %v", len(tc.expectedDecodedBEEF.CMPSlice), len(decodedBEEF.CMPSlice))

				assert.NotNil(t, decodedBEEF.ProcessedTxData.Transaction, "expected original transaction to be not nil")

				if tc.expectedCMPForTheOldestInput {
					assert.NotNil(t, tc.expectedDecodedBEEF.InputsTxData[0].PathIndex, "expected %v, but got %v", "expected PathIndex for oldest input to be not nil")
				}
			}
		})
	}
}
