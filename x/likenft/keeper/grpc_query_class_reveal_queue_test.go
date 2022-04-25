package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestClassRevealQueueQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNClassRevealQueue(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetClassRevealQueueRequest
		response *types.QueryGetClassRevealQueueResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetClassRevealQueueRequest{
				RevealTime: msgs[0].RevealTime,
				ClassId:    msgs[0].ClassId,
			},
			response: &types.QueryGetClassRevealQueueResponse{ClassRevealQueue: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetClassRevealQueueRequest{
				RevealTime: msgs[1].RevealTime,
				ClassId:    msgs[1].ClassId,
			},
			response: &types.QueryGetClassRevealQueueResponse{ClassRevealQueue: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetClassRevealQueueRequest{
				RevealTime: strconv.Itoa(100000),
				ClassId:    strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ClassRevealQueue(wctx, tc.request)
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

func TestClassRevealQueueQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNClassRevealQueue(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllClassRevealQueueRequest {
		return &types.QueryAllClassRevealQueueRequest{
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
			resp, err := keeper.ClassRevealQueueAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassRevealQueue), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ClassRevealQueue),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ClassRevealQueueAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassRevealQueue), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ClassRevealQueue),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ClassRevealQueueAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ClassRevealQueue),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ClassRevealQueueAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
