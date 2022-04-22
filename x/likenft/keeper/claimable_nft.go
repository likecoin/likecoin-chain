package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetClaimableNFT set a specific claimableNFT in the store from its index
func (k Keeper) SetClaimableNFT(ctx sdk.Context, claimableNFT types.ClaimableNFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimableNFTKeyPrefix))
	b := k.cdc.MustMarshal(&claimableNFT)
	store.Set(types.ClaimableNFTKey(
		claimableNFT.ClassId,
		claimableNFT.Id,
	), b)
}

// GetClaimableNFT returns a claimableNFT from its index
func (k Keeper) GetClaimableNFT(
	ctx sdk.Context,
	classId string,
	id string,

) (val types.ClaimableNFT, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimableNFTKeyPrefix))

	b := store.Get(types.ClaimableNFTKey(
		classId,
		id,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveClaimableNFT removes a claimableNFT from the store
func (k Keeper) RemoveClaimableNFT(
	ctx sdk.Context,
	classId string,
	id string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimableNFTKeyPrefix))
	store.Delete(types.ClaimableNFTKey(
		classId,
		id,
	))
}

// GetClaimableNFTs returns all claimableNFT of a class
func (k Keeper) GetClaimableNFTs(ctx sdk.Context, classId string) (list []types.ClaimableNFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimableNFTKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.ClaimableNFTsKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClaimableNFT
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllClaimableNFT returns all claimableNFT
func (k Keeper) GetAllClaimableNFT(ctx sdk.Context) (list []types.ClaimableNFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimableNFTKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClaimableNFT
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
