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

func createNClassesByAccount(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ClassesByAccount {
	items := make([]types.ClassesByAccount, n)
	for i := range items {
		items[i].Account = strconv.Itoa(i)

		keeper.SetClassesByAccount(ctx, items[i])
	}
	return items
}

func TestClassesByAccountGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassesByAccount(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetClassesByAccount(ctx,
			item.Account,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestClassesByAccountRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassesByAccount(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveClassesByAccount(ctx,
			item.Account,
		)
		_, found := keeper.GetClassesByAccount(ctx,
			item.Account,
		)
		require.False(t, found)
	}
}

func TestClassesByAccountGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassesByAccount(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllClassesByAccount(ctx)),
	)
}
