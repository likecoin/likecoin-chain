package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// OfferKeyPrefix is the prefix to retrieve all Offer
	OfferKeyPrefix = "Offer/value/"
)

// OfferKey returns the store key to retrieve a Offer from the index fields
func OfferKey(
	classId string,
	nftId string,
	buyer string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	nftIdBytes := []byte(nftId)
	key = append(key, nftIdBytes...)
	key = append(key, []byte("/")...)

	buyerBytes := []byte(buyer)
	key = append(key, buyerBytes...)
	key = append(key, []byte("/")...)

	return key
}
