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

func createNListing(keeper *keeper.Keeper, ctx sdk.Context, nClass int, nNFT int, nListing int) ([]types.Listing, []sdk.AccAddress) {
	var items []types.Listing
	accounts := testutil.CreateIncrementalAccounts(nListing)
	for i := 0; i < nClass; i++ {
		for j := 0; j < nNFT; j++ {
			for k := 0; k < nListing; k++ {
				listing := types.Listing{
					ClassId:    strconv.Itoa(i),
					NftId:      strconv.Itoa(j),
					Seller:     accounts[k].String(),
					Price:      uint64(k),
					Expiration: time.Date(2022, 1, 1+k, 0, 0, 0, 0, time.UTC),
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
		acc, err := sdk.AccAddressFromBech32(item.Seller)
		require.NoError(t, err)
		rst, found := keeper.GetListing(ctx,
			item.ClassId,
			item.NftId,
			acc,
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
		acc, err := sdk.AccAddressFromBech32(item.Seller)
		require.NoError(t, err)
		keeper.RemoveListing(ctx,
			item.ClassId,
			item.NftId,
			acc,
		)
		_, found := keeper.GetListing(ctx,
			item.ClassId,
			item.NftId,
			acc,
		)
		require.False(t, found)
	}
}

func TestListingGetAll(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	items, _ := createNListing(keeper, ctx, 2, 2, 2)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllListing(ctx)),
	)
}
