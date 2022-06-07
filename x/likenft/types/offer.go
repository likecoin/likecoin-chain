package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (o Offer) ToStoreRecord() OfferStoreRecord {
	buyer, err := sdk.AccAddressFromBech32(o.Buyer)
	if err != nil {
		panic(err)
	}

	return OfferStoreRecord{
		ClassId:    o.ClassId,
		NftId:      o.NftId,
		Buyer:      buyer,
		Price:      o.Price,
		Expiration: o.Expiration,
	}
}

func (r OfferStoreRecord) ToPublicRecord() Offer {
	return Offer{
		ClassId:    r.ClassId,
		NftId:      r.NftId,
		Buyer:      sdk.AccAddress(r.Buyer).String(),
		Price:      r.Price,
		Expiration: r.Expiration,
	}
}

func MapOffersToStoreRecords(offers []Offer) (records []OfferStoreRecord) {
	for _, offer := range offers {
		records = append(records, offer.ToStoreRecord())
	}
	return
}

func MapOffersToPublicRecords(records []OfferStoreRecord) (offers []Offer) {
	for _, record := range records {
		offers = append(offers, record.ToPublicRecord())
	}
	return
}
