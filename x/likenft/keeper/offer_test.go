package keeper_test

import (
	"strconv"
	"testing"
	"time"

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

func createNOffer(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.Offer, []sdk.AccAddress) {
	items := make([]types.Offer, n)
	accounts := testutil.CreateIncrementalAccounts(n)
	for i := range items {
		items[i] = types.Offer{
			ClassId:    strconv.Itoa(i),
			NftId:      strconv.Itoa(i),
			Buyer:      accounts[i].String(),
			Price:      uint64(i),
			Expiration: time.Date(2022, 1, 1+i, 0, 0, 0, 0, time.UTC),
		}
		keeper.SetOffer(ctx, items[i])
	}
	return items, accounts
}

func TestOfferGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, accs := createNOffer(keeper, ctx, 10)
	for i, item := range items {
		rst, found := keeper.GetOffer(ctx,
			item.ClassId,
			item.NftId,
			accs[i],
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
	items, accs := createNOffer(keeper, ctx, 10)
	for i, item := range items {
		keeper.RemoveOffer(ctx,
			item.ClassId,
			item.NftId,
			accs[i],
		)
		_, found := keeper.GetOffer(ctx,
			item.ClassId,
			item.NftId,
			accs[i],
		)
		require.False(t, found)
	}
}

func TestOfferGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNOffer(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOffer(ctx)),
	)
}
