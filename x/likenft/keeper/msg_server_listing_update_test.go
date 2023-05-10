package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likecoin-chain/v4/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateListingNormal(t *testing.T) {
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
	fullPayToRoyalty := true

	// Seed listing
	prevExpiration := time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           userAddressBytes,
		Price:            uint64(987654),
		Expiration:       prevExpiration,
		FullPayToRoyalty: false,
	})
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: prevExpiration,
		ListingKey: types.ListingKey(classId, nftId, userAddressBytes),
	})

	// Call
	res, err := msgServer.UpdateListing(goCtx, &types.MsgUpdateListing{
		Creator:          userAddress,
		ClassId:          classId,
		NftId:            nftId,
		Price:            price,
		Expiration:       expiration,
		FullPayToRoyalty: fullPayToRoyalty,
	})
	require.NoError(t, err)
	expectedListing := types.Listing{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           userAddress,
		Price:            price,
		Expiration:       expiration,
		FullPayToRoyalty: fullPayToRoyalty,
	}
	require.Equal(t, &types.MsgUpdateListingResponse{
		Listing: expectedListing,
	}, res)

	// Check state
	// expect listing updated
	listing, found := k.GetListing(ctx, classId, nftId, sdk.AccAddress(userAddressBytes))
	require.True(t, found)
	require.Equal(t, expectedListing.ToStoreRecord(), listing)
	// expect queue updated
	_, found = k.GetListingExpireQueueEntry(ctx, prevExpiration, types.ListingKey(classId, nftId, userAddressBytes))
	require.False(t, found)
	_, found = k.GetListingExpireQueueEntry(ctx, expiration, types.ListingKey(classId, nftId, userAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}

// Test old listing not exists
func TestUpdateListingNotExist(t *testing.T) {
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
	fullPayToRoyalty := true

	// Call
	res, err := msgServer.UpdateListing(goCtx, &types.MsgUpdateListing{
		Creator:          userAddress,
		ClassId:          classId,
		NftId:            nftId,
		Price:            price,
		Expiration:       expiration,
		FullPayToRoyalty: fullPayToRoyalty,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrListingNotFound.Error())

	ctrl.Finish()
}

// Test expiration out of range
func TestUpdateListingExpirationInPast(t *testing.T) {
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
	expiration := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)
	fullPayToRoyalty := true

	// Seed listing
	prevListing := types.ListingStoreRecord{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           userAddressBytes,
		Price:            uint64(987654),
		Expiration:       time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
		FullPayToRoyalty: false,
	}
	k.SetListing(ctx, prevListing)
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: prevListing.Expiration,
		ListingKey: types.ListingKey(classId, nftId, userAddressBytes),
	})

	// Call
	res, err := msgServer.UpdateListing(goCtx, &types.MsgUpdateListing{
		Creator:          userAddress,
		ClassId:          classId,
		NftId:            nftId,
		Price:            price,
		Expiration:       expiration,
		FullPayToRoyalty: fullPayToRoyalty,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidRequest.Error())

	// Check state
	// expect listing not updated
	listing, found := k.GetListing(ctx, classId, nftId, sdk.AccAddress(userAddressBytes))
	require.True(t, found)
	require.Equal(t, prevListing, listing)
	// expect queue not updated
	_, found = k.GetListingExpireQueueEntry(ctx, prevListing.Expiration, types.ListingKey(classId, nftId, userAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}
