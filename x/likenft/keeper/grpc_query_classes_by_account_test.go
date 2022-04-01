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

func TestClassesByAccountQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs, _ := createNClassesByAccount(keeper, ctx, 2)
	classesByMsgs := testutil.BatchMakeDummyNFTClassesForAccount(msgs)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryClassesByAccountRequest
		response *types.QueryClassesByAccountResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryClassesByAccountRequest{
				Account: msgs[0].Account,
				Pagination: &query.PageRequest{
					Limit: uint64(len(classesByMsgs[0])),
				},
			},
			response: &types.QueryClassesByAccountResponse{
				Account: msgs[0].Account,
				Classes: classesByMsgs[0],
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   uint64(len(classesByMsgs[0])),
				},
			},
		},
		{
			desc: "Second",
			request: &types.QueryClassesByAccountRequest{
				Account: msgs[1].Account,
				Pagination: &query.PageRequest{
					Limit: uint64(len(classesByMsgs[0])),
				},
			},
			response: &types.QueryClassesByAccountResponse{
				Account: msgs[1].Account,
				Classes: classesByMsgs[1],
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   uint64(len(classesByMsgs[0])),
				},
			},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryClassesByAccountRequest{
				Account: "cosmos1kznrznww4pd6gx0zwrpthjk68fdmqypjpkj5hp",
			},
			err: status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ClassesByAccount(wctx, tc.request)
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

func TestClassesByAccountQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs, _ := createNClassesByAccount(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryClassesByAccountIndexRequest {
		return &types.QueryClassesByAccountIndexRequest{
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
			resp, err := keeper.ClassesByAccountIndex(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassesByAccount), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ClassesByAccount),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ClassesByAccountIndex(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ClassesByAccount), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ClassesByAccount),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ClassesByAccountIndex(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ClassesByAccount),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ClassesByAccountIndex(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
