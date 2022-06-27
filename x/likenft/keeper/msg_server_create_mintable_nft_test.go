package keeper_test

import (
	"testing"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestCreateBlindBoxContentNormal(t *testing.T) {
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
	nftInput := types.NFTInput{
		Uri:      "ipfs://123456",
		UriHash:  "123456",
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
	}

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

	accountKeeper.EXPECT().GetAccount(gomock.Any(), ownerAddressBytes).Return(authtypes.NewBaseAccountWithAddress(ownerAddressBytes))
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), ownerAddressBytes, authtypes.FeeCollectorName, gomock.Any()).Return(nil)

	nftKeeper.EXPECT().UpdateClass(gomock.Any(), gomock.Any()).Return(nil)

	// Run
	res, err := msgServer.CreateBlindBoxContent(goCtx, &types.MsgCreateBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
		Input:   nftInput,
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, &types.MsgCreateBlindBoxContentResponse{
		BlindBoxContent: types.BlindBoxContent{
			ClassId: classId,
			Id:      mintableId,
			Input:   nftInput,
		},
	}, res)

	created, found := keeper.GetBlindBoxContent(ctx, classId, mintableId)
	require.True(t, found)
	require.Equal(t, types.BlindBoxContent{
		ClassId: classId,
		Id:      mintableId,
		Input:   nftInput,
	}, created)

	ctrl.Finish()
}

func TestCreateBlindBoxContentClassNotFound(t *testing.T) {
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
	nftInput := types.NFTInput{
		Uri:      "ipfs://123456",
		UriHash:  "123456",
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
	}

	// Mock calls
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{}, false).MinTimes(1)

	// Run
	res, err := msgServer.CreateBlindBoxContent(goCtx, &types.MsgCreateBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
		Input:   nftInput,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotFound.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestCreateBlindBoxContentBadRelation(t *testing.T) {
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
	nftInput := types.NFTInput{
		Uri:      "ipfs://123456",
		UriHash:  "123456",
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
	}

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
	res, err := msgServer.CreateBlindBoxContent(goCtx, &types.MsgCreateBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
		Input:   nftInput,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyAccount.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestCreateBlindBoxContentAlreadyMinted(t *testing.T) {
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
	nftInput := types.NFTInput{
		Uri:      "ipfs://123456",
		UriHash:  "123456",
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
	}

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
	res, err := msgServer.CreateBlindBoxContent(goCtx, &types.MsgCreateBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
		Input:   nftInput,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrCannotUpdateClassWithMintedTokens.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestCreateBlindBoxContentMaxSupplyReached(t *testing.T) {
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
	nftInput := types.NFTInput{
		Uri:      "ipfs://123456",
		UriHash:  "123456",
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
	}

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
			ContentCount: uint64(5),
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
	res, err := msgServer.CreateBlindBoxContent(goCtx, &types.MsgCreateBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
		Input:   nftInput,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftNoSupply.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestCreateBlindBoxContentNotOwner(t *testing.T) {
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
	nftInput := types.NFTInput{
		Uri:      "ipfs://123456",
		UriHash:  "123456",
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
	}

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
	res, err := msgServer.CreateBlindBoxContent(goCtx, &types.MsgCreateBlindBoxContent{
		Creator: notOwnerAddress,
		ClassId: classId,
		Id:      mintableId,
		Input:   nftInput,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())
	require.Nil(t, res)

	ctrl.Finish()
}

func TestCreateBlindBoxContentAlreadyExist(t *testing.T) {
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
	nftInput := types.NFTInput{
		Uri:      "ipfs://123456",
		UriHash:  "123456",
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
	}

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

	// called when seeding existing mintable
	nftKeeper.EXPECT().UpdateClass(gomock.Any(), gomock.Any()).Return(nil)

	keeper.SetBlindBoxContent(ctx, types.BlindBoxContent{
		ClassId: classId,
		Id:      mintableId,
	})

	// Run
	res, err := msgServer.CreateBlindBoxContent(goCtx, &types.MsgCreateBlindBoxContent{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      mintableId,
		Input:   nftInput,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrMintableNftAlreadyExists.Error())
	require.Nil(t, res)

	ctrl.Finish()
}
