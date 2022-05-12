package keeper_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// Note: only tests control flow & external call counts here, token / class content to be tested in queue or e2e test case

func TestRevealNormalMintToOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	keeper, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, keepertest.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	// hash from mainnet block 1
	hash, err := hex.DecodeString("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash,
	})

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	supply := 100
	mintableCount := 99
	totalSupply := 90
	mintToOwnerCount := mintableCount - totalSupply

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(supply),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        uint64(0),
					},
				},
			},
		},
		MintableCount: uint64(mintableCount),
		ToBeRevealed:  true,
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	class := nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}
	// mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(class, true).AnyTimes()
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})
	nftKeeper.EXPECT().GetTotalSupply(ctx, classId).Return(uint64(totalSupply))
	nftKeeper.EXPECT().Mint(gomock.Any(), gomock.Any(), ownerAddressBytes).Return(nil).Times(mintToOwnerCount)
	nftKeeper.EXPECT().UpdateClass(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	for i := 0; i < mintableCount; i++ {
		keeper.SetMintableNFT(ctx, types.MintableNFT{
			ClassId: classId,
			Id:      fmt.Sprintf("mintable%d", i),
			Input: types.NFTInput{
				Uri: fmt.Sprintf("mintable%d", i),
			},
		})
	}
	var blindTokens []nft.NFT
	for i := 0; i < totalSupply+mintToOwnerCount; i++ {
		nftData := types.NFTData{
			ClassParent:  classData.Parent,
			ToBeRevealed: true,
		}
		nftDataInAny, err := cdctypes.NewAnyWithValue(&nftData)
		require.NoError(t, err)
		blindTokens = append(blindTokens, nft.NFT{
			ClassId: classId,
			Id:      fmt.Sprintf("nft%d", i),
			Data:    nftDataInAny,
		})
	}
	nftKeeper.EXPECT().GetNFTsOfClass(ctx, classId).Return(blindTokens)
	nftKeeper.EXPECT().Update(ctx, gomock.Any()).Return(nil).Times(totalSupply + mintToOwnerCount)

	// call
	err = keeper.RevealMintableNFTs(ctx, classId)
	require.NoError(t, err)

	ctrl.Finish()
}

func TestRevealNormalNoMintToOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	keeper, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, keepertest.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	// hash from mainnet block 1
	hash, err := hex.DecodeString("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash,
	})

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	supply := 100
	mintableCount := 99
	totalSupply := 99

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(supply),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        uint64(0),
					},
				},
			},
		},
		MintableCount: uint64(mintableCount),
		ToBeRevealed:  true,
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	class := nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}
	// mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(class, true).AnyTimes()
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})
	nftKeeper.EXPECT().GetTotalSupply(ctx, classId).Return(uint64(totalSupply))
	nftKeeper.EXPECT().UpdateClass(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	for i := 0; i < mintableCount; i++ {
		keeper.SetMintableNFT(ctx, types.MintableNFT{
			ClassId: classId,
			Id:      fmt.Sprintf("mintable%d", i),
			Input: types.NFTInput{
				Uri: fmt.Sprintf("mintable%d", i),
			},
		})
	}
	var blindTokens []nft.NFT
	for i := 0; i < totalSupply; i++ {
		nftData := types.NFTData{
			ClassParent:  classData.Parent,
			ToBeRevealed: true,
		}
		nftDataInAny, err := cdctypes.NewAnyWithValue(&nftData)
		require.NoError(t, err)
		blindTokens = append(blindTokens, nft.NFT{
			ClassId: classId,
			Id:      fmt.Sprintf("nft%d", i),
			Data:    nftDataInAny,
		})
	}
	nftKeeper.EXPECT().GetNFTsOfClass(ctx, classId).Return(blindTokens)
	nftKeeper.EXPECT().Update(ctx, gomock.Any()).Return(nil).Times(totalSupply)

	// call
	err = keeper.RevealMintableNFTs(ctx, classId)
	require.NoError(t, err)

	ctrl.Finish()
}

func TestRevealClassNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	keeper, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, keepertest.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	// hash from mainnet block 1
	hash, err := hex.DecodeString("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash,
	})

	classId := "likenft1aabbccddeeff"

	// mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{}, false)

	// call
	err = keeper.RevealMintableNFTs(ctx, classId)
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotFound.Error())

	ctrl.Finish()
}

func TestRevealNotBlindBox(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	keeper, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, keepertest.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	// hash from mainnet block 1
	hash, err := hex.DecodeString("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash,
	})

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	supply := 100

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(supply),
			BlindBoxConfig: nil,
		},
		MintableCount: uint64(0),
		ToBeRevealed:  false,
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	class := nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}
	// mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(class, true).AnyTimes()
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	// call
	err = keeper.RevealMintableNFTs(ctx, classId)
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrClassIsNotBlindBox.Error())

	ctrl.Finish()
}

