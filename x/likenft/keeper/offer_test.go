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

func createNOffer(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Offer {
	items := make([]types.Offer, n)
	for i := range items {
		items[i].ClassId = strconv.Itoa(i)
		items[i].NftId = strconv.Itoa(i)
		items[i].Buyer = strconv.Itoa(i)

		keeper.SetOffer(ctx, items[i])
	}
	return items
}

func TestOfferGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNOffer(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetOffer(ctx,
			item.ClassId,
			item.NftId,
			item.Buyer,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestOfferRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNOffer(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveOffer(ctx,
			item.ClassId,
			item.NftId,
			item.Buyer,
		)
		_, found := keeper.GetOffer(ctx,
			item.ClassId,
			item.NftId,
			item.Buyer,
		)
		require.False(t, found)
	}
}

func TestOfferGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items := createNOffer(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOffer(ctx)),
	)
}
