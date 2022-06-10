package likenft

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
)

func tryRevealClassCatchPanic(ctx sdk.Context, keeper keeper.Keeper, classId string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	err = keeper.RevealMintableNFTs(ctx, classId)
	return
}

func processClassRevealQueue(ctx sdk.Context, keeper keeper.Keeper) {
	// Reveal classes with reveal time < current block header time
	keeper.IterateClassRevealQueueByTime(ctx, ctx.BlockHeader().Time, func(entry types.ClassRevealQueueEntry) (stop bool) {
		err := tryRevealClassCatchPanic(ctx, keeper, entry.ClassId)

		if err != nil {
			ctx.EventManager().EmitTypedEvent(&types.EventRevealClass{
				ClassId: entry.ClassId,
				Success: false,
				Error:   err.Error(),
			})
		} else {
			ctx.EventManager().EmitTypedEvent(&types.EventRevealClass{
				ClassId: entry.ClassId,
				Success: true,
			})
		}

		keeper.RemoveClassRevealQueueEntry(ctx, entry.RevealTime, entry.ClassId)
		return false
	})
}

func tryExpireOfferCatchPanic(ctx sdk.Context, keeper keeper.Keeper, offer types.OfferStoreRecord) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	err = keeper.ExpireOffer(ctx, offer)
	return
}

func processOfferExpireQueue(ctx sdk.Context, keeper keeper.Keeper) {
	// Expire offers with expiration time < current block header time
	keeper.IterateOfferExpireQueueByTime(ctx, ctx.BlockHeader().Time, func(val types.OfferExpireQueueEntry) (stop bool) {
		// Get offer
		offer, found := keeper.GetOfferByKeyBytes(ctx, val.OfferKey)
		if !found {
			// offer not found, dequeue and continue
			keeper.RemoveOfferExpireQueueEntry(ctx, val.ExpireTime, val.OfferKey)
			return false
		}

		err := tryExpireOfferCatchPanic(ctx, keeper, offer)
		if err != nil {
			ctx.EventManager().EmitTypedEvent(&types.EventExpireOffer{
				ClassId: offer.ClassId,
				NftId:   offer.NftId,
				Buyer:   offer.Buyer.String(),
				Success: false,
				Error:   err.Error(),
			})
		} else {
			ctx.EventManager().EmitTypedEvent(&types.EventExpireOffer{
				ClassId: offer.ClassId,
				NftId:   offer.NftId,
				Buyer:   offer.Buyer.String(),
				Success: true,
			})
		}

		keeper.RemoveOfferExpireQueueEntry(ctx, val.ExpireTime, val.OfferKey)
		return false
	})
}

func tryExpireListingCatchPanic(ctx sdk.Context, keeper keeper.Keeper, listing types.ListingStoreRecord) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	err = keeper.ExpireListing(ctx, listing)
	return
}

func processListingExpireQueue(ctx sdk.Context, keeper keeper.Keeper) {
	// Expire listings with expiration time < current block header time
	keeper.IterateListingExpireQueueByTime(ctx, ctx.BlockHeader().Time, func(val types.ListingExpireQueueEntry) (stop bool) {
		// Get listing
		listing, found := keeper.GetListingByKeyBytes(ctx, val.ListingKey)
		if !found {
			// listing not found, dequeue and continue
			keeper.RemoveListingExpireQueueEntry(ctx, val.ExpireTime, val.ListingKey)
			return false
		}

		err := tryExpireListingCatchPanic(ctx, keeper, listing)
		if err != nil {
			ctx.EventManager().EmitTypedEvent(&types.EventExpireListing{
				ClassId: listing.ClassId,
				NftId:   listing.NftId,
				Seller:  listing.Seller.String(),
				Success: false,
				Error:   err.Error(),
			})
		} else {
			ctx.EventManager().EmitTypedEvent(&types.EventExpireListing{
				ClassId: listing.ClassId,
				NftId:   listing.NftId,
				Seller:  listing.Seller.String(),
				Success: true,
			})
		}

		keeper.RemoveOfferExpireQueueEntry(ctx, val.ExpireTime, val.ListingKey)
		return false
	})
}

// EndBlocker called every block, process class reveal queue.
func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	processClassRevealQueue(ctx, keeper)
	processOfferExpireQueue(ctx, keeper)
	processListingExpireQueue(ctx, keeper)
}
