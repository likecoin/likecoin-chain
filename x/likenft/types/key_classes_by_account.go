package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ClassesByAccountKeyPrefix is the prefix to retrieve all ClassesByAccount
	ClassesByAccountKeyPrefix = "ClassesByAccount/value/"
)

// ClassesByAccountKey returns the store key to retrieve a ClassesByAccount from the index fields
func ClassesByAccountKey(
	account string,
) []byte {
	var key []byte

	accountBytes := []byte(account)
	key = append(key, accountBytes...)
	key = append(key, []byte("/")...)

	return key
}
