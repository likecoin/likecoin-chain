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
func (tx *SimpleTransferTransaction) GenerateJSONMap() JSONMap {
	return JSONMap{
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

// SimpleTransferEIP712Signature is the EIP-712 form of RegisterSignature
type SimpleTransferEIP712Signature struct {
	EIP712Signature
}

// GenerateEIP712SignData generates the EIP-712 sign data, which is used for generating and verifying EIP-712 signature
func (tx *SimpleTransferTransaction) GenerateEIP712SignData() EIP712SignData {
	return EIP712SignData{
		Name: "SimpleTransfer",
		Fields: []EIP712Field{
			{"identity", EIP712Identifier{tx.From}},
			{"to", EIP712Identifier{tx.To}},
			{"value", EIP712Uint256(tx.Value)},
			{"remark", EIP712String(tx.Remark)},
			{"fee", EIP712Uint256(tx.Fee)},
			{"nonce", EIP712Uint64(tx.Nonce)},
		},
	}
}

// RecoverAddress recovers the signing address
func (sig *SimpleTransferEIP712Signature) RecoverAddress(tx *SimpleTransferTransaction) (*types.Address, error) {
	return sig.EIP712Signature.RecoverAddress(tx.GenerateEIP712SignData())
}
