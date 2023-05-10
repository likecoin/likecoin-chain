package keeper_test

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/cosmos/cosmos-sdk/x/nft"
	keepertest "github.com/likecoin/likecoin-chain/v4/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestAccountByClassNormal(t *testing.T) {
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
	goCtx := sdk.WrapSDKContext(ctx)

	// Test input
	classId := "likenft1aabbccddeeff"
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	// Run
	res, err := keeper.AccountByClass(goCtx, &types.QueryAccountByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, ownerAddress, res.Address)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestAccountByClassNotFound(t *testing.T) {
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
	goCtx := sdk.WrapSDKContext(ctx)

	// Test input
	classId := "likenft1aabbccddeeff"

	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{}, false)

	// Run
	res, err := keeper.AccountByClass(goCtx, &types.QueryAccountByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotFound.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestAccountByClassMissingRelationRecord(t *testing.T) {
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
	goCtx := sdk.WrapSDKContext(ctx)

	// Test input
	classId := "likenft1aabbccddeeff"
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

	// Run
	res, err := keeper.AccountByClass(goCtx, &types.QueryAccountByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyAccount.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestISCNByClassNotRelatedToAccount(t *testing.T) {
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
	goCtx := sdk.WrapSDKContext(ctx)

	// Test input
	classId := "likenft1aabbccddeeff"
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{"likenft199887766"},
	})

	// Run
	res, err := keeper.AccountByClass(goCtx, &types.QueryAccountByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyAccount.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}
