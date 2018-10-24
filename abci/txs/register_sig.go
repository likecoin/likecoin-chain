package txs

import "github.com/likecoin/likechain/abci/types"

// RegisterSignature is the signature of a RegisterTransaction
type RegisterSignature interface {
	RecoverAddress(*RegisterTransaction) (*types.Address, error)
}

// RegisterJSONSignature is the JSON form of RegisterSignature
type RegisterJSONSignature struct {
	JSONSignature
}

// RecoverAddress recovers the signing address
func (sig *RegisterJSONSignature) RecoverAddress(tx *RegisterTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(map[string]interface{}{
		"addr": tx.Addr.String(),
	})
}
