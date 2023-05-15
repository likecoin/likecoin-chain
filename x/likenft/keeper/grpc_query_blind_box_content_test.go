package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/likecoin/likecoin-chain/v4/testutil/nullify"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestBlindBoxContentQuerySingle(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBlindBoxContent(keeper, ctx, 3, 3)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryBlindBoxContentRequest
		response *types.QueryBlindBoxContentResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryBlindBoxContentRequest{
				ClassId: msgs[0].ClassId,
				Id:      msgs[0].Id,
			},
			response: &types.QueryBlindBoxContentResponse{BlindBoxContent: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryBlindBoxContentRequest{
				ClassId: msgs[1].ClassId,
				Id:      msgs[1].Id,
			},
			response: &types.QueryBlindBoxContentResponse{BlindBoxContent: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryBlindBoxContentRequest{
				ClassId: strconv.Itoa(100000),
				Id:      strconv.Itoa(100000),
			},
			err: status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.BlindBoxContent(wctx, tc.request)
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

func TestBlindBoxContentQueryByClass(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	wctx := sdk.WrapSDKContext(ctx)
	nClass := 3
	nItem := 3
	msgs := createNBlindBoxContent(keeper, ctx, nClass, nItem)
	request := func(classId string, next []byte, offset, limit uint64, total bool) *types.QueryBlindBoxContentsRequest {
		return &types.QueryBlindBoxContentsRequest{
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
		for classId := 0; classId < nClass; classId++ {
			step := 2
			for i := 0; i < nItem; i += step {
				resp, err := keeper.BlindBoxContents(wctx, request(strconv.Itoa(classId), nil, uint64(i), uint64(step), false))
				require.NoError(t, err)
				require.LessOrEqual(t, len(resp.BlindBoxContents), step)
				require.Subset(t,
					nullify.Fill(msgs),
					nullify.Fill(resp.BlindBoxContents),
				)
			}
		}

	})
	t.Run("ByKey", func(t *testing.T) {
		for classId := 0; classId < nClass; classId++ {
			step := 2
			var next []byte
			for i := 0; i < nItem; i += step {
				resp, err := keeper.BlindBoxContents(wctx, request(strconv.Itoa(classId), next, 0, uint64(step), false))
				require.NoError(t, err)
				require.LessOrEqual(t, len(resp.BlindBoxContents), step)
				require.Subset(t,
					nullify.Fill(msgs),
					nullify.Fill(resp.BlindBoxContents),
				)
				next = resp.Pagination.NextKey
			}
		}
	})
	t.Run("Total", func(t *testing.T) {
		for classId := 0; classId < nClass; classId++ {
			resp, err := keeper.BlindBoxContents(wctx, request(strconv.Itoa(classId), nil, 0, 0, true))
			require.NoError(t, err)
			require.Equal(t, nItem, int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(msgs[classId*nItem:classId*nItem+nItem]),
				nullify.Fill(resp.BlindBoxContents),
			)
		}
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.BlindBoxContents(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestBlindBoxContentQueryPaginated(t *testing.T) {
	keeper, ctx, ctrl := testutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBlindBoxContent(keeper, ctx, 3, 3)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryBlindBoxContentIndexRequest {
		return &types.QueryBlindBoxContentIndexRequest{
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
			resp, err := keeper.BlindBoxContentIndex(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.BlindBoxContents), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.BlindBoxContents),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.BlindBoxContentIndex(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.BlindBoxContents), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.BlindBoxContents),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.BlindBoxContentIndex(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.BlindBoxContents),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.BlindBoxContentIndex(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
