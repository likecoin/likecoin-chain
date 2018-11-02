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

// RecoverAddress recovers the signing address
func (sig *DepositApprovalJSONSignature) RecoverAddress(tx *DepositApprovalTransaction) (*types.Address, error) {
	return sig.JSONSignature.RecoverAddress(map[string]interface{}{
		"deposit_tx_hash": "0x" + strings.ToLower(common.Bytes2Hex(tx.DepositTxHash)),
		"identity":        tx.Approver,
		"nonce":           tx.Nonce,
	})
}
