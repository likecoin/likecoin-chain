package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) BuyNFT(goCtx context.Context, msg *types.MsgBuyNFT) (*types.MsgBuyNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	buyerAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// check listing exists
	sellerAddress, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return nil, types.ErrListingNotFound
	}
	listing, isFound := k.GetListing(ctx, msg.ClassId, msg.NftId, sellerAddress)
	if !isFound {
		return nil, types.ErrListingNotFound
	}

	// check listing owner is still valid
	if !k.nftKeeper.GetOwner(ctx, msg.ClassId, msg.NftId).Equals(sellerAddress) {
		return nil, types.ErrListingExpired.Wrapf("Listing owner is no longer valid")
	}

	// check listing not expired
	if listing.Expiration.Before(ctx.BlockTime()) {
		return nil, types.ErrListingExpired
	}

	// check price >= listing price
	if msg.Price < listing.Price {
		return nil, types.ErrFailedToBuyNFT.Wrapf("Price is too low. Listing price was %d", listing.Price)
	}

	// check user has enough balance
	if k.bankKeeper.GetBalance(ctx, buyerAddress, k.GetParams(ctx).PriceDenom).Amount.Uint64() < msg.Price {
		return nil, types.ErrFailedToBuyNFT.Wrapf("User does not have enough balance")
	}

	// transact
	// calculate royalty
	_, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}
	if classData.Config.RoyaltyBasisPoints > types.MaxRoyaltyBasisPoints {
		return nil, types.ErrInvalidNftClassConfig.Wrapf("Royalty basis points cannot be greater than %s", types.MaxRoyaltyBasisPointsText)
	}
	royaltyAmount := msg.Price / 10000 * classData.Config.RoyaltyBasisPoints
	// pay royalty if needed, could be 0 if price < 10000
	if royaltyAmount > 0 {
		classParent, err := k.ValidateAndRefreshClassParent(ctx, msg.ClassId, classData.Parent)
		if err != nil {
			return nil, err
		}
		royaltyAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(royaltyAmount))))
		err = k.bankKeeper.SendCoins(ctx, buyerAddress, classParent.Owner, royaltyAmountCoins)
		if err != nil {
			return nil, types.ErrFailedToBuyNFT.Wrapf(err.Error())
		}
	}
	// pay seller
	netAmount := msg.Price - royaltyAmount
	netAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(netAmount))))
	err = k.bankKeeper.SendCoins(ctx, buyerAddress, sellerAddress, netAmountCoins)
	if err != nil {
		return nil, types.ErrFailedToBuyNFT.Wrapf(err.Error())
	}
	// sanity check
	if royaltyAmount+netAmount != msg.Price {
		return nil, types.ErrFailedToBuyNFT.Wrapf("Price split calculation error")
	}
	// transfer nft to buyer
	err = k.nftKeeper.Transfer(ctx, msg.ClassId, msg.NftId, buyerAddress)
	if err != nil {
		return nil, types.ErrFailedToBuyNFT.Wrapf(err.Error())
	}

	// remove listing
	k.RemoveListing(ctx, msg.ClassId, msg.NftId, sellerAddress)

	// owner changed, prune invalid listings
	k.PruneInvalidListingsForNFT(ctx, msg.ClassId, msg.NftId)

	// emit event
	ctx.EventManager().EmitTypedEvent(&types.EventBuyNFT{
		ClassId: msg.ClassId,
		NftId:   msg.NftId,
		Buyer:   buyerAddress.String(),
		Seller:  sellerAddress.String(),
		Price:   msg.Price,
	})

	return &types.MsgBuyNFTResponse{}, nil
}
