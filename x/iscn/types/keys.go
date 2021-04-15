package types

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	ModuleName   = "iscn"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
)

var (
	CidBlockKey         = []byte{0x01}
	IscnIdToCidKey      = []byte{0x02}
	IscnCountKey        = []byte{0x03}
	CidToIscnIdKey      = []byte{0x04}
	FingerprintToCidKey = []byte{0x05}
	IscnIdOwnerKey      = []byte{0x06}
	IscnIdVersionKey    = []byte{0x07}
)

func GetIscnIdToCidKey(cdc codec.BinaryMarshaler, iscnId IscnId) []byte {
	iscnIdBytes := cdc.MustMarshalBinaryBare(&iscnId)
	return append(IscnIdToCidKey, iscnIdBytes...)
}

func GetCidBlockKey(cid CID) []byte {
	bz := cid.Bytes()
	return append(CidBlockKey, bz...)
}

func GetCidToIscnIdKey(cid CID) []byte {
	bz := cid.Bytes()
	return append(CidToIscnIdKey, bz...)
}

func GetIscnIdVersionKey(cdc codec.BinaryMarshaler, iscnId IscnId) []byte {
	iscnId.Version = 0
	iscnIdBytes := cdc.MustMarshalBinaryBare(&iscnId)
	return append(IscnIdVersionKey, iscnIdBytes...)
}

func GetIscnIdOwnerKey(cdc codec.BinaryMarshaler, iscnId IscnId) []byte {
	iscnId.Version = 0
	iscnIdBytes := cdc.MustMarshalBinaryBare(&iscnId)
	return append(IscnIdOwnerKey, iscnIdBytes...)
}

func GetFingerprintToCidKey(fingerprint string) []byte {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(fingerprint)))
	fingerprintBytes := []byte(fingerprint)
	output := make([]byte, len(FingerprintToCidKey)+4+len(fingerprintBytes))
	copy(output, FingerprintToCidKey)
	copy(output[len(FingerprintToCidKey):], lenBytes)
	copy(output[len(FingerprintToCidKey)+4:], fingerprintBytes)
	return output
}

func GetFingerprintCidRecordKey(fingerprint string, cid CID) []byte {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(fingerprint)))
	fingerprintBytes := []byte(fingerprint)
	cidBytes := cid.Bytes()
	output := make([]byte, len(FingerprintToCidKey)+4+len(fingerprintBytes)+len(cidBytes))
	copy(output, FingerprintToCidKey)
	copy(output[len(FingerprintToCidKey):], lenBytes)
	copy(output[len(FingerprintToCidKey)+4:], fingerprintBytes)
	copy(output[len(FingerprintToCidKey)+4+len(fingerprintBytes):], cidBytes)
	return output
}
