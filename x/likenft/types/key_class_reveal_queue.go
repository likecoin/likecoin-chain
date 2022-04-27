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

// ClassRevealQueueKey returns the store key to retrieve a ClassRevealQueueEntry from the index fields
func ClassRevealQueueKey(
	revealTime time.Time,
	classId string,
) []byte {
	var key []byte

	revealTimeBytes := sdk.FormatTimeBytes(revealTime)
	key = append(key, revealTimeBytes...)
	key = append(key, []byte("/")...)

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
