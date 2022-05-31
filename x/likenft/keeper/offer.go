package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetOffer set a specific offer in the store from its index
func (k Keeper) SetOffer(ctx sdk.Context, offer types.Offer) {
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
	buyer string,

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

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOffer removes a offer from the store
func (k Keeper) RemoveOffer(
	ctx sdk.Context,
	classId string,
	nftId string,
	buyer string,

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
		var val types.Offer
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
