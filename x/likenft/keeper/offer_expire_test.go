package keeper_test

import (
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	apptestutil "github.com/likecoin/likechain/testutil"
	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// Test feature

func parseOfferExpireEvent(t *testing.T, ctx sdk.Context) *types.EventExpireOffer {
	for _, event := range ctx.EventManager().Events() {
		if event.Type != "likechain.likenft.EventExpireOffer" {
			continue
		}
		ev := types.EventExpireOffer{}
		for _, attr := range event.Attributes {
			val := strings.Trim(string(attr.Value), "\"")
			if string(attr.Key) == "class_id" {
				ev.ClassId = val
				continue
			}
			if string(attr.Key) == "nft_id" {
				ev.NftId = val
				continue
			}
			if string(attr.Key) == "buyer" {
				ev.Buyer = val
				continue
			}
			if string(attr.Key) == "success" {
				ev.Success = val == "true"
				continue
			}
			if string(attr.Key) == "error" {
				ev.Error = val
				continue
			}
		}
		return &ev
	}
	return nil
}

func TestExpireOfferFeature(t *testing.T) {
	app := apptestutil.SetupTestAppWithDefaultState()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// seed offer to be expired
	classId := "likenft11"
	nftId := "nft1"
	buyer := sdk.AccAddress([]byte{0, 1, 0, 1, 0, 1, 0, 1})
	expireTime := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	app.LikeNftKeeper.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyer,
		Price:      uint64(0), // refund to be covered by unit test
		Expiration: expireTime,
	})
	app.LikeNftKeeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expireTime,
		OfferKey:   types.OfferKey(classId, nftId, buyer),
	})

	// seed another with future expire time
	buyer2 := sdk.AccAddress([]byte{1, 1, 1, 1, 0, 0, 0, 0})
	expireTime2 := time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)
	app.LikeNftKeeper.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyer2,
		Price:      uint64(0), // refund to be covered by unit test
		Expiration: expireTime2,
	})
	app.LikeNftKeeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expireTime2,
		OfferKey:   types.OfferKey(classId, nftId, buyer2),
	})

	// increase height to after expiry time
	newHeader := ctx.BlockHeader()
	newHeader.Time = time.Date(2022, 2, 1, 0, 0, 1, 0, time.UTC)
	ctx = ctx.WithBlockHeader(newHeader)

	require.NotPanics(t, func() {
		likenft.EndBlocker(ctx, app.LikeNftKeeper)
	})

	// check result from event
	event := parseOfferExpireEvent(t, ctx)
	require.Equal(t, &types.EventExpireOffer{
		ClassId: classId,
		NftId:   nftId,
		Buyer:   buyer.String(),
		Success: true,
		Error:   "",
	}, event)

	// check offer deleted
	_, found := app.LikeNftKeeper.GetOffer(ctx, classId, nftId, buyer)
	require.False(t, found)
	_, found = app.LikeNftKeeper.GetOfferExpireQueueEntry(ctx, expireTime, types.OfferKey(classId, nftId, buyer))
	require.False(t, found)

	// check second offer not deleted
	_, found = app.LikeNftKeeper.GetOffer(ctx, classId, nftId, buyer2)
	require.True(t, found)
	_, found = app.LikeNftKeeper.GetOfferExpireQueueEntry(ctx, expireTime2, types.OfferKey(classId, nftId, buyer2))
	require.True(t, found)
}

