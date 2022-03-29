package keeper_test

import (
	"fmt"
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
	msgs := createNClassesByISCN(keeper, ctx, 2)
	// Note: currently does not test pagination (class id arrays are empty)
	// Effort needed to seed data via mock (similar to e2e test)
	// Pagination is already covered by unit test and cli test
	classesByMsgs := testutil.BatchMakeDummyNFTClassesForISCN(msgs)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryClassesByISCNRequest
		response *types.QueryClassesByISCNResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryClassesByISCNRequest{
				IscnIdPrefix: msgs[0].IscnIdPrefix,
				Pagination: &query.PageRequest{
					Limit: uint64(len(classesByMsgs[0])),
				},
			},
			response: &types.QueryClassesByISCNResponse{
				IscnIdPrefix: msgs[0].IscnIdPrefix,
				Classes:      classesByMsgs[0],
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   uint64(len(classesByMsgs[0])),
				},
			},
		},
		{
			desc: "Second",
			request: &types.QueryClassesByISCNRequest{
				IscnIdPrefix: msgs[1].IscnIdPrefix,
				Pagination: &query.PageRequest{
					Limit: uint64(len(classesByMsgs[0])),
				},
			},
			response: &types.QueryClassesByISCNResponse{
				IscnIdPrefix: msgs[1].IscnIdPrefix,
				Classes:      classesByMsgs[1],
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   uint64(len(classesByMsgs[0])),
				},
			},
		},
		{
			desc: "Full id",
			request: &types.QueryClassesByISCNRequest{
				IscnIdPrefix: fmt.Sprintf("%s/1", msgs[1].IscnIdPrefix),
				Pagination: &query.PageRequest{
					Limit: uint64(len(classesByMsgs[0])),
				},
			},
			response: &types.QueryClassesByISCNResponse{
				IscnIdPrefix: msgs[1].IscnIdPrefix,
				Classes:      classesByMsgs[1],
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   uint64(len(classesByMsgs[0])),
				},
			},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryClassesByISCNRequest{
				IscnIdPrefix: "iscn://likecoin-chain/100000",
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
	msgs := createNClassesByISCN(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryClassesByISCNIndexRequest {
		return &types.QueryClassesByISCNIndexRequest{
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
			resp, err := keeper.ClassesByISCNIndex(wctx, request(nil, uint64(i), uint64(step), false))
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
			resp, err := keeper.ClassesByISCNIndex(wctx, request(next, 0, uint64(step), false))
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
		resp, err := keeper.ClassesByISCNIndex(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ClassesByISCN),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ClassesByISCNIndex(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
