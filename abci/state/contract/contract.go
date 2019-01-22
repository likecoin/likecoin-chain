package contract

import (
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

// GetContractUpdaterInfo checks if a LikeChain ID is one of the ContractUpdaters
func GetContractUpdaterInfo(state context.IImmutableState, id *types.LikeChainID) *Updater {
	updaters := GetContractUpdaters(state)
	for _, updater := range updaters {
		if id.Equals(updater.ID) {
			return &updater
		}
	}
	return nil
}

// CheckUpdate checks if an update proposal could be proposed in the current context
func CheckUpdate(state context.IImmutableState, proposal *Proposal, proposer *types.LikeChainID) response.R {
	contractIndex := proposal.ContractIndex
	currentIndex := GetContractIndex(state)
	if GetContractUpdaterInfo(state, proposer) == nil {
		return response.ContractUpdateNotUpdater
	}
	if HasApprovedUpdate(state, proposer, proposal.Bytes()) {
		return response.ContractUpdateDoubleApproval
	}
	if contractIndex != currentIndex+1 {
		return response.ContractUpdateInvalidIndex
	}
	return response.Success
}

// ProcessUpdate creates an update proposal and related entries in the context, returns whether the proposal is executed
func ProcessUpdate(state context.IMutableState, proposal *Proposal, proposer *types.LikeChainID) bool {
	proposalBytes := proposal.Bytes()
	setUpdateApproval(state, proposer, proposalBytes)
	updaterInfo := GetContractUpdaterInfo(state, proposer)
	newWeight := IncreaseUpdateProposalWeight(state, proposalBytes, updaterInfo.Weight)
	weightSum := GetContractUpdatersWeightSum(state)
	if newWeight*3 <= weightSum*2 {
		return false
	}
	setUpdateExecution(state, proposal)
	increaseContractIndex(state)
	return true
}
