package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// SetListing set a specific listing in the store from its index
func (k Keeper) SetListing(ctx sdk.Context, listing types.Listing) {
	storeRecord := listing.ToStoreRecord()
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	b := k.cdc.MustMarshal(&storeRecord)
	store.Set(types.ListingKey(
		storeRecord.ClassId,
		storeRecord.NftId,
		storeRecord.Seller,
	), b)
}

// GetListing returns a listing from its index
func (k Keeper) GetListing(
	ctx sdk.Context,
	classId string,
	nftId string,
	seller sdk.AccAddress,

) (val types.Listing, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))

	b := store.Get(types.ListingKey(
		classId,
		nftId,
		seller,
	))
	if b == nil {
		return val, false
	}

	var storeRecord types.ListingStoreRecord
	k.cdc.MustUnmarshal(b, &storeRecord)
	return storeRecord.ToPublicRecord(), true
}

func (k Keeper) GetListingsByClass(
	ctx sdk.Context,
	classId string,
) (list []types.Listing) {
	k.IterateListingsByClass(ctx, classId, func(l types.Listing) {
		list = append(list, l)
	})

	return
}

func (k Keeper) IterateListingsByClass(
	ctx sdk.Context,
	classId string,
	callback func(types.Listing),
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.ListingsByClassKey(classId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ListingStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val.ToPublicRecord())
	}

	return
}

func (k Keeper) GetListingsByNFT(
	ctx sdk.Context,
	classId string,
	nftId string,
) (list []types.Listing) {
	k.IterateListingsByNFT(ctx, classId, nftId, func(l types.Listing) {
		list = append(list, l)
	})

	return
}

func (k Keeper) IterateListingsByNFT(
	ctx sdk.Context,
	classId string,
	nftId string,
	callback func(types.Listing),
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, types.ListingsByNFTKey(classId, nftId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ListingStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		callback(val.ToPublicRecord())
	}

	return
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
func (k Keeper) GetAllListing(ctx sdk.Context) (list []types.Listing) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ListingKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ListingStoreRecord
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val.ToPublicRecord())
	}

	return
}

func (k Keeper) PruneInvalidListingsForNFT(ctx sdk.Context, classId string, nftId string) {
	nftOwner := k.nftKeeper.GetOwner(ctx, classId, nftId)

	k.IterateListingsByNFT(ctx, classId, nftId, func(l types.Listing) {
		seller, err := sdk.AccAddressFromBech32(l.Seller)
		if err != nil {
			return
		}
		if !seller.Equals(nftOwner) {
			k.RemoveListing(ctx, l.ClassId, l.NftId, seller)
			// TODO dequeue listing as well
		}
	})
}
