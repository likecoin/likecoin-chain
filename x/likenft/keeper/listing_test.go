package keeper_test

import (
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/testutil"
	keepertest "github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/testutil/nullify"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNListing(keeper *keeper.Keeper, ctx sdk.Context, nClass int, nNFT int, nListing int) ([]types.ListingStoreRecord, []sdk.AccAddress) {
	var items []types.ListingStoreRecord
	accounts := testutil.CreateIncrementalAccounts(nListing)
	for i := 0; i < nClass; i++ {
		for j := 0; j < nNFT; j++ {
			for k := 0; k < nListing; k++ {
				listing := types.ListingStoreRecord{
					ClassId:          strconv.Itoa(i),
					NftId:            strconv.Itoa(j),
					Seller:           accounts[k],
					Price:            uint64(k),
					Expiration:       time.Date(2022, 1, 1+k, 0, 0, 0, 0, time.UTC),
					FullPayToRoyalty: false,
				}
				items = append(items, listing)
				keeper.SetListing(ctx, listing)
			}
		}
	}
	return items, accounts
}

func TestListingGet(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNListing(keeper, ctx, 2, 2, 2)
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
	items, _ := createNListing(keeper, ctx, 2, 2, 2)
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

func TestListingGetByClass(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNListing(keeper, ctx, 3, 3, 1)
	require.ElementsMatch(t,
		nullify.Fill(items[0:3]),
		nullify.Fill(keeper.GetListingsByClass(ctx, "0")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[3:6]),
		nullify.Fill(keeper.GetListingsByClass(ctx, "1")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[6:9]),
		nullify.Fill(keeper.GetListingsByClass(ctx, "2")),
	)
}

func TestListingGetByNFT(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNListing(keeper, ctx, 1, 3, 1)
	require.ElementsMatch(t,
		nullify.Fill([]types.ListingStoreRecord{items[0]}),
		nullify.Fill(keeper.GetListingsByNFT(ctx, "0", "0")),
	)
	require.ElementsMatch(t,
		nullify.Fill([]types.ListingStoreRecord{items[1]}),
		nullify.Fill(keeper.GetListingsByNFT(ctx, "0", "1")),
	)
	require.ElementsMatch(t,
		nullify.Fill([]types.ListingStoreRecord{items[2]}),
		nullify.Fill(keeper.GetListingsByNFT(ctx, "0", "2")),
	)
}

func TestListingGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNListing(keeper, ctx, 2, 2, 2)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllListing(ctx)),
	)
}
