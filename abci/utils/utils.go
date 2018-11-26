package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

// DbRawKey composes a key with prefix and suffix for IAVL tree
func DbRawKey(content []byte, prefix string, suffix string) []byte {
	var buf bytes.Buffer

	if len(prefix) > 0 {
		buf.WriteString(prefix)
		buf.WriteString("_")
	}

	buf.Write(content)

	if len(suffix) > 0 {
		buf.WriteString("_")
		buf.WriteString(suffix)
	}

	return buf.Bytes()
}

// DbTxHashKey returns a key with txHash
func DbTxHashKey(txHash []byte, suffix string) []byte {
	return DbRawKey(txHash, "tx:hash:", suffix)
}

// HashRawTx hash a rawTx in byte
func HashRawTx(rawTx []byte) []byte {
	return tmhash.Sum(rawTx)
}

// IsValidBigIntegerString verifies big integer in string
func IsValidBigIntegerString(s string) bool {
	i, ok := new(big.Int).SetString(s, 10)
	if i != nil {
		j := new(big.Int).SetBytes(i.Bytes())
		ok = i.Cmp(j) == 0
	}
	return ok
}

// Hex2Bytes strips the prefix "0x" (if any), then decode the hex string into bytes
func Hex2Bytes(s string) ([]byte, error) {
	if len(s) >= 2 && s[0:2] == "0x" {
		s = s[2:]
	}
	return hex.DecodeString(s)
}

// EncodeUint64 encodes uint64 into big-endian 8 bytes array
func EncodeUint64(n uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, n)
	return bs
}

// DecodeUint64 decodes uint64 from big-endian 8 bytes array
func DecodeUint64(bs []byte) uint64 {
	return binary.BigEndian.Uint64(bs)
}

// JoinKeys joins the keys and returns a new key, usually used for tree keys with multiple components
func JoinKeys(keys [][]byte) []byte {
	if len(keys) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	for i := 0; i < len(keys)-1; i++ {
		buf.Write(keys[i])
		buf.WriteByte('_')
	}
	buf.Write(keys[len(keys)-1])
	return buf.Bytes()
}

// PrefixKey returns the prefix of a JoinKeys result, logically equivalent to JoinKeys(...) + "_"
func PrefixKey(keys [][]byte) []byte {
	buf := new(bytes.Buffer)
	for i := 0; i < len(keys); i++ {
		buf.Write(keys[i])
		buf.WriteByte('_')
	}
	return buf.Bytes()
}
