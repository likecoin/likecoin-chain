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
	// calculate royalty
	_, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}
	if classData.Config.RoyaltyBasisPoints > k.MaxRoyaltyBasisPoints(ctx) {
		return nil, types.ErrInvalidNftClassConfig.Wrapf("Royalty basis points cannot be greater than %s", k.MaxRoyaltyBasisPointsText(ctx))
	}
	royaltyAmount := msg.Price / 10000 * classData.Config.RoyaltyBasisPoints
	// pay royalty if needed, could be 0 if price < 10000
	if royaltyAmount > 0 {
		classParent, err := k.ValidateAndRefreshClassParent(ctx, msg.ClassId, classData.Parent)
		if err != nil {
			return nil, err
		}
		royaltyAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(royaltyAmount))))
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, classParent.Owner, royaltyAmountCoins)
		if err != nil {
			return nil, types.ErrFailedToSellNFT.Wrapf(err.Error())
		}
	}
	// pay seller
	netAmount := msg.Price - royaltyAmount
	netAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(netAmount))))
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sellerAddress, netAmountCoins)
	if err != nil {
		return nil, types.ErrFailedToSellNFT.Wrapf(err.Error())
	}
	// refund remainder to buyer
	remainder := offer.Price - msg.Price
	if remainder > 0 {
		remainCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(remainder))))
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyerAddress, remainCoins)
		if err != nil {
			return nil, types.ErrFailedToSellNFT.Wrapf(err.Error())
		}
	}
	// sanity check
	if royaltyAmount+netAmount+remainder != offer.Price {
		return nil, types.ErrFailedToSellNFT.Wrapf("Price split calculation error")
	}
	// transfer nft to buyer
	err = k.nftKeeper.Transfer(ctx, msg.ClassId, msg.NftId, buyerAddress)
	if err != nil {
		return nil, types.ErrFailedToSellNFT.Wrapf(err.Error())
	}
	// remove offer
	k.RemoveOffer(
		ctx,
		offer.ClassId,
		offer.NftId,
		offer.Buyer,
	)
	k.RemoveOfferExpireQueueEntry(
		ctx,
		offer.Expiration,
		types.OfferKey(offer.ClassId, offer.NftId, offer.Buyer),
	)

	// owner changed, remove all listings
	k.PruneAllListingsForNFT(ctx, msg.ClassId, msg.NftId)

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
