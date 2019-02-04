package txs

import (
	"strconv"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/types"

	cmn "github.com/tendermint/tendermint/libs/common"
)

// DepositTransaction represents a Deposit transaction
type DepositTransaction struct {
	Proposer types.Identifier
	Proposal deposit.Proposal
	Nonce    uint64
	Sig      DepositSignature
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *DepositTransaction) ValidateFormat() bool {
	if tx.Proposer == nil || tx.Sig == nil {
		return false
	}
	return tx.Proposal.Validate()
}

func (tx *DepositTransaction) checkTx(state context.IImmutableState) (
	r response.R, senderID *types.LikeChainID,
) {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.DepositInvalidFormat.Info)
		return response.DepositInvalidFormat, nil
	}

	senderID = account.IdentifierToLikeChainID(state, tx.Proposer)
	if senderID == nil {
		logTx(tx).Info(response.DepositSenderNotRegistered.Info)
		return response.DepositSenderNotRegistered, nil
	}

	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !account.IsLikeChainIDHasAddress(state, senderID, addr) {
		logTx(tx).
			WithField("recovered_addr", addr).
			WithError(err).
			Info(response.DepositInvalidSignature.Info)
		return response.DepositInvalidSignature, senderID
	}

	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.DepositInvalidNonce.Info)
		return response.DepositInvalidNonce, senderID
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.DepositDuplicated.Info)
		return response.DepositDuplicated, senderID
	}

	return deposit.CheckDeposit(state, tx.Proposal, senderID), senderID
}

// CheckTx checks the transaction to see if it should be executed
func (tx *DepositTransaction) CheckTx(state context.IImmutableState) response.R {
	r, _ := tx.checkTx(state)
	return r
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *DepositTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxRes, senderID := tx.checkTx(state)
	if checkTxRes.Code != response.Success.Code {
		if checkTxRes.ShouldIncrementNonce {
			account.IncrementNextNonce(state, senderID)
		}
		return checkTxRes
	}

	account.IncrementNextNonce(state, senderID)
	height := state.GetHeight() + 1
	tags := []cmn.KVPair{
		{
			Key:   []byte("deposit.height"),
			Value: []byte(strconv.FormatInt(height, 10)),
		},
	}
	executed := deposit.ProcessDeposit(state, tx.Proposal, senderID)
	if executed {
		tags = append(tags, cmn.KVPair{
			Key:   []byte("deposit_execution.height"),
			Value: []byte(strconv.FormatInt(height, 10)),
		})
		tags = append(tags, cmn.KVPair{
			Key:   []byte("deposit_execution.eth_block"),
			Value: []byte(strconv.FormatInt(int64(tx.Proposal.BlockNumber), 10)),
		})
	}

	return response.Success.Merge(response.R{
		Tags: tags,
	})
}

// DepositTx returns raw bytes of a DepositTransaction
func DepositTx(proposer types.Identifier, blockNumber uint64, inputs []deposit.Input, nonce uint64, sigHex string) *DepositTransaction {
	sig := &DepositJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &DepositTransaction{
		Proposer: proposer,
		Proposal: deposit.Proposal{
			BlockNumber: blockNumber,
			Inputs:      inputs,
		},
		Nonce: nonce,
		Sig:   sig,
	}
}

// RawDepositTx returns raw bytes of a DepositTransaction
func RawDepositTx(proposer types.Identifier, blockNumber uint64, inputs []deposit.Input, nonce uint64, sigHex string) []byte {
	return EncodeTx(DepositTx(proposer, blockNumber, inputs, nonce, sigHex))
}
