package txs

import (
	"strings"

	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
)

// ClaimHashedTransferSignature is the signature of a ClaimHashedTransferTransaction
type ClaimHashedTransferSignature interface {
	RecoverAddress(*ClaimHashedTransferTransaction) (*types.Address, error)
}

// ClaimHashedTransferJSONSignature is the JSON form of ClaimHashedTransferSignature
type ClaimHashedTransferJSONSignature struct {
	JSONSignature
}

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *ClaimHashedTransferTransaction) GenerateJSONMap() JSONMap {
	secretStr := ""
	if len(tx.Secret) != 0 {
		secretStr = "0x" + strings.ToLower(common.Bytes2Hex(tx.Secret))
	}
	return JSONMap{
		"identity":     tx.From.String(),
		"htlc_tx_hash": "0x" + strings.ToLower(common.Bytes2Hex(tx.HTLCTxHash)),
		"secret":       secretStr,
		"nonce":        tx.Nonce,
	}
}

// RecoverAddress recovers the signing address
func (sig *ClaimHashedTransferJSONSignature) RecoverAddress(tx *ClaimHashedTransferTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}

// ClaimHashedTransferEIP712Signature is the EIP-712 form of RegisterSignature
type ClaimHashedTransferEIP712Signature struct {
	EIP712Signature
}

// GenerateEIP712SignData generates the EIP-712 sign data, which is used for generating and verifying EIP-712 signature
func (tx *ClaimHashedTransferTransaction) GenerateEIP712SignData() EIP712SignData {
	return EIP712SignData{
		Name: "ClaimHashedTransfer",
		Fields: []EIP712Field{
			{"identity", EIP712Identifier{tx.From}},
			{"htlc_tx_hash", EIP712Bytes32(tx.HTLCTxHash)},
			{"secret", EIP712Bytes32(tx.Secret[:])},
			{"nonce", EIP712Uint64(tx.Nonce)},
		},
	}
}

// RecoverAddress recovers the signing address
func (sig *ClaimHashedTransferEIP712Signature) RecoverAddress(tx *ClaimHashedTransferTransaction) (*types.Address, error) {
	return sig.EIP712Signature.RecoverAddress(tx.GenerateEIP712SignData())
}
