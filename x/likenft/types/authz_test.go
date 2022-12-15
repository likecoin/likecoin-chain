package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
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
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewUpdateClassAuthorization(classId1)
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

func TestUpdateOfferAuthorizat(t *testing.T) {
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
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgUpdateOffer{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewUpdateOfferAuthorization(classId1, nftId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)
}

func TestUpdateListingAuthorization(t *testing.T) {
	var auth *types.UpdateListingAuthorization
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

	auth = types.NewUpdateListingAuthorization(classId1, nftId1)
	msg = &types.MsgUpdateListing{
		ClassId: classId1,
		NftId:   nftId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateListing{
		ClassId: classId1,
		NftId:   nftId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgUpdateListing{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewUpdateListingAuthorization(classId1, nftId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)
}
