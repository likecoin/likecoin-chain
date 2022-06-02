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

func createNListing(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Listing {
	items := make([]types.Listing, n)
	for i := range items {
		items[i].ClassId = strconv.Itoa(i)
		items[i].NftId = strconv.Itoa(i)
		items[i].Seller = strconv.Itoa(i)

		keeper.SetListing(ctx, items[i])
	}
	return items
}

func TestListingGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNListing(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetListing(ctx,
			item.ClassId,
			item.NftId,
			item.Seller,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestListingRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNListing(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveListing(ctx,
			item.ClassId,
			item.NftId,
			item.Seller,
		)
		_, found := keeper.GetListing(ctx,
			item.ClassId,
			item.NftId,
			item.Seller,
		)
		require.False(t, found)
	}
}

func TestListingGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNListing(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllListing(ctx)),
	)
}
