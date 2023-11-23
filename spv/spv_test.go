package spv

import (
	"context"
	"errors"
	"testing"

	"github.com/bitcoin-sv/go-paymail/beef"
	"github.com/stretchr/testify/require"
)

func TestExecuteSimplifiedPaymentVerification_ValidTransaction_Success(t *testing.T) {
	t.Parallel()

	const validBeefHex = "0100beef04fe4c7c0c000e02fd460302d2e6252f996ab1a6fcbe8911e8f865bb719c0e11787397fd818b5bb1ff554a3cfd470300d67d0a32df0ecd5976479f068f19ff671e0a285570b915bc6e6658e98a9c0e3401fda20100eafe36c6f1adef584b1d199286c5024a4791eb16f35161c4b69cc119c1c9493c01d00077605d020124b1f6e039f1b2ccebf3e46acf45584bc559b654b57179b99d906a0169002c4fe7d1d2b2990df28e5132b059e10cc6ce9be85a79f834546eee881a8f968d013500a187befab2d1cc22004f96045afecd929e2edc7017a044aaf0f530847afa1644011b00bf5108b8176d4abba8c80f935607abc6a6bc38b9461c86a812369975ad2abb1e010c002efe8c4f9a630b0d0f4a410d93523d7003085dd166248d6e76f405fe8fc606720107007043c5b1be3e6952859bca6d6320f5a0f6c5bba5d5fbfd8a195f1d9ec4e0964c010200822378b29fd273a3bf38a4b852089e8f5a1251b5fedeee6ee2f548c0ee93a14e010000788ff2bf42a41c2e32bd59504fa59fd8b2ab0058430f9d4d2988ee8341ccac1b010100470c98a5fb48ed000ca746a05c7d398d2b76efe70a65c67c1cec14a8c8c84cce010100de821cfed42c83ab4dfc3fef452ec5f4086d18fb7324c046826cd51e0fea2c8b0101006422ad7eb4999dc4ca89a3488862f03b4101c32ed49c19e3df7563958e6e480201010087ce3715a94573ac7746f3b66e4070d313b4a8df6114509678ab08c425250b06fe137c0c000e04fde00a02827f1758c64b4a0b1226c54316ddeed4618500f026602ebca3cd4b96174a690afde10a001a906867a2cbd98095453400107f814e5380640443bf4504eb6fe412acba674bfdac0d001f3da43722e72c749846913dbce0623b1d22fa8499bbca49f3f9aec4cd5b3d70fdad0d02213e4fca3103f812ffcba253caf452c6811947ff6f2fb99b4e18baa1233e84b602fd7105001399fab1150a4f8013dca738bec4ae8c5b215e301697d6453f648230b36cd431fdd70600282c33fe4ec70f91818781daa6a97e85022bc2031a370202a385f2aeb9086a1a02fdb90200fd35a389307434d05891122dc81f3a5246f577334de368ebcfc15836648e73a8fd6a0300ed1732da6c0666a7270d63ca0d4af8ba1c9a80c149bccd5e1ac039159418740a02fd5d010019dc6865859cc67245c5d9b04bc3a9141a4e6fda17ba33d24271600244b2dd35fdb40100ccb07e295f3bbac75957015c58c5f7f64be4904865c95ebf4181c40d0089429602af002d31340de572c9cdbe9c75c62474ad3130e417d7a61b370aab7345901f69badfdb006c6b6e5958eb5f2ae497e4b7d8ecc90ac614856ada2458f35e1f2880fd155530025600efb95e533fbed11730d7eb9f2c891a885cb658f4ca058fdf6bed8d348c9415de6c0056e17aded7bd785e6a44cf023a294440bd89886dfa5d48ce28386199361e4b47022a003c5a68cac12c161d6a3c5a1e8ba77d9df937954a3c5ebb669fc5d3f04883c22237003978958e578f7326e8271b058ad60c85f1d2d404bdf75611ef9192de27ce644d021400692e79f2246540d6526e8261bda1b430f218be172a2672ae0a58e5fe20c359331a00595cbb11c5591f7630c31a779fe71b5a48526b525078a37492b087dcff4f1314020b00d6a30d5bccc4f3d9387e56e9e1f32de9db91c52c8fb4179aab7732a38df8df600c009fe2e3698677771d5fd4819796ee8342277ae3bdcd68e4e38ad3c04d207df06d0204007635530acf3aff406ca2ae825aa7253c3c1fb92cadf0124145a4565477c9198f0700f1f84a4c5efaeb4d5660968aa9234db2a98c5bca3f5942ec1d007dc2d996822f0202003578890ece8c5028080f613c01b8662df60c4c5ea6cc858e35ec266f3e8c8bf80300d2fcdb70b5967e0b53341b6bcc517511aa977ec6c98e587035e734bafb2df40201000050e8890913d743a6a49e244dab275bcafb7ec2fd8828163e6a07560351b77b2e01010039787d0d5252e1347883ebcad07c141448e366d2b40d3241211482ef52f34d97010100278ecd209c80ea49e6581c8b9e03b118e66cbc39e4acb14b26bda7d378fe2117feeb7b0c000f02fdd438020dacf934645c462a155ca35453ab578e7d510687fda689a565e000fc4df11cd5fdd53800de97bf2bdde019a90f6c784c3731b71a45ade39a58c436e2c70a4c9f83371e1d01fd6b1c00aaccb26fe28312165b6a2fd6414c651dd178a8053324d25008b8f533af7ad30a01fd340e0010cd00845ba8b3e5e08a057de308885b01ec1cb2f4cc56da6583ff1f3f9d98ff01fd1b0700d712139ba223278e95809bd1aebd6d334288ef66b989116a344485c603e9998d01fd8c0300097e0e68b5f5aab8ef48a8983b4b4fa992d80bd4cc2fd2c3069be71407d73b0001fdc701008be1cc85f6bb068b86ad60cfbd321878c2774c99274e953530cf4d7611227f3b01e2000453717a6dee5b33d9c379e7b262d17ddcd6267b159fd8a0b980603d03ce31b601700050992dd58d2332d87df4874852f9275eeb9c7813e1a6c1e1cad1a74a9e5560a10139003fdc83bc63d0af296cd6a71e45355eab21634d829f6b17e769fe0d2308a99484011d001eb98912db04c69f675a4cfb64a57167d9b136d174c0d0e1f2635840a27a2819010f00b0a43e31d356c1ccd7cd1733f8bf5885013361d7f02627629136b360c55aa43a010600a7b2280eb6bc6f398f7f5423033bd011e819caf0beeb248f26a1498fbc7e4bb1010200152e881046cdb491d33c93d799f28e07686708db2cdc8a3011c375f7b9bef87a010000851add544d85300e318ec60be010f38f2bea85cc86ce5a972df28b94e53f5801010100a7308e947ed378dcf6e6f79ba0bb8d062f7c7ffcc025e5d08fc30a3288e87ef5fe1c7c0c000e02fd5a1902e5232220aee5069017d31cc30818bcc971de3e6418f6e62b8cc9a3430d64f3e9fd5b19004359d0c4ca1b47399960cb57ff4b0d2c850c825ba56a2119c3ff875fece11a4001fdac0c003c0fd4b7a8621c97bf7b6adad9569fe1f07a52c41d472d1e3407f9e37b3400d701fd5706006400dbcd6c1ea8a5c8cccfa153eb02076bf03b4ac7fef38db8e3a5032d908f3101fd2a0300ceac7150203519610ec78f93d64bbdc2e76e9cfadb72a83070d5fc72c8c0291a01fd94010058dde4e9d746faed13b29bd0f44ea86cdab712bc3d430a88e8cec22364e19cfd01cb001f77543e1a81dc8d4fdef4a4b132cf907b8a22de393400035a40391904fcb5c20164002cc86783ce17e19b7a56e0d08ba4c5d94cfd586771dd28e2fe1284d887db1285013300ceb2afce84f3e257bce737ebaef9ded331348d750d17f408c22d79dc8e9de9e9011800d939fb970c2044768981bea25d19be80302773a4b0d469566315258d39db5059010d00c815d2c629a1d419d6ef3c33a670da1bc64f7f23a840365fc292aab7d7b4f770010700828ced5731867ba622104fe64bd18a2132097623f8fe4a834537b329dbb6d06a010200f9df4e876faa4ddd9f909bfafed085ad24297f80003a3700bd13250498e252190100004834f3f8c1be8c4f43bb61dd6d0eea6ab996bae149aefd5f9aced05c54b3395901010036679c2ca45cd1d37b4b0d1d4a9a0e923d80f6241957fd61969b641ca1a2b62c060100000001aa5df6cc8e5f41e8770d222d06d2dfdea647769871044df714c99b59d3528858000000006b483045022100d65cb97677ba806e752d63ea988faf5fb00cabd0df78bb78fe9f42a079e7ae540220158a2d51ca8b03f18691eab1c4fa5e85b6f21846df9bfaf797c995b1c9f34e47412103597c97b42c5f880e27eaf2eb2de0fb0fe5a81d689c646b2c6527acbc9c2509c7ffffffff0264000000000000001976a914e5c8e12ec010e7bf318ebc94efcb887ff751069c88ac8e010000000000001976a914ae791b7e9d0108b4000ac4138c2b47d8cfea404288ac0000000001010100000001a25525c35101ad7ad7653e16d5d525762f3826e4c7b87077a096aebabe292ad9000000006b4830450221009fa60fe23a004850fd671a834534b278d4a2a9a4e7d8dfc0e4f9ffcb2899c1ef022010a8e718b608dfd308cdbc382e6f666583c8148b5f70819fbb7df1f1edd8cc164121023d72bbb613c7e8a4230c080ffaff1abc038df9008046b2157f888f36f427d2f0ffffffff0228000000000000001976a914596ce4dbbefd792d782ccb7ce5652726fd0c676288acca010000000000001976a9143e8a2536d8bc8f2a995800f16e068ee85aa1939288ac00000000010301000000017a7123cf30268df225315c47b6ccbfd3e4f5c19f765bb9abc8fbeed134eb0b7f010000006a47304402201715f13c3456310f7b73386ca947621b53f71f27d2a23c0d1f09deb4b8931ca7022000d9c8daa400b1c2ad34df4c22e024e0d9191ecefe83160f47fbf53026fb4d53412102b6f13a7f2599f7312cf55ca1524bf9ebefeed22daa1f80ffed82df038db0167affffffff012d000000000000001976a9148e4221cdcf3f547ffd0aee46411b40a3248fdb3688ac0000000001020100000001a25525c35101ad7ad7653e16d5d525762f3826e4c7b87077a096aebabe292ad9010000006a47304402202ed74a7c67dd15fa9594ae600775e3a289511ceda7994b7914bd94435a269c7202205767d4deaf8db2c7c3a862a2b0e8f43f377a597bbbbed10409819bc2239064804121024ffdb4c74d0f3e92dbfe41348ee5295a8a8aa91f7eea1c17c879baf3d31492c3ffffffff01e7030000000000001976a914424cc7a221a4bf81790a850014a5e1a1663d3b5888ac0000000001010100000001213e4fca3103f812ffcba253caf452c6811947ff6f2fb99b4e18baa1233e84b6000000006a47304402202b577840e0fbab66a150a79feeb9ea2e0469df83f44a4aebb4779c35ea1d7b5402204764d20a211b112e8907a3a80e418dcc933ba4e49dea7072bd5814cca8a76e664121027db94b451dafb7368da30ed73abbf064a622ebcffc1bd3362802fd8eaaf5fc6dffffffff0163000000000000001976a9144451354c29626cbeb4f130e6bb7ddd54151c6b3b88ac0000000001000100000005d2e6252f996ab1a6fcbe8911e8f865bb719c0e11787397fd818b5bb1ff554a3c000000006a4730440220179e09c88dc5d83e6db5abcf2a46c989ba766612f3be664f66f9d5932681b42e02207ea73be5317362c511e76a40f3b416ca2ec99ce891bd55fbbfc60f3e2598c1be41210277e25e8e4ab96e46d94a037411d195b537810b1d7301c2d72e9eb40bb47aae34ffffffff827f1758c64b4a0b1226c54316ddeed4618500f026602ebca3cd4b96174a690a000000006b4830450221008a9418e0be1b65b45ac27449071bf8a7cb55d9008ae5c8d28da47bce9bc2cc28022056aa8ffd170eb4b86b2be7da81b1d01dd0f026d0488645c3b27f6e521eebb51a412103138a3aac623a5fc9789cd9476efad0605dea61f3e7cd5089195eb0b446b6ab4effffffff0dacf934645c462a155ca35453ab578e7d510687fda689a565e000fc4df11cd5000000006a4730440220147729079aef1449d0c1b08f79c1fa4749d685190d0fdd25801725845d0186740220796be61c7a059e65b4e418d1fec99a596ea565b179683d334dc7a2e2312ed466412103ba49c495254796f5ccf0e28f205f62965fafc33367b2b8d6609e5de30c206ad4ffffffff213e4fca3103f812ffcba253caf452c6811947ff6f2fb99b4e18baa1233e84b6010000006b483045022100e5374c185ffb992bcabec480572427f402911cdaeda92246f2a94813a32f33ab0220676121e6bd049c981593a8093f779927cc95bdf56d50069b7786ad9b9fe528c84121035d1d732dbe247c0886753c84dc3d2fc96a9eac26662e8664fe9ce8f67ab6dd98ffffffffe5232220aee5069017d31cc30818bcc971de3e6418f6e62b8cc9a3430d64f3e9010000006a473044022048e3532181d848dcf69a4da9b8ae296b90b4660a27f880b928ce15e9d3c5979e02205cef48653b1defbc2a87d93b268693c8d078199b69eb5bba32ba10f301727bf641210343caa07997898400cefe7a28445b233d30463d13359c1d87ac42ea5da61432a0ffffffff01ce070000000000001976a9145c21ae83ea5892dea33b9b002eb1d8450b581ea888ac0000000000"
	validDecodedBeef, err := beef.DecodeBEEF(validBeefHex)
	require.Nil(t, err)

	t.Run("SPV on valid beef", func(t *testing.T) {
		require.Nil(t, ExecuteSimplifiedPaymentVerification(context.Background(), validDecodedBeef, new(mockServiceProvider)))
	})
}

