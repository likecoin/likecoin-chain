package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ClassRevealQueueKeyPrefix is the prefix to retrieve all ClassRevealQueue
	ClassRevealQueueKeyPrefix = "ClassRevealQueue/value/"
)

// ClassRevealQueueKey returns the store key to retrieve a ClassRevealQueue from the index fields
func ClassRevealQueueKey(
	revealTime string,
	classId string,
) []byte {
	var key []byte

	revealTimeBytes := []byte(revealTime)
	key = append(key, revealTimeBytes...)
	key = append(key, []byte("/")...)

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
