package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateOffer(goCtx context.Context, msg *types.MsgCreateOffer) (*types.MsgCreateOfferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// Check if the value already exists
	_, isFound := k.GetOffer(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if isFound {
		return nil, types.ErrOfferAlreadyExists
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

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// Check if the value exists
	_, isFound := k.GetOffer(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if !isFound {
		return nil, types.ErrOfferNotFound
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

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// Check if the value exists
	_, isFound := k.GetOffer(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if !isFound {
		return nil, types.ErrOfferNotFound
	}

	k.RemoveOffer(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)

	return &types.MsgDeleteOfferResponse{}, nil
}
