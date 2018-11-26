package txs

import "github.com/likecoin/likechain/abci/types"

// WithdrawSignature is the signature of a WithdrawTransaction
type WithdrawSignature interface {
	RecoverAddress(*WithdrawTransaction) (*types.Address, error)
}

// WithdrawJSONSignature is the JSON form of WithdrawSignature
type WithdrawJSONSignature struct {
	JSONSignature
}

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *WithdrawTransaction) GenerateJSONMap() map[string]interface{} {
	return map[string]interface{}{
		"fee":      tx.Fee.String(),
		"identity": tx.From.String(),
		"nonce":    tx.Nonce,
		"to_addr":  tx.ToAddr.String(),
		"value":    tx.Value.String(),
	}
}

// RecoverAddress recovers the signing address
func (sig *WithdrawJSONSignature) RecoverAddress(tx *WithdrawTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}
