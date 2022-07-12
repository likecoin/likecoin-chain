package types

import (
	"encoding/binary"
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// OfferExpireQueueKeyPrefix is the prefix to retrieve all OfferExpireQueueEntry
	OfferExpireQueueKeyPrefix = "OfferExpireQueueEntry/value/"
)

func OfferExpireByTimeKey(
	expireTime time.Time,
) []byte {
	var key []byte
	expireTimeBytes := sdk.FormatTimeBytes(expireTime)
	key = append(key, expireTimeBytes...)
	key = append(key, []byte("/")...)

	return key
}

// OfferExpireQueueKey returns the store key to retrieve a OfferExpireQueueEntry from the index fields
func OfferExpireQueueKey(
	expireTime time.Time,
	offerKey []byte,
) []byte {
	key := OfferExpireByTimeKey(expireTime)

	offerKeyBytes := offerKey
	key = append(key, offerKeyBytes...)
	key = append(key, []byte("/")...)

	return key
}
