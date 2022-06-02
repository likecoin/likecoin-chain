package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateListing(goCtx context.Context, msg *types.MsgCreateListing) (*types.MsgCreateListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		msg.Creator,
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

	// Check if the value exists
	valFound, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		msg.Creator,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Seller {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
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

	// Check if the value exists
	valFound, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		msg.Creator,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Seller {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		msg.Creator,
	)

	return &types.MsgDeleteListingResponse{}, nil
}
