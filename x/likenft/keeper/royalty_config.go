package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetRoyaltyConfig set a specific royaltyConfigByClass in the store from its index
func (k Keeper) SetRoyaltyConfig(ctx sdk.Context, royaltyConfigByClass types.RoyaltyConfigByClass) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoyaltyConfigByClassKeyPrefix))
	b := k.cdc.MustMarshal(&royaltyConfigByClass)
	store.Set(types.RoyaltyConfigByClassKey(
		royaltyConfigByClass.ClassId,
	), b)
}

// GetRoyaltyConfig returns a royaltyConfigByClass from its index
func (k Keeper) GetRoyaltyConfig(
	ctx sdk.Context,
	classId string,

) (config types.RoyaltyConfig, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoyaltyConfigByClassKeyPrefix))

	var val types.RoyaltyConfigByClass
	b := store.Get(types.RoyaltyConfigByClassKey(
		classId,
	))
	if b == nil {
		return config, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val.RoyaltyConfig, true
}

// RemoveRoyaltyConfig removes a royaltyConfigByClass from the store
func (k Keeper) RemoveRoyaltyConfig(
	ctx sdk.Context,
	classId string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoyaltyConfigByClassKeyPrefix))
	store.Delete(types.RoyaltyConfigByClassKey(
		classId,
	))
}

// GetAllRoyaltyConfig returns all royaltyConfigByClass
func (k Keeper) GetAllRoyaltyConfig(ctx sdk.Context) (list []types.RoyaltyConfigByClass) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoyaltyConfigByClassKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RoyaltyConfigByClass
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
