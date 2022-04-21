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

func createNClaimableNFT(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ClaimableNFT {
	items := make([]types.ClaimableNFT, n)
	for i := range items {
		items[i].ClassId = strconv.Itoa(i)
		items[i].Id = strconv.Itoa(i)

		keeper.SetClaimableNFT(ctx, items[i])
	}
	return items
}

func TestClaimableNFTGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClaimableNFT(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetClaimableNFT(ctx,
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
func TestClaimableNFTRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClaimableNFT(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveClaimableNFT(ctx,
			item.ClassId,
			item.Id,
		)
		_, found := keeper.GetClaimableNFT(ctx,
			item.ClassId,
			item.Id,
		)
		require.False(t, found)
	}
}

func TestClaimableNFTGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNClaimableNFT(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllClaimableNFT(ctx)),
	)
}
