package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetListing set a specific listing in the store from its index
func (k Keeper) SetListing(ctx sdk.Context, listing types.ListingStoreRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	b := k.cdc.MustMarshal(&listing)
	store.Set(types.ListingKey(
		listing.ClassId,
		listing.NftId,
		listing.Seller,
	), b)
}

// GetListing returns a listing from its index
func (k Keeper) GetListing(
	ctx sdk.Context,
	classId string,
	nftId string,
	seller sdk.AccAddress,

) (val types.ListingStoreRecord, found bool) {
	return k.GetListingByKeyBytes(ctx, types.ListingKey(
		classId,
		nftId,
		seller,
	))
}

func (k Keeper) GetListingByKeyBytes(
	ctx sdk.Context,
	key []byte,
) (val types.ListingStoreRecord, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))

	b := store.Get(key)
	if b == nil {
		return val, false
	}

	var storeRecord types.ListingStoreRecord
	k.cdc.MustUnmarshal(b, &storeRecord)
	return storeRecord, true
}

func (k Keeper) GetListingsByClass(
	ctx sdk.Context,
	classId string,
) (list []types.ListingStoreRecord) {
	k.IterateListingsByClass(ctx, classId, func(l types.ListingStoreRecord) {
		list = append(list, l)
	})

	return
}

func (k Keeper) IterateListingsByClass(
	ctx sdk.Context,
	classId string,
	callback func(types.ListingStoreRecord),
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.ListingsByClassKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ListingStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val)
	}
}

func (k Keeper) GetListingsByNFT(
	ctx sdk.Context,
	classId string,
	nftId string,
) (list []types.ListingStoreRecord) {
	k.IterateListingsByNFT(ctx, classId, nftId, func(l types.ListingStoreRecord) {
		list = append(list, l)
	})

	return
}

func (k Keeper) IterateListingsByNFT(
	ctx sdk.Context,
	classId string,
	nftId string,
	callback func(types.ListingStoreRecord),
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.ListingsByNFTKey(classId, nftId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ListingStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val)
	}
}

// RemoveListing removes a listing from the store
func (k Keeper) RemoveListing(
	ctx sdk.Context,
	classId string,
	nftId string,
	seller sdk.AccAddress,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	store.Delete(types.ListingKey(
		classId,
		nftId,
		seller,
	))
}

// GetAllListing returns all listing
func (k Keeper) GetAllListing(ctx sdk.Context) (list []types.ListingStoreRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ListingStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) PruneInvalidListingsForNFT(ctx sdk.Context, classId string, nftId string) {
	nftOwner := k.nftKeeper.GetOwner(ctx, classId, nftId)

	k.IterateListingsByNFT(ctx, classId, nftId, func(l types.ListingStoreRecord) {
		if !l.Seller.Equals(nftOwner) {
			k.RemoveListing(ctx, l.ClassId, l.NftId, l.Seller)
			// TODO dequeue listing as well
			ctx.EventManager().EmitTypedEvent(&types.EventDeleteListing{
				ClassId: l.ClassId,
				NftId:   l.NftId,
				Seller:  l.Seller.String(),
			})
		}
	})
}
