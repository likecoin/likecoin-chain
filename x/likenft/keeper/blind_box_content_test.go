package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNBlindBoxContent(keeper *keeper.Keeper, ctx sdk.Context, nClass int, nNFT int) []types.BlindBoxContent {
	var items []types.BlindBoxContent
	for i := 0; i < nClass; i++ {
		for j := 0; j < nNFT; j++ {
			item := types.BlindBoxContent{
				ClassId: strconv.Itoa(i),
				Id:      strconv.Itoa(j),
			}
			items = append(items, item)
			keeper.SetBlindBoxContent(ctx, item)
		}
	}
	return items
}

func TestBlindBoxContentGet(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	items := createNBlindBoxContent(keeper, ctx, 3, 3)
	for _, item := range items {
		rst, found := keeper.GetBlindBoxContent(ctx,
			item.ClassId,
			item.Id,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestBlindBoxContentRemove(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	items := createNBlindBoxContent(keeper, ctx, 3, 3)
	for _, item := range items {
		keeper.RemoveBlindBoxContent(ctx,
			item.ClassId,
			item.Id,
		)
		_, found := keeper.GetBlindBoxContent(ctx,
			item.ClassId,
			item.Id,
		)
		require.False(t, found)
	}
}

func TestBlindBoxContentRemoveMultiple(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	items := createNBlindBoxContent(keeper, ctx, 1, 5)
	keeper.RemoveBlindBoxContents(ctx,
		items[0].ClassId,
	)
	foundItems := keeper.GetBlindBoxContents(ctx,
		items[0].ClassId,
	)
	require.Empty(t, foundItems)
}

func TestBlindBoxContentGetByClass(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	items := createNBlindBoxContent(keeper, ctx, 3, 3)
	require.ElementsMatch(t,
		nullify.Fill(items[0:3]),
		nullify.Fill(keeper.GetBlindBoxContents(ctx, "0")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[3:6]),
		nullify.Fill(keeper.GetBlindBoxContents(ctx, "1")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[6:9]),
		nullify.Fill(keeper.GetBlindBoxContents(ctx, "2")),
	)
}

func TestBlindBoxContentGetAll(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	items := createNBlindBoxContent(keeper, ctx, 3, 3)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllBlindBoxContent(ctx)),
	)
}
