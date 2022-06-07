package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) SellNFT(goCtx context.Context, msg *types.MsgSellNFT) (*types.MsgSellNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check user is current owner
	sellerAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}
	if !k.nftKeeper.GetOwner(ctx, msg.ClassId, msg.NftId).Equals(sellerAddress) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("User do not own the NFT")
	}

	// check offer exists
	buyerAddress, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return nil, types.ErrOfferNotFound
	}
	offer, isFound := k.GetOffer(ctx, msg.ClassId, msg.NftId, buyerAddress)
	if !isFound {
		return nil, types.ErrOfferNotFound
	}

	// check offer is not expired
	if offer.Expiration.Before(ctx.BlockHeader().Time) {
		return nil, types.ErrOfferExpired
	}

	// check price <= offer price
	if msg.Price > offer.Price {
		return nil, types.ErrFailedToSellNFT.Wrapf("Price is too high. Offered price was %d", offer.Price)
	}

	// transact
	// pay seller
	priceCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).MintPriceDenom, sdk.NewInt(int64(msg.Price))))
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sellerAddress, priceCoins)
	if err != nil {
		return nil, types.ErrFailedToSellNFT.Wrapf(err.Error())
	}
	// TODO: pay royalty to class parent owner
	// refund remainder to buyer
	remainder := offer.Price - msg.Price
	if remainder > 0 {
		remainCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).MintPriceDenom, sdk.NewInt(int64(remainder))))
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyerAddress, remainCoins)
		if err != nil {
			return nil, types.ErrFailedToSellNFT.Wrapf(err.Error())
		}
	}
	// transfer nft to buyer
	err = k.nftKeeper.Transfer(ctx, msg.ClassId, msg.NftId, buyerAddress)
	if err != nil {
		return nil, types.ErrFailedToSellNFT.Wrapf(err.Error())
	}
	// remove offer
	k.RemoveOffer(ctx, msg.ClassId, msg.NftId, buyerAddress)

	// owner changed, prune invalid listings
	k.PruneInvalidListingsForNFT(ctx, msg.ClassId, msg.NftId)

	// emit event
	ctx.EventManager().EmitTypedEvent(&types.EventSellNFT{
		ClassId: msg.ClassId,
		NftId:   msg.NftId,
		Seller:  sellerAddress.String(),
		Buyer:   buyerAddress.String(),
		Price:   msg.Price,
	})

	return &types.MsgSellNFTResponse{}, nil
}
