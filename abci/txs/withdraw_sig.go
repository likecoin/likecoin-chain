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
func (tx *WithdrawTransaction) GenerateJSONMap() JSONMap {
	return JSONMap{
		"identity": tx.From.String(),
		"to_addr":  tx.ToAddr.String(),
		"value":    tx.Value.String(),
		"fee":      tx.Fee.String(),
		"nonce":    tx.Nonce,
	}
}

// RecoverAddress recovers the signing address
func (sig *WithdrawJSONSignature) RecoverAddress(tx *WithdrawTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}

// WithdrawEIP712Signature is the EIP-712 form of RegisterSignature
type WithdrawEIP712Signature struct {
	EIP712Signature
}

// GenerateEIP712SignData generates the EIP-712 sign data, which is used for generating and verifying EIP-712 signature
func (tx *WithdrawTransaction) GenerateEIP712SignData() EIP712SignData {
	return EIP712SignData{
		Name: "Withdraw",
		Fields: []EIP712Field{
			{"identity", EIP712Identifier{tx.From}},
			{"to_addr", EIP712Address(tx.ToAddr)},
			{"value", EIP712Uint256(tx.Value)},
			{"fee", EIP712Uint256(tx.Fee)},
			{"nonce", EIP712Uint64(tx.Nonce)},
		},
	}
}

// RecoverAddress recovers the signing address
func (sig *WithdrawEIP712Signature) RecoverAddress(tx *WithdrawTransaction) (*types.Address, error) {
	return sig.EIP712Signature.RecoverAddress(tx.GenerateEIP712SignData())
}
