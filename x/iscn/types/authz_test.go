package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/likecoin/likecoin-chain/v3/x/iscn/types"
)

func TestUpdateAuthorization(t *testing.T) {
	var auth *types.UpdateAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	iscnId1v1 := types.NewIscnId("test", "1111", 1)
	iscnId2v1 := types.NewIscnId("test", "2222", 1)

	auth = types.NewUpdateAuthorization(iscnId1v1.Prefix.String())
	msg = &types.MsgUpdateIscnRecord{
		IscnId: iscnId1v1.String(),
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateIscnRecord{
		IscnId: iscnId2v1.String(),
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgUpdateIscnRecord{
		IscnId: "invalid",
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, types.ErrInvalidIscnId)

	msg = &types.MsgCreateIscnRecord{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewUpdateAuthorization(iscnId1v1.Prefix.String())
	err = auth.ValidateBasic()
	require.NoError(t, err)
	auth = types.NewUpdateAuthorization(iscnId1v1.String())
	err = auth.ValidateBasic()
	require.NoError(t, err)
	auth = types.NewUpdateAuthorization("invalid")
	err = auth.ValidateBasic()
	require.ErrorIs(t, err, types.ErrInvalidIscnId)
}
