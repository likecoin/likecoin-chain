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

func TestCreateOfferAuthorization(t *testing.T) {
	var auth *types.CreateOfferAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)
	nftId1 := "nft-id-1"
	nftId2 := "nft-id-2"

	auth = types.NewCreateOfferAuthorization(classId1, nftId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgCreateOffer{
		ClassId: classId1,
		NftId:   nftId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgCreateOffer{
		ClassId: classId1,
		NftId:   nftId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgCreateOffer{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewCreateOfferAuthorization(classId1, "")
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgCreateOffer{
		ClassId: classId1,
		NftId:   nftId1,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgCreateOffer{
		ClassId: classId1,
		NftId:   nftId2,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgCreateOffer{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}

func TestUpdateOfferAuthorization(t *testing.T) {
	var auth *types.UpdateOfferAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)
	nftId1 := "nft-id-1"
	nftId2 := "nft-id-2"

	auth = types.NewUpdateOfferAuthorization(classId1, nftId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgUpdateOffer{
		ClassId: classId1,
		NftId:   nftId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateOffer{
		ClassId: classId1,
		NftId:   nftId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgUpdateOffer{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewUpdateOfferAuthorization(classId1, "")
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgUpdateOffer{
		ClassId: classId1,
		NftId:   nftId1,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateOffer{
		ClassId: classId1,
		NftId:   nftId2,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateOffer{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}

func TestdeleteOfferAuthorization(t *testing.T) {
	var auth *types.DeleteOfferAuthorization
	var msg sdk.Msg
	var res authz.AcceptResponse
	var err error

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	classId1, err := types.NewClassIdForAccount(sdk.AccAddress{1}, 1)
	require.NoError(t, err)
	classId2, err := types.NewClassIdForAccount(sdk.AccAddress{2}, 1)
	require.NoError(t, err)
	nftId1 := "nft-id-1"
	nftId2 := "nft-id-2"

	auth = types.NewDeleteOfferAuthorization(classId1, nftId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgDeleteOffer{
		ClassId: classId1,
		NftId:   nftId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgDeleteOffer{
		ClassId: classId1,
		NftId:   nftId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgDeleteOffer{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewDeleteOfferAuthorization(classId1, "")
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgDeleteOffer{
		ClassId: classId1,
		NftId:   nftId1,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgDeleteOffer{
		ClassId: classId1,
		NftId:   nftId2,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgDeleteOffer{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}
