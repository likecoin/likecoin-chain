package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ListingKeyPrefix is the prefix to retrieve all Listing
	ListingKeyPrefix = "Listing/value/"
)

// ListingKey returns the store key to retrieve a Listing from the index fields
func ListingKey(
	classId string,
	nftId string,
	seller string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	nftIdBytes := []byte(nftId)
	key = append(key, nftIdBytes...)
	key = append(key, []byte("/")...)

	sellerBytes := []byte(seller)
	key = append(key, sellerBytes...)
	key = append(key, []byte("/")...)

	return key
}
