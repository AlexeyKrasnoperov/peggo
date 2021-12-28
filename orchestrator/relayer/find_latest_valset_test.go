package relayer

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/umee-network/peggo/mocks"
	peggyMocks "github.com/umee-network/peggo/mocks/peggy"
	"github.com/umee-network/umee/x/peggy/types"
)

func TestFindLatestValset(t *testing.T) {
	t.Run("ok. 1 member", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
		mockQClient := mocks.NewMockQueryClient(mockCtrl)
		ethProvider := mocks.NewMockEVMProviderWithRet(mockCtrl)
		mockPeggyContract := peggyMocks.NewMockContract(mockCtrl)

		peggyAddress := ethcmn.HexToAddress("0x3bdf8428734244c9e5d82c95d125081939d6d42d")
		fromAddress := ethcmn.HexToAddress("0xd8da6bf26964af9d7eed9e03e53415d37aa96045")

		ethProvider.EXPECT().HeaderByNumber(gomock.Any(), nil).Return(&ethtypes.Header{
			Number: big.NewInt(112),
		}, nil)
		ethProvider.EXPECT().PendingNonceAt(gomock.Any(), fromAddress).Return(uint64(0), nil).AnyTimes()

		mockPeggyContract.EXPECT().FromAddress().Return(fromAddress).AnyTimes()
		mockPeggyContract.EXPECT().Address().Return(peggyAddress).AnyTimes()
		mockPeggyContract.EXPECT().GetValsetNonce(gomock.Any(), fromAddress).Return(big.NewInt(2), nil)

		mockQClient.EXPECT().ValsetRequest(gomock.Any(), &types.QueryValsetRequestRequest{
			Nonce: uint64(2),
		}).Return(&types.QueryValsetRequestResponse{Valset: &types.Valset{
			Nonce: 2,
		}}, nil)

		// FilterValsetUpdatedEvent
		ethProvider.EXPECT().FilterLogs(
			gomock.Any(),
			MatchFilterQuery(ethereum.FilterQuery{
				FromBlock: new(big.Int).SetUint64(0),
				ToBlock:   new(big.Int).SetUint64(112),
				Addresses: []ethcmn.Address{peggyAddress},
				Topics:    [][]ethcmn.Hash{{ethcmn.HexToHash("0x76d08978c024a4bf8cbb30c67fd78fcaa1827cbc533e4e175f36d07e64ccf96a")}, {}},
			})).
			Return(
				// The test data is from a real tx: https://goerli.etherscan.io/tx/0x79a63e4fdcadb35bc89d6aab9ca2a2c80916817744f472901375290c548e0022#eventlog
				[]ethtypes.Log{
					{
						Address:     peggyAddress,
						Topics:      []ethcmn.Hash{ethcmn.HexToHash("0x76d08978c024a4bf8cbb30c67fd78fcaa1827cbc533e4e175f36d07e64ccf96a"), ethcmn.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000019")},
						Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000005a64fe82628217900ced80bf3747b5ef88bfa21000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000ffffffff"),
						BlockNumber: 3,
						TxHash:      ethcmn.HexToHash("0x0"),
						TxIndex:     2,
						BlockHash:   ethcmn.HexToHash("0x0"),
						Index:       1,
						Removed:     false,
					},
				},
				nil,
			).Times(1)

		relayer := peggyRelayer{
			logger:            logger,
			cosmosQueryClient: mockQClient,
			peggyContract:     mockPeggyContract,
			ethProvider:       ethProvider,
		}

		valset, err := relayer.FindLatestValset(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, valset)
		assert.Len(t, valset.Members, 1)
	})

	t.Run("ok. 99 members", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
		mockQClient := mocks.NewMockQueryClient(mockCtrl)
		ethProvider := mocks.NewMockEVMProviderWithRet(mockCtrl)
		mockPeggyContract := peggyMocks.NewMockContract(mockCtrl)

		peggyAddress := ethcmn.HexToAddress("0x3bdf8428734244c9e5d82c95d125081939d6d42d")
		fromAddress := ethcmn.HexToAddress("0xd8da6bf26964af9d7eed9e03e53415d37aa96045")

		ethProvider.EXPECT().HeaderByNumber(gomock.Any(), nil).Return(&ethtypes.Header{
			Number: big.NewInt(112),
		}, nil)
		ethProvider.EXPECT().PendingNonceAt(gomock.Any(), fromAddress).Return(uint64(0), nil).AnyTimes()

		mockPeggyContract.EXPECT().FromAddress().Return(fromAddress).AnyTimes()
		mockPeggyContract.EXPECT().Address().Return(peggyAddress).AnyTimes()
		mockPeggyContract.EXPECT().GetValsetNonce(gomock.Any(), fromAddress).Return(big.NewInt(2), nil)

		mockQClient.EXPECT().ValsetRequest(gomock.Any(), &types.QueryValsetRequestRequest{
			Nonce: uint64(2),
		}).Return(&types.QueryValsetRequestResponse{Valset: &types.Valset{
			Nonce: 2,
		}}, nil)

		// FilterValsetUpdatedEvent
		ethProvider.EXPECT().FilterLogs(
			gomock.Any(),
			MatchFilterQuery(ethereum.FilterQuery{
				FromBlock: new(big.Int).SetUint64(0),
				ToBlock:   new(big.Int).SetUint64(112),
				Addresses: []ethcmn.Address{peggyAddress},
				Topics:    [][]ethcmn.Hash{{ethcmn.HexToHash("0x76d08978c024a4bf8cbb30c67fd78fcaa1827cbc533e4e175f36d07e64ccf96a")}, {}},
			})).
			Return(
				// The test data is from a real tx: https://goerli.etherscan.io/tx/0x4714abe3e48c4f730dd6e851cff83ab4baed33f7ef1991e504722ef4d28fd30f#eventlog
				[]ethtypes.Log{
					{
						Address:     peggyAddress,
						Topics:      []ethcmn.Hash{ethcmn.HexToHash("0x76d08978c024a4bf8cbb30c67fd78fcaa1827cbc533e4e175f36d07e64ccf96a"), ethcmn.HexToHash("0x0000000000000000000000000000000000000000000000000000000000007cac")},
						Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000290000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000d200000000000000000000000000000000000000000000000000000000000000063000000000000000000000000005362286505c44486f2091db1113c6b97ffb56900000000000000000000000001f71f2a59d19030bf79961e3e57c82922feecca00000000000000000000000003246e835bb5b45d2058a79c21f5bfd632327b610000000000000000000000000c2b90c5a2026426444da9015ded175e6a90f89300000000000000000000000016f2c385b2158916d319b271db31feff1b2d32b30000000000000000000000001b80db8b22369344aa13996df8497e2ccad0774a0000000000000000000000002043fc83d3e031e80680a6e2fa52c5b1c85112b90000000000000000000000002065f9aed8e2cafb8d02dcf9a664994bf39a80c7000000000000000000000000222e7efe38c218cbc251a11d4f95e1cfb2630a8e00000000000000000000000024c44aa9ad5171eb62923ab7cbef9f927e94e00800000000000000000000000026a643429a2f52f9988d3bcea17fdbcd12bee18800000000000000000000000029774aed3c6b6dbc0d4c6652cd6197df1e65299200000000000000000000000029b63a7bd276f9df7e87efba07b34b92b9418eb90000000000000000000000002cd8aee645a981455ee7ea1889531a1eaf733a2f0000000000000000000000002a59dd215a91338b2490a9c6036d597cab3eaa4900000000000000000000000034737d904e18637fc40034d8dd4955b3180c688900000000000000000000000037d5bd3a9fd787790a851874338fe168d5d1cdae000000000000000000000000384ca1d336dca6497ee5e7020980dbde2c611aee0000000000000000000000003995a1e56ef395b9427633292939cc19a0e6e5220000000000000000000000003ab20d8fc60bdeddc1195c82efd971024e7f90f40000000000000000000000003b686e73f3c7d5e9e5f1fdfa4ef41ca0e5e1f60a0000000000000000000000003ae6c9b11602ea986f4e6f9cadb0d2014e4aedc10000000000000000000000003c921221ba740847862bc90aeb627de5805996d100000000000000000000000041841f9260d50496361ca1218bb3dd0b6574ccb7000000000000000000000000447df38586787d3cb145f4a1367c9aea294ca9be00000000000000000000000049101b60abfa293c14c129939a1b747a9daf54ac0000000000000000000000004b11ba0561fc03249b90d6e81aaebb82ca117b850000000000000000000000004c478f19f6330eaa2c9879c9f516c0f3af3177f20000000000000000000000004e561711f16f9d317beb61491e9b0d0cb6f202d6000000000000000000000000517593dcc43b2ac25463a6c874d315a714f644e100000000000000000000000054c800f2570fe8b3ec07f3aa92dbec448acc2efe00000000000000000000000059c4fc02d62e6a6380df5569bd1f50f2b56d99640000000000000000000000006129a4b88af189df6dc7702381e0a8e9243c31aa00000000000000000000000062e536be93d62c581c0f4fca4d026cf4b3021567000000000000000000000000653888070bddbd9c3bb07359ddf6df37498ac6fd0000000000000000000000006993a4a66d4b5dbbab0e609582df3c38e2d2d5a50000000000000000000000006adaf8552678e6f9a97a83ec99fbb6192b1680ca0000000000000000000000006a40b4a7db61a1c64bc922a30d0dd402f973689b0000000000000000000000007042590d89fce92a68e4af39d6ebc79fb0f77be70000000000000000000000007173feab4811282bf22b865420c4cdb0263a52fc000000000000000000000000743802032d3547b4254ce176fb791a81c8d2288b0000000000000000000000007d31889469ad6c7ee742e30c5de593384a3f160e0000000000000000000000007d3f251784091edd12e99aff2047977f89a8639900000000000000000000000081e2609303add0150a1e7c4b885a8f619ba1ca37000000000000000000000000839809d5eb842b7644c34954c6aa23808781cdd50000000000000000000000008d9b999c276cd0248dbb198a81badc4d295167c90000000000000000000000008fc2755b6b5aea0f984e5ab75d1523423a11d2c300000000000000000000000096ea1f6e078cef2d1a327d72c01de2d51a7757d300000000000000000000000098529d938abf84075a5e17b1dc7c628c1bb805df0000000000000000000000009ec26a78f4d850d41269ad23f328227019bb35050000000000000000000000009fb438f1a292fbdaf4f569a8f2e2e2972c0172ec0000000000000000000000009e2d1dd7551385c94bbe59858c29ea8d4c98ffb60000000000000000000000009e9959fee01da59b42853329422bd20355f90a4f000000000000000000000000a602d18eec7565766df301bf4ef18253dde19ec1000000000000000000000000a925763841c5e9ebcfaf442f055ec3371f64318d000000000000000000000000b0ce763127cd1de2a01d7831c12f98104e58bdcf000000000000000000000000b4a680180a82f51881a477deb46e1c802cb6ab64000000000000000000000000bd670f22933a17f4747f88cc9d6ca0e586226c31000000000000000000000000baba0ea6d8b8d9cd19f6f1aed2cdb76888d273d4000000000000000000000000c2b9671db9d00a8136ca1b2faa5c5a46cfa40e6b000000000000000000000000c5bc8523c94d7c87a03976f9ab0867d7fdf458ee000000000000000000000000c9aa59bf68ff97fc67cdade3f20ed8220bf6762b000000000000000000000000ce2daae917fe14ae0fef173aaab9e4fc955b8336000000000000000000000000efc37a45a5515ffa8682194b3db53254673ff38d000000000000000000000000f0b5cf17d48500e0aefc92aa1671dc8bd0fb421f000000000000000000000000f51ac58e5faaab86c7ef93796530db94e36a9977000000000000000000000000f71df0d9e715ff9c337f84cf57963aea2362fe08000000000000000000000000f8e581270fb2f6c54da86037ffe19848d7b57235000000000000000000000000f96bfcf9eb509eb407784fdb56305eaa3671d74a000000000000000000000000fd662bb8aea2c35efe9135749727c58e6218112e000000000000000000000000ff518c080fe35bbbf853be2ec8311ecd61d14c6b000000000000000000000000a163f6355abe610edf8360861ee6e429719fccb7000000000000000000000000a7705d98610c146438e04fed78e6c7d2e1d8dcdd000000000000000000000000ac1f897344e0038f0af2e7a318dd125e73d46358000000000000000000000000b3c8d8d0513bf9878bee148aa0d46ba59a37baaa000000000000000000000000b4ce3a61f9c4bbf8009b02a15cb8ef5f32391980000000000000000000000000b5270d7d8b7e6d371e01950ab0268f0d559c3ae5000000000000000000000000b799c871ba410f97b8b3ebc12d774834f4dfaf22000000000000000000000000bbc0916f2933e734d54ea57749ae5f5980f123c5000000000000000000000000c01bc2728960decdb6633bc5d7f441972cd41caa000000000000000000000000c2f5a910a8f831808b959d2084f19f71ea5eed57000000000000000000000000cb0303deaa03fde8cbc21a320fa1a8d22e96ea4d000000000000000000000000cc1c7e2974e3203def047e7cd1868709b6c61e8d000000000000000000000000d0bcb91c88b640dde5d0476382e074b5b6e0bf3d000000000000000000000000d3d77cb7207d4a40c19bd0068d9919bf324d2744000000000000000000000000dcf93a4cd868545752b0c6fdd9467a36e85fc56e000000000000000000000000dcf402433caf65548524873727d5cfbb27c8e7e0000000000000000000000000e1d08d0f195e15f6dd237f87b4a775185db05067000000000000000000000000ea572284a3fdee1b237ebdffb444cacc39e232d4000000000000000000000000f02b8a3a0e9439b79ed54b1ec7fa64f5a4905688000000000000000000000000f2e678bf92dff213836aea6e6d261f78fad034ad0000000000000000000000001aacd42b608626953964fca1a0a4ec924a82407b000000000000000000000000678db247391df12fd83c0f4fc4a29d09a13de70d0000000000000000000000005101d13dfb7153823ca85e0b665685f57640c60f000000000000000000000000b4f8ccaa731c491b962933f4e9b6e0e10241177f000000000000000000000000ab395a1086ef9fdc2167fe520208aa443995bff40000000000000000000000008f534ae38cf21340c493e1873aab9ba1074e7d9a0000000000000000000000003d5fc27024960ff814c4a83ae7c8395bf2e1fd45000000000000000000000000c1c517f6c7c236235cdd2438576362346162eabf00000000000000000000000000000000000000000000000000000000000000630000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3b44d0000000000000000000000000000000000000000000000000000000002a3a4bb0000000000000000000000000000000000000000000000000000000002a311b70000000000000000000000000000000000000000000000000000000002a307560000000000000000000000000000000000000000000000000000000002a3059b0000000000000000000000000000000000000000000000000000000002a3059b00000000000000000000000000000000000000000000000000000000029efd04000000000000000000000000000000000000000000000000000000000000114b000000000000000000000000000000000000000000000000000000000000114b"),
						BlockNumber: 3,
						TxHash:      ethcmn.HexToHash("0x0"),
						TxIndex:     2,
						BlockHash:   ethcmn.HexToHash("0x0"),
						Index:       1,
						Removed:     false,
					},
				},
				nil,
			).Times(1)

		relayer := peggyRelayer{
			logger:            logger,
			cosmosQueryClient: mockQClient,
			peggyContract:     mockPeggyContract,
			ethProvider:       ethProvider,
		}

		valset, err := relayer.FindLatestValset(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, valset)
		assert.Len(t, valset.Members, 99)
	})
}