func TestExecuteSimplifiedPaymentVerification_CorruptedTransaction_ReturnError(t *testing.T) {
	t.Parallel()

	const someoneElse = "0100beef01fe4e6d0c001002fd909002088a382ec07a8cf47c6158b68e5822852362102d8571482d1257e0b7527e1882fd91900065cb01218f2506bb51155d243e4d6b32d69d1b5f2221c52e26963cfd8cf7283201fd4948008d7a44ae384797b0ae84db0c857e8c1083425d64d09ef8bc5e2e9d270677260501fd25240060f38aa33631c8d70adbac1213e7a5b418c90414e919e3a12ced63dd152fd85a01fd1312005ff132ee64a7a0c79150a29f66ef861e552d3a05b47d6303f5d8a2b2a09bc61501fd080900cc0baf21cf06b9439dfe05dce9bdb14ddc2ca2d560b1138296ef5769851a84b301fd85040063ccb26232a6e1d3becdb47a0f19a67a562b754e8894155b3ae7bba10335ce5101fd430200e153fc455a0f2c8372885c11af70af904dcf44740b9ebf3b3e5b2234cce550bc01fd20010077d5ea69d1dcc379dde65d6adcebde1838190118a8fae928c037275e78bd87910191000263e4f31684a25169857f2788aeef603504931f92585f02c4c9e023b2aa43d1014900de72292e0b3e5eeacfa2b657bf4d46c885559b081ee78632a99b318c1148d85c01250068a5f831ca99b9e7f3720920d6ea977fd2ab52b83d1a6567dafa4c8cafd941ed0113006a0b91d83f9056b702d6a8056af6365c7da626fc3818b815dd4b0de22d05450f0108009876ce56b68545a75859e93d200bdde7880d46f39384818b259ed847a9664ddf010500990bc5e95cacbc927b5786ec39a183f983fe160d52829cf47521c7eb369771c30103004fe794e50305f590b6010a51d050bf47dfeaabfdb949c5ee0673f577a59537d70100004dad44a358aea4d8bc1917912539901f5ae44e07a4748e1a9d3018814b0759d0020100000002704273c86298166ac351c3aa9ac90a8029e4213b5f1b03c3bbf4bc5fb09cdd43010000006a4730440220398d6389e8a156a3c6c1ca355e446d844fd480193a93af832afd1c87d0f04784022050091076b8f7405b37ce6e795d1b92526396ac2b14f08e91649b908e711e2b044121030ef6975d46dbab4b632ef62fdbe97de56d183be1acc0be641d2c400ae01cf136ffffffff2f41ed6a2488ac3ba4a3c330a15fa8193af87f0192aa59935e6c6401d92dc3a00a0000006a47304402200ad9cf0dc9c90a4c58b08910740b4a8b3e1a7e37db1bc5f656361b93f412883d0220380b6b3d587103fc8bf3fe7bed19ab375766984c67ebb7d43c993bcd199f32a441210205ef4171f58213b5a2ddf16ac6038c10a2a8c3edc1e6275cb943af4bb3a58182ffffffff03e8030000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988acf4010000000000001976a914c7662da5e0a6a179141a7872045538126f1e954288acf5000000000000001976a914765bdf10934f5aac894cf8a3795c9eeb494c013488ac0000000001000100000001088a382ec07a8cf47c6158b68e5822852362102d8571482d1257e0b7527e1882000000006a4730440220610bba9ed83a47641c34bbbcf8eeb536d2ae6cfddc7644a8c520bb747f798c3702206a23c9f45273772dd7e80ba21a5c4613d6ffe7ba1c75b729eae0cdd484fee2bd412103c0cd91af135d09f98d57e34af28e307daf36bccd4764708e8a3f7ea5cebf01a9ffffffff01c8000000000000001976a9148ce2d21f9a75e98600be76b25b91c4fef6b40bcd88ac0000000000"
	const tooMuch = "0100beef01fe4e6d0c001002fd9c67028ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160cfd9d6700db1332728830a58c83a5970dcd111a575a585b43b0492361ea8082f41668f8bd01fdcf3300e568706954aae516ef6df7b5db7828771a1f3fcf1b6d65389ec8be8c46057a3c01fde6190001a6028d13cc988f55c8765e3ffcdcfc7d5185a8ebd68709c0adbe37b528557b01fdf20c001cc64f09a217e1971cabe751b925f246e3c2a8e145c49be7b831eaea3e064d7501fd7806009ccf122626a20cdb054877ef3f8ae2d0503bb7a8704fdb6295b3001b5e8876a101fd3d0300aeea966733175ff60b55bc77edcb83c0fce3453329f51195e5cbc7a874ee47ad01fd9f0100f67f50b53d73ffd6e84c02ee1903074b9a5b2ac64c508f7f26349b73cca9d7e901ce006ce74c7beed0c61c50dda8b578f0c0dc5a393e1f8758af2fb65edf483afcaa68016600e32475e17bdd141d62524d0005989dd1db6ca92c6af70791b0e4802be4c5c8c1013200b88162f494f26cc3a1a4a7dcf2829a295064e93b3dbb2f72e21a73522869277a011800a938d3f80dd25b6a3a80e450403bf7d62a1068e2e4b13f0656c83f764c55bb77010d006feac6e4fea41c37c508b5bfdc00d582f6e462e6754b338c95b448df37bd342c010700bf5448356be23b2b9afe53d00cee047065bbc16d0bbcc5f80aa8c1b509e45678010200c2e37431a437ee311a737aecd3caae1213db353847f33792fd539e380bdb4d440100005d5aef298770e2702448af2ce014f8bfcded5896df5006a44b5f1b6020007aeb01010091484f513003fcdb25f336b9b56dafcb05fbc739593ab573a2c6516b344ca5320201000000027b0a1b12c7c9e48015e78d3a08a4d62e439387df7e0d7a810ebd4af37661daaa000000006a47304402207d972759afba7c0ffa6cfbbf39a31c2aeede1dae28d8841db56c6dd1197d56a20220076a390948c235ba8e72b8e43a7b4d4119f1a81a77032aa6e7b7a51be5e13845412103f78ec31cf94ca8d75fb1333ad9fc884e2d489422034a1efc9d66a3b72eddca0fffffffff7f36874f858fb43ffcf4f9e3047825619bad0e92d4b9ad4ba5111d1101cbddfe010000006a473044022043f048043d56eb6f75024808b78f18808b7ab45609e4c4c319e3a27f8246fc3002204b67766b62f58bf6f30ea608eaba76b8524ed49f67a90f80ac08a9b96a6922cd41210254a583c1c51a06e10fab79ddf922915da5f5c1791ef87739f40cb68638397248ffffffff03e8030000000000001976a914b08f70bc5010fb026de018f19e7792385a146b4a88acf3010000000000001976a9147d48635f889372c3da12d75ce246c59f4ab907ed88acf7000000000000001976a914b8fbd58685b6920d8f9a8f1b274d8696708b51b088ac00000000010001000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006b483045022100e47fbd96b59e2c22be273dcacea74a4be568b3e61da7eddddb6ce43d459c4cf202201a580f3d9442d5dce3f2ced03256ca147bcd230975a6067954e22415715f4490412102b0c8980f5d2cab77c92c68ac46442feba163a9d306913f6a34911fc618c3c4e7ffffffff0188130000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988ac0000000000"
	const nlockTime = "0100beef01fe4e6d0c001002fd9c67028ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160cfd9d6700db1332728830a58c83a5970dcd111a575a585b43b0492361ea8082f41668f8bd01fdcf3300e568706954aae516ef6df7b5db7828771a1f3fcf1b6d65389ec8be8c46057a3c01fde6190001a6028d13cc988f55c8765e3ffcdcfc7d5185a8ebd68709c0adbe37b528557b01fdf20c001cc64f09a217e1971cabe751b925f246e3c2a8e145c49be7b831eaea3e064d7501fd7806009ccf122626a20cdb054877ef3f8ae2d0503bb7a8704fdb6295b3001b5e8876a101fd3d0300aeea966733175ff60b55bc77edcb83c0fce3453329f51195e5cbc7a874ee47ad01fd9f0100f67f50b53d73ffd6e84c02ee1903074b9a5b2ac64c508f7f26349b73cca9d7e901ce006ce74c7beed0c61c50dda8b578f0c0dc5a393e1f8758af2fb65edf483afcaa68016600e32475e17bdd141d62524d0005989dd1db6ca92c6af70791b0e4802be4c5c8c1013200b88162f494f26cc3a1a4a7dcf2829a295064e93b3dbb2f72e21a73522869277a011800a938d3f80dd25b6a3a80e450403bf7d62a1068e2e4b13f0656c83f764c55bb77010d006feac6e4fea41c37c508b5bfdc00d582f6e462e6754b338c95b448df37bd342c010700bf5448356be23b2b9afe53d00cee047065bbc16d0bbcc5f80aa8c1b509e45678010200c2e37431a437ee311a737aecd3caae1213db353847f33792fd539e380bdb4d440100005d5aef298770e2702448af2ce014f8bfcded5896df5006a44b5f1b6020007aeb01010091484f513003fcdb25f336b9b56dafcb05fbc739593ab573a2c6516b344ca5320201000000027b0a1b12c7c9e48015e78d3a08a4d62e439387df7e0d7a810ebd4af37661daaa000000006a47304402207d972759afba7c0ffa6cfbbf39a31c2aeede1dae28d8841db56c6dd1197d56a20220076a390948c235ba8e72b8e43a7b4d4119f1a81a77032aa6e7b7a51be5e13845412103f78ec31cf94ca8d75fb1333ad9fc884e2d489422034a1efc9d66a3b72eddca0fffffffff7f36874f858fb43ffcf4f9e3047825619bad0e92d4b9ad4ba5111d1101cbddfe010000006a473044022043f048043d56eb6f75024808b78f18808b7ab45609e4c4c319e3a27f8246fc3002204b67766b62f58bf6f30ea608eaba76b8524ed49f67a90f80ac08a9b96a6922cd41210254a583c1c51a06e10fab79ddf922915da5f5c1791ef87739f40cb68638397248ffffffff03e8030000000000001976a914b08f70bc5010fb026de018f19e7792385a146b4a88acf3010000000000001976a9147d48635f889372c3da12d75ce246c59f4ab907ed88acf7000000000000001976a914b8fbd58685b6920d8f9a8f1b274d8696708b51b088ac00000000010001000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006a473044022052b71d3f9701e29419a0feb77f4eed2d1eeba113806c66956a4516531a5c8e7d022056ac7694a79ad45c105d28954034e71d30a7ec9d4dd0782098cedd32f8952a4a412102b0c8980f5d2cab77c92c68ac46442feba163a9d306913f6a34911fc618c3c4e7ffffffff01c8000000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988ac9f86010000"
	const nlockNseq = "0100beef01fe4e6d0c001002fd9c67028ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160cfd9d6700db1332728830a58c83a5970dcd111a575a585b43b0492361ea8082f41668f8bd01fdcf3300e568706954aae516ef6df7b5db7828771a1f3fcf1b6d65389ec8be8c46057a3c01fde6190001a6028d13cc988f55c8765e3ffcdcfc7d5185a8ebd68709c0adbe37b528557b01fdf20c001cc64f09a217e1971cabe751b925f246e3c2a8e145c49be7b831eaea3e064d7501fd7806009ccf122626a20cdb054877ef3f8ae2d0503bb7a8704fdb6295b3001b5e8876a101fd3d0300aeea966733175ff60b55bc77edcb83c0fce3453329f51195e5cbc7a874ee47ad01fd9f0100f67f50b53d73ffd6e84c02ee1903074b9a5b2ac64c508f7f26349b73cca9d7e901ce006ce74c7beed0c61c50dda8b578f0c0dc5a393e1f8758af2fb65edf483afcaa68016600e32475e17bdd141d62524d0005989dd1db6ca92c6af70791b0e4802be4c5c8c1013200b88162f494f26cc3a1a4a7dcf2829a295064e93b3dbb2f72e21a73522869277a011800a938d3f80dd25b6a3a80e450403bf7d62a1068e2e4b13f0656c83f764c55bb77010d006feac6e4fea41c37c508b5bfdc00d582f6e462e6754b338c95b448df37bd342c010700bf5448356be23b2b9afe53d00cee047065bbc16d0bbcc5f80aa8c1b509e45678010200c2e37431a437ee311a737aecd3caae1213db353847f33792fd539e380bdb4d440100005d5aef298770e2702448af2ce014f8bfcded5896df5006a44b5f1b6020007aeb01010091484f513003fcdb25f336b9b56dafcb05fbc739593ab573a2c6516b344ca5320201000000027b0a1b12c7c9e48015e78d3a08a4d62e439387df7e0d7a810ebd4af37661daaa000000006a47304402207d972759afba7c0ffa6cfbbf39a31c2aeede1dae28d8841db56c6dd1197d56a20220076a390948c235ba8e72b8e43a7b4d4119f1a81a77032aa6e7b7a51be5e13845412103f78ec31cf94ca8d75fb1333ad9fc884e2d489422034a1efc9d66a3b72eddca0fffffffff7f36874f858fb43ffcf4f9e3047825619bad0e92d4b9ad4ba5111d1101cbddfe010000006a473044022043f048043d56eb6f75024808b78f18808b7ab45609e4c4c319e3a27f8246fc3002204b67766b62f58bf6f30ea608eaba76b8524ed49f67a90f80ac08a9b96a6922cd41210254a583c1c51a06e10fab79ddf922915da5f5c1791ef87739f40cb68638397248ffffffff03e8030000000000001976a914b08f70bc5010fb026de018f19e7792385a146b4a88acf3010000000000001976a9147d48635f889372c3da12d75ce246c59f4ab907ed88acf7000000000000001976a914b8fbd58685b6920d8f9a8f1b274d8696708b51b088ac00000000010001000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006b483045022100e3da33b9f50ee8c4a2383d92bbfc9cc19367f1c227a8b7c08d11eecfc8ec01f702202ba54cdbf766dfd314c87e29557d18db2c451502a6a23acafa1c970b816f431c412102b0c8980f5d2cab77c92c68ac46442feba163a9d306913f6a34911fc618c3c4e70f27000001c8000000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988ac9f86010000"
	const nseq = "0100beef01fe4e6d0c001002fd9c67028ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160cfd9d6700db1332728830a58c83a5970dcd111a575a585b43b0492361ea8082f41668f8bd01fdcf3300e568706954aae516ef6df7b5db7828771a1f3fcf1b6d65389ec8be8c46057a3c01fde6190001a6028d13cc988f55c8765e3ffcdcfc7d5185a8ebd68709c0adbe37b528557b01fdf20c001cc64f09a217e1971cabe751b925f246e3c2a8e145c49be7b831eaea3e064d7501fd7806009ccf122626a20cdb054877ef3f8ae2d0503bb7a8704fdb6295b3001b5e8876a101fd3d0300aeea966733175ff60b55bc77edcb83c0fce3453329f51195e5cbc7a874ee47ad01fd9f0100f67f50b53d73ffd6e84c02ee1903074b9a5b2ac64c508f7f26349b73cca9d7e901ce006ce74c7beed0c61c50dda8b578f0c0dc5a393e1f8758af2fb65edf483afcaa68016600e32475e17bdd141d62524d0005989dd1db6ca92c6af70791b0e4802be4c5c8c1013200b88162f494f26cc3a1a4a7dcf2829a295064e93b3dbb2f72e21a73522869277a011800a938d3f80dd25b6a3a80e450403bf7d62a1068e2e4b13f0656c83f764c55bb77010d006feac6e4fea41c37c508b5bfdc00d582f6e462e6754b338c95b448df37bd342c010700bf5448356be23b2b9afe53d00cee047065bbc16d0bbcc5f80aa8c1b509e45678010200c2e37431a437ee311a737aecd3caae1213db353847f33792fd539e380bdb4d440100005d5aef298770e2702448af2ce014f8bfcded5896df5006a44b5f1b6020007aeb01010091484f513003fcdb25f336b9b56dafcb05fbc739593ab573a2c6516b344ca5320201000000027b0a1b12c7c9e48015e78d3a08a4d62e439387df7e0d7a810ebd4af37661daaa000000006a47304402207d972759afba7c0ffa6cfbbf39a31c2aeede1dae28d8841db56c6dd1197d56a20220076a390948c235ba8e72b8e43a7b4d4119f1a81a77032aa6e7b7a51be5e13845412103f78ec31cf94ca8d75fb1333ad9fc884e2d489422034a1efc9d66a3b72eddca0fffffffff7f36874f858fb43ffcf4f9e3047825619bad0e92d4b9ad4ba5111d1101cbddfe010000006a473044022043f048043d56eb6f75024808b78f18808b7ab45609e4c4c319e3a27f8246fc3002204b67766b62f58bf6f30ea608eaba76b8524ed49f67a90f80ac08a9b96a6922cd41210254a583c1c51a06e10fab79ddf922915da5f5c1791ef87739f40cb68638397248ffffffff03e8030000000000001976a914b08f70bc5010fb026de018f19e7792385a146b4a88acf3010000000000001976a9147d48635f889372c3da12d75ce246c59f4ab907ed88acf7000000000000001976a914b8fbd58685b6920d8f9a8f1b274d8696708b51b088ac00000000010001000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006b483045022100bcae2985c92e421e25b271d565398bf0fe0ef870d62d4945a1589933ad9f0af102207f9fd4ad14ac0c5507ca317e64231b9ba1f2de1e8a7b13520cccee3e5f0650e4412102b0c8980f5d2cab77c92c68ac46442feba163a9d306913f6a34911fc618c3c4e70f27000001c8000000000000001976a9148a8c4546a95e6fc8d18076a9980d59fd882b4e6988ac0000000000"
	const bump = "0100beef01fe4e6d0c001002fd909002088a382ec07a8cf47c6158b68e5822852362102d8571482d1257e0b7527e1882fd91900065cb01218f2506bb51155d243e4d6b32d69d1b5f2221c52e26963cfd8cf7283201fd4948008d7a44ae384797b0ae84db0c857e8c1083425d64d09ef8bc5e2e9d270677260501fd25240060f38aa33631c8d70adbac1213e7a5b418c90414e919e3a12ced63dd152fd85a01fd1312005ff132ee64a7a0c79150a29f66ef861e552d3a05b47d6303f5d8a2b2a09bc61501fd080900cc0baf21cf06b9439dfe05dce9bdb14ddc2ca2d560b1138296ef5769851a84b301fd85040063ccb26232a6e1d3becdb47a0f19a67a562b754e8894155b3ae7bba10335ce5101fd430200e153fc455a0f2c8372885c11af70af904dcf44740b9ebf3b3e5b2234cce550bc01fd20010077d5ea69d1dcc379dde65d6adcebde1838190118a8fae928c037275e78bd87910191000263e4f31684a25169857f2788aeef603504931f92585f02c4c9e023b2aa43d1014900de72292e0b3e5eeacfa2b657bf4d46c885559b081ee78632a99b318c1148d85c01250068a5f831ca99b9e7f3720920d6ea977fd2ab52b83d1a6567dafa4c8cafd941ed0113006a0b91d83f9056b702d6a8056af6365c7da626fc3818b815dd4b0de22d05450f0108009876ce56b68545a75859e93d200bdde7880d46f39384818b259ed847a9664ddf010500990bc5e95cacbc927b5786ec39a183f983fe160d52829cf47521c7eb369771c30103004fe794e50305f590b6010a51d050bf47dfeaabfdb949c5ee0673f577a59537d70100004dad44a358aea4d8bc1917912539901f5ae44e07a4748e1a9d3018814b0759d00201000000027b0a1b12c7c9e48015e78d3a08a4d62e439387df7e0d7a810ebd4af37661daaa000000006a47304402207d972759afba7c0ffa6cfbbf39a31c2aeede1dae28d8841db56c6dd1197d56a20220076a390948c235ba8e72b8e43a7b4d4119f1a81a77032aa6e7b7a51be5e13845412103f78ec31cf94ca8d75fb1333ad9fc884e2d489422034a1efc9d66a3b72eddca0fffffffff7f36874f858fb43ffcf4f9e3047825619bad0e92d4b9ad4ba5111d1101cbddfe010000006a473044022043f048043d56eb6f75024808b78f18808b7ab45609e4c4c319e3a27f8246fc3002204b67766b62f58bf6f30ea608eaba76b8524ed49f67a90f80ac08a9b96a6922cd41210254a583c1c51a06e10fab79ddf922915da5f5c1791ef87739f40cb68638397248ffffffff03e8030000000000001976a914b08f70bc5010fb026de018f19e7792385a146b4a88acf3010000000000001976a9147d48635f889372c3da12d75ce246c59f4ab907ed88acf7000000000000001976a914b8fbd58685b6920d8f9a8f1b274d8696708b51b088ac00000000010001000000018ae36502fdc82837319362c488fb9cb978e064daf600bbfc48389663fc5c160c000000006a47304402204a04841f6f626d30e21200e1c404ea80e319b643fe86f08e709413a89a493a4b022038a2e3e25a813d8d540c1a572fa8ec5fa2d2434bcea78d17902dcccddcc1c9484121028fd1afeee81361e801800afb264e35cdce3037ec6f7dc4f1d1eaba7ad519c948ffffffff01c8000000000000001976a9148ce2d21f9a75e98600be76b25b91c4fef6b40bcd88ac0000000000"

	tcs := []struct {
		name          string
		beef          string
		expectedError error
	}{

		{
			name:          "SPV on someone else UTXOs",
			beef:          someoneElse,
			expectedError: errors.New("invalid script"),
		},
		{
			name:          "SPV on trying to spend more satoshis in outputs then in inputs",
			beef:          tooMuch,
			expectedError: errors.New("invalid input and output sum, outputs can not be larger than inputs"),
		},
		{
			name:          "SPV on unsupported LockTime",
			beef:          nlockTime,
			expectedError: errors.New("unexpected transaction with nLockTime"),
		},
		{
			name:          "SPV on unsupported LockTime and unsupported Sequence",
			beef:          nlockNseq,
			expectedError: errors.New("unexpected transaction with nLockTime"),
		},
		{
			name:          "SPV on supported LockTime but unsupported Sequence",
			beef:          nseq,
			expectedError: errors.New("unexpected transaction with nSequence"),
		},
		{
			name:          "SPV on valid BUMP from other tx",
			beef:          bump,
			expectedError: errors.New("invalid BUMP - input mined ancestor is not present in BUMPs"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// given
			validDecodedBeef, err := beef.DecodeBEEF(tc.beef)
			require.Nil(t, err)

			//when
			err = ExecuteSimplifiedPaymentVerification(context.Background(), validDecodedBeef, new(mockServiceProvider))
			require.NotNil(t, err)

			//then
			require.Equal(t, tc.expectedError, err)
		})
	}

}

// Mock implementation of a service provider
type mockServiceProvider struct{}

// VerifyMerkleRoots is a mock implementation of this interface
func (m *mockServiceProvider) VerifyMerkleRoots(_ context.Context, _ []*MerkleRootConfirmationRequestItem) error {
	// Verify the merkle roots
	return nil
}