package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNMintableNFT(keeper *keeper.Keeper, ctx sdk.Context, nClass int, nNFT int) []types.MintableNFT {
	var items []types.MintableNFT
	for i := 0; i < nClass; i++ {
		for j := 0; j < nNFT; j++ {
			item := types.MintableNFT{
				ClassId: strconv.Itoa(i),
				Id:      strconv.Itoa(j),
			}
			items = append(items, item)
			keeper.SetMintableNFT(ctx, item)
		}
	}
	return items
}

func TestMintableNFTGet(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForMintableTest(t)
	defer ctrl.Finish()

	items := createNMintableNFT(keeper, ctx, 3, 3)
	for _, item := range items {
		rst, found := keeper.GetMintableNFT(ctx,
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
func TestMintableNFTRemove(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForMintableTest(t)
	defer ctrl.Finish()

	items := createNMintableNFT(keeper, ctx, 3, 3)
	for _, item := range items {
		keeper.RemoveMintableNFT(ctx,
			item.ClassId,
			item.Id,
		)
		_, found := keeper.GetMintableNFT(ctx,
			item.ClassId,
			item.Id,
		)
		require.False(t, found)
	}
}

func TestMintableNFTGetByClass(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForMintableTest(t)
	defer ctrl.Finish()

	items := createNMintableNFT(keeper, ctx, 3, 3)
	require.ElementsMatch(t,
		nullify.Fill(items[0:3]),
		nullify.Fill(keeper.GetMintableNFTs(ctx, "0")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[3:6]),
		nullify.Fill(keeper.GetMintableNFTs(ctx, "1")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[6:9]),
		nullify.Fill(keeper.GetMintableNFTs(ctx, "2")),
	)
}

func TestMintableNFTGetAll(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForMintableTest(t)
	defer ctrl.Finish()

	items := createNMintableNFT(keeper, ctx, 3, 3)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllMintableNFT(ctx)),
	)
}