func TestCheckIfValsetsDiffer(t *testing.T) {
	// this function doesn't return a value. Running different scenarios just to increase code coverage.

	t.Run("ok", func(t *testing.T) {
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

		relayer := peggyRelayer{
			logger: logger,
		}

		relayer.checkIfValsetsDiffer(&types.Valset{}, &types.Valset{})
		relayer.checkIfValsetsDiffer(nil, &types.Valset{})
		relayer.checkIfValsetsDiffer(nil, &types.Valset{Nonce: 2})
		relayer.checkIfValsetsDiffer(&types.Valset{Nonce: 12}, &types.Valset{Nonce: 11})
		relayer.checkIfValsetsDiffer(&types.Valset{}, &types.Valset{Members: []*types.BridgeValidator{{EthereumAddress: "0x0"}}})
	})

}

func TestBridgeValidator(t *testing.T) {
	var bridgeValidators BridgeValidators = []*types.BridgeValidator{
		{
			EthereumAddress: "0x0",
			Power:           10000,
		},
		{
			EthereumAddress: "0x1",
			Power:           20000,
		},
		{
			EthereumAddress: "0x2",
			Power:           20000,
		},
	}
	bridgeValidators.Sort()
	assert.Equal(t, bridgeValidators[0].EthereumAddress, "0x1")
	assert.False(t, bridgeValidators.HasDuplicates())
	assert.Equal(t, []uint64{20000, 20000, 10000}, bridgeValidators.GetPowers())

}

type matchFilterQuery struct {
	q ethereum.FilterQuery
}

func (m *matchFilterQuery) Matches(input interface{}) bool {
	q, ok := input.(ethereum.FilterQuery)
	if ok {

		if q.BlockHash != m.q.BlockHash {
			return false
		}

		if q.FromBlock.Int64() != m.q.FromBlock.Int64() {
			return false
		}

		if q.ToBlock.Int64() != m.q.ToBlock.Int64() {
			return false
		}

		if !assert.ObjectsAreEqual(q.Addresses, m.q.Addresses) {
			return false
		}

		// Comparing 2 slices of slices seems to be a bit tricky.

		if len(q.Topics) != len(m.q.Topics) {
			return false
		}

		for i := range q.Topics {
			if len(q.Topics[i]) != len(m.q.Topics[i]) {
				return false
			}

			for j := range q.Topics[i] {
				if q.Topics[i][j] != m.q.Topics[i][j] {
					return false
				}
			}
		}
		return true
	}

	return false
}

func (m *matchFilterQuery) String() string {
	return fmt.Sprintf("is equal to %v (%T)", m.q, m.q)
}

func MatchFilterQuery(q ethereum.FilterQuery) gomock.Matcher {
	return &matchFilterQuery{q: q}
}
