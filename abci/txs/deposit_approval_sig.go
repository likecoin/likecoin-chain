package txs

import (
	"strings"

	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
)

// DepositApprovalSignature is the signature of a DepositApprovalTransaction
type DepositApprovalSignature interface {
	RecoverAddress(*DepositApprovalTransaction) (*types.Address, error)
}

// DepositApprovalJSONSignature is the JSON form of DepositApprovalSignature
type DepositApprovalJSONSignature struct {
	JSONSignature
}

// GenerateJSONMap generates the JSON map from the transaction, which is used for generating and verifying JSON signature
func (tx *DepositApprovalTransaction) GenerateJSONMap() map[string]interface{} {
	return map[string]interface{}{
		"deposit_tx_hash": "0x" + strings.ToLower(common.Bytes2Hex(tx.DepositTxHash)),
		"identity":        tx.Approver,
		"nonce":           tx.Nonce,
	}
}

// RecoverAddress recovers the signing address
func (sig *DepositApprovalJSONSignature) RecoverAddress(tx *DepositApprovalTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(tx.GenerateJSONMap())
}
