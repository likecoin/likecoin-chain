package deposit

import (
	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

// GetDepositApproverInfo checks if a LikeChain ID is one of the DepositApprovers
func GetDepositApproverInfo(state context.IImmutableState, id *types.LikeChainID) *Approver {
	depositApprovers := GetDepositApprovers(state)
	for _, approver := range depositApprovers {
		if id.Equals(approver.ID) {
			return &approver
		}
	}
	return nil
}

// CheckDeposit checks if a proposer could be proposed in the current context
func CheckDeposit(state context.IImmutableState, proposal Proposal, proposer *types.LikeChainID) response.R {
	blockNumber := proposal.BlockNumber
	if GetDepositExecution(state, blockNumber) != nil {
		return response.DepositAlreadyExecuted
	}
	if GetDepositApproverInfo(state, proposer) == nil {
		return response.DepositNotApprover
	}
	if HasApprovedDeposit(state, proposer, proposal.Hash()) {
		return response.DepositDoubleApproval
	}
	return response.Success
}

// ProcessDeposit creates a deposit proposal and related entries in the context, returns whether the proposal is executed
func ProcessDeposit(state context.IMutableState, proposal Proposal, proposer *types.LikeChainID) bool {
	blockNumber := proposal.BlockNumber
	proposalHash := proposal.Hash()
	setDepositApproval(state, proposer, proposalHash)
	approverInfo := GetDepositApproverInfo(state, proposer)
	newWeight := IncreaseDepositProposalWeight(state, proposalHash, approverInfo.Weight)
	weightSum := GetDepositApproversWeightSum(state)
	if newWeight*3 <= weightSum*2 {
		return false
	}
	setDepositExecution(state, blockNumber, proposalHash)
	for _, input := range proposal.Inputs {
		addr := input.FromAddr
		id := account.NormalizeIdentifier(state, &addr)
		account.AddBalance(state, id, input.Value.Int)
	}
	return true
}
