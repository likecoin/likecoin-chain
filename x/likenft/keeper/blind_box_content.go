package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

// SetBlindBoxContent set a specific blind box content in the store from its index
func (k Keeper) SetBlindBoxContent(ctx sdk.Context, content types.BlindBoxContent) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlindBoxContentKeyPrefix))
	b := k.cdc.MustMarshal(&content)
	key := types.BlindBoxContentKey(
		content.ClassId,
		content.Id,
	)
	if !store.Has(key) {
		// new content, increment count
		if err := k.incBlindBoxContentCount(ctx, content.ClassId); err != nil {
			panic(fmt.Errorf("Failed to increase blind box content count: %s", err.Error()))
		}
	}
	store.Set(key, b)
}

// GetBlindBoxContent returns a blind box content from its index
func (k Keeper) GetBlindBoxContent(
	ctx sdk.Context,
	classId string,
	id string,

) (val types.BlindBoxContent, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlindBoxContentKeyPrefix))

	b := store.Get(types.BlindBoxContentKey(
		classId,
		id,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBlindBoxContent removes a blind box content from the store
func (k Keeper) RemoveBlindBoxContent(
	ctx sdk.Context,
	classId string,
	id string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlindBoxContentKeyPrefix))
	key := types.BlindBoxContentKey(
		classId,
		id,
	)
	if store.Has(key) {
		// remove existing blind box content, decrement count
		if err := k.decBlindBoxContentCount(ctx, classId); err != nil {
			panic(fmt.Errorf("Failed to decrease blind box content count: %s", err.Error()))
		}
	}
	store.Delete(key)
}

// RemoveBlindBoxContent removes a blind box content from the store
func (k Keeper) RemoveBlindBoxContents(
	ctx sdk.Context,
	classId string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlindBoxContentKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.BlindBoxContentsKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	// reset count to 0
	if err := k.setBlindBoxContentCount(ctx, classId, 0); err != nil {
		panic(fmt.Errorf("Failed to reset blind box content count: %s", err.Error()))
	}
}

// GetBlindBoxContents returns all blind box content of a class
func (k Keeper) GetBlindBoxContents(ctx sdk.Context, classId string) (list []types.BlindBoxContent) {
	k.IterateBlindBoxContents(ctx, classId, func(mn types.BlindBoxContent) {
		list = append(list, mn)
	})

	return
}

func (k Keeper) IterateBlindBoxContents(ctx sdk.Context, classId string, callback func(types.BlindBoxContent)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlindBoxContentKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.BlindBoxContentsKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.BlindBoxContent
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val)
	}
}

func (k Keeper) IterateAllBlindBoxContent(ctx sdk.Context, callback func(types.BlindBoxContent)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlindBoxContentKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.BlindBoxContent
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val)
	}
}

// GetAllBlindBoxContent returns all BlindBoxContent
func (k Keeper) GetAllBlindBoxContent(ctx sdk.Context) (list []types.BlindBoxContent) {
	k.IterateAllBlindBoxContent(ctx, func(val types.BlindBoxContent) {
		list = append(list, val)
	})
	return
}
