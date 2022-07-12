package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/testutil/nullify"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestRoyaltyConfigByClassQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNRoyaltyConfig(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryRoyaltyConfigRequest
		response *types.QueryRoyaltyConfigResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryRoyaltyConfigRequest{
				ClassId: msgs[0].ClassId,
			},
			response: &types.QueryRoyaltyConfigResponse{RoyaltyConfig: msgs[0].RoyaltyConfig},
		},
		{
			desc: "Second",
			request: &types.QueryRoyaltyConfigRequest{
				ClassId: msgs[1].ClassId,
			},
			response: &types.QueryRoyaltyConfigResponse{RoyaltyConfig: msgs[1].RoyaltyConfig},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryRoyaltyConfigRequest{
				ClassId: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.RoyaltyConfig(wctx, tc.request)
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

func TestRoyaltyConfigByClassQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNRoyaltyConfig(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryRoyaltyConfigIndexRequest {
		return &types.QueryRoyaltyConfigIndexRequest{
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
			resp, err := keeper.RoyaltyConfigIndex(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RoyaltyConfigByClass), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.RoyaltyConfigByClass),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.RoyaltyConfigIndex(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RoyaltyConfigByClass), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.RoyaltyConfigByClass),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.RoyaltyConfigIndex(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.RoyaltyConfigByClass),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.RoyaltyConfigIndex(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
