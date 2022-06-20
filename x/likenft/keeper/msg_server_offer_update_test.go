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

// normal raise price
func TestUpdateOfferNormalRaisePrice(t *testing.T) {
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
	newPrice := uint64(200000)
	newExpiration := time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC)

	// Seed existing offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddressBytes,
		Price:      price,
		Expiration: expiration,
	})
	k.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expiration,
		OfferKey:   types.OfferKey(classId, nftId, userAddressBytes),
	})

	// Mock
	priceDiff := newPrice - price
	bankKeeper.EXPECT().GetBalance(gomock.Any(), userAddressBytes, "nanolike").Return(sdk.NewCoin("nanolike", sdk.NewInt(int64(priceDiff))))
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, userAddressBytes, sdk.NewCoins(sdk.NewCoin("nanolike", sdk.NewInt(int64(price))))).Return(nil)
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), userAddressBytes, types.ModuleName, sdk.NewCoins(sdk.NewCoin("nanolike", sdk.NewInt(int64(newPrice))))).Return(nil)

	// Call
	res, err := msgServer.UpdateOffer(goCtx, &types.MsgUpdateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      newPrice,
		Expiration: newExpiration,
	})
	require.NoError(t, err)
	expectedOffer := types.Offer{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddress,
		Price:      newPrice,
		Expiration: newExpiration,
	}
	require.Equal(t, &types.MsgUpdateOfferResponse{
		Offer: expectedOffer,
	}, res)

	// check state
	// expect offer updated
	offer, found := k.GetOffer(ctx, classId, nftId, userAddressBytes)
	require.True(t, found)
	require.Equal(t, expectedOffer.ToStoreRecord(), offer)
	// expect offer expire queue entry updated
	_, found = k.GetOfferExpireQueueEntry(ctx, expiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.False(t, found)
	_, found = k.GetOfferExpireQueueEntry(ctx, newExpiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}

// normal drop price
func TestUpdateOfferNormalDropPrice(t *testing.T) {
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
	newPrice := uint64(100000)
	newExpiration := time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC)

	// Seed existing offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddressBytes,
		Price:      price,
		Expiration: expiration,
	})
	k.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expiration,
		OfferKey:   types.OfferKey(classId, nftId, userAddressBytes),
	})

	// Mock
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, userAddressBytes, sdk.NewCoins(sdk.NewCoin("nanolike", sdk.NewInt(int64(price))))).Return(nil)
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), userAddressBytes, types.ModuleName, sdk.NewCoins(sdk.NewCoin("nanolike", sdk.NewInt(int64(newPrice))))).Return(nil)

	// Call
	res, err := msgServer.UpdateOffer(goCtx, &types.MsgUpdateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      newPrice,
		Expiration: newExpiration,
	})
	require.NoError(t, err)
	expectedOffer := types.Offer{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddress,
		Price:      newPrice,
		Expiration: newExpiration,
	}
	require.Equal(t, &types.MsgUpdateOfferResponse{
		Offer: expectedOffer,
	}, res)

	// check state
	// expect offer updated
	offer, found := k.GetOffer(ctx, classId, nftId, userAddressBytes)
	require.True(t, found)
	require.Equal(t, expectedOffer.ToStoreRecord(), offer)
	// expect offer expire queue entry updated
	_, found = k.GetOfferExpireQueueEntry(ctx, expiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.False(t, found)
	_, found = k.GetOfferExpireQueueEntry(ctx, newExpiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}

// normal same price
func TestUpdateOfferNormalSamePrice(t *testing.T) {
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
	newPrice := uint64(123456)
	newExpiration := time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC)

	// Seed existing offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddressBytes,
		Price:      price,
		Expiration: expiration,
	})
	k.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expiration,
		OfferKey:   types.OfferKey(classId, nftId, userAddressBytes),
	})

	// Call
	res, err := msgServer.UpdateOffer(goCtx, &types.MsgUpdateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      newPrice,
		Expiration: newExpiration,
	})
	require.NoError(t, err)
	expectedOffer := types.Offer{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddress,
		Price:      newPrice,
		Expiration: newExpiration,
	}
	require.Equal(t, &types.MsgUpdateOfferResponse{
		Offer: expectedOffer,
	}, res)

	// check state
	// expect offer updated
	offer, found := k.GetOffer(ctx, classId, nftId, userAddressBytes)
	require.True(t, found)
	require.Equal(t, expectedOffer.ToStoreRecord(), offer)
	// expect offer expire queue entry updated
	_, found = k.GetOfferExpireQueueEntry(ctx, expiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.False(t, found)
	_, found = k.GetOfferExpireQueueEntry(ctx, newExpiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}

// not exists
func TestUpdateOfferNotExists(t *testing.T) {
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
	newPrice := uint64(123456)
	newExpiration := time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC)

	// Call
	res, err := msgServer.UpdateOffer(goCtx, &types.MsgUpdateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      newPrice,
		Expiration: newExpiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrOfferNotFound.Error())

	ctrl.Finish()
}

// expiration out of range
func TestUpdateOfferExpirationInPast(t *testing.T) {
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
	newPrice := uint64(123456)
	newExpiration := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

	// Seed existing offer
	offer := types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddressBytes,
		Price:      price,
		Expiration: expiration,
	}
	k.SetOffer(ctx, offer)
	k.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expiration,
		OfferKey:   types.OfferKey(classId, nftId, userAddressBytes),
	})

	// Call
	res, err := msgServer.UpdateOffer(goCtx, &types.MsgUpdateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      newPrice,
		Expiration: newExpiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidRequest.Error())

	// check state
	// expect offer not updated
	_offer, found := k.GetOffer(ctx, classId, nftId, userAddressBytes)
	require.True(t, found)
	require.Equal(t, _offer, offer)
	// expect offer expire queue entry not updated
	_, found = k.GetOfferExpireQueueEntry(ctx, expiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}

func TestUpdateOfferExpirationTooFarInFuture(t *testing.T) {
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
	newPrice := uint64(123456)
	newExpiration := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)

	// Seed existing offer
	offer := types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddressBytes,
		Price:      price,
		Expiration: expiration,
	}
	k.SetOffer(ctx, offer)
	k.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expiration,
		OfferKey:   types.OfferKey(classId, nftId, userAddressBytes),
	})

	// Call
	res, err := msgServer.UpdateOffer(goCtx, &types.MsgUpdateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      newPrice,
		Expiration: newExpiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidRequest.Error())

	// check state
	// expect offer not updated
	_offer, found := k.GetOffer(ctx, classId, nftId, userAddressBytes)
	require.True(t, found)
	require.Equal(t, _offer, offer)
	// expect offer expire queue entry not updated
	_, found = k.GetOfferExpireQueueEntry(ctx, expiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.True(t, found)

	ctrl.Finish()
}
