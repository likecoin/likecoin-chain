package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (l Listing) ToStoreRecord() ListingStoreRecord {
	seller, err := sdk.AccAddressFromBech32(l.Seller)
	if err != nil {
		panic(err)
	}

	return ListingStoreRecord{
		ClassId:    l.ClassId,
		NftId:      l.NftId,
		Seller:     seller,
		Price:      l.Price,
		Expiration: l.Expiration,
	}
}

func (r ListingStoreRecord) ToPublicRecord() Listing {
	return Listing{
		ClassId:    r.ClassId,
		NftId:      r.NftId,
		Seller:     sdk.AccAddress(r.Seller).String(),
		Price:      r.Price,
		Expiration: r.Expiration,
	}
}
