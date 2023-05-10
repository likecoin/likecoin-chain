package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func TestMintNFTAuthorization(t *testing.T) {
	var auth *types.MintNFTAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)

	auth = types.NewMintNFTAuthorization(classId1)
	msg = &types.MsgMintNFT{
		ClassId: classId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgMintNFT{
		ClassId: classId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgBurnNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewMintNFTAuthorization(classId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)
}

func TestSendNFTAuthorization(t *testing.T) {
	var auth *types.SendNFTAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)

	nftId1 := "testing-id-1"
	nftId2 := "testing-id-2"

	auth = types.NewSendNFTAuthorization(classId1, nftId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &nft.MsgSend{
		ClassId: classId1,
		Id:      nftId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())

	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &nft.MsgSend{
		ClassId: classId1,
		Id:      nftId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &nft.MsgSend{
		ClassId: classId2,
		Id:      nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgBurnNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewSendNFTAuthorization(classId1, "")
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &nft.MsgSend{
		ClassId: classId1,
		Id:      nftId1,
	}

	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &nft.MsgSend{
		ClassId: classId1,
		Id:      nftId2,
	}

	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &nft.MsgSend{
		ClassId: classId2,
		Id:      nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}
