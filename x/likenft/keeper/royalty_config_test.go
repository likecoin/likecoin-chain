package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/likecoin/likecoin-chain/v4/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v4/testutil/nullify"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNRoyaltyConfig(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.RoyaltyConfigByClass {
	items := make([]types.RoyaltyConfigByClass, n)
	for i := range items {
		items[i].ClassId = strconv.Itoa(i)

		keeper.SetRoyaltyConfig(ctx, items[i])
	}
	return items
}

func TestRoyaltyConfigGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNRoyaltyConfig(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetRoyaltyConfig(ctx,
			item.ClassId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item.RoyaltyConfig),
			nullify.Fill(&rst),
		)
	}
}
func TestRoyaltyRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNRoyaltyConfig(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveRoyaltyConfig(ctx,
			item.ClassId,
		)
		_, found := keeper.GetRoyaltyConfig(ctx,
			item.ClassId,
		)
		require.False(t, found)
	}
}

func TestRoyaltyGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNRoyaltyConfig(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllRoyaltyConfig(ctx)),
	)
}
