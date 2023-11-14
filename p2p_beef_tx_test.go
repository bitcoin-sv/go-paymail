package paymail

import (
	"context"
	"errors"
	"testing"

	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestDecodingBeef will test methods on the DecodedBEEF struct
func TestDecodedBeef(t *testing.T) {
	t.Parallel()

	const validBeefHex = "0100beef04fe4c7c0c000e02fd460302d2e6252f996ab1a6fcbe8911e8f865bb719c0e11787397fd818b5bb1ff554a3cfd470300d67d0a32df0ecd5976479f068f19ff671e0a285570b915bc6e6658e98a9c0e3401fda20100eafe36c6f1adef584b1d199286c5024a4791eb16f35161c4b69cc119c1c9493c01d00077605d020124b1f6e039f1b2ccebf3e46acf45584bc559b654b57179b99d906a0169002c4fe7d1d2b2990df28e5132b059e10cc6ce9be85a79f834546eee881a8f968d013500a187befab2d1cc22004f96045afecd929e2edc7017a044aaf0f530847afa1644011b00bf5108b8176d4abba8c80f935607abc6a6bc38b9461c86a812369975ad2abb1e010c002efe8c4f9a630b0d0f4a410d93523d7003085dd166248d6e76f405fe8fc606720107007043c5b1be3e6952859bca6d6320f5a0f6c5bba5d5fbfd8a195f1d9ec4e0964c010200822378b29fd273a3bf38a4b852089e8f5a1251b5fedeee6ee2f548c0ee93a14e010000788ff2bf42a41c2e32bd59504fa59fd8b2ab0058430f9d4d2988ee8341ccac1b010100470c98a5fb48ed000ca746a05c7d398d2b76efe70a65c67c1cec14a8c8c84cce010100de821cfed42c83ab4dfc3fef452ec5f4086d18fb7324c046826cd51e0fea2c8b0101006422ad7eb4999dc4ca89a3488862f03b4101c32ed49c19e3df7563958e6e480201010087ce3715a94573ac7746f3b66e4070d313b4a8df6114509678ab08c425250b06fe137c0c000e04fde00a02827f1758c64b4a0b1226c54316ddeed4618500f026602ebca3cd4b96174a690afde10a001a906867a2cbd98095453400107f814e5380640443bf4504eb6fe412acba674bfdac0d001f3da43722e72c749846913dbce0623b1d22fa8499bbca49f3f9aec4cd5b3d70fdad0d02213e4fca3103f812ffcba253caf452c6811947ff6f2fb99b4e18baa1233e84b602fd7105001399fab1150a4f8013dca738bec4ae8c5b215e301697d6453f648230b36cd431fdd70600282c33fe4ec70f91818781daa6a97e85022bc2031a370202a385f2aeb9086a1a02fdb90200fd35a389307434d05891122dc81f3a5246f577334de368ebcfc15836648e73a8fd6a0300ed1732da6c0666a7270d63ca0d4af8ba1c9a80c149bccd5e1ac039159418740a02fd5d010019dc6865859cc67245c5d9b04bc3a9141a4e6fda17ba33d24271600244b2dd35fdb40100ccb07e295f3bbac75957015c58c5f7f64be4904865c95ebf4181c40d0089429602af002d31340de572c9cdbe9c75c62474ad3130e417d7a61b370aab7345901f69badfdb006c6b6e5958eb5f2ae497e4b7d8ecc90ac614856ada2458f35e1f2880fd155530025600efb95e533fbed11730d7eb9f2c891a885cb658f4ca058fdf6bed8d348c9415de6c0056e17aded7bd785e6a44cf023a294440bd89886dfa5d48ce28386199361e4b47022a003c5a68cac12c161d6a3c5a1e8ba77d9df937954a3c5ebb669fc5d3f04883c22237003978958e578f7326e8271b058ad60c85f1d2d404bdf75611ef9192de27ce644d021400692e79f2246540d6526e8261bda1b430f218be172a2672ae0a58e5fe20c359331a00595cbb11c5591f7630c31a779fe71b5a48526b525078a37492b087dcff4f1314020b00d6a30d5bccc4f3d9387e56e9e1f32de9db91c52c8fb4179aab7732a38df8df600c009fe2e3698677771d5fd4819796ee8342277ae3bdcd68e4e38ad3c04d207df06d0204007635530acf3aff406ca2ae825aa7253c3c1fb92cadf0124145a4565477c9198f0700f1f84a4c5efaeb4d5660968aa9234db2a98c5bca3f5942ec1d007dc2d996822f0202003578890ece8c5028080f613c01b8662df60c4c5ea6cc858e35ec266f3e8c8bf80300d2fcdb70b5967e0b53341b6bcc517511aa977ec6c98e587035e734bafb2df40201000050e8890913d743a6a49e244dab275bcafb7ec2fd8828163e6a07560351b77b2e01010039787d0d5252e1347883ebcad07c141448e366d2b40d3241211482ef52f34d97010100278ecd209c80ea49e6581c8b9e03b118e66cbc39e4acb14b26bda7d378fe2117feeb7b0c000f02fdd438020dacf934645c462a155ca35453ab578e7d510687fda689a565e000fc4df11cd5fdd53800de97bf2bdde019a90f6c784c3731b71a45ade39a58c436e2c70a4c9f83371e1d01fd6b1c00aaccb26fe28312165b6a2fd6414c651dd178a8053324d25008b8f533af7ad30a01fd340e0010cd00845ba8b3e5e08a057de308885b01ec1cb2f4cc56da6583ff1f3f9d98ff01fd1b0700d712139ba223278e95809bd1aebd6d334288ef66b989116a344485c603e9998d01fd8c0300097e0e68b5f5aab8ef48a8983b4b4fa992d80bd4cc2fd2c3069be71407d73b0001fdc701008be1cc85f6bb068b86ad60cfbd321878c2774c99274e953530cf4d7611227f3b01e2000453717a6dee5b33d9c379e7b262d17ddcd6267b159fd8a0b980603d03ce31b601700050992dd58d2332d87df4874852f9275eeb9c7813e1a6c1e1cad1a74a9e5560a10139003fdc83bc63d0af296cd6a71e45355eab21634d829f6b17e769fe0d2308a99484011d001eb98912db04c69f675a4cfb64a57167d9b136d174c0d0e1f2635840a27a2819010f00b0a43e31d356c1ccd7cd1733f8bf5885013361d7f02627629136b360c55aa43a010600a7b2280eb6bc6f398f7f5423033bd011e819caf0beeb248f26a1498fbc7e4bb1010200152e881046cdb491d33c93d799f28e07686708db2cdc8a3011c375f7b9bef87a010000851add544d85300e318ec60be010f38f2bea85cc86ce5a972df28b94e53f5801010100a7308e947ed378dcf6e6f79ba0bb8d062f7c7ffcc025e5d08fc30a3288e87ef5fe1c7c0c000e02fd5a1902e5232220aee5069017d31cc30818bcc971de3e6418f6e62b8cc9a3430d64f3e9fd5b19004359d0c4ca1b47399960cb57ff4b0d2c850c825ba56a2119c3ff875fece11a4001fdac0c003c0fd4b7a8621c97bf7b6adad9569fe1f07a52c41d472d1e3407f9e37b3400d701fd5706006400dbcd6c1ea8a5c8cccfa153eb02076bf03b4ac7fef38db8e3a5032d908f3101fd2a0300ceac7150203519610ec78f93d64bbdc2e76e9cfadb72a83070d5fc72c8c0291a01fd94010058dde4e9d746faed13b29bd0f44ea86cdab712bc3d430a88e8cec22364e19cfd01cb001f77543e1a81dc8d4fdef4a4b132cf907b8a22de393400035a40391904fcb5c20164002cc86783ce17e19b7a56e0d08ba4c5d94cfd586771dd28e2fe1284d887db1285013300ceb2afce84f3e257bce737ebaef9ded331348d750d17f408c22d79dc8e9de9e9011800d939fb970c2044768981bea25d19be80302773a4b0d469566315258d39db5059010d00c815d2c629a1d419d6ef3c33a670da1bc64f7f23a840365fc292aab7d7b4f770010700828ced5731867ba622104fe64bd18a2132097623f8fe4a834537b329dbb6d06a010200f9df4e876faa4ddd9f909bfafed085ad24297f80003a3700bd13250498e252190100004834f3f8c1be8c4f43bb61dd6d0eea6ab996bae149aefd5f9aced05c54b3395901010036679c2ca45cd1d37b4b0d1d4a9a0e923d80f6241957fd61969b641ca1a2b62c060100000001aa5df6cc8e5f41e8770d222d06d2dfdea647769871044df714c99b59d3528858000000006b483045022100d65cb97677ba806e752d63ea988faf5fb00cabd0df78bb78fe9f42a079e7ae540220158a2d51ca8b03f18691eab1c4fa5e85b6f21846df9bfaf797c995b1c9f34e47412103597c97b42c5f880e27eaf2eb2de0fb0fe5a81d689c646b2c6527acbc9c2509c7ffffffff0264000000000000001976a914e5c8e12ec010e7bf318ebc94efcb887ff751069c88ac8e010000000000001976a914ae791b7e9d0108b4000ac4138c2b47d8cfea404288ac0000000001010100000001a25525c35101ad7ad7653e16d5d525762f3826e4c7b87077a096aebabe292ad9000000006b4830450221009fa60fe23a004850fd671a834534b278d4a2a9a4e7d8dfc0e4f9ffcb2899c1ef022010a8e718b608dfd308cdbc382e6f666583c8148b5f70819fbb7df1f1edd8cc164121023d72bbb613c7e8a4230c080ffaff1abc038df9008046b2157f888f36f427d2f0ffffffff0228000000000000001976a914596ce4dbbefd792d782ccb7ce5652726fd0c676288acca010000000000001976a9143e8a2536d8bc8f2a995800f16e068ee85aa1939288ac00000000010301000000017a7123cf30268df225315c47b6ccbfd3e4f5c19f765bb9abc8fbeed134eb0b7f010000006a47304402201715f13c3456310f7b73386ca947621b53f71f27d2a23c0d1f09deb4b8931ca7022000d9c8daa400b1c2ad34df4c22e024e0d9191ecefe83160f47fbf53026fb4d53412102b6f13a7f2599f7312cf55ca1524bf9ebefeed22daa1f80ffed82df038db0167affffffff012d000000000000001976a9148e4221cdcf3f547ffd0aee46411b40a3248fdb3688ac0000000001020100000001a25525c35101ad7ad7653e16d5d525762f3826e4c7b87077a096aebabe292ad9010000006a47304402202ed74a7c67dd15fa9594ae600775e3a289511ceda7994b7914bd94435a269c7202205767d4deaf8db2c7c3a862a2b0e8f43f377a597bbbbed10409819bc2239064804121024ffdb4c74d0f3e92dbfe41348ee5295a8a8aa91f7eea1c17c879baf3d31492c3ffffffff01e7030000000000001976a914424cc7a221a4bf81790a850014a5e1a1663d3b5888ac0000000001010100000001213e4fca3103f812ffcba253caf452c6811947ff6f2fb99b4e18baa1233e84b6000000006a47304402202b577840e0fbab66a150a79feeb9ea2e0469df83f44a4aebb4779c35ea1d7b5402204764d20a211b112e8907a3a80e418dcc933ba4e49dea7072bd5814cca8a76e664121027db94b451dafb7368da30ed73abbf064a622ebcffc1bd3362802fd8eaaf5fc6dffffffff0163000000000000001976a9144451354c29626cbeb4f130e6bb7ddd54151c6b3b88ac0000000001000100000005d2e6252f996ab1a6fcbe8911e8f865bb719c0e11787397fd818b5bb1ff554a3c000000006a4730440220179e09c88dc5d83e6db5abcf2a46c989ba766612f3be664f66f9d5932681b42e02207ea73be5317362c511e76a40f3b416ca2ec99ce891bd55fbbfc60f3e2598c1be41210277e25e8e4ab96e46d94a037411d195b537810b1d7301c2d72e9eb40bb47aae34ffffffff827f1758c64b4a0b1226c54316ddeed4618500f026602ebca3cd4b96174a690a000000006b4830450221008a9418e0be1b65b45ac27449071bf8a7cb55d9008ae5c8d28da47bce9bc2cc28022056aa8ffd170eb4b86b2be7da81b1d01dd0f026d0488645c3b27f6e521eebb51a412103138a3aac623a5fc9789cd9476efad0605dea61f3e7cd5089195eb0b446b6ab4effffffff0dacf934645c462a155ca35453ab578e7d510687fda689a565e000fc4df11cd5000000006a4730440220147729079aef1449d0c1b08f79c1fa4749d685190d0fdd25801725845d0186740220796be61c7a059e65b4e418d1fec99a596ea565b179683d334dc7a2e2312ed466412103ba49c495254796f5ccf0e28f205f62965fafc33367b2b8d6609e5de30c206ad4ffffffff213e4fca3103f812ffcba253caf452c6811947ff6f2fb99b4e18baa1233e84b6010000006b483045022100e5374c185ffb992bcabec480572427f402911cdaeda92246f2a94813a32f33ab0220676121e6bd049c981593a8093f779927cc95bdf56d50069b7786ad9b9fe528c84121035d1d732dbe247c0886753c84dc3d2fc96a9eac26662e8664fe9ce8f67ab6dd98ffffffffe5232220aee5069017d31cc30818bcc971de3e6418f6e62b8cc9a3430d64f3e9010000006a473044022048e3532181d848dcf69a4da9b8ae296b90b4660a27f880b928ce15e9d3c5979e02205cef48653b1defbc2a87d93b268693c8d078199b69eb5bba32ba10f301727bf641210343caa07997898400cefe7a28445b233d30463d13359c1d87ac42ea5da61432a0ffffffff01ce070000000000001976a9145c21ae83ea5892dea33b9b002eb1d8450b581ea888ac0000000000"
	validDecodedBeef, err := DecodeBEEF(validBeefHex)
	require.Nil(t, err)

	t.Run("SPV on valid beef", func(t *testing.T) {
		require.Nil(t, ExecuteSimplifiedPaymentVerification(validDecodedBeef, new(mockServiceProvider)))
	})
}
