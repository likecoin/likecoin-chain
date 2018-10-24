package txs

import "github.com/likecoin/likechain/abci/types"

// TransferSignature is the signature of a TransferTransaction
type TransferSignature interface {
	RecoverAddress(*TransferTransaction) (*types.Address, error)
}

// TransferJSONSignature is the JSON form of TransferSignature
type TransferJSONSignature struct {
	JSONSignature
}

// RecoverAddress recovers the signing address
func (sig *TransferJSONSignature) RecoverAddress(tx *TransferTransaction) (*types.Address, error) {
	outputs := make([]map[string]interface{}, len(tx.Outputs))
	for i, output := range tx.Outputs {
		outputs[i] = map[string]interface{}{
			"identity": output.To.String(),
			"remark":   output.Remark.String(),
			"value":    output.Value.String(),
		}
	}
	return sig.JSONSignature.RecoverAddress(map[string]interface{}{
		"fee":      tx.Fee.String(),
		"identity": tx.From.String(),
		"nonce":    tx.Nonce,
		"outputs":  outputs,
	})
}
