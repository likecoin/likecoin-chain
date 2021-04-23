package types

import (
	"encoding/binary"
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
)

// one fingerprint points to many sequence
// key structure:
//  - 4 bytes fingerprint length
//  - fpLen bytes fingerprint
//  - 8 bytes index
func GetFingerprintPrefix(fingerprint string) []byte {
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

func ParseFingerprintSequenceBytes(key []byte) (fingerprint string, index uint64) {
	fpLen := binary.BigEndian.Uint32(key[:4])
	fpBytes := key[4 : 4+fpLen]
	fingerprint = string(fpBytes)
	index = binary.BigEndian.Uint64(key[4+fpLen:])
	return fingerprint, index
}
