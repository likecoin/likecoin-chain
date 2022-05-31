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

	// Check nft exists
	if isFound := k.nftKeeper.HasNFT(ctx, msg.ClassId, msg.NftId); !isFound {
		return nil, types.ErrNftNotFound
	}

	offer := types.Offer{
		ClassId:    msg.ClassId,
		NftId:      msg.NftId,
		Buyer:      msg.Creator,
		Price:      msg.Price,
		Expiration: msg.Expiration,
	}

	// Take deposit if needed
	if offer.Price > 0 {
		// TODO: rename mint price denom
		denom := k.MintPriceDenom(ctx)
		coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(int64(offer.Price))))
		if k.bankKeeper.GetBalance(ctx, userAddress, denom).Amount.Uint64() < offer.Price {
			return nil, types.ErrInsufficientFunds
		}
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, userAddress, types.ModuleName, coins); err != nil {
			return nil, types.ErrFailedToCreateOffer.Wrapf(err.Error())
		}
	}

	k.SetOffer(
		ctx,
		offer,
	)

	// TODO emit event

	return &types.MsgCreateOfferResponse{
		Offer: offer,
	}, nil
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
