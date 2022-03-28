package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetClassesByAccount set a specific classesByAccount in the store from its index
func (k Keeper) SetClassesByAccount(ctx sdk.Context, classesByAccount types.ClassesByAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassesByAccountKeyPrefix))
	storeRecord := classesByAccount.ToStoreRecord()
	b := k.cdc.MustMarshal(&storeRecord)
	store.Set(types.ClassesByAccountKey(
		storeRecord.AccAddress,
	), b)
}

// GetClassesByAccount returns a classesByAccount from its index
func (k Keeper) GetClassesByAccount(
	ctx sdk.Context,
	account sdk.AccAddress,

) (val types.ClassesByAccount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassesByAccountKeyPrefix))

	b := store.Get(types.ClassesByAccountKey(
		account,
	))
	if b == nil {
		return val, false
	}

	var storeRecord types.ClassesByAccountStoreRecord
	k.cdc.MustUnmarshal(b, &storeRecord)
	return storeRecord.ToPublicRecord(), true
}

// RemoveClassesByAccount removes a classesByAccount from the store
func (k Keeper) RemoveClassesByAccount(
	ctx sdk.Context,
	account sdk.AccAddress,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassesByAccountKeyPrefix))
	store.Delete(types.ClassesByAccountKey(
		account,
	))
}

// GetAllClassesByAccount returns all classesByAccount
func (k Keeper) GetAllClassesByAccount(ctx sdk.Context) (list []types.ClassesByAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClassesByAccountKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClassesByAccountStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val.ToPublicRecord())
	}

	return
}
