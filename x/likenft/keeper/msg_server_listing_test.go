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

func TestListingMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.LikenftKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateListing{Creator: creator,
			ClassId: strconv.Itoa(i),
			NftId:   strconv.Itoa(i),
		}
		_, err := srv.CreateListing(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetListing(ctx,
			expected.ClassId,
			expected.NftId,
			expected.Creator,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Seller)
	}
}

func TestListingMsgServerUpdate(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateListing
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateListing{Creator: creator,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateListing{Creator: "B",
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateListing{Creator: creator,
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
			expected := &types.MsgCreateListing{Creator: creator,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			}
			_, err := srv.CreateListing(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateListing(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetListing(ctx,
					expected.ClassId,
					expected.NftId,
					expected.Creator,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Seller)
			}
		})
	}
}

func TestListingMsgServerDelete(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteListing
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteListing{Creator: creator,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteListing{Creator: "B",
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteListing{Creator: creator,
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

			_, err := srv.CreateListing(wctx, &types.MsgCreateListing{Creator: creator,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteListing(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetListing(ctx,
					tc.request.ClassId,
					tc.request.NftId,
					tc.request.Creator,
				)
				require.False(t, found)
			}
		})
	}
}
