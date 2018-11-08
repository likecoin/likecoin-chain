package txs

import "github.com/likecoin/likechain/abci/types"

// SimpleTransferSignature is the signature of a SimpleTransferTransaction
type SimpleTransferSignature interface {
	RecoverAddress(*SimpleTransferTransaction) (*types.Address, error)
}

// SimpleTransferJSONSignature is the JSON form of SimpleTransferSignature
type SimpleTransferJSONSignature struct {
	JSONSignature
}

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *SimpleTransferTransaction) GenerateJSONMap() map[string]interface{} {
	return map[string]interface{}{
		"identity": tx.From.String(),
		"to":       tx.To.String(),
		"value":    tx.Value.String(),
		"remark":   tx.Remark,
		"fee":      tx.Fee.String(),
		"nonce":    tx.Nonce,
	}
}

// RecoverAddress recovers the signing address
func (sig *SimpleTransferJSONSignature) RecoverAddress(tx *SimpleTransferTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}
