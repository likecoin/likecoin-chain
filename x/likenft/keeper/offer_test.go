package keeper_test

import (
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v4/testutil"
	keepertest "github.com/likecoin/likecoin-chain/v4/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v4/testutil/nullify"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNOffer(keeper *keeper.Keeper, ctx sdk.Context, nClass int, nNFT int, nOffer int) ([]types.OfferStoreRecord, []sdk.AccAddress) {
	var items []types.OfferStoreRecord
	accounts := testutil.CreateIncrementalAccounts(nOffer)
	for i := 0; i < nClass; i++ {
		for j := 0; j < nNFT; j++ {
			for k := 0; k < nOffer; k++ {
				offer := types.OfferStoreRecord{
					ClassId:    strconv.Itoa(i),
					NftId:      strconv.Itoa(j),
					Buyer:      accounts[k],
					Price:      uint64(k),
					Expiration: time.Date(2022, 1, 1+k, 0, 0, 0, 0, time.UTC),
				}
				items = append(items, offer)
				keeper.SetOffer(ctx, offer)
			}
		}
	}
	return items, accounts
}

func TestOfferGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNOffer(keeper, ctx, 3, 3, 1)
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
	items, _ := createNOffer(keeper, ctx, 3, 3, 1)
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

func TestOfferGetByClass(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNOffer(keeper, ctx, 3, 3, 1)
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
	items, _ := createNOffer(keeper, ctx, 1, 3, 1)
	require.ElementsMatch(t,
		nullify.Fill([]types.OfferStoreRecord{items[0]}),
		nullify.Fill(keeper.GetOffersByNFT(ctx, "0", "0")),
	)
	require.ElementsMatch(t,
		nullify.Fill([]types.OfferStoreRecord{items[1]}),
		nullify.Fill(keeper.GetOffersByNFT(ctx, "0", "1")),
	)
	require.ElementsMatch(t,
		nullify.Fill([]types.OfferStoreRecord{items[2]}),
		nullify.Fill(keeper.GetOffersByNFT(ctx, "0", "2")),
	)
}

func TestOfferGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNOffer(keeper, ctx, 3, 3, 1)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOffer(ctx)),
	)
}
