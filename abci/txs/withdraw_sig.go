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

// RecoverAddress recovers the signing address
func (sig *WithdrawJSONSignature) RecoverAddress(tx *WithdrawTransaction) (*types.Address, error) {
	jsonMap := map[string]interface{}{
		"fee":      tx.Fee.String(),
		"identity": tx.From.String(),
		"nonce":    tx.Nonce,
		"to_addr":  tx.ToAddr.String(),
		"value":    tx.Value.String(),
	}
	return sig.JSONSignature.RecoverAddress(jsonMap)
}
