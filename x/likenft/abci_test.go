package likenft_test

import (
	"testing"
	"time"

	"github.com/likecoin/likechain/testutil"
	"github.com/likecoin/likechain/x/likenft"
	"github.com/likecoin/likechain/x/likenft/types"
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
	pendingQueue := app.LikeNftKeeper.ClassRevealQueueByTimeIterator(ctx, header.Time)
	require.False(t, pendingQueue.Valid())
	pendingQueue.Close()

	// Seed queue entries, expect to cause class not found / bad class id error
	entry0 := types.ClassRevealQueueEntry{
		RevealTime: header.Time,
		ClassId:    "1",
	}
	entry1 := types.ClassRevealQueueEntry{
		RevealTime: header.Time.Add(time.Duration(1) * time.Second),
		ClassId:    "2",
	}
	entry2 := types.ClassRevealQueueEntry{
		RevealTime: header.Time.Add(time.Duration(2) * time.Second),
		ClassId:    "3",
	}
	app.LikeNftKeeper.SetClassRevealQueueEntry(ctx, entry0)
	app.LikeNftKeeper.SetClassRevealQueueEntry(ctx, entry1)
	app.LikeNftKeeper.SetClassRevealQueueEntry(ctx, entry2)

	// Double check queue is still not valid after insert since no height increase
	pendingQueue = app.LikeNftKeeper.ClassRevealQueueByTimeIterator(ctx, header.Time)
	require.False(t, pendingQueue.Valid())
	pendingQueue.Close()
	pendingEntries := app.LikeNftKeeper.GetClassRevealQueueByTime(ctx, header.Time)
	require.Empty(t, pendingEntries)

	// Increase height
	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	// Check pending entries after height increase
	pendingQueue = app.LikeNftKeeper.ClassRevealQueueByTimeIterator(ctx, newHeader.Time)
	require.True(t, pendingQueue.Valid())
	pendingQueue.Close()
	pendingEntries = app.LikeNftKeeper.GetClassRevealQueueByTime(ctx, newHeader.Time)
	require.Equal(t, []types.ClassRevealQueueEntry{entry0}, pendingEntries)

	// Process queue, should process entry 0 only
	require.NotPanics(t, func() {
		likenft.EndBlocker(ctx, app.LikeNftKeeper)
	})

	// Check if queue has been digested
	pendingQueue = app.LikeNftKeeper.ClassRevealQueueByTimeIterator(ctx, newHeader.Time)
	require.False(t, pendingQueue.Valid())
	pendingQueue.Close()
	pendingEntries = app.LikeNftKeeper.GetClassRevealQueueByTime(ctx, newHeader.Time)
	require.Empty(t, pendingEntries)

	// Check future entries still on queue after digest
	remainingEntries := app.LikeNftKeeper.GetClassRevealQueue(ctx)
	require.Equal(t, []types.ClassRevealQueueEntry{entry1, entry2}, remainingEntries)
}
