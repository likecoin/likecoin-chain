package keeper_test

import (
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	apptestutil "github.com/likecoin/likechain/testutil"
	"github.com/likecoin/likechain/x/likenft"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// Test feature

func parseListingExpireEvent(t *testing.T, ctx sdk.Context) *types.EventExpireListing {
	for _, event := range ctx.EventManager().Events() {
		if event.Type != "likechain.likenft.v1.EventExpireListing" {
			continue
		}
		ev := types.EventExpireListing{}
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
			if string(attr.Key) == "seller" {
				ev.Seller = val
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

func TestListingExpireFeatureNormal(t *testing.T) {
	app := apptestutil.SetupTestAppWithDefaultState()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// seed listing to be expired
	classId := "likenft11"
	nftId := "nft1"
	seller := sdk.AccAddress([]byte{0, 1, 0, 1, 0, 1, 0, 1})
	expireTime := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	app.LikeNftKeeper.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     seller,
		Price:      uint64(123456),
		Expiration: expireTime,
	})
	app.LikeNftKeeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: expireTime,
		ListingKey: types.ListingKey(classId, nftId, seller),
	})

	// seed another with future expire time
	seller2 := sdk.AccAddress([]byte{1, 1, 1, 1, 0, 0, 0, 0})
	expireTime2 := time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)
	app.LikeNftKeeper.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     seller2,
		Price:      uint64(987654),
		Expiration: expireTime2,
	})
	app.LikeNftKeeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: expireTime2,
		ListingKey: types.ListingKey(classId, nftId, seller2),
	})

	// increase height to after expiry time
	newHeader := ctx.BlockHeader()
	newHeader.Time = time.Date(2022, 2, 1, 0, 0, 1, 0, time.UTC)
	ctx = ctx.WithBlockHeader(newHeader)

	require.NotPanics(t, func() {
		likenft.EndBlocker(ctx, app.LikeNftKeeper)
	})

	// check result from event
	event := parseListingExpireEvent(t, ctx)
	require.Equal(t, &types.EventExpireListing{
		ClassId: classId,
		NftId:   nftId,
		Seller:  seller.String(),
		Success: true,
		Error:   "",
	}, event)

	// check listing deleted
	_, found := app.LikeNftKeeper.GetListing(ctx, classId, nftId, seller)
	require.False(t, found)
	_, found = app.LikeNftKeeper.GetListingExpireQueueEntry(ctx, expireTime, types.ListingKey(classId, nftId, seller))
	require.False(t, found)

	// check second listing not deleted
	_, found = app.LikeNftKeeper.GetListing(ctx, classId, nftId, seller2)
	require.True(t, found)
	_, found = app.LikeNftKeeper.GetListingExpireQueueEntry(ctx, expireTime2, types.ListingKey(classId, nftId, seller2))
	require.True(t, found)
}

func TestListingExpireFeatureQueueMismatch(t *testing.T) {
	app := apptestutil.SetupTestAppWithDefaultState()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// seed listing
	classId := "likenft11"
	nftId := "nft1"
	seller := sdk.AccAddress([]byte{0, 1, 0, 1, 0, 1, 0, 1})
	expireTime := time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)
	app.LikeNftKeeper.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     seller,
		Price:      uint64(123456),
		Expiration: expireTime,
	})
	app.LikeNftKeeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: expireTime,
		ListingKey: types.ListingKey(classId, nftId, seller),
	})

	// seed wrong, earlier queue entry
	badExpireTime := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
	app.LikeNftKeeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: badExpireTime,
		ListingKey: types.ListingKey(classId, nftId, seller),
	})

	// increase height to after wrong expiry time
	newHeader := ctx.BlockHeader()
	newHeader.Time = time.Date(2022, 2, 1, 0, 0, 1, 0, time.UTC)
	ctx = ctx.WithBlockHeader(newHeader)

	require.NotPanics(t, func() {
		likenft.EndBlocker(ctx, app.LikeNftKeeper)
	})

	// check result from event
	event := parseListingExpireEvent(t, ctx)
	require.Equal(t, classId, event.ClassId)
	require.Equal(t, nftId, event.NftId)
	require.Equal(t, seller.String(), event.Seller)
	require.False(t, event.Success)
	require.Contains(t, event.Error, types.ErrFailedToExpireListing.Error())

	// check listing not deleted
	_, found := app.LikeNftKeeper.GetListing(ctx, classId, nftId, seller)
	require.True(t, found)
	_, found = app.LikeNftKeeper.GetListingExpireQueueEntry(ctx, expireTime, types.ListingKey(classId, nftId, seller))
	require.True(t, found)

	// check wrong queue entry deleted
	_, found = app.LikeNftKeeper.GetListingExpireQueueEntry(ctx, badExpireTime, types.ListingKey(classId, nftId, seller))
	require.False(t, found)
}
