package types

import (
	"encoding/binary"
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// ListingExpireQueueKeyPrefix is the prefix to retrieve all ListingExpireQueueEntry
	ListingExpireQueueKeyPrefix = "ListingExpireQueueEntry/value/"
)

func ListingExpireByTimeKey(
	expireTime time.Time,
) []byte {
	var key []byte
	expireTimeBytes := sdk.FormatTimeBytes(expireTime)
	key = append(key, expireTimeBytes...)
	key = append(key, []byte("/")...)

	return key
}

// ListingExpireQueueKey returns the store key to retrieve a ListingExpireQueueEntry from the index fields
func ListingExpireQueueKey(
	expireTime time.Time,
	listingKey []byte,
) []byte {
	key := ListingExpireByTimeKey(expireTime)

	listingKeyBytes := listingKey
	key = append(key, listingKeyBytes...)
	key = append(key, []byte("/")...)

	return key
}
