package keeper_test

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/golang/mock/gomock"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/likecoin/likecoin-chain/v4/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// test normal
func TestDeleteRoyaltyConfigNormal(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"

	// Mock
	k.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  userAddress,
		ClassIds: []string{classId},
	})
	classData := types.ClassData{
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: userAddress,
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	require.NoError(t, err)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:   classId,
		Data: classDataInAny,
	}, true)

	// Seed old record
	k.SetRoyaltyConfig(ctx, types.RoyaltyConfigByClass{
		ClassId:       classId,
		RoyaltyConfig: types.RoyaltyConfig{},
	})

	// Call
	res, err := msgServer.DeleteRoyaltyConfig(goCtx, &types.MsgDeleteRoyaltyConfig{
		Creator: userAddress,
		ClassId: classId,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgDeleteRoyaltyConfigResponse{}, res)

	// check state
	_, found := k.GetRoyaltyConfig(ctx, classId)
	require.False(t, found)

	ctrl.Finish()
}

// test not exist
func TestDeleteRoyaltyConfigNotExist(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"

	// Mock
	k.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  userAddress,
		ClassIds: []string{classId},
	})
	classData := types.ClassData{
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: userAddress,
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	require.NoError(t, err)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:   classId,
		Data: classDataInAny,
	}, true)

	// do not seed existing record

	// Call
	res, err := msgServer.DeleteRoyaltyConfig(goCtx, &types.MsgDeleteRoyaltyConfig{
		Creator: userAddress,
		ClassId: classId,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrRoyaltyConfigNotFound.Error())

	// check not found
	_, found := k.GetRoyaltyConfig(ctx, classId)
	require.False(t, found)

	ctrl.Finish()
}

// test user not class owner
func TestDeleteRoyaltyConfigNotOwner(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	notUserAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	notUserAddress, _ := sdk.Bech32ifyAddressBytes("like", notUserAddressBytes)

	// Mock
	k.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  notUserAddress,
		ClassIds: []string{classId},
	})
	classData := types.ClassData{
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: notUserAddress,
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	require.NoError(t, err)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:   classId,
		Data: classDataInAny,
	}, true)

	// Seed old record
	k.SetRoyaltyConfig(ctx, types.RoyaltyConfigByClass{
		ClassId:       classId,
		RoyaltyConfig: types.RoyaltyConfig{},
	})

	// Call
	res, err := msgServer.DeleteRoyaltyConfig(goCtx, &types.MsgDeleteRoyaltyConfig{
		Creator: userAddress,
		ClassId: classId,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())

	// Check not changed
	config, found := k.GetRoyaltyConfig(ctx, classId)
	require.True(t, found)
	require.Equal(t, types.RoyaltyConfig{}, config)

	ctrl.Finish()
}
