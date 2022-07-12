package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func (k Keeper) ExpireListing(ctx sdk.Context, listing types.ListingStoreRecord) error {
	// Check listing is actually expired
	if !listing.Expiration.Before(ctx.BlockTime()) {
		return types.ErrFailedToExpireListing.Wrap("Listing is not expired on record")
	}

	// Delete offer
	k.RemoveListing(ctx, listing.ClassId, listing.NftId, listing.Seller)

	return nil
}
