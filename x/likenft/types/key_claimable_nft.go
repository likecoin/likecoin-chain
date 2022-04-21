package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ClaimableNFTKeyPrefix is the prefix to retrieve all ClaimableNFT
	ClaimableNFTKeyPrefix = "ClaimableNFT/value/"
)

// ClaimableNFTKey returns the store key to retrieve a ClaimableNFT from the index fields
func ClaimableNFTKey(
	classId string,
	id string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	idBytes := []byte(id)
	key = append(key, idBytes...)
	key = append(key, []byte("/")...)

	return key
}
