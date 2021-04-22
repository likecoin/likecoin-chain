package types

import "encoding/binary"

func EncodeUint64(n uint64) []byte {
	output := make([]byte, 8)
	binary.BigEndian.PutUint64(output, n)
	return output
}

func DecodeUint64(bz []byte) uint64 {
	if len(bz) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func EncodeUint32(n uint32) []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, n)
	return output
}

func DecodeUint32(bz []byte) uint32 {
	if len(bz) == 0 {
		return 0
	}
	return binary.BigEndian.Uint32(bz)
}
