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

func MapListingsToStoreRecords(listings []Listing) (records []ListingStoreRecord) {
	for _, listing := range listings {
		records = append(records, listing.ToStoreRecord())
	}
	return
}

func MapListingsToPublicRecords(records []ListingStoreRecord) (listings []Listing) {
	for _, record := range records {
		listings = append(listings, record.ToPublicRecord())
	}
	return
}
