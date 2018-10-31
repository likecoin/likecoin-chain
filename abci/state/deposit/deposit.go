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

// CheckDepositProposal checks if a proposer could be proposed in the current context
func CheckDepositProposal(state context.IImmutableState, proposal Proposal, proposer *types.LikeChainID) response.R {
	blockNumber := proposal.BlockNumber
	if GetDepositExecution(state, blockNumber) != nil {
		return response.DepositAlreadyExecuted
	}
	if GetDepositApproverInfo(state, proposer) == nil {
		return response.DepositNotApprover
	}
	if GetDepositApproval(state, proposer, blockNumber) != nil {
		return response.DepositDoubleApproval
	}
	return response.Success
}

// CreateDepositProposal creates a deposit proposal and related entries in the context, returns the weight of the proposal
func CreateDepositProposal(state context.IMutableState, txHash []byte, proposal Proposal, proposer *types.LikeChainID) uint64 {
	blockNumber := proposal.BlockNumber
	SetDepositProposal(state, txHash, proposal)
	SetDepositApproval(state, proposer, blockNumber, txHash)
	approverInfo := GetDepositApproverInfo(state, proposer)
	return IncreaseDepositProposalWeight(state, txHash, approverInfo.Weight)
}

// CheckDepositApproval checks if a deposit approval is valid in the current context
func CheckDepositApproval(state context.IImmutableState, txHash []byte, approver *types.LikeChainID) response.R {
	proposal := GetDepositProposal(state, txHash)
	if proposal == nil {
		return response.DepositApprovalProposalNotExist
	}
	blockNumber := proposal.BlockNumber
	if GetDepositExecution(state, blockNumber) != nil {
		return response.DepositApprovalAlreadyExecuted
	}
	if GetDepositApproverInfo(state, approver) == nil {
		return response.DepositApprovalNotApprover
	}
	if GetDepositApproval(state, approver, blockNumber) != nil {
		return response.DepositApprovalDoubleApproval
	}
	return response.Success
}

// CreateDepositApproval creates a deposit approval and related entries in the context, returns the new proposal weight
func CreateDepositApproval(state context.IMutableState, txHash []byte, approver *types.LikeChainID) uint64 {
	proposal := GetDepositProposal(state, txHash)
	SetDepositApproval(state, approver, proposal.BlockNumber, txHash)
	approverInfo := GetDepositApproverInfo(state, approver)
	newWeight := IncreaseDepositProposalWeight(state, txHash, approverInfo.Weight)
	return newWeight
}

// ExecuteDepositProposal executes a deposit proposal and clean up the related entries in the context
func ExecuteDepositProposal(state context.IMutableState, txHash []byte) {
	proposal := GetDepositProposal(state, txHash)
	blockNumber := proposal.BlockNumber
	SetDepositExecution(state, blockNumber, txHash)
	for _, input := range proposal.Inputs {
		addr := input.FromAddr
		id := account.NormalizeIdentifier(state, &addr)
		account.AddBalance(state, id, input.Value.Int)
	}
}
