package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetOffer set a specific offer in the store from its index
func (k Keeper) SetOffer(ctx sdk.Context, offer types.OfferStoreRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferKeyPrefix))
	b := k.cdc.MustMarshal(&offer)
	store.Set(types.OfferKey(
		offer.ClassId,
		offer.NftId,
		offer.Buyer,
	), b)
}

// GetOffer returns a offer from its index
func (k Keeper) GetOffer(
	ctx sdk.Context,
	classId string,
	nftId string,
	buyer sdk.AccAddress,

) (val types.OfferStoreRecord, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferKeyPrefix))

	b := store.Get(types.OfferKey(
		classId,
		nftId,
		buyer,
	))
	if b == nil {
		return val, false
	}

	var storeRecord types.OfferStoreRecord
	k.cdc.MustUnmarshal(b, &storeRecord)
	return storeRecord, true
}

func (k Keeper) GetOffersByClass(
	ctx sdk.Context,
	classId string,
) (list []types.OfferStoreRecord) {
	k.IterateOffersByClass(ctx, classId, func(o types.OfferStoreRecord) {
		list = append(list, o)
	})

	return
}

func (k Keeper) IterateOffersByClass(
	ctx sdk.Context,
	classId string,
	callback func(types.OfferStoreRecord),
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.OffersByClassKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OfferStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val)
	}

	return
}

func (k Keeper) GetOffersByNFT(
	ctx sdk.Context,
	classId string,
	nftId string,
) (list []types.OfferStoreRecord) {
	k.IterateOffersByNFT(ctx, classId, nftId, func(o types.OfferStoreRecord) {
		list = append(list, o)
	})

	return
}

func (k Keeper) IterateOffersByNFT(
	ctx sdk.Context,
	classId string,
	nftId string,
	callback func(types.OfferStoreRecord),
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.OffersByNFTKey(classId, nftId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OfferStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val)
	}

	return
}

// RemoveOffer removes a offer from the store
func (k Keeper) RemoveOffer(
	ctx sdk.Context,
	classId string,
	nftId string,
	buyer sdk.AccAddress,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferKeyPrefix))
	store.Delete(types.OfferKey(
		classId,
		nftId,
		buyer,
	))
}

// GetAllOffer returns all offer
func (k Keeper) GetAllOffer(ctx sdk.Context) (list []types.OfferStoreRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OfferStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
