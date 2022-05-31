package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateOffer(goCtx context.Context, msg *types.MsgCreateOffer) (*types.MsgCreateOfferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetOffer(
		ctx,
		msg.ClassId,
		msg.NftId,
		msg.Creator,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var offer = types.Offer{
		ClassId:    msg.ClassId,
		NftId:      msg.NftId,
		Buyer:      msg.Creator,
		Price:      msg.Price,
		Expiration: msg.Expiration,
	}

	k.SetOffer(
		ctx,
		offer,
	)
	return &types.MsgCreateOfferResponse{}, nil
}

func (k msgServer) UpdateOffer(goCtx context.Context, msg *types.MsgUpdateOffer) (*types.MsgUpdateOfferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetOffer(
		ctx,
		msg.ClassId,
		msg.NftId,
		msg.Creator,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Buyer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var offer = types.Offer{
		ClassId:    msg.ClassId,
		NftId:      msg.NftId,
		Buyer:      msg.Creator,
		Price:      msg.Price,
		Expiration: msg.Expiration,
	}

	k.SetOffer(ctx, offer)

	return &types.MsgUpdateOfferResponse{}, nil
}

func (k msgServer) DeleteOffer(goCtx context.Context, msg *types.MsgDeleteOffer) (*types.MsgDeleteOfferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetOffer(
		ctx,
		msg.ClassId,
		msg.NftId,
		msg.Creator,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Buyer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveOffer(
		ctx,
		msg.ClassId,
		msg.NftId,
		msg.Creator,
	)

	return &types.MsgDeleteOfferResponse{}, nil
}
