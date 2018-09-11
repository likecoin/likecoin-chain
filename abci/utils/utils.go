package utils

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/likecoin/likechain/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
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

// DbIDKey composes a key with LikeChain ID in
// `{prefix}:id:_{id}_{suffix}` format
func DbIDKey(id *types.LikeChainID, prefix string, suffix string) []byte {
	var buf bytes.Buffer
	buf.WriteString(prefix)
	buf.WriteString(":id:")
	return DbRawKey(id.Content, buf.String(), suffix)
}

// DbAddrKey returns a key with Ethereum address in
// `{prefix}:addr:_{addr}_{suffix}` format
func DbAddrKey(addr *types.Address, prefix string, suffix string) []byte {
	var buf bytes.Buffer
	buf.WriteString(prefix)
	buf.WriteString(":addr:")
	return DbRawKey(addr.Content, buf.String(), suffix)
}

// DbIdentifierKey returns a key either with Ethereum address or LikeChain ID
func DbIdentifierKey(
	identifier *types.Identifier,
	prefix string,
	suffix string,
) (key []byte) {
	if addr := identifier.GetAddr(); addr != nil {
		key = DbAddrKey(addr, prefix, suffix)
	} else if id := identifier.GetLikeChainID(); id != nil {
		key = DbIDKey(id, prefix, suffix)
	}
	return key
}

// DbTxHashKey returns a key with txHash
func DbTxHashKey(txHash []byte, suffix string) []byte {
	return DbRawKey(txHash, "tx:hash:", suffix)
}

// HashRawTx hash a rawTx in byte
func HashRawTx(rawTx []byte) []byte {
	return tmhash.Sum(rawTx)
}
