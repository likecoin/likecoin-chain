package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/likecoin/likecoin-chain/v4/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v4/testutil/nullify"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestListingQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	_msgs, _ := createNListing(keeper, ctx, 2, 1, 1)
	msgs := types.MapListingsToPublicRecords(_msgs)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryListingRequest
		response *types.QueryListingResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryListingRequest{
				ClassId: msgs[0].ClassId,
				NftId:   msgs[0].NftId,
				Seller:  msgs[0].Seller,
			},
			response: &types.QueryListingResponse{Listing: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryListingRequest{
				ClassId: msgs[1].ClassId,
				NftId:   msgs[1].NftId,
				Seller:  msgs[1].Seller,
			},
			response: &types.QueryListingResponse{Listing: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryListingRequest{
				ClassId: strconv.Itoa(100000),
				NftId:   strconv.Itoa(100000),
				Seller:  msgs[0].Seller,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Listing(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestListingQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	_msgs, _ := createNListing(keeper, ctx, 2, 2, 2)
	msgs := types.MapListingsToPublicRecords(_msgs)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryListingIndexRequest {
		return &types.QueryListingIndexRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ListingIndex(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Listings), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Listings),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ListingIndex(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Listings), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Listings),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ListingIndex(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Listings),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ListingIndex(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestListingByClassQuery(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	_msgs, _ := createNListing(keeper, ctx, 3, 3, 1)
	msgs := types.MapListingsToPublicRecords(_msgs)

	request := func(classId string, next []byte, offset, limit uint64, total bool) *types.QueryListingsByClassRequest {
		return &types.QueryListingsByClassRequest{
			ClassId: classId,
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs[3:6]); i += step {
			resp, err := keeper.ListingsByClass(wctx, request("1", nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Listings), step)
			require.Subset(t,
				nullify.Fill(msgs[3:6]),
				nullify.Fill(resp.Listings),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs[3:6]); i += step {
			resp, err := keeper.ListingsByClass(wctx, request("1", next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Listings), step)
			require.Subset(t,
				nullify.Fill(msgs[3:6]),
				nullify.Fill(resp.Listings),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ListingsByClass(wctx, request("1", nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs[3:6]), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs[3:6]),
			nullify.Fill(resp.Listings),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ListingsByClass(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestListingByNftQuery(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	_msgs, _ := createNListing(keeper, ctx, 1, 3, 3)
	msgs := types.MapListingsToPublicRecords(_msgs)

	request := func(classId string, nftId string, next []byte, offset, limit uint64, total bool) *types.QueryListingsByNFTRequest {
		return &types.QueryListingsByNFTRequest{
			ClassId: classId,
			NftId:   nftId,
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs[3:6]); i += step {
			resp, err := keeper.ListingsByNFT(wctx, request("0", "1", nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Listings), step)
			require.Subset(t,
				nullify.Fill(msgs[3:6]),
				nullify.Fill(resp.Listings),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs[3:6]); i += step {
			resp, err := keeper.ListingsByNFT(wctx, request("0", "1", next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Listings), step)
			require.Subset(t,
				nullify.Fill(msgs[3:6]),
				nullify.Fill(resp.Listings),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ListingsByNFT(wctx, request("0", "1", nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs[3:6]), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs[3:6]),
			nullify.Fill(resp.Listings),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ListingsByClass(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
