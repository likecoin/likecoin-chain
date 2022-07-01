package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

// SetClassesByISCN set a specific classesByISCN in the store from its index
func (k Keeper) SetClassesByISCN(ctx sdk.Context, classesByISCN types.ClassesByISCN) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassesByISCNKeyPrefix))
	b := k.cdc.MustMarshal(&classesByISCN)
	store.Set(types.ClassesByISCNKey(
		classesByISCN.IscnIdPrefix,
	), b)
}

// GetClassesByISCN returns a classesByISCN from its index
func (k Keeper) GetClassesByISCN(
	ctx sdk.Context,
	iscnIdPrefix string,

) (val types.ClassesByISCN, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassesByISCNKeyPrefix))

	b := store.Get(types.ClassesByISCNKey(
		iscnIdPrefix,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveClassesByISCN removes a classesByISCN from the store
func (k Keeper) RemoveClassesByISCN(
	ctx sdk.Context,
	iscnIdPrefix string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassesByISCNKeyPrefix))
	store.Delete(types.ClassesByISCNKey(
		iscnIdPrefix,
	))
}

// GetAllClassesByISCN returns all classesByISCN
func (k Keeper) GetAllClassesByISCN(ctx sdk.Context) (list []types.ClassesByISCN) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassesByISCNKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClassesByISCN
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
