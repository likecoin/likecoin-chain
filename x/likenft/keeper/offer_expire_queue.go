package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetOfferExpireQueueEntry set a specific offerExpireQueueEntry in the store from its index
func (k Keeper) SetOfferExpireQueueEntry(ctx sdk.Context, offerExpireQueueEntry types.OfferExpireQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferExpireQueueKeyPrefix))
	b := k.cdc.MustMarshal(&offerExpireQueueEntry)
	store.Set(types.OfferExpireQueueKey(
		offerExpireQueueEntry.ExpireTime,
		offerExpireQueueEntry.OfferKey,
	), b)
}

// GetOfferExpireQueueEntry returns a offerExpireQueueEntry from its index
func (k Keeper) GetOfferExpireQueueEntry(
	ctx sdk.Context,
	expireTime time.Time,
	offerKey []byte,
) (val types.OfferExpireQueueEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferExpireQueueKeyPrefix))

	b := store.Get(types.OfferExpireQueueKey(
		expireTime,
		offerKey,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOfferExpireQueueEntry removes a offerExpireQueueEntry from the store
func (k Keeper) RemoveOfferExpireQueueEntry(
	ctx sdk.Context,
	expireTime time.Time,
	offerKey []byte,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferExpireQueueKeyPrefix))
	store.Delete(types.OfferExpireQueueKey(
		expireTime,
		offerKey,
	))
}

func (k Keeper) UpdateOfferExpireQueueEntry(ctx sdk.Context, originalExpireTime time.Time, offerKey []byte, updatedExpireTime time.Time) {
	k.RemoveOfferExpireQueueEntry(ctx, originalExpireTime, offerKey)
	k.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: updatedExpireTime,
		OfferKey:   offerKey,
	})
}

func (k Keeper) OfferExpireQueueByTimeIterator(ctx sdk.Context, expireTime time.Time) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferExpireQueueKeyPrefix))
	iterator := store.Iterator(types.OfferExpireByTimeKey(time.Time{}), types.OfferExpireByTimeKey(expireTime))
	return iterator
}

func (k Keeper) IterateOfferExpireQueueByTime(ctx sdk.Context, endTime time.Time, cb func(val types.OfferExpireQueueEntry) (stop bool)) {
	iterator := k.OfferExpireQueueByTimeIterator(ctx, endTime)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OfferExpireQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if cb(val) {
			break
		}
	}
}

func (k Keeper) GetOfferExpireQueueByTime(ctx sdk.Context, endTime time.Time) (list []types.OfferExpireQueueEntry) {
	k.IterateOfferExpireQueueByTime(ctx, endTime, func(val types.OfferExpireQueueEntry) bool {
		list = append(list, val)
		return false
	})
	return
}

func (k Keeper) OfferExpireQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferExpireQueueKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	return iterator
}

func (k Keeper) IterateOfferExpireQueue(ctx sdk.Context, cb func(val types.OfferExpireQueueEntry) (stop bool)) {
	iterator := k.OfferExpireQueueIterator(ctx)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OfferExpireQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if cb(val) {
			break
		}
	}
}

// GetOfferExpireQueue returns all offerExpireQueueEntry
func (k Keeper) GetOfferExpireQueue(ctx sdk.Context) (list []types.OfferExpireQueueEntry) {
	k.IterateOfferExpireQueue(ctx, func(val types.OfferExpireQueueEntry) (stop bool) {
		list = append(list, val)
		return false
	})
	return
}