// Note: validateAndGetClassParentAndOwner covered by other cases
// FIXME: refactor to test that utils separately

func TestRevealFailedToMint(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	keeper, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, keepertest.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	// hash from mainnet block 1
	hash, err := hex.DecodeString("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash,
	})

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	supply := 100
	mintableCount := 99
	totalSupply := 90

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(supply),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        uint64(0),
					},
				},
			},
		},
		MintableCount: uint64(mintableCount),
		ToBeRevealed:  true,
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	class := nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}
	// mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(class, true).AnyTimes()
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})
	nftKeeper.EXPECT().GetTotalSupply(ctx, classId).Return(uint64(totalSupply))
	nftKeeper.EXPECT().Mint(gomock.Any(), gomock.Any(), ownerAddressBytes).Return(fmt.Errorf("Failed to mint"))

	// call
	err = keeper.RevealMintableNFTs(ctx, classId)
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrFailedToMintNFT.Error())

	ctrl.Finish()
}

func TestRevealMintableMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	keeper, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, keepertest.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	// hash from mainnet block 1
	hash, err := hex.DecodeString("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash,
	})

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	supply := 100
	mintableCount := 99
	totalSupply := 90
	mintToOwnerCount := mintableCount - totalSupply

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(supply),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        uint64(0),
					},
				},
			},
		},
		MintableCount: uint64(mintableCount),
		ToBeRevealed:  true,
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	class := nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}
	// mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(class, true).AnyTimes()
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})
	nftKeeper.EXPECT().GetTotalSupply(ctx, classId).Return(uint64(totalSupply))
	nftKeeper.EXPECT().Mint(gomock.Any(), gomock.Any(), ownerAddressBytes).Return(nil).Times(mintToOwnerCount)
	nftKeeper.EXPECT().UpdateClass(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	for i := 0; i < mintableCount; i++ {
		keeper.SetMintableNFT(ctx, types.MintableNFT{
			ClassId: classId,
			Id:      fmt.Sprintf("mintable%d", i),
			Input: types.NFTInput{
				Uri: fmt.Sprintf("mintable%d", i),
			},
		})
	}
	nftKeeper.EXPECT().GetNFTsOfClass(ctx, classId).Return([]nft.NFT{})

	// call
	err = keeper.RevealMintableNFTs(ctx, classId)
	require.Error(t, err)
	require.Contains(t, err.Error(), "length mismatch")

	ctrl.Finish()
}

func TestRevealFailToUpdateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	keeper, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, keepertest.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	// hash from mainnet block 1
	hash, err := hex.DecodeString("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash,
	})

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	supply := 100
	mintableCount := 99
	totalSupply := 90
	mintToOwnerCount := mintableCount - totalSupply

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(supply),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        uint64(0),
					},
				},
			},
		},
		MintableCount: uint64(mintableCount),
		ToBeRevealed:  true,
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	class := nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}
	// mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(class, true).AnyTimes()
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})
	nftKeeper.EXPECT().GetTotalSupply(ctx, classId).Return(uint64(totalSupply))
	nftKeeper.EXPECT().Mint(gomock.Any(), gomock.Any(), ownerAddressBytes).Return(nil).Times(mintToOwnerCount)
	nftKeeper.EXPECT().UpdateClass(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	for i := 0; i < mintableCount; i++ {
		keeper.SetMintableNFT(ctx, types.MintableNFT{
			ClassId: classId,
			Id:      fmt.Sprintf("mintable%d", i),
			Input: types.NFTInput{
				Uri: fmt.Sprintf("mintable%d", i),
			},
		})
	}
	var blindTokens []nft.NFT
	for i := 0; i < totalSupply+mintToOwnerCount; i++ {
		nftData := types.NFTData{
			ClassParent:  classData.Parent,
			ToBeRevealed: true,
		}
		nftDataInAny, err := cdctypes.NewAnyWithValue(&nftData)
		require.NoError(t, err)
		blindTokens = append(blindTokens, nft.NFT{
			ClassId: classId,
			Id:      fmt.Sprintf("nft%d", i),
			Data:    nftDataInAny,
		})
	}
	nftKeeper.EXPECT().GetNFTsOfClass(ctx, classId).Return(blindTokens)
	nftKeeper.EXPECT().Update(ctx, gomock.Any()).Return(fmt.Errorf("Failed to update"))

	// call
	err = keeper.RevealMintableNFTs(ctx, classId)
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrFailedToUpdateNFT.Error())

	ctrl.Finish()
}
