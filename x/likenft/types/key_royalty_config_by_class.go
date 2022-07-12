package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// RoyaltyConfigByClassKeyPrefix is the prefix to retrieve all RoyaltyConfigByClass
	RoyaltyConfigByClassKeyPrefix = "RoyaltyConfigByClass/value/"
)

// RoyaltyConfigByClassKey returns the store key to retrieve a RoyaltyConfigByClass from the index fields
func RoyaltyConfigByClassKey(
	classId string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
