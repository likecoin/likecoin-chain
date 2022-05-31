package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestOfferMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.LikenftKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateOffer{Creator: creator,
			ClassId: strconv.Itoa(i),
			NftId:   strconv.Itoa(i),
		}
		_, err := srv.CreateOffer(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetOffer(ctx,
			expected.ClassId,
			expected.NftId,
			expected.Creator,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Buyer)
	}
}

func TestOfferMsgServerUpdate(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateOffer
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateOffer{Creator: creator,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateOffer{Creator: "B",
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateOffer{Creator: creator,
				ClassId: strconv.Itoa(100000),
				NftId:   strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.LikenftKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateOffer{Creator: creator,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			}
			_, err := srv.CreateOffer(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateOffer(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetOffer(ctx,
					expected.ClassId,
					expected.NftId,
					expected.Creator,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Buyer)
			}
		})
	}
}

func TestOfferMsgServerDelete(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteOffer
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteOffer{Creator: creator,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteOffer{Creator: "B",
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteOffer{Creator: creator,
				ClassId: strconv.Itoa(100000),
				NftId:   strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.LikenftKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateOffer(wctx, &types.MsgCreateOffer{Creator: creator,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteOffer(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetOffer(ctx,
					tc.request.ClassId,
					tc.request.NftId,
					tc.request.Creator,
				)
				require.False(t, found)
			}
		})
	}
}
