package txs

import (
	"strconv"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"

	cmn "github.com/tendermint/tendermint/libs/common"
)

// DepositApprovalTransaction represents a DepositApproval transaction
type DepositApprovalTransaction struct {
	Approver      types.Identifier
	DepositTxHash []byte
	Nonce         uint64
	Sig           DepositApprovalSignature
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *DepositApprovalTransaction) ValidateFormat() bool {
	if tx.Approver == nil || tx.Sig == nil {
		return false
	}
	if len(tx.DepositTxHash) != 20 {
		return false
	}
	return true
}

func (tx *DepositApprovalTransaction) checkTx(state context.IImmutableState) (
	r response.R, senderID *types.LikeChainID,
) {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.DepositApprovalInvalidFormat.Info)
		return response.DepositApprovalInvalidFormat, nil
	}

	senderID = account.IdentifierToLikeChainID(state, tx.Approver)
	if senderID == nil {
		logTx(tx).Info(response.DepositApprovalSenderNotRegistered.Info)
		return response.DepositApprovalSenderNotRegistered, nil
	}

	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !account.IsLikeChainIDHasAddress(state, senderID, addr) {
		logTx(tx).
			WithField("recovered_addr", addr).
			WithError(err).
			Info(response.DepositApprovalInvalidSignature.Info)
		return response.DepositApprovalInvalidSignature, senderID
	}

	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.DepositApprovalInvalidNonce.Info)
		return response.DepositApprovalInvalidNonce, senderID
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.DepositApprovalDuplicated.Info)
		return response.DepositApprovalDuplicated, senderID
	}

	return deposit.CheckDepositApproval(state, tx.DepositTxHash, senderID), senderID
}

// CheckTx checks the transaction to see if it should be executed
func (tx *DepositApprovalTransaction) CheckTx(state context.IImmutableState) response.R {
	r, _ := tx.checkTx(state)
	return r
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *DepositApprovalTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxRes, senderID := tx.checkTx(state)
	if checkTxRes.Code != 0 {
		if checkTxRes.ShouldIncrementNonce {
			account.IncrementNextNonce(state, senderID)
		}
		return checkTxRes
	}

	account.IncrementNextNonce(state, senderID)
	weight := deposit.CreateDepositApproval(state, tx.DepositTxHash, senderID)

	height := state.GetHeight() + 1
	tags := []cmn.KVPair{
		{
			Key:   []byte("deposit_approval.height"),
			Value: []byte(strconv.FormatInt(height, 10)),
		},
	}

	weightSum := deposit.GetDepositApproversWeightSum(state)
	if weight*3 > weightSum*2 {
		deposit.ExecuteDepositProposal(state, tx.DepositTxHash)
		txstatus.SetStatus(state, tx.DepositTxHash, txstatus.TxStatusSuccess)
		tags = append(tags, cmn.KVPair{
			Key:   []byte("deposit_execution.height"),
			Value: []byte(strconv.FormatInt(height, 10)),
		})
	}

	return response.Success.Merge(response.R{
		Tags: tags,
	})
}

// DepositApprovalTx returns raw bytes of a DepositApprovalTransaction
func DepositApprovalTx(approver types.Identifier, depositTxHash []byte, nonce uint64, sigHex string) *DepositApprovalTransaction {
	sig := &DepositApprovalJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &DepositApprovalTransaction{
		Approver:      approver,
		DepositTxHash: depositTxHash,
		Nonce:         nonce,
		Sig:           sig,
	}
}

// RawDepositApprovalTx returns raw bytes of a DepositApprovalTransaction
func RawDepositApprovalTx(approver types.Identifier, depositTxHash []byte, nonce uint64, sigHex string) []byte {
	return EncodeTx(DepositApprovalTx(approver, depositTxHash, nonce, sigHex))
}
