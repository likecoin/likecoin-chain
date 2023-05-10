package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func TestNewClassAuthorization(t *testing.T) {
	var auth *types.NewClassAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	iscnIdPrefix1 := "iscn://testing/prefix1"
	iscnIdPrefix2 := "iscn://testing/prefix2"

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)

	auth = types.NewNewClassAuthorization(iscnIdPrefix1)
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgNewClass{
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnIdPrefix1,
		},
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgNewClass{
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnIdPrefix2,
		},
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgNewClass{
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ACCOUNT,
			IscnIdPrefix: iscnIdPrefix2,
		},
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)
}

func TestUpdateClassAuthorization(t *testing.T) {
	var auth *types.UpdateClassAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)

	auth = types.NewUpdateClassAuthorization(classId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgUpdateClass{
		ClassId: classId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateClass{
		ClassId: classId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)
}
