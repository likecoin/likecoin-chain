package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// ListingKeyPrefix is the prefix to retrieve all Listing
	ListingKeyPrefix = "Listing/value/"
)

func ListingsByClassKey(
	classId string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func ListingsByNFTKey(
	classId string,
	nftId string,
) []byte {
	key := ListingsByClassKey(classId)

	nftIdBytes := []byte(nftId)
	key = append(key, nftIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

// ListingKey returns the store key to retrieve a Listing from the index fields
func ListingKey(
	classId string,
	nftId string,
	seller sdk.AccAddress,
) []byte {
	key := ListingsByNFTKey(classId, nftId)

	sellerBytes := seller
	key = append(key, sellerBytes...)
	key = append(key, []byte("/")...)

	return key
}
