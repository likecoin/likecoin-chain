package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ClaimableNFTKeyPrefix is the prefix to retrieve all ClaimableNFT
	ClaimableNFTKeyPrefix = "ClaimableNFT/value/"
)

// ClaimableNFTsKey gets the first part of the ClaimableNFT key based on the classID
func ClaimableNFTsKey(
	classId string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

// ClaimableNFTKey returns the store key to retrieve a ClaimableNFT from the index fields
func ClaimableNFTKey(
	classId string,
	id string,
) []byte {
	key := ClaimableNFTsKey(classId)

	idBytes := []byte(id)
	key = append(key, idBytes...)
	key = append(key, []byte("/")...)

	return key
}
