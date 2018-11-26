package txs

import "github.com/likecoin/likechain/abci/types"

// DepositSignature is the signature of a DepositTransaction
type DepositSignature interface {
	RecoverAddress(*DepositTransaction) (*types.Address, error)
}

// DepositJSONSignature is the JSON form of DepositSignature
type DepositJSONSignature struct {
	JSONSignature
}

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *DepositTransaction) GenerateJSONMap() map[string]interface{} {
	inputs := make([]map[string]interface{}, len(tx.Proposal.Inputs))
	for i, input := range tx.Proposal.Inputs {
		inputs[i] = map[string]interface{}{
			"value":     input.Value.String(),
			"from_addr": input.FromAddr.String(),
		}
	}
	return map[string]interface{}{
		"block_number": tx.Proposal.BlockNumber,
		"identity":     tx.Proposer.String(),
		"nonce":        tx.Nonce,
		"inputs":       inputs,
	}
}

// RecoverAddress recovers the signing address
func (sig *DepositJSONSignature) RecoverAddress(tx *DepositTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}
