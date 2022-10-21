package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
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
	royaltyConfig, found := k.GetRoyaltyConfig(ctx, msg.ClassId)
	var royaltyAmount uint64
	if found {
		_royaltyAmount, allocations, err := k.ComputeRoyaltyAllocation(ctx, msg.Price, listing.FullPayToRoyalty, royaltyConfig)
		if err != nil {
			return nil, err
		}
		royaltyAmount = _royaltyAmount
		for _, allocation := range allocations {
			coins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(allocation.Amount))))
			err = k.bankKeeper.SendCoins(ctx, buyerAddress, allocation.Account, coins)
			if err != nil {
				return nil, types.ErrFailedToBuyNFT.Wrapf(err.Error())
			}
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

	// owner changed, remove all listings
	k.PruneAllListingsForNFT(ctx, msg.ClassId, msg.NftId)

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
