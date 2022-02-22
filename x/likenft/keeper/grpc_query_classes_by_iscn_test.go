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
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestClassesByISCNQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	_msgs := createNClassesByISCN(keeper, ctx, 2)
	msgs := testutil.BatchDummyConcretizeClassesByISCN(_msgs)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetClassesByISCNRequest
		response *types.QueryGetClassesByISCNResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetClassesByISCNRequest{
				IscnIdPrefix: msgs[0].IscnIdPrefix,
			},
			response: &types.QueryGetClassesByISCNResponse{ClassesByISCN: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetClassesByISCNRequest{
				IscnIdPrefix: msgs[1].IscnIdPrefix,
			},
			response: &types.QueryGetClassesByISCNResponse{ClassesByISCN: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetClassesByISCNRequest{
				IscnIdPrefix: strconv.Itoa(100000),
			},
			err: status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ClassesByISCN(wctx, tc.request)
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

func TestClassesByISCNQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	_msgs := createNClassesByISCN(keeper, ctx, 5)
	msgs := testutil.BatchDummyConcretizeClassesByISCN(_msgs)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllClassesByISCNRequest {
		return &types.QueryAllClassesByISCNRequest{
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
			resp, err := keeper.ClassesByISCNAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassesByISCN), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ClassesByISCN),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ClassesByISCNAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassesByISCN), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ClassesByISCN),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ClassesByISCNAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ClassesByISCN),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ClassesByISCNAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
