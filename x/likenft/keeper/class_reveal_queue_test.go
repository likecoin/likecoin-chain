package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNClassRevealQueue(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ClassRevealQueue {
	items := make([]types.ClassRevealQueue, n)
	for i := range items {
		items[i].RevealTime = strconv.Itoa(i)
		items[i].ClassId = strconv.Itoa(i)

		keeper.SetClassRevealQueue(ctx, items[i])
	}
	return items
}

func TestClassRevealQueueGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassRevealQueue(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetClassRevealQueue(ctx,
			item.RevealTime,
			item.ClassId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestClassRevealQueueRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassRevealQueue(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveClassRevealQueue(ctx,
			item.RevealTime,
			item.ClassId,
		)
		_, found := keeper.GetClassRevealQueue(ctx,
			item.RevealTime,
			item.ClassId,
		)
		require.False(t, found)
	}
}

func TestClassRevealQueueGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassRevealQueue(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllClassRevealQueue(ctx)),
	)
}
