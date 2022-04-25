package likenft_test

import (
	"testing"
	"time"

	"github.com/likecoin/likechain/testutil"
	"github.com/likecoin/likechain/x/likenft"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestQueueProcess(t *testing.T) {
	app := testutil.SetupTestAppWithDefaultState()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// Check queue is not a valid one without any entries at height 0
	revealQueue := app.LikeNftKeeper.ClassRevealQueueIterator(ctx)
	require.False(t, revealQueue.Valid())
	revealQueue.Close()

	// TODO: Test queue insert via class create/update

	// Double check queue is still not valid after insert since no height increase
	revealQueue = app.LikeNftKeeper.ClassRevealQueueIterator(ctx)
	require.False(t, revealQueue.Valid())
	revealQueue.Close()

	// Increase height
	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	// Check queue is valid after height increase
	// TODO: Uncomment when queue insert above is implemented

	// revealQueue = app.LikeNftKeeper.ClassRevealQueueIterator(ctx)
	// require.True(t, revealQueue.Valid())
	// revealQueue.Close()

	// Process queue
	likenft.EndBlocker(ctx, app.LikeNftKeeper)

	// Check if queue has been digested
	revealQueue = app.LikeNftKeeper.ClassRevealQueueIterator(ctx)
	require.False(t, revealQueue.Valid())
	revealQueue.Close()
}
