package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// ClassRevealQueueKeyPrefix is the prefix to retrieve all ClassRevealQueueEntry
	ClassRevealQueueKeyPrefix = "ClassRevealQueue/value/"
)

func ClassRevealByTimeKey(
	revealTime time.Time,
) []byte {
	var key []byte

	revealTimeBytes := sdk.FormatTimeBytes(revealTime)
	key = append(key, revealTimeBytes...)
	key = append(key, []byte("/")...)

	return key
}

// ClassRevealQueueKey returns the store key to retrieve a ClassRevealQueueEntry from the index fields
func ClassRevealQueueKey(
	revealTime time.Time,
	classId string,
) []byte {
	key := ClassRevealByTimeKey(revealTime)

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
