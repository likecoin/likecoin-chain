package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateListing(goCtx context.Context, msg *types.MsgCreateListing) (*types.MsgCreateListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// Check if the value already exists
	_, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var listing = types.Listing{
		ClassId:    msg.ClassId,
		NftId:      msg.NftId,
		Seller:     msg.Creator,
		Price:      msg.Price,
		Expiration: msg.Expiration,
	}

	k.SetListing(
		ctx,
		listing,
	)
	return &types.MsgCreateListingResponse{}, nil
}

func (k msgServer) UpdateListing(goCtx context.Context, msg *types.MsgUpdateListing) (*types.MsgUpdateListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// Check if the value exists
	_, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	var listing = types.Listing{
		ClassId:    msg.ClassId,
		NftId:      msg.NftId,
		Seller:     msg.Creator,
		Price:      msg.Price,
		Expiration: msg.Expiration,
	}

	k.SetListing(ctx, listing)

	return &types.MsgUpdateListingResponse{}, nil
}

func (k msgServer) DeleteListing(goCtx context.Context, msg *types.MsgDeleteListing) (*types.MsgDeleteListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// Check if the value exists
	_, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	k.RemoveListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)

	return &types.MsgDeleteListingResponse{}, nil
}
