package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func (k Keeper) ExpireOffer(ctx sdk.Context, offer types.OfferStoreRecord) error {
	// Check offer is actually expired
	if !offer.Expiration.Before(ctx.BlockTime()) {
		return types.ErrFailedToExpireOffer.Wrap("Offer is not expired on record")
	}

	// Refund deposit if needed
	if offer.Price > 0 {
		denom := k.PriceDenom(ctx)
		coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(int64(offer.Price))))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, offer.Buyer, coins); err != nil {
			return types.ErrFailedToExpireOffer.Wrapf(err.Error())
		}
	}

	// Delete offer
	k.RemoveOffer(ctx, offer.ClassId, offer.NftId, offer.Buyer)

	return nil
}
