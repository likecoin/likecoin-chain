package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// MintableNFTKeyPrefix is the prefix to retrieve all MintableNFT
	MintableNFTKeyPrefix = "MintableNFT/value/"
)

// MintableNFTsKey gets the first part of the MintableNFT key based on the classID
func MintableNFTsKey(
	classId string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

// MintableNFTKey returns the store key to retrieve a MintableNFT from the index fields
func MintableNFTKey(
	classId string,
	id string,
) []byte {
	key := MintableNFTsKey(classId)

	idBytes := []byte(id)
	key = append(key, idBytes...)
	key = append(key, []byte("/")...)

	return key
}
