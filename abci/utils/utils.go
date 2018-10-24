package utils

import (
	"bytes"
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
	if s[0:2] == "0x" {
		s = s[2:]
	}
	return hex.DecodeString(s)
}
