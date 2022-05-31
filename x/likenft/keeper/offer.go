package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetOffer set a specific offer in the store from its index
func (k Keeper) SetOffer(ctx sdk.Context, offer types.Offer) {
	storeRecord := offer.ToStoreRecord()
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferKeyPrefix))
	b := k.cdc.MustMarshal(&storeRecord)
	store.Set(types.OfferKey(
		storeRecord.ClassId,
		storeRecord.NftId,
		storeRecord.Buyer,
	), b)
}

// GetOffer returns a offer from its index
func (k Keeper) GetOffer(
	ctx sdk.Context,
	classId string,
	nftId string,
	buyer sdk.AccAddress,

) (val types.Offer, found bool) {
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
	return storeRecord.ToPublicRecord(), true
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
func (k Keeper) GetAllOffer(ctx sdk.Context) (list []types.Offer) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OfferKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OfferStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val.ToPublicRecord())
	}

	return
}
