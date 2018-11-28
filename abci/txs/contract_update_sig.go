package txs

import "github.com/likecoin/likechain/abci/types"

// ContractUpdateSignature is the signature of a ContractUpdateTransaction
type ContractUpdateSignature interface {
	RecoverAddress(*ContractUpdateTransaction) (*types.Address, error)
}

// ContractUpdateJSONSignature is the JSON form of ContractUpdateSignature
type ContractUpdateJSONSignature struct {
	JSONSignature
}

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *ContractUpdateTransaction) GenerateJSONMap() JSONMap {
	return JSONMap{
		"contract_index": tx.Proposal.ContractIndex,
		"contract_addr":  tx.Proposal.ContractAddress.String(),
		"identity":       tx.Proposer.String(),
		"nonce":          tx.Nonce,
	}
}

// RecoverAddress recovers the signing address
func (sig *ContractUpdateJSONSignature) RecoverAddress(tx *ContractUpdateTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}
