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

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *TransferTransaction) GenerateJSONMap() JSONMap {
	outputs := make([]JSONMap, len(tx.Outputs))
	for i, output := range tx.Outputs {
		outputs[i] = JSONMap{
			"identity": output.To.String(),
			"remark":   output.Remark.String(),
			"value":    output.Value.String(),
		}
	}
	return JSONMap{
		"fee":      tx.Fee.String(),
		"identity": tx.From.String(),
		"nonce":    tx.Nonce,
		"outputs":  outputs,
	}
}

// RecoverAddress recovers the signing address
func (sig *TransferJSONSignature) RecoverAddress(tx *TransferTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}
