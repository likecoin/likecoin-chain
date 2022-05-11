package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetMintableNFT set a specific mintableNFT in the store from its index
func (k Keeper) SetMintableNFT(ctx sdk.Context, mintableNFT types.MintableNFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintableNFTKeyPrefix))
	b := k.cdc.MustMarshal(&mintableNFT)
	key := types.MintableNFTKey(
		mintableNFT.ClassId,
		mintableNFT.Id,
	)
	if !store.Has(key) {
		// new mintable, increment count
		if err := k.incMintableCount(ctx, mintableNFT.ClassId); err != nil {
			panic(fmt.Errorf("Failed to increase mintable count: %s", err.Error()))
		}
	}
	store.Set(key, b)
}

// GetMintableNFT returns a mintableNFT from its index
func (k Keeper) GetMintableNFT(
	ctx sdk.Context,
	classId string,
	id string,

) (val types.MintableNFT, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintableNFTKeyPrefix))

	b := store.Get(types.MintableNFTKey(
		classId,
		id,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMintableNFT removes a mintableNFT from the store
func (k Keeper) RemoveMintableNFT(
	ctx sdk.Context,
	classId string,
	id string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintableNFTKeyPrefix))
	key := types.MintableNFTKey(
		classId,
		id,
	)
	if store.Has(key) {
		// remove existing mintable, decrement count
		if err := k.decMintableCount(ctx, classId); err != nil {
			panic(fmt.Errorf("Failed to decrease mintable count: %s", err.Error()))
		}
	}
	store.Delete(key)
}

// RemoveMintableNFT removes a mintableNFT from the store
func (k Keeper) RemoveMintableNFTs(
	ctx sdk.Context,
	classId string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintableNFTKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.MintableNFTsKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	// reset count to 0
	if err := k.setMintableCount(ctx, classId, 0); err != nil {
		panic(fmt.Errorf("Failed to reset mintable count: %s", err.Error()))
	}
}

// GetMintableNFTs returns all mintableNFT of a class
func (k Keeper) GetMintableNFTs(ctx sdk.Context, classId string) (list []types.MintableNFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintableNFTKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.MintableNFTsKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MintableNFT
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) IterateMintableNFTs(ctx sdk.Context, classId string, callback func(types.MintableNFT)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintableNFTKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.MintableNFTsKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MintableNFT
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val)
	}
}

func (k Keeper) IterateAllMintableNFT(ctx sdk.Context, callback func(types.MintableNFT)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintableNFTKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MintableNFT
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val)
	}
}

// GetAllMintableNFT returns all MintableNFT
func (k Keeper) GetAllMintableNFT(ctx sdk.Context) (list []types.MintableNFT) {
	k.IterateAllMintableNFT(ctx, func(val types.MintableNFT) {
		list = append(list, val)
	})
	return
}
