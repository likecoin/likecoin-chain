package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/testutil/nullify"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNClassesByISCN(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ClassesByISCN {
	items := make([]types.ClassesByISCN, n)
	for i := range items {
		items[i].IscnIdPrefix = fmt.Sprintf("iscn://likecoin-chain/%s", strconv.Itoa(i))

		keeper.SetClassesByISCN(ctx, items[i])
	}
	return items
}

func TestClassesByISCNGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassesByISCN(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetClassesByISCN(ctx,
			item.IscnIdPrefix,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestClassesByISCNRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassesByISCN(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveClassesByISCN(ctx,
			item.IscnIdPrefix,
		)
		_, found := keeper.GetClassesByISCN(ctx,
			item.IscnIdPrefix,
		)
		require.False(t, found)
	}
}

func TestClassesByISCNGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClassesByISCN(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllClassesByISCN(ctx)),
	)
}
