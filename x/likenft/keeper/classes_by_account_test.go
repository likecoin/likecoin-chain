package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/testutil"
	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNClassesByAccount(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.ClassesByAccount, []sdk.AccAddress) {
	items := make([]types.ClassesByAccount, n)
	accounts := testutil.CreateIncrementalAccounts(n)
	for i := range items {
		items[i].Account = accounts[i].String()
		keeper.SetClassesByAccount(ctx, items[i])
	}
	return items, accounts
}

func TestClassesByAccountGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, accounts := createNClassesByAccount(keeper, ctx, 10)
	for i, item := range items {
		rst, found := keeper.GetClassesByAccount(ctx,
			accounts[i],
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
	items, accounts := createNClassesByAccount(keeper, ctx, 10)
	for i, _ := range items {
		keeper.RemoveClassesByAccount(ctx,
			accounts[i],
		)
		_, found := keeper.GetClassesByAccount(ctx,
			accounts[i],
		)
		require.False(t, found)
	}
}

func TestClassesByAccountGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNClassesByAccount(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllClassesByAccount(ctx)),
	)
}
