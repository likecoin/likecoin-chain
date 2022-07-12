package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ClassesByISCNKeyPrefix is the prefix to retrieve all ClassesByISCN
	ClassesByISCNKeyPrefix = "ClassesByISCN/value/"
)

// ClassesByISCNKey returns the store key to retrieve a ClassesByISCN from the index fields
func ClassesByISCNKey(
	iscnIdPrefix string,
) []byte {
	var key []byte

	iscnIdPrefixBytes := []byte(iscnIdPrefix)
	key = append(key, iscnIdPrefixBytes...)
	key = append(key, []byte("/")...)

	return key
}
