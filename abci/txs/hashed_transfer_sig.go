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
func (tx *HashedTransferTransaction) GenerateJSONMap() map[string]interface{} {
	return map[string]interface{}{
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
