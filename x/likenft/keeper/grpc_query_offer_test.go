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

func TestOfferQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs, _ := createNOffer(keeper, ctx, 3, 3)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryOfferRequest
		response *types.QueryOfferResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryOfferRequest{
				ClassId: msgs[0].ClassId,
				NftId:   msgs[0].NftId,
				Buyer:   msgs[0].Buyer,
			},
			response: &types.QueryOfferResponse{Offer: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryOfferRequest{
				ClassId: msgs[1].ClassId,
				NftId:   msgs[1].NftId,
				Buyer:   msgs[1].Buyer,
			},
			response: &types.QueryOfferResponse{Offer: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryOfferRequest{
				ClassId: strconv.Itoa(100000),
				NftId:   strconv.Itoa(100000),
				Buyer:   msgs[0].Buyer,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Offer(wctx, tc.request)
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

func TestOfferQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs, _ := createNOffer(keeper, ctx, 3, 3)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryOfferIndexRequest {
		return &types.QueryOfferIndexRequest{
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
			resp, err := keeper.OfferIndex(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Offers), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Offers),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.OfferIndex(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Offers), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Offers),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.OfferIndex(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Offers),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.OfferIndex(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestOfferByClassQuery(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs, _ := createNOffer(keeper, ctx, 3, 3)

	request := func(classId string, next []byte, offset, limit uint64, total bool) *types.QueryOffersByClassRequest {
		return &types.QueryOffersByClassRequest{
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
			resp, err := keeper.OffersByClass(wctx, request("1", nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Offers), step)
			require.Subset(t,
				nullify.Fill(msgs[3:6]),
				nullify.Fill(resp.Offers),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs[3:6]); i += step {
			resp, err := keeper.OffersByClass(wctx, request("1", next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Offers), step)
			require.Subset(t,
				nullify.Fill(msgs[3:6]),
				nullify.Fill(resp.Offers),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.OffersByClass(wctx, request("1", nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs[3:6]), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs[3:6]),
			nullify.Fill(resp.Offers),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.OfferIndex(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
