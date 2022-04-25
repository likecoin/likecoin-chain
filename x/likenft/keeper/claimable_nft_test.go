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

func createNClaimableNFT(keeper *keeper.Keeper, ctx sdk.Context, nClass int, nNFT int) []types.ClaimableNFT {
	var items []types.ClaimableNFT
	for i := 0; i < nClass; i++ {
		for j := 0; j < nNFT; j++ {
			item := types.ClaimableNFT{
				ClassId: strconv.Itoa(i),
				Id:      strconv.Itoa(j),
			}
			items = append(items, item)
			keeper.SetClaimableNFT(ctx, item)
		}
	}
	return items
}

func TestClaimableNFTGet(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForClaimableTest(t)
	defer ctrl.Finish()

	items := createNClaimableNFT(keeper, ctx, 3, 3)
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
	keeper, ctx, ctrl := testutil.LikenftKeeperForClaimableTest(t)
	defer ctrl.Finish()

	items := createNClaimableNFT(keeper, ctx, 3, 3)
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

func TestClaimableNFTGetByClass(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForClaimableTest(t)
	defer ctrl.Finish()

	items := createNClaimableNFT(keeper, ctx, 3, 3)
	require.ElementsMatch(t,
		nullify.Fill(items[0:3]),
		nullify.Fill(keeper.GetClaimableNFTs(ctx, "0")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[3:6]),
		nullify.Fill(keeper.GetClaimableNFTs(ctx, "1")),
	)
	require.ElementsMatch(t,
		nullify.Fill(items[6:9]),
		nullify.Fill(keeper.GetClaimableNFTs(ctx, "2")),
	)
}

func TestClaimableNFTGetAll(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForClaimableTest(t)
	defer ctrl.Finish()

	items := createNClaimableNFT(keeper, ctx, 3, 3)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllClaimableNFT(ctx)),
	)
}
