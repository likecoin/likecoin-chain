package txs

import (
	"strconv"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/contract"
	"github.com/likecoin/likechain/abci/types"

	cmn "github.com/tendermint/tendermint/libs/common"
)

// ContractUpdateTransaction represents a ContractUpdate transaction
type ContractUpdateTransaction struct {
	Proposer types.Identifier
	Proposal contract.Proposal
	Nonce    uint64
	Sig      ContractUpdateSignature
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *ContractUpdateTransaction) ValidateFormat() bool {
	return tx.Proposer != nil && tx.Sig != nil
}

func (tx *ContractUpdateTransaction) checkTx(state context.IImmutableState) (
	r response.R, senderID *types.LikeChainID,
) {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.ContractUpdateInvalidFormat.Info)
		return response.ContractUpdateInvalidFormat, nil
	}

	senderID = account.IdentifierToLikeChainID(state, tx.Proposer)
	if senderID == nil {
		logTx(tx).Info(response.ContractUpdateSenderNotRegistered.Info)
		return response.ContractUpdateSenderNotRegistered, nil
	}

	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !account.IsLikeChainIDHasAddress(state, senderID, addr) {
		logTx(tx).
			WithField("recovered_addr", addr).
			WithError(err).
			Info(response.ContractUpdateInvalidSignature.Info)
		return response.ContractUpdateInvalidSignature, senderID
	}

	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.ContractUpdateInvalidNonce.Info)
		return response.ContractUpdateInvalidNonce, senderID
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.ContractUpdateDuplicated.Info)
		return response.ContractUpdateDuplicated, senderID
	}

	return contract.CheckUpdate(state, &tx.Proposal, senderID), senderID
}

// CheckTx checks the transaction to see if it should be executed
func (tx *ContractUpdateTransaction) CheckTx(state context.IImmutableState) response.R {
	r, _ := tx.checkTx(state)
	return r
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *ContractUpdateTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxRes, senderID := tx.checkTx(state)
	if checkTxRes.Code != 0 {
		if checkTxRes.ShouldIncrementNonce {
			account.IncrementNextNonce(state, senderID)
		}
		return checkTxRes
	}

	account.IncrementNextNonce(state, senderID)
	height := state.GetHeight() + 1
	tags := []cmn.KVPair{
		{
			Key:   []byte("contract_update.height"),
			Value: []byte(strconv.FormatInt(height, 10)),
		},
	}
	executed := contract.ProcessUpdate(state, &tx.Proposal, senderID)
	if executed {
		tags = append(tags, cmn.KVPair{
			Key:   []byte("contract_update_execution.height"),
			Value: []byte(strconv.FormatInt(height, 10)),
		})
	}

	return response.Success.Merge(response.R{
		Tags: tags,
	})
}

// ContractUpdateTx returns raw bytes of a ContractUpdateTransaction
func ContractUpdateTx(proposer types.Identifier, contractIndex uint64, contractAddr *types.Address, nonce uint64, sigHex string) *ContractUpdateTransaction {
	sig := &ContractUpdateJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &ContractUpdateTransaction{
		Proposer: proposer,
		Proposal: contract.Proposal{
			ContractIndex:   contractIndex,
			ContractAddress: *contractAddr,
		},
		Nonce: nonce,
		Sig:   sig,
	}
}

// RawContractUpdateTx returns raw bytes of a ContractUpdateTransaction
func RawContractUpdateTx(proposer types.Identifier, contractIndex uint64, contractAddr *types.Address, nonce uint64, sigHex string) []byte {
	return EncodeTx(ContractUpdateTx(proposer, contractIndex, contractAddr, nonce, sigHex))
}
