package types

import (
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const CLASS_ID_PREFIX = "likenft"

// deterministically create bech32 style class id
func NewClassId(iscnIdPrefix string, serial int) (*string, error) {
	data := fmt.Sprintf("%s:%d", iscnIdPrefix, serial)
	hash := sha256.Sum256([]byte(data))

	classId, err := sdk.Bech32ifyAddressBytes(CLASS_ID_PREFIX, hash[:])
	if err != nil {
		return nil, err
	}
	return &classId, nil
}
