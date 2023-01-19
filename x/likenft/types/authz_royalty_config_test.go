package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func TestCreateRoyaltyConfigAuthorization(t *testing.T) {
	var auth *types.CreateRoyaltyConfigAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)

	auth = types.NewCreateRoyaltyConfigAuthorization(classId1)
	msg = &types.MsgCreateRoyaltyConfig{
		ClassId: classId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgCreateRoyaltyConfig{
		ClassId: classId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewCreateRoyaltyConfigAuthorization(classId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)
}

func TestUpdateRoyaltyConfigAuthorization(t *testing.T) {
	var auth *types.UpdateRoyaltyConfigAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)

	auth = types.NewUpdateRoyaltyConfigAuthorization(classId1)
	msg = &types.MsgUpdateRoyaltyConfig{
		ClassId: classId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateRoyaltyConfig{
		ClassId: classId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewUpdateRoyaltyConfigAuthorization(classId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)
}

func TestDeleteRoyaltyConfigAuthorization(t *testing.T) {
	var auth *types.DeleteRoyaltyConfigAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)

	auth = types.NewDeleteRoyaltyConfigAuthorization(classId1)
	msg = &types.MsgDeleteRoyaltyConfig{
		ClassId: classId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgDeleteRoyaltyConfig{
		ClassId: classId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewDeleteRoyaltyConfigAuthorization(classId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)
}
