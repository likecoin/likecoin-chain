package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetClassRevealQueue set a specific classRevealQueue in the store from its index
func (k Keeper) SetClassRevealQueue(ctx sdk.Context, classRevealQueue types.ClassRevealQueue) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	b := k.cdc.MustMarshal(&classRevealQueue)
	store.Set(types.ClassRevealQueueKey(
		classRevealQueue.RevealTime,
		classRevealQueue.ClassId,
	), b)
}

// RemoveClassRevealQueue removes a classRevealQueue from the store
func (k Keeper) RemoveFromClassRevealQueue(
	ctx sdk.Context,
	revealTime time.Time,
	classId string,

) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	if !store.Has(types.ClassRevealQueueKey(revealTime, classId)) {
		return sdkerrors.ErrKeyNotFound.Wrapf("classRevealQueue entry not found: classId: %s, revealTime: %s", classId, revealTime)
	}
	store.Delete(types.ClassRevealQueueKey(
		revealTime,
		classId,
	))
	return nil
}

// UpdateClassRevealQueue updates a classRevealQueue in the store
func (k Keeper) UpdateClassRevealQueue(ctx sdk.Context, originalRevealTime time.Time, classId string, updatedRevealTime time.Time) error {
	err := k.RemoveFromClassRevealQueue(ctx, originalRevealTime, classId)
	if err != nil {
		return err
	}
	k.SetClassRevealQueue(ctx, types.ClassRevealQueue{
		RevealTime: updatedRevealTime,
		ClassId:    classId,
	})
	return nil
}

func (k Keeper) ClassRevealQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	return iterator
}

// IterateClassRevealQueue iterates over all classRevealQueue
func (k Keeper) IterateClassRevealQueue(ctx sdk.Context, cb func(val types.ClassRevealQueue) (stop bool)) {
	iterator := k.ClassRevealQueueIterator(ctx)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClassRevealQueue
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		if cb(val) {
			break
		}
	}

	return
}

// GetAllClassRevealQueue returns all classRevealQueue
func (k Keeper) GetAllClassRevealQueue(ctx sdk.Context) (list []types.ClassRevealQueue) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClassRevealQueue
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
