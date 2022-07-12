package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	// ClassesByAccountKeyPrefix is the prefix to retrieve all ClassesByAccount
	ClassesByAccountKeyPrefix = "ClassesByAccount/value/"
)

// ClassesByAccountKey returns the store key to retrieve a ClassesByAccount from the index fields
func ClassesByAccountKey(
	account sdk.AccAddress,
) []byte {
	var key []byte

	key = append(key, account...)
	key = append(key, []byte("/")...)

	return key
}
