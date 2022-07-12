package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

// SetListingExpireQueueEntry set a specific listingExpireQueueEntry in the store from its index
func (k Keeper) SetListingExpireQueueEntry(ctx sdk.Context, listingExpireQueueEntry types.ListingExpireQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingExpireQueueKeyPrefix))
	b := k.cdc.MustMarshal(&listingExpireQueueEntry)
	store.Set(types.ListingExpireQueueKey(
		listingExpireQueueEntry.ExpireTime,
		listingExpireQueueEntry.ListingKey,
	), b)
}

// GetListingExpireQueueEntry returns a listingExpireQueueEntry from its index
func (k Keeper) GetListingExpireQueueEntry(
	ctx sdk.Context,
	expireTime time.Time,
	listingKey []byte,
) (val types.ListingExpireQueueEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingExpireQueueKeyPrefix))

	b := store.Get(types.ListingExpireQueueKey(
		expireTime,
		listingKey,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveListingExpireQueueEntry removes a listingExpireQueueEntry from the store
func (k Keeper) RemoveListingExpireQueueEntry(
	ctx sdk.Context,
	expireTime time.Time,
	listingKey []byte,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingExpireQueueKeyPrefix))
	store.Delete(types.ListingExpireQueueKey(
		expireTime,
		listingKey,
	))
}

func (k Keeper) UpdateListingExpireQueueEntry(ctx sdk.Context, originalExpireTime time.Time, listingKey []byte, updatedExpireTime time.Time) {
	k.RemoveListingExpireQueueEntry(ctx, originalExpireTime, listingKey)
	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: updatedExpireTime,
		ListingKey: listingKey,
	})
}

func (k Keeper) ListingExpireQueueByTimeIterator(ctx sdk.Context, expireTime time.Time) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingExpireQueueKeyPrefix))
	iterator := store.Iterator(types.ListingExpireByTimeKey(time.Time{}), types.ListingExpireByTimeKey(expireTime))
	return iterator
}

func (k Keeper) IterateListingExpireQueueByTime(ctx sdk.Context, endTime time.Time, cb func(val types.ListingExpireQueueEntry) (stop bool)) {
	iterator := k.ListingExpireQueueByTimeIterator(ctx, endTime)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ListingExpireQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if cb(val) {
			break
		}
	}
}

func (k Keeper) GetListingExpireQueueByTime(ctx sdk.Context, endTime time.Time) (list []types.ListingExpireQueueEntry) {
	k.IterateListingExpireQueueByTime(ctx, endTime, func(val types.ListingExpireQueueEntry) bool {
		list = append(list, val)
		return false
	})
	return
}

func (k Keeper) ListingExpireQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingExpireQueueKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	return iterator
}

func (k Keeper) IterateListingExpireQueue(ctx sdk.Context, cb func(val types.ListingExpireQueueEntry) (stop bool)) {
	iterator := k.ListingExpireQueueIterator(ctx)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ListingExpireQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if cb(val) {
			break
		}
	}
}

// GetListingExpireQueue returns all listingExpireQueueEntry
func (k Keeper) GetListingExpireQueue(ctx sdk.Context) (list []types.ListingExpireQueueEntry) {
	k.IterateListingExpireQueue(ctx, func(val types.ListingExpireQueueEntry) (stop bool) {
		list = append(list, val)
		return false
	})
	return
}
