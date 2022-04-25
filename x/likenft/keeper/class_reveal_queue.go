package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// GetClassRevealQueue returns a classRevealQueue from its index
func (k Keeper) GetClassRevealQueue(
	ctx sdk.Context,
	revealTime string,
	classId string,

) (val types.ClassRevealQueue, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))

	b := store.Get(types.ClassRevealQueueKey(
		revealTime,
		classId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveClassRevealQueue removes a classRevealQueue from the store
func (k Keeper) RemoveClassRevealQueue(
	ctx sdk.Context,
	revealTime string,
	classId string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassRevealQueueKeyPrefix))
	store.Delete(types.ClassRevealQueueKey(
		revealTime,
		classId,
	))
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
