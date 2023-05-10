package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likecoin-chain/v4/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestDeleteListingNormalOwner(t *testing.T) {
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
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	notUserAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}

	// Seed listing
	prevUserListing := types.ListingStoreRecord{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           userAddressBytes,
		Price:            uint64(987654),
		Expiration:       time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
		FullPayToRoyalty: false,
	}
	k.SetListing(ctx, prevUserListing)
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: prevUserListing.Expiration,
		ListingKey: types.ListingKey(classId, nftId, userAddressBytes),
	})
	prevNotUserListing := types.ListingStoreRecord{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           notUserAddressBytes,
		Price:            uint64(987654),
		Expiration:       time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
		FullPayToRoyalty: false,
	}
	k.SetListing(ctx, prevNotUserListing)
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: prevNotUserListing.Expiration,
		ListingKey: types.ListingKey(classId, nftId, notUserAddressBytes),
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(userAddressBytes).MinTimes(1)

	// Call
	res, err := msgServer.DeleteListing(goCtx, &types.MsgDeleteListing{
		Creator: userAddress,
		ClassId: classId,
		NftId:   nftId,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgDeleteListingResponse{}, res)

	// Check state
	// expect both listing deleted
	_, found := k.GetListing(ctx, classId, nftId, sdk.AccAddress(userAddressBytes))
	require.False(t, found)
	_, found = k.GetListing(ctx, classId, nftId, sdk.AccAddress(notUserAddressBytes))
	require.False(t, found)
	// expect both listing dequeued
	_, found = k.GetListingExpireQueueEntry(ctx, prevUserListing.Expiration, types.ListingKey(classId, nftId, userAddressBytes))
	require.False(t, found)
	_, found = k.GetListingExpireQueueEntry(ctx, prevNotUserListing.Expiration, types.ListingKey(classId, nftId, notUserAddressBytes))
	require.False(t, found)

	ctrl.Finish()
}

func TestDeleteListingNormalNotOwner(t *testing.T) {
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
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	notUserAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}

	// Seed listing
	prevUserListing := types.ListingStoreRecord{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           userAddressBytes,
		Price:            uint64(987654),
		Expiration:       time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
		FullPayToRoyalty: false,
	}
	k.SetListing(ctx, prevUserListing)
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: prevUserListing.Expiration,
		ListingKey: types.ListingKey(classId, nftId, userAddressBytes),
	})
	prevNotUserListing := types.ListingStoreRecord{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           notUserAddressBytes,
		Price:            uint64(987654),
		Expiration:       time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
		FullPayToRoyalty: false,
	}
	k.SetListing(ctx, prevNotUserListing)
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: prevNotUserListing.Expiration,
		ListingKey: types.ListingKey(classId, nftId, notUserAddressBytes),
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(notUserAddressBytes).MinTimes(1)

	// Call
	res, err := msgServer.DeleteListing(goCtx, &types.MsgDeleteListing{
		Creator: userAddress,
		ClassId: classId,
		NftId:   nftId,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgDeleteListingResponse{}, res)

	// Check state
	// expect own listing deleted
	_, found := k.GetListing(ctx, classId, nftId, sdk.AccAddress(userAddressBytes))
	require.False(t, found)
	_, found = k.GetListing(ctx, classId, nftId, sdk.AccAddress(notUserAddressBytes))
	require.True(t, found)
	// expect own listing dequeued
	_, found = k.GetListingExpireQueueEntry(ctx, prevUserListing.Expiration, types.ListingKey(classId, nftId, userAddressBytes))
	require.False(t, found)
	_, found = k.GetListingExpireQueueEntry(ctx, prevNotUserListing.Expiration, types.ListingKey(classId, nftId, notUserAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}

func TestDeleteListingNormalNotOwnerNotFound(t *testing.T) {
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
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	notUserAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}

	// Seed listing
	prevNotUserListing := types.ListingStoreRecord{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           notUserAddressBytes,
		Price:            uint64(987654),
		Expiration:       time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
		FullPayToRoyalty: false,
	}
	k.SetListing(ctx, prevNotUserListing)
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: prevNotUserListing.Expiration,
		ListingKey: types.ListingKey(classId, nftId, notUserAddressBytes),
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(notUserAddressBytes).MinTimes(1)

	// Call
	res, err := msgServer.DeleteListing(goCtx, &types.MsgDeleteListing{
		Creator: userAddress,
		ClassId: classId,
		NftId:   nftId,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrListingNotFound.Error())

	// Check state
	// expect no listing deleted
	_, found := k.GetListing(ctx, classId, nftId, sdk.AccAddress(notUserAddressBytes))
	require.True(t, found)
	// expect no listing dequeued
	_, found = k.GetListingExpireQueueEntry(ctx, prevNotUserListing.Expiration, types.ListingKey(classId, nftId, notUserAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}
