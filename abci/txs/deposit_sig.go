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
func (tx *DepositTransaction) GenerateJSONMap() JSONMap {
	inputs := make([]JSONMap, len(tx.Proposal.Inputs))
	for i, input := range tx.Proposal.Inputs {
		inputs[i] = JSONMap{
			"value":     input.Value.String(),
			"from_addr": input.FromAddr.String(),
		}
	}
	return JSONMap{
		"block_number": tx.Proposal.BlockNumber,
		"identity":     tx.Proposer.String(),
		"nonce":        tx.Nonce,
		"inputs":       inputs,
	}
}

// RecoverAddress recovers the signing address
func (sig *DepositJSONSignature) RecoverAddress(tx *DepositTransaction) (*types.Address, error) {
	tx.Proposal.Sort()
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}
