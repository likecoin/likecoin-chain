package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "iscn"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
)

var (
	SequenceCountKey            = []byte{0x01}
	SequenceToStoreRecordPrefix = []byte{0x02}
	CidToSequencePrefix         = []byte{0x03}
	IscnIdToSequencePrefix      = []byte{0x04}
	ContentIdRecordPrefix       = []byte{0x05}
	FingerprintSequencePrefix   = []byte{0x06}
	OwnerSequencePrefix         = []byte{0x07}
)

// one fingerprint points to many sequence
// key structure:
//  - 4 bytes fingerprint length
//  - fpLen bytes fingerprint
//  - 8 bytes sequence
func GetFingerprintStorePrefix(fingerprint string) []byte {
	fpBytes := []byte(fingerprint)
	fpLen := len(fpBytes)
	output := make([]byte, len(FingerprintSequencePrefix)+4+fpLen)
	copy(output, FingerprintSequencePrefix)
	binary.BigEndian.PutUint32(output[len(FingerprintSequencePrefix):], uint32(fpLen))
	copy(output[len(FingerprintSequencePrefix)+4:], fpBytes)
	return output
}

func GetFingerprintSequenceKey(fingerprint string, seq uint64) []byte {
	fpBytes := []byte(fingerprint)
	fpLen := len(fpBytes)
	output := make([]byte, 4+fpLen+8)
	binary.BigEndian.PutUint32(output, uint32(fpLen))
	copy(output[4:], fpBytes)
	binary.BigEndian.PutUint64(output[4+fpLen:], seq)
	return output
}

func ParseFingerprintSequenceBytes(key []byte) (fingerprint string, seq uint64) {
	fpLen := binary.BigEndian.Uint32(key[:4])
	fpBytes := key[4 : 4+fpLen]
	fingerprint = string(fpBytes)
	seq = binary.BigEndian.Uint64(key[4+fpLen:])
	return fingerprint, seq
}

// one owner points to many sequence
// key structure:
//  - 4 bytes owner address bytes length
//  - addrLen bytes owner address
//  - 8 bytes sequence
func GetOwnerStorePrefix(owner sdk.AccAddress) []byte {
	addrBytes := []byte(owner)
	addrLen := len(addrBytes)
	output := make([]byte, len(OwnerSequencePrefix)+4+addrLen)
	copy(output, OwnerSequencePrefix)
	binary.BigEndian.PutUint32(output[len(OwnerSequencePrefix):], uint32(addrLen))
	copy(output[len(OwnerSequencePrefix)+4:], addrBytes)
	return output
}

func GetOwnerSequenceKey(owner sdk.AccAddress, seq uint64) []byte {
	addrBytes := []byte(owner)
	addrLen := len(addrBytes)
	output := make([]byte, 4+addrLen+8)
	binary.BigEndian.PutUint32(output, uint32(addrLen))
	copy(output[4:], addrBytes)
	binary.BigEndian.PutUint64(output[4+addrLen:], seq)
	return output
}

func ParseOwnerSequenceBytes(key []byte) (owner sdk.AccAddress, seq uint64) {
	addrLen := binary.BigEndian.Uint32(key[:4])
	addrBytes := key[4 : 4+addrLen]
	owner = sdk.AccAddress(addrBytes)
	seq = binary.BigEndian.Uint64(key[4+addrLen:])
	return owner, seq
}
