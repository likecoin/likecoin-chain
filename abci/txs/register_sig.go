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

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *RegisterTransaction) GenerateJSONMap() JSONMap {
	return JSONMap{
		"addr": tx.Addr.String(),
	}
}

// RecoverAddress recovers the signing address
func (sig *RegisterJSONSignature) RecoverAddress(tx *RegisterTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}

// RegisterEIP712Signature is the EIP-712 form of RegisterSignature
type RegisterEIP712Signature struct {
	EIP712Signature
}

// GenerateEIP712SignData generates the EIP-712 sign data, which is used for generating and verifying EIP-712 signature
func (tx *RegisterTransaction) GenerateEIP712SignData() EIP712SignData {
	return EIP712SignData{
		Name: "Register",
		Fields: []EIP712Field{
			{"addr", EIP712Address(tx.Addr)},
		},
	}
}

// RecoverAddress recovers the signing address
func (sig *RegisterEIP712Signature) RecoverAddress(tx *RegisterTransaction) (*types.Address, error) {
	return sig.EIP712Signature.RecoverAddress(tx.GenerateEIP712SignData())
}
