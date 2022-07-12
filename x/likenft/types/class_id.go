package types

import (
	"crypto/sha256"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const CLASS_ID_PREFIX = "likenft"

func newClassId(prefix []byte, serial int) (string, error) {
	data := append(prefix, []byte(strconv.Itoa(serial))...)
	hash := sha256.Sum256(data)

	classId, err := sdk.Bech32ifyAddressBytes(CLASS_ID_PREFIX, hash[:])
	if err != nil {
		return "", err
	}
	return classId, nil
}

func NewClassIdForISCN(iscnIdPrefix string, serial int) (string, error) {
	return newClassId([]byte(iscnIdPrefix), serial)
}

func NewClassIdForAccount(accAddress sdk.AccAddress, serial int) (string, error) {
	return newClassId(accAddress, serial)
}
