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

func createNOffer(keeper *keeper.Keeper, ctx sdk.Context, nClass int, nNFT int) ([]types.Offer, []sdk.AccAddress) {
	var items []types.Offer
	accounts := testutil.CreateIncrementalAccounts(nClass * nNFT)
	for i := 0; i < nClass; i++ {
		for j := 0; j < nNFT; j++ {
			offer := types.Offer{
				ClassId:    strconv.Itoa(i),
				NftId:      strconv.Itoa(j),
				Buyer:      accounts[i*j].String(),
				Price:      uint64(i * j),
				Expiration: time.Date(2022, 1, 1+i*j, 0, 0, 0, 0, time.UTC),
			}
			items = append(items, offer)
			keeper.SetOffer(ctx, offer)
		}
	}
	return items, accounts
}

func TestOfferGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, accs := createNOffer(keeper, ctx, 3, 3)
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
	items, accs := createNOffer(keeper, ctx, 3, 3)
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

func TestOfferGetByClass(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNOffer(keeper, ctx, 3, 3)
	require.ElementsMatch(t,
		nullify.Fill(items[0:3]),
		nullify.Fill(keeper.GetOffersByClass(ctx, "0")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[3:6]),
		nullify.Fill(keeper.GetOffersByClass(ctx, "1")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[6:9]),
		nullify.Fill(keeper.GetOffersByClass(ctx, "2")),
	)
}

func TestOfferGetByNFT(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNOffer(keeper, ctx, 1, 3)
	require.ElementsMatch(t,
		nullify.Fill([]types.Offer{items[0]}),
		nullify.Fill(keeper.GetOffersByNFT(ctx, "0", "0")),
	)
	require.ElementsMatch(t,
		nullify.Fill([]types.Offer{items[1]}),
		nullify.Fill(keeper.GetOffersByNFT(ctx, "0", "1")),
	)
	require.ElementsMatch(t,
		nullify.Fill([]types.Offer{items[2]}),
		nullify.Fill(keeper.GetOffersByNFT(ctx, "0", "2")),
	)
}

func TestOfferGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNOffer(keeper, ctx, 3, 3)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOffer(ctx)),
	)
}
