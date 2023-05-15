package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

// SetClassRevealQueueEntry set a specific classRevealQueueEntry in the store from its index
func (k Keeper) SetClassRevealQueueEntry(ctx sdk.Context, classRevealQueueEntry types.ClassRevealQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	b := k.cdc.MustMarshal(&classRevealQueueEntry)
	store.Set(types.ClassRevealQueueKey(
		classRevealQueueEntry.RevealTime,
		classRevealQueueEntry.ClassId,
	), b)
}

// RemoveClassRevealQueueEntry removes a classRevealQueueEntry from the store
func (k Keeper) RemoveClassRevealQueueEntry(
	ctx sdk.Context,
	revealTime time.Time,
	classId string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	store.Delete(types.ClassRevealQueueKey(
		revealTime,
		classId,
	))
}

// UpdateClassRevealQueueEntry updates a classRevealQueueEntry in the store
func (k Keeper) UpdateClassRevealQueueEntry(ctx sdk.Context, originalRevealTime time.Time, classId string, updatedRevealTime time.Time) {
	k.RemoveClassRevealQueueEntry(ctx, originalRevealTime, classId)
	k.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: updatedRevealTime,
		ClassId:    classId,
	})
}

func (k Keeper) ClassRevealQueueByTimeIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	iterator := store.Iterator(types.ClassRevealByTimeKey(time.Time{}), types.ClassRevealByTimeKey(endTime))
	return iterator
}

func (k Keeper) IterateClassRevealQueueByTime(ctx sdk.Context, endTime time.Time, cb func(val types.ClassRevealQueueEntry) (stop bool)) {
	iterator := k.ClassRevealQueueByTimeIterator(ctx, endTime)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClassRevealQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if cb(val) {
			break
		}
	}

	return
}

func (k Keeper) GetClassRevealQueueByTime(ctx sdk.Context, endTime time.Time) (list []types.ClassRevealQueueEntry) {
	k.IterateClassRevealQueueByTime(ctx, endTime, func(val types.ClassRevealQueueEntry) (stop bool) {
		list = append(list, val)
		return false
	})
	return
}

func (k Keeper) ClassRevealQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	return iterator
}

// IterateClassRevealQueue iterates over all classRevealQueueEntry
func (k Keeper) IterateClassRevealQueue(ctx sdk.Context, cb func(val types.ClassRevealQueueEntry) (stop bool)) {
	iterator := k.ClassRevealQueueIterator(ctx)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClassRevealQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if cb(val) {
			break
		}
	}

	return
}

// GetClassRevealQueue returns all classRevealQueueEntry
func (k Keeper) GetClassRevealQueue(ctx sdk.Context) (list []types.ClassRevealQueueEntry) {
	k.IterateClassRevealQueue(ctx, func(val types.ClassRevealQueueEntry) (stop bool) {
		list = append(list, val)
		return false
	})
	return
}
