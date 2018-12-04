package txs

import (
	"strings"

	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
)

// ClaimHashedTransferSignature is the signature of a ClaimHashedTransferTransaction
type ClaimHashedTransferSignature interface {
	RecoverAddress(*ClaimHashedTransferTransaction) (*types.Address, error)
}

// ClaimHashedTransferJSONSignature is the JSON form of ClaimHashedTransferSignature
type ClaimHashedTransferJSONSignature struct {
	JSONSignature
}

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *ClaimHashedTransferTransaction) GenerateJSONMap() map[string]interface{} {
	secretStr := ""
	if len(tx.Secret) != 0 {
		secretStr = "0x" + strings.ToLower(common.Bytes2Hex(tx.Secret))
	}
	return map[string]interface{}{
		"identity":     tx.From.String(),
		"htlc_tx_hash": "0x" + strings.ToLower(common.Bytes2Hex(tx.HTLCTxHash)),
		"secret":       secretStr,
		"nonce":        tx.Nonce,
	}
}

// RecoverAddress recovers the signing address
func (sig *ClaimHashedTransferJSONSignature) RecoverAddress(tx *ClaimHashedTransferTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}
