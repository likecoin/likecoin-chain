package utils

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/likecoin/likechain/abci/types"
)

// RecoverSignature returns the address for the account that was used to create the signature.
func RecoverSignature(hash []byte, sig *types.Signature) (common.Address, error) {
	// Transform yellow paper V from 27/28 to 0/1
	sigContent := make([]byte, len(sig.Content))
	copy(sigContent, sig.Content)

	if len(sigContent) == 65 {
		sigContent[64] -= 27
	}
	pubKeyBytes, err := crypto.Ecrecover(hash, sigContent)
	if err != nil {
		return common.Address{}, err
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

// DbKeyRaw composes a key with prefix and suffix for IAVL tree
func DbKeyRaw(content []byte, prefix string, suffix string) []byte {
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

// DbIDKey composes a key with LikeChain ID in `{prefix}_{id}_{suffix}` format
func DbIDKey(id types.LikeChainID, prefix string, suffix string) []byte {
	return DbKeyRaw(id.Content, prefix, suffix)
}

// DbAddrKey returns a key with Ethereum address in `addr_{addr}_id` format
func DbAddrKey(ethAddr common.Address) []byte {
	return DbRawAddrKey(ethAddr.Bytes())
}

// DbRawAddrKey returns a key with protobuf address in `addr_{addr}_id` format
func DbRawAddrKey(addr []byte) []byte {
	return DbKeyRaw(addr, "addr", "id")
}
