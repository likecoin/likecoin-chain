package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// BlindBoxContentKeyPrefix is the prefix to retrieve all BlindBoxContent
	BlindBoxContentKeyPrefix = "BlindBoxContent/value/"
)

// BlindBoxContentsKey gets the first part of the BlindBoxContent key based on the classID
func BlindBoxContentsKey(
	classId string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

// BlindBoxContentKey returns the store key to retrieve a BlindBoxContent from the index fields
func BlindBoxContentKey(
	classId string,
	id string,
) []byte {
	key := BlindBoxContentsKey(classId)

	idBytes := []byte(id)
	key = append(key, idBytes...)
	key = append(key, []byte("/")...)

	return key
}
