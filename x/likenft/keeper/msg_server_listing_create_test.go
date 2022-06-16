package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestCreateListingNormal(t *testing.T) {
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
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	notUserAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}

	// Seed listing by previous owner
	prevOwnerListing := types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     notUserAddressBytes,
		Price:      uint64(987654),
		Expiration: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
	}
	k.SetListing(ctx, prevOwnerListing)
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: prevOwnerListing.Expiration,
		ListingKey: types.ListingKey(prevOwnerListing.ClassId, prevOwnerListing.NftId, prevOwnerListing.Seller),
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(userAddressBytes).MinTimes(1)

	// Call
	res, err := msgServer.CreateListing(goCtx, &types.MsgCreateListing{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.NoError(t, err)
	expectedListing := types.Listing{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     userAddress,
		Price:      price,
		Expiration: expiration,
	}
	require.Equal(t, &types.MsgCreateListingResponse{
		Listing: expectedListing,
	}, res)

	// Check state
	// expect new listing
	listing, found := k.GetListing(ctx, classId, nftId, sdk.AccAddress(userAddressBytes))
	require.True(t, found)
	require.Equal(t, expectedListing.ToStoreRecord(), listing)
	// expect enqueued listing
	_, found = k.GetListingExpireQueueEntry(ctx, expiration, types.ListingKey(classId, nftId, userAddressBytes))
	require.True(t, found)
	// expect previous listing to be deleted
	_, found = k.GetListing(ctx, prevOwnerListing.ClassId, prevOwnerListing.NftId, sdk.AccAddress(prevOwnerListing.Seller))
	require.False(t, found)
	_, found = k.GetListingExpireQueueEntry(ctx, prevOwnerListing.Expiration, types.ListingKey(prevOwnerListing.ClassId, prevOwnerListing.NftId, prevOwnerListing.Seller))
	require.False(t, found)

	ctrl.Finish()
}

// Test user not owner
func TestCreateListingUserNotOwner(t *testing.T) {
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
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	notUserAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(notUserAddressBytes).MinTimes(1)

	// Call
	res, err := msgServer.CreateListing(goCtx, &types.MsgCreateListing{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())
	ctrl.Finish()
}

// Test expiration out of range
func TestCreateListingExpirationInPast(t *testing.T) {
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
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(userAddressBytes).MinTimes(1)

	// Call
	res, err := msgServer.CreateListing(goCtx, &types.MsgCreateListing{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidRequest.Error())

	ctrl.Finish()
}

func TestCreateListingExpirationTooFarInFuture(t *testing.T) {
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
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(userAddressBytes).MinTimes(1)

	// Call
	res, err := msgServer.CreateListing(goCtx, &types.MsgCreateListing{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidRequest.Error())

	ctrl.Finish()
}

// Test listing already exists
func TestCreateListingAlreadyExists(t *testing.T) {
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
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)

	// Seed listing by current owner
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     userAddressBytes,
		Price:      uint64(987654),
		Expiration: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(userAddressBytes).MinTimes(1)

	// Call
	res, err := msgServer.CreateListing(goCtx, &types.MsgCreateListing{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrListingAlreadyExists.Error())

	ctrl.Finish()
}
