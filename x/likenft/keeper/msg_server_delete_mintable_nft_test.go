package keeper_test

import (
	"testing"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestDeleteBlindBoxContentNormal(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	mintableId := "mintable1"

	// Mock calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"1234": "5678"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(5),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        1000,
					},
				},
				RevealTime: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		BlindBoxState: types.BlindBoxState{
			ToBeRevealed: true,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}, true).MinTimes(1)
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	nftKeeper.EXPECT().GetTotalSupply(gomock.Any(), classId).Return(uint64(0))

	// once at seeding, once at delete
	nftKeeper.EXPECT().UpdateClass(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	keeper.SetBlindBoxContent(ctx, types.BlindBoxContent{
		ClassId: classId,
		Id:      mintableId,
	})

	// Run
	res, err := msgServer.DeleteBlindBoxContent(goCtx, &types.MsgDeleteBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
	})

	// Check output
	require.NoError(t, err)
	require.NotNil(t, res)

	_, found := keeper.GetBlindBoxContent(ctx, classId, mintableId)
	require.False(t, found)

	ctrl.Finish()
}

func TestDeleteBlindBoxContentClassNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	mintableId := "mintable1"

	// Mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{}, false).MinTimes(1)
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	// Run
	res, err := msgServer.DeleteBlindBoxContent(goCtx, &types.MsgDeleteBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotFound.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestDeleteBlindBoxContentBadRelation(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, _ := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	mintableId := "mintable1"

	// Mock calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"1234": "5678"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(5),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        1000,
					},
				},
				RevealTime: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		BlindBoxState: types.BlindBoxState{
			ToBeRevealed: true,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}, true).MinTimes(1)

	// Run
	res, err := msgServer.DeleteBlindBoxContent(goCtx, &types.MsgDeleteBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyAccount.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestDeleteBlindBoxContentAlreadyMinted(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	mintableId := "mintable1"

	// Mock calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"1234": "5678"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(5),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        1000,
					},
				},
				RevealTime: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		BlindBoxState: types.BlindBoxState{
			ToBeRevealed: true,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}, true).MinTimes(1)
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	nftKeeper.EXPECT().GetTotalSupply(gomock.Any(), classId).Return(uint64(1))

	// Run
	res, err := msgServer.DeleteBlindBoxContent(goCtx, &types.MsgDeleteBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrCannotUpdateClassWithMintedTokens.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestDeleteBlindBoxContentNotOwner(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	mintableId := "mintable1"

	// Mock calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"1234": "5678"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(5),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        1000,
					},
				},
				RevealTime: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		BlindBoxState: types.BlindBoxState{
			ToBeRevealed: true,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}, true).MinTimes(1)
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	nftKeeper.EXPECT().GetTotalSupply(gomock.Any(), classId).Return(uint64(0))

	// Run
	notOwnerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	notOwnerAddress, _ := sdk.Bech32ifyAddressBytes("like", notOwnerAddressBytes)
	res, err := msgServer.DeleteBlindBoxContent(goCtx, &types.MsgDeleteBlindBoxContent{
		Creator: notOwnerAddress,
		ClassId: classId,
		Id:      mintableId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestDeleteBlindBoxContentDoNotExist(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
	classId := "likenft1aabbccddeeff"
	mintableId := "mintable1"

	// Mock calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"1234": "5678"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(5),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
					{
						StartTime:        time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						AllowedAddresses: []string{},
						MintPrice:        1000,
					},
				},
				RevealTime: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		BlindBoxState: types.BlindBoxState{
			ToBeRevealed: true,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:          classId,
		Name:        "Class Name",
		Symbol:      "ABC",
		Description: "Testing Class 123",
		Uri:         "ipfs://abcdef",
		UriHash:     "abcdef",
		Data:        classDataInAny,
	}, true).MinTimes(1)
	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	nftKeeper.EXPECT().GetTotalSupply(gomock.Any(), classId).Return(uint64(0))

	// Run
	res, err := msgServer.DeleteBlindBoxContent(goCtx, &types.MsgDeleteBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrMintableNftNotFound.Error())
	require.Nil(t, res)

	ctrl.Finish()
}
