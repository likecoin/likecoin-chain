package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// normal refund
func TestDeleteOfferNormalRefund(t *testing.T) {
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

	// Call
	res, err := msgServer.DeleteOffer(goCtx, &types.MsgDeleteOffer{
		Creator: userAddress,
		ClassId: classId,
		NftId:   nftId,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgDeleteOfferResponse{}, res)

	// check state
	// expect offer deleted
	_, found := k.GetOffer(ctx, classId, nftId, userAddressBytes)
	require.False(t, found)
	// expect offer expire queue entry deleted
	_, found = k.GetOfferExpireQueueEntry(ctx, expiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.False(t, found)

	ctrl.Finish()
}

// normal no refund
func TestDeleteOfferNormalNoRefund(t *testing.T) {
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
	price := uint64(0)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)

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
	res, err := msgServer.DeleteOffer(goCtx, &types.MsgDeleteOffer{
		Creator: userAddress,
		ClassId: classId,
		NftId:   nftId,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgDeleteOfferResponse{}, res)

	// check state
	// expect offer deleted
	_, found := k.GetOffer(ctx, classId, nftId, userAddressBytes)
	require.False(t, found)
	// expect offer expire queue entry deleted
	_, found = k.GetOfferExpireQueueEntry(ctx, expiration, types.OfferKey(classId, nftId, userAddressBytes))
	require.False(t, found)

	ctrl.Finish()
}

// not exist
func TestDeleteOfferNotExist(t *testing.T) {
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

	// Call
	res, err := msgServer.DeleteOffer(goCtx, &types.MsgDeleteOffer{
		Creator: userAddress,
		ClassId: classId,
		NftId:   nftId,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrOfferNotFound.Error())

	ctrl.Finish()
}
