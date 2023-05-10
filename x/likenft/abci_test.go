package likenft_test

import (
	"testing"
	"time"

	"github.com/likecoin/likecoin-chain/v4/testutil"
	"github.com/likecoin/likecoin-chain/v4/x/likenft"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestClassRevealQueueProcessing(t *testing.T) {
	app := testutil.SetupTestAppWithDefaultState()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// Check queue is not a valid one without any entries at height 0
	pendingClassQueue := app.LikeNftKeeper.ClassRevealQueueByTimeIterator(ctx, header.Time)
	require.False(t, pendingClassQueue.Valid())
	pendingClassQueue.Close()

	pendingOfferQueue := app.LikeNftKeeper.OfferExpireQueueByTimeIterator(ctx, header.Time)
	require.False(t, pendingOfferQueue.Valid())
	pendingOfferQueue.Close()

	pendingListingQueue := app.LikeNftKeeper.ListingExpireQueueByTimeIterator(ctx, header.Time)
	require.False(t, pendingListingQueue.Valid())
	pendingListingQueue.Close()

	// Seed queue entries, expect to cause class not found / bad class id error
	classEntry0 := types.ClassRevealQueueEntry{
		RevealTime: header.Time,
		ClassId:    "1",
	}
	classEntry1 := types.ClassRevealQueueEntry{
		RevealTime: header.Time.Add(time.Duration(1) * time.Second),
		ClassId:    "2",
	}
	classEntry2 := types.ClassRevealQueueEntry{
		RevealTime: header.Time.Add(time.Duration(2) * time.Second),
		ClassId:    "3",
	}
	app.LikeNftKeeper.SetClassRevealQueueEntry(ctx, classEntry0)
	app.LikeNftKeeper.SetClassRevealQueueEntry(ctx, classEntry1)
	app.LikeNftKeeper.SetClassRevealQueueEntry(ctx, classEntry2)

	offerEntry0 := types.OfferExpireQueueEntry{
		ExpireTime: header.Time,
		OfferKey:   []byte("1"),
	}
	offerEntry1 := types.OfferExpireQueueEntry{
		ExpireTime: header.Time.Add(time.Duration(1) * time.Second),
		OfferKey:   []byte("2"),
	}
	offerEntry2 := types.OfferExpireQueueEntry{
		ExpireTime: header.Time.Add(time.Duration(2) * time.Second),
		OfferKey:   []byte("3"),
	}
	app.LikeNftKeeper.SetOfferExpireQueueEntry(ctx, offerEntry0)
	app.LikeNftKeeper.SetOfferExpireQueueEntry(ctx, offerEntry1)
	app.LikeNftKeeper.SetOfferExpireQueueEntry(ctx, offerEntry2)

	listingEntry0 := types.ListingExpireQueueEntry{
		ExpireTime: header.Time,
		ListingKey: []byte("1"),
	}
	listingEntry1 := types.ListingExpireQueueEntry{
		ExpireTime: header.Time.Add(time.Duration(1) * time.Second),
		ListingKey: []byte("2"),
	}
	listingEntry2 := types.ListingExpireQueueEntry{
		ExpireTime: header.Time.Add(time.Duration(2) * time.Second),
		ListingKey: []byte("3"),
	}
	app.LikeNftKeeper.SetListingExpireQueueEntry(ctx, listingEntry0)
	app.LikeNftKeeper.SetListingExpireQueueEntry(ctx, listingEntry1)
	app.LikeNftKeeper.SetListingExpireQueueEntry(ctx, listingEntry2)

	// Double check queue is still not valid after insert since no height increase
	pendingClassQueue = app.LikeNftKeeper.ClassRevealQueueByTimeIterator(ctx, header.Time)
	require.False(t, pendingClassQueue.Valid())
	pendingClassQueue.Close()
	pendingClassEntries := app.LikeNftKeeper.GetClassRevealQueueByTime(ctx, header.Time)
	require.Empty(t, pendingClassEntries)

	pendingOfferQueue = app.LikeNftKeeper.OfferExpireQueueByTimeIterator(ctx, header.Time)
	require.False(t, pendingOfferQueue.Valid())
	pendingOfferQueue.Close()
	pendingOfferEntries := app.LikeNftKeeper.GetOfferExpireQueueByTime(ctx, header.Time)
	require.Empty(t, pendingOfferEntries)

	pendingListingQueue = app.LikeNftKeeper.ListingExpireQueueByTimeIterator(ctx, header.Time)
	require.False(t, pendingListingQueue.Valid())
	pendingListingQueue.Close()
	pendingListingEntries := app.LikeNftKeeper.GetListingExpireQueueByTime(ctx, header.Time)
	require.Empty(t, pendingListingEntries)

	// Increase height
	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	// Check pending entries after height increase
	pendingClassQueue = app.LikeNftKeeper.ClassRevealQueueByTimeIterator(ctx, newHeader.Time)
	require.True(t, pendingClassQueue.Valid())
	pendingClassQueue.Close()
	pendingClassEntries = app.LikeNftKeeper.GetClassRevealQueueByTime(ctx, newHeader.Time)
	require.Equal(t, []types.ClassRevealQueueEntry{classEntry0}, pendingClassEntries)

	pendingOfferQueue = app.LikeNftKeeper.OfferExpireQueueByTimeIterator(ctx, newHeader.Time)
	require.True(t, pendingOfferQueue.Valid())
	pendingOfferQueue.Close()
	pendingOfferEntries = app.LikeNftKeeper.GetOfferExpireQueueByTime(ctx, newHeader.Time)
	require.Equal(t, []types.OfferExpireQueueEntry{offerEntry0}, pendingOfferEntries)

	pendingListingQueue = app.LikeNftKeeper.ListingExpireQueueByTimeIterator(ctx, newHeader.Time)
	require.True(t, pendingListingQueue.Valid())
	pendingListingQueue.Close()
	pendingListingEntries = app.LikeNftKeeper.GetListingExpireQueueByTime(ctx, newHeader.Time)
	require.Equal(t, []types.ListingExpireQueueEntry{listingEntry0}, pendingListingEntries)

	// Process queue, should process entry 0 only
	require.NotPanics(t, func() {
		likenft.EndBlocker(ctx, app.LikeNftKeeper)
	})

	// Check if queue has been digested
	pendingClassQueue = app.LikeNftKeeper.ClassRevealQueueByTimeIterator(ctx, newHeader.Time)
	require.False(t, pendingClassQueue.Valid())
	pendingClassQueue.Close()
	pendingClassEntries = app.LikeNftKeeper.GetClassRevealQueueByTime(ctx, newHeader.Time)
	require.Empty(t, pendingClassEntries)

	pendingOfferQueue = app.LikeNftKeeper.OfferExpireQueueByTimeIterator(ctx, newHeader.Time)
	require.False(t, pendingOfferQueue.Valid())
	pendingOfferQueue.Close()
	pendingOfferEntries = app.LikeNftKeeper.GetOfferExpireQueueByTime(ctx, newHeader.Time)
	require.Empty(t, pendingOfferEntries)

	pendingListingQueue = app.LikeNftKeeper.ListingExpireQueueByTimeIterator(ctx, newHeader.Time)
	require.False(t, pendingListingQueue.Valid())
	pendingListingQueue.Close()
	pendingListingEntries = app.LikeNftKeeper.GetListingExpireQueueByTime(ctx, newHeader.Time)
	require.Empty(t, pendingListingEntries)

	// Check future entries still on queue after digest
	remainingClassEntries := app.LikeNftKeeper.GetClassRevealQueue(ctx)
	require.Equal(t, []types.ClassRevealQueueEntry{classEntry1, classEntry2}, remainingClassEntries)

	remainingOfferEntries := app.LikeNftKeeper.GetOfferExpireQueue(ctx)
	require.Equal(t, []types.OfferExpireQueueEntry{offerEntry1, offerEntry2}, remainingOfferEntries)

	remainingListingEntries := app.LikeNftKeeper.GetListingExpireQueue(ctx)
	require.Equal(t, []types.ListingExpireQueueEntry{listingEntry1, listingEntry2}, remainingListingEntries)
}
