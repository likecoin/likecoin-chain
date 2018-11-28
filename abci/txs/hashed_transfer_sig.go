package txs

import (
	"strings"

	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
)

// HashedTransferSignature is the signature of a HashedTransferTransaction
type HashedTransferSignature interface {
	RecoverAddress(*HashedTransferTransaction) (*types.Address, error)
}

// HashedTransferJSONSignature is the JSON form of HashedTransferSignature
type HashedTransferJSONSignature struct {
	JSONSignature
}

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *HashedTransferTransaction) GenerateJSONMap() JSONMap {
	return JSONMap{
		"identity":    tx.HashedTransfer.From.String(),
		"to":          tx.HashedTransfer.To.String(),
		"value":       tx.HashedTransfer.Value.String(),
		"hash_commit": "0x" + strings.ToLower(common.Bytes2Hex(tx.HashedTransfer.HashCommit[:])),
		"expiry":      tx.HashedTransfer.Expiry,
		"fee":         tx.Fee.String(),
		"nonce":       tx.Nonce,
	}
}

// RecoverAddress recovers the signing address
func (sig *HashedTransferJSONSignature) RecoverAddress(tx *HashedTransferTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}

// HashedTransferEIP712Signature is the EIP-712 form of RegisterSignature
type HashedTransferEIP712Signature struct {
	EIP712Signature
}

// GenerateEIP712SignData generates the EIP-712 sign data, which is used for generating and verifying EIP-712 signature
func (tx *HashedTransferTransaction) GenerateEIP712SignData() EIP712SignData {
	return EIP712SignData{
		Name: "HashedTransfer",
		Fields: []EIP712Field{
			{"identity", EIP712Identifier{tx.HashedTransfer.From}},
			{"to", EIP712Identifier{tx.HashedTransfer.To}},
			{"value", EIP712Uint256(tx.HashedTransfer.Value)},
			{"hash_commit", EIP712Bytes32(tx.HashedTransfer.HashCommit[:])},
			{"expiry", EIP712Uint256(types.NewBigInt(int64(tx.HashedTransfer.Expiry)))},
			{"fee", EIP712Uint256(tx.Fee)},
			{"nonce", EIP712Uint64(tx.Nonce)},
		},
	}
}

// RecoverAddress recovers the signing address
func (sig *HashedTransferEIP712Signature) RecoverAddress(tx *HashedTransferTransaction) (*types.Address, error) {
	return sig.EIP712Signature.RecoverAddress(tx.GenerateEIP712SignData())
}
