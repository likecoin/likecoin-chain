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

func TestCreateListingAuthorization(t *testing.T) {
	var auth *types.CreateListingAuthorization
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

	auth = types.NewCreateListingAuthorization(classId1, nftId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgCreateListing{
		ClassId: classId1,
		NftId:   nftId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgCreateListing{
		ClassId: classId1,
		NftId:   nftId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgCreateListing{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewCreateListingAuthorization(classId1, "")
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgCreateListing{
		ClassId: classId1,
		NftId:   nftId1,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgCreateListing{
		ClassId: classId1,
		NftId:   nftId2,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgCreateListing{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
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
	err = auth.ValidateBasic()
	require.NoError(t, err)

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
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgUpdateListing{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewUpdateListingAuthorization(classId1, "")
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgUpdateListing{
		ClassId: classId1,
		NftId:   nftId1,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateListing{
		ClassId: classId1,
		NftId:   nftId2,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgUpdateListing{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}

func TestdeleteListingAuthorization(t *testing.T) {
	var auth *types.DeleteListingAuthorization
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

	auth = types.NewDeleteListingAuthorization(classId1, nftId1)
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgDeleteListing{
		ClassId: classId1,
		NftId:   nftId1,
	}
	require.Equal(t, sdk.MsgTypeURL(msg), auth.MsgTypeURL())
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgDeleteListing{
		ClassId: classId1,
		NftId:   nftId2,
	}
	_, err = auth.Accept(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgDeleteListing{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	msg = &types.MsgMintNFT{}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidType)

	auth = types.NewDeleteListingAuthorization(classId1, "")
	err = auth.ValidateBasic()
	require.NoError(t, err)

	msg = &types.MsgDeleteListing{
		ClassId: classId1,
		NftId:   nftId1,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgDeleteListing{
		ClassId: classId1,
		NftId:   nftId2,
	}
	res, err = auth.Accept(ctx, msg)
	require.NoError(t, err)
	require.True(t, res.Accept)
	require.False(t, res.Delete)
	require.Nil(t, res.Updated)

	msg = &types.MsgDeleteListing{
		ClassId: classId2,
		NftId:   nftId1,
	}
	_, err = auth.Accept(ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}
