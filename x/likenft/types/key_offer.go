package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// OfferKeyPrefix is the prefix to retrieve all Offer
	OfferKeyPrefix = "Offer/value/"
)

func OffersByClassKey(
	classId string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func OffersByNFTKey(
	classId string,
	nftId string,
) []byte {
	key := OffersByClassKey(classId)

	nftIdBytes := []byte(nftId)
	key = append(key, nftIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

// OfferKey returns the store key to retrieve a Offer from the index fields
func OfferKey(
	classId string,
	nftId string,
	buyer sdk.AccAddress,
) []byte {
	key := OffersByNFTKey(classId, nftId)

	buyerBytes := buyer
	key = append(key, buyerBytes...)
	key = append(key, []byte("/")...)

	return key
}
