package paymail

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBEEF_DecodeCMPSliceFromStream(t *testing.T) {
	testCases := []struct {
		name             string
		hexString        string
		expectedCMPSlice CMPSlice
		expectedError    error
	}{
		{
			name:      "Valid CMPSlice",
			hexString: "02020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b020101cd73c0c6bb645581816fa960fd2f1636062fcbf23cb57981074ab8d708a76e3b02003470d882cf556a4b943639eba15dc795dffdbebdc98b9a98e3637fda96e3811e01c58e40f22b9e9fcd05a09689a9b19e6e62dbfd3335c5253d09a7a7cd755d9a3c04008c00bb9360e93fb822c84b2e579fa4ce75c8378ae87f67730a49552f73c56ee801da256f78ae0ad74bbf539662cdb9122aa02ba9a9d883f1d52468d96290515adb02b4c8d919190a090e77b73ffcd52b85babaaeeb62da000473102aca7f070facef03e5b331f4961d764373f3a4e2751954e75489fb17902aad583eedbb41dc165a3b02020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac000000000100020000000158cb8b052fded9a6c450c4212562df8820359ec34da41286421e0d0f2b7eefee000000006a47304402206b1255cb23454c63b22833de25a3a3ecbdb8d8645ad129d3269cdddf10b2ec98022034cadf46e5bfecc38940e5497ddf5fa9aeb37ff5ec3fe8e21b19cbb64a45ec324121029a82bfce319faccc34095c8405896e1223921917501a4f736a04f126d6a01c12ffffffff0101000000000000001976a914d866ec5ebb0f4e3840351ee61887101e5407562988ac0000000000",
			expectedCMPSlice: CMPSlice{
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
			expectedError: nil,
		},
		{
			name:             "Empty hexString",
			hexString:        "",
			expectedCMPSlice: nil,
			expectedError:    errors.New("provided hexStream is empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			bytes, _ := hex.DecodeString(tc.hexString)

			// when
			result, _, err := DecodeCMPSliceFromStream(bytes)

			// then
			assert.Equal(t, tc.expectedError, err, "expected error %v, but got %v", tc.expectedError, err)
			assert.Equal(t, tc.expectedCMPSlice, result, "expected result %v, but got %v", tc.expectedCMPSlice, result)
		})
	}
}
