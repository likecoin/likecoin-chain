package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateOffer{Creator: ownerAddress,
			ClassId: strconv.Itoa(i),
			NftId:   strconv.Itoa(i),
		}
		_, err := srv.CreateOffer(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetOffer(ctx,
			expected.ClassId,
			expected.NftId,
			ownerAddressBytes,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Buyer)
	}
}

func TestOfferMsgServerUpdate(t *testing.T) {
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)

	notOwnerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	notOwnerAddress, _ := sdk.Bech32ifyAddressBytes("like", notOwnerAddressBytes)

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateOffer
		acc     sdk.AccAddress
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateOffer{Creator: ownerAddress,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
			acc: ownerAddressBytes,
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateOffer{Creator: notOwnerAddress,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
			acc: notOwnerAddressBytes,
			err: types.ErrOfferNotFound,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateOffer{Creator: ownerAddress,
				ClassId: strconv.Itoa(100000),
				NftId:   strconv.Itoa(100000),
			},
			acc: ownerAddressBytes,
			err: types.ErrOfferNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.LikenftKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateOffer{Creator: ownerAddress,
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
					tc.acc,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Buyer)
			}
		})
	}
}

func TestOfferMsgServerDelete(t *testing.T) {
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)

	notOwnerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	notOwnerAddress, _ := sdk.Bech32ifyAddressBytes("like", notOwnerAddressBytes)

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteOffer
		acc     sdk.AccAddress
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteOffer{Creator: ownerAddress,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
			acc: ownerAddressBytes,
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteOffer{Creator: notOwnerAddress,
				ClassId: strconv.Itoa(0),
				NftId:   strconv.Itoa(0),
			},
			acc: notOwnerAddressBytes,
			err: types.ErrOfferNotFound,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteOffer{Creator: ownerAddress,
				ClassId: strconv.Itoa(100000),
				NftId:   strconv.Itoa(100000),
			},
			acc: notOwnerAddressBytes,
			err: types.ErrOfferNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.LikenftKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateOffer(wctx, &types.MsgCreateOffer{Creator: ownerAddress,
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
					tc.acc,
				)
				require.False(t, found)
			}
		})
	}
}
