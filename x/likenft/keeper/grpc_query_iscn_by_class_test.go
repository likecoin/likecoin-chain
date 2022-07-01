package keeper_test

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likecoin-chain/v3/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	keepertest "github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	iscntypes "github.com/likecoin/likecoin-chain/v3/x/iscn/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestISCNByClassNormal(t *testing.T) {
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
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
	iscnSequence := uint64(1)
	iscnCidBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	iscnData := iscntypes.IscnInput(`{"qwer": "asdf"}`)
	iscnLatestVersion := uint64(1)
	iscnStoreRecord := iscntypes.StoreRecord{
		IscnId:   iscnId,
		CidBytes: iscnCidBytes,
		Data:     iscnData,
	}

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
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

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(&iscntypes.ContentIdRecord{
			OwnerAddressBytes: ownerAddressBytes,
			LatestVersion:     iscnLatestVersion,
		})

	iscnKeeper.
		EXPECT().
		GetIscnIdSequence(gomock.Any(), gomock.Eq(iscnId)).
		Return(iscnSequence)

	iscnKeeper.
		EXPECT().
		GetStoreRecord(gomock.Any(), gomock.Eq(iscnSequence)).
		Return(
			&iscnStoreRecord,
		)

	// Run
	res, err := keeper.ISCNByClass(goCtx, &types.QueryISCNByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, iscnId.Prefix.String(), res.IscnIdPrefix)
	require.Equal(t, ownerAddress, res.Owner)
	require.Equal(t, iscnLatestVersion, res.LatestVersion)
	require.Equal(t, iscnStoreRecord.Cid().String(), res.LatestRecord.Ipld)
	require.Equal(t, iscnStoreRecord.Data, res.LatestRecord.Data)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestISCNByClassNotFound(t *testing.T) {
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
	res, err := keeper.ISCNByClass(goCtx, &types.QueryISCNByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotFound.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestISCNByClassMissingRelationRecord(t *testing.T) {
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
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
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
	res, err := keeper.ISCNByClass(goCtx, &types.QueryISCNByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyIscn.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestISCNByClassNotRelatedToISCN(t *testing.T) {
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
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
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

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{"likenft199887766"},
	})

	// Run
	res, err := keeper.ISCNByClass(goCtx, &types.QueryISCNByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyIscn.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestISCNByClassISCNNotFound(t *testing.T) {
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
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)

	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
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

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(nil)

	// Run
	res, err := keeper.ISCNByClass(goCtx, &types.QueryISCNByClassRequest{
		ClassId: classId,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrIscnRecordNotFound.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}
