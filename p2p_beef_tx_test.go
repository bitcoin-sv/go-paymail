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

func TestDecodeBEEF_InvalidBeef_RetunrError(t *testing.T) {
	const rawTx = "01000000027b0a1b12c7c9e48015e78d3a08a4d62e439387df7e0d7a810ebd4af37661daaa000000006a47304402207d972759afba7c0ffa6cfbbf39a31c2aeede1dae28d8841db56c6dd1197d56a20220076a390948c235ba8e72b8e43a7b4d4119f1a81a77032aa6e7b7a51be5e13845412103f78ec31cf94ca8d75fb1333ad9fc884e2d489422034a1efc9d66a3b72eddca0fffffffff7f36874f858fb43ffcf4f9e3047825619bad0e92d4b9ad4ba5111d1101cbddfe010000006a473044022043f048043d56eb6f75024808b78f18808b7ab45609e4c4c319e3a27f8246fc3002204b67766b62f58bf6f30ea608eaba76b8524ed49f67a90f80ac08a9b96a6922cd41210254a583c1c51a06e10fab79ddf922915da5f5c1791ef87739f40cb68638397248ffffffff03e8030000000000001976a914b08f70bc5010fb026de018f19e7792385a146b4a88acf3010000000000001976a9147d48635f889372c3da12d75ce246c59f4ab907ed88acf7000000000000001976a914b8fbd58685b6920d8f9a8f1b274d8696708b51b088ac00000000"
	const emptyBumps = "0100beef000201000000027b0a1b12c7c9e48015e78d3a08a4d62e439387df7e0d7a810ebd4af37661daaa000000006a47304402207d972759afba7c0ffa6cfbbf39a31c2aeede1dae28d8841db56c6dd1197d56a20220076a390948c235ba8e72b8e43a7b4d4119f1a81a77032aa6e7b7a51be5e13845412103f78ec31cf94ca8d75fb1333ad9fc884e2d489422034a1efc9d66a3b72eddca0fffffffff7f36874f858fb43ffcf4f9e3047825619bad0e92d4b9ad4ba5111d1101cbddfe010000006a473044022043f048043d56eb6f75024808b78f18808b7ab45609e4c4c319e3a27f8246fc3002204b67766b62f58bf6f30ea608eaba76b8524ed49f67a90f80ac08a9b96a6922cd41210254a583c1c51a06e10fab79ddf922915da5f5c1791ef87739f40cb68638397248ffffffff03e8030000000000001976a914b08f70bc5010fb026de018f19e7792385a146b4a88acf3010000000000001976a9147d48635f889372c3da12d75ce246c59f4ab907ed88acf7000000000000001976a914b8fbd58685b6920d8f9a8f1b274d8696708b51b088ac000000000001000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006b48304502210095ea412e82881f81de63764f38f2562bd8a184b06686b3a9e4a5a8b47b9ea1cf022018c76a08b46168c8beb3f7e508cd21c208882b1ec801a9ad7c1b791995092429412102b0c8980f5d2cab77c92c68ac46442feba163a9d306913f6a34911fc618c3c4e7ffffffff01f4010000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988ac0000000000"
	const withoutBumps = "0100beef000201000000027b0a1b12c7c9e48015e78d3a08a4d62e439387df7e0d7a810ebd4af37661daaa000000006a47304402207d972759afba7c0ffa6cfbbf39a31c2aeede1dae28d8841db56c6dd1197d56a20220076a390948c235ba8e72b8e43a7b4d4119f1a81a77032aa6e7b7a51be5e13845412103f78ec31cf94ca8d75fb1333ad9fc884e2d489422034a1efc9d66a3b72eddca0fffffffff7f36874f858fb43ffcf4f9e3047825619bad0e92d4b9ad4ba5111d1101cbddfe010000006a473044022043f048043d56eb6f75024808b78f18808b7ab45609e4c4c319e3a27f8246fc3002204b67766b62f58bf6f30ea608eaba76b8524ed49f67a90f80ac08a9b96a6922cd41210254a583c1c51a06e10fab79ddf922915da5f5c1791ef87739f40cb68638397248ffffffff03e8030000000000001976a914b08f70bc5010fb026de018f19e7792385a146b4a88acf3010000000000001976a9147d48635f889372c3da12d75ce246c59f4ab907ed88acf7000000000000001976a914b8fbd58685b6920d8f9a8f1b274d8696708b51b088ac000000000001000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006b48304502210095ea412e82881f81de63764f38f2562bd8a184b06686b3a9e4a5a8b47b9ea1cf022018c76a08b46168c8beb3f7e508cd21c208882b1ec801a9ad7c1b791995092429412102b0c8980f5d2cab77c92c68ac46442feba163a9d306913f6a34911fc618c3c4e7ffffffff01f4010000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988ac0000000000"
	const withoutParents = "0100beef01fe4e6d0c001002fd9c67028ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160cfd9d6700db1332728830a58c83a5970dcd111a575a585b43b0492361ea8082f41668f8bd01fdcf3300e568706954aae516ef6df7b5db7828771a1f3fcf1b6d65389ec8be8c46057a3c01fde6190001a6028d13cc988f55c8765e3ffcdcfc7d5185a8ebd68709c0adbe37b528557b01fdf20c001cc64f09a217e1971cabe751b925f246e3c2a8e145c49be7b831eaea3e064d7501fd7806009ccf122626a20cdb054877ef3f8ae2d0503bb7a8704fdb6295b3001b5e8876a101fd3d0300aeea966733175ff60b55bc77edcb83c0fce3453329f51195e5cbc7a874ee47ad01fd9f0100f67f50b53d73ffd6e84c02ee1903074b9a5b2ac64c508f7f26349b73cca9d7e901ce006ce74c7beed0c61c50dda8b578f0c0dc5a393e1f8758af2fb65edf483afcaa68016600e32475e17bdd141d62524d0005989dd1db6ca92c6af70791b0e4802be4c5c8c1013200b88162f494f26cc3a1a4a7dcf2829a295064e93b3dbb2f72e21a73522869277a011800a938d3f80dd25b6a3a80e450403bf7d62a1068e2e4b13f0656c83f764c55bb77010d006feac6e4fea41c37c508b5bfdc00d582f6e462e6754b338c95b448df37bd342c010700bf5448356be23b2b9afe53d00cee047065bbc16d0bbcc5f80aa8c1b509e45678010200c2e37431a437ee311a737aecd3caae1213db353847f33792fd539e380bdb4d440100005d5aef298770e2702448af2ce014f8bfcded5896df5006a44b5f1b6020007aeb01010091484f513003fcdb25f336b9b56dafcb05fbc739593ab573a2c6516b344ca5320101000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006b48304502210095ea412e82881f81de63764f38f2562bd8a184b06686b3a9e4a5a8b47b9ea1cf022018c76a08b46168c8beb3f7e508cd21c208882b1ec801a9ad7c1b791995092429412102b0c8980f5d2cab77c92c68ac46442feba163a9d306913f6a34911fc618c3c4e7ffffffff01f4010000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988ac0000000000"
	const withBumpTreeHeightEq65 = "0100beef01fe4e6d0c004102fd9c67028ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160cfd9d6700db1332728830a58c83a5970dcd111a575a585b43b0492361ea8082f41668f8bd01fdcf3300e568706954aae516ef6df7b5db7828771a1f3fcf1b6d65389ec8be8c46057a3c01fde6190001a6028d13cc988f55c8765e3ffcdcfc7d5185a8ebd68709c0adbe37b528557b01fdf20c001cc64f09a217e1971cabe751b925f246e3c2a8e145c49be7b831eaea3e064d7501fd7806009ccf122626a20cdb054877ef3f8ae2d0503bb7a8704fdb6295b3001b5e8876a101fd3d0300aeea966733175ff60b55bc77edcb83c0fce3453329f51195e5cbc7a874ee47ad01fd9f0100f67f50b53d73ffd6e84c02ee1903074b9a5b2ac64c508f7f26349b73cca9d7e901ce006ce74c7beed0c61c50dda8b578f0c0dc5a393e1f8758af2fb65edf483afcaa68016600e32475e17bdd141d62524d0005989dd1db6ca92c6af70791b0e4802be4c5c8c1013200b88162f494f26cc3a1a4a7dcf2829a295064e93b3dbb2f72e21a73522869277a011800a938d3f80dd25b6a3a80e450403bf7d62a1068e2e4b13f0656c83f764c55bb77010d006feac6e4fea41c37c508b5bfdc00d582f6e462e6754b338c95b448df37bd342c010700bf5448356be23b2b9afe53d00cee047065bbc16d0bbcc5f80aa8c1b509e45678010200c2e37431a437ee311a737aecd3caae1213db353847f33792fd539e380bdb4d440100005d5aef298770e2702448af2ce014f8bfcded5896df5006a44b5f1b6020007aeb01010091484f513003fcdb25f336b9b56dafcb05fbc739593ab573a2c6516b344ca5320101000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006b48304502210095ea412e82881f81de63764f38f2562bd8a184b06686b3a9e4a5a8b47b9ea1cf022018c76a08b46168c8beb3f7e508cd21c208882b1ec801a9ad7c1b791995092429412102b0c8980f5d2cab77c92c68ac46442feba163a9d306913f6a34911fc618c3c4e7ffffffff01f4010000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988ac0000000000"

	tcs := []struct {
		name          string
		beef          string
		expectedError error
	}{
		{
			name:          "DecodeBEEF - rawTx",
			beef:          rawTx,
			expectedError: errors.New("invalid format of transaction, BEEF marker not found"),
		},
		{
			name:          "DecodeBEEF - empty BUMPs",
			beef:          emptyBumps,
			expectedError: errors.New("invalid BEEF- lack of BUMPs"),
		},
		{
			name:          "DecodeBEEF - without  BUMPs",
			beef:          withoutBumps,
			expectedError: errors.New("invalid BEEF- lack of BUMPs"),
		},
		{
			name:          "DecodeBEEF - without  input parent transactions",
			beef:          withoutParents,
			expectedError: errors.New("invalid BEEF- not enough transactions provided to decode BEEF"),
		},
		{
			name:          "DecodeBEEF - with a bump tree higher than 64",
			beef:          withBumpTreeHeightEq65,
			expectedError: errors.New("invalid BEEF - treeHeight cannot be grater than 64"),
		},
	}

	t.Parallel()
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			result, err := DecodeBEEF(tc.beef)

			// then
			assert.Equal(t, tc.expectedError, err, "expected error %v, but got %v", tc.expectedError, err)
			assert.Nil(t, result, "expected nil result, but got %v", result)
		})
	}
}