func TestExpireOfferFeatureErrorOccured(t *testing.T) {
	app := apptestutil.SetupTestAppWithDefaultState()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// seed offer to be expired
	classId := "likenft11"
	nftId := "nft1"
	buyer := sdk.AccAddress([]byte{0, 1, 0, 1, 0, 1, 0, 1})
	expireTime := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	app.LikeNftKeeper.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyer,
		Price:      uint64(123456), // will failed due to no balance in module acc
		Expiration: expireTime,
	})
	app.LikeNftKeeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expireTime,
		OfferKey:   types.OfferKey(classId, nftId, buyer),
	})

	// increase height to after expiry time
	newHeader := ctx.BlockHeader()
	newHeader.Time = time.Date(2022, 2, 1, 0, 0, 1, 0, time.UTC)
	ctx = ctx.WithBlockHeader(newHeader)

	require.NotPanics(t, func() {
		likenft.EndBlocker(ctx, app.LikeNftKeeper)
	})

	// check result from event
	event := parseOfferExpireEvent(t, ctx)
	require.Equal(t, classId, event.ClassId)
	require.Equal(t, nftId, event.NftId)
	require.Equal(t, buyer.String(), event.Buyer)
	require.False(t, event.Success)
	require.Contains(t, event.Error, types.ErrFailedToExpireOffer.Error())

	// check offer not deleted
	_, found := app.LikeNftKeeper.GetOffer(ctx, classId, nftId, buyer)
	require.True(t, found)
	// expect queue entry deleted
	_, found = app.LikeNftKeeper.GetOfferExpireQueueEntry(ctx, expireTime, types.OfferKey(classId, nftId, buyer))
	require.False(t, found)
}

// Unit tests

func TestExpireOfferNormalRefund(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: time.Date(2022, 2, 1, 0, 0, 0, 1, time.UTC),
	})

	// seed offer to be expired
	classId := "likenft11"
	nftId := "nft1"
	buyer := sdk.AccAddress([]byte{0, 1, 0, 1, 0, 1, 0, 1})
	expireTime := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	offer := types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyer,
		Price:      uint64(123456),
		Expiration: expireTime,
	}
	keeper.SetOffer(ctx, offer)
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expireTime,
		OfferKey:   types.OfferKey(classId, nftId, buyer),
	})

	// Mock
	denom := keeper.PriceDenom(ctx)
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(int64(offer.Price))))
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, buyer, coins).Return(nil)

	// Call
	err := keeper.ExpireOffer(ctx, offer)
	require.NoError(t, err)

	// Check state
	// Expect offer deleted
	_, found := keeper.GetOffer(ctx, classId, nftId, buyer)
	require.False(t, found)
	// queue entry to be deleted by abci queue logic

	ctrl.Finish()
}

// test no refund
func TestExpireOfferNormalNoRefund(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: time.Date(2022, 2, 1, 0, 0, 0, 1, time.UTC),
	})

	// seed offer to be expired
	classId := "likenft11"
	nftId := "nft1"
	buyer := sdk.AccAddress([]byte{0, 1, 0, 1, 0, 1, 0, 1})
	expireTime := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	offer := types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyer,
		Price:      uint64(0),
		Expiration: expireTime,
	}
	keeper.SetOffer(ctx, offer)
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expireTime,
		OfferKey:   types.OfferKey(classId, nftId, buyer),
	})

	// Call
	err := keeper.ExpireOffer(ctx, offer)
	require.NoError(t, err)

	// Check state
	// Expect offer deleted
	_, found := keeper.GetOffer(ctx, classId, nftId, buyer)
	require.False(t, found)
	// queue entry to be deleted by abci queue logic

	ctrl.Finish()
}

// test not expired on record
func TestExpireOfferBadQueueEntry(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: time.Date(2022, 2, 1, 0, 0, 0, 1, time.UTC),
	})

	// seed offer to be expired
	classId := "likenft11"
	nftId := "nft1"
	buyer := sdk.AccAddress([]byte{0, 1, 0, 1, 0, 1, 0, 1})
	expireTime := time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)
	offer := types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyer,
		Price:      uint64(0),
		Expiration: expireTime,
	}
	keeper.SetOffer(ctx, offer)
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: expireTime,
		OfferKey:   types.OfferKey(classId, nftId, buyer),
	})
	// seed bad queue entry with earlier expire time
	badExpireTime := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: badExpireTime,
		OfferKey:   types.OfferKey(classId, nftId, buyer),
	})

	// Call
	err := keeper.ExpireOffer(ctx, offer)
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrFailedToExpireOffer.Error())

	// Check state
	// Expect offer not deleted
	_, found := keeper.GetOffer(ctx, classId, nftId, buyer)
	require.True(t, found)
	_, found = keeper.GetOfferExpireQueueEntry(ctx, expireTime, types.OfferKey(classId, nftId, buyer))
	require.True(t, found)

	ctrl.Finish()
}
