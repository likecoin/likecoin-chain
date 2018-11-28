package contract

import (
	"bytes"

	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	"github.com/tendermint/iavl"
	cmn "github.com/tendermint/tendermint/libs/common"
)

var (
	contractUpdatersKey          = []byte("contractUpdaters")
	contractUpdatersWeightSumKey = []byte("contractUpdatersWeightSum")

	currentContractIndexKey = []byte("currentContractIndex")

	updateWeightKey   = []byte("updateWeight")
	updateApprovalKey = []byte("updateApproval")

	log = logger.L
)

// Updater represents the identity and weight of a contract updater
type Updater struct {
	ID     *types.LikeChainID
	Weight uint32
}

// Proposal represents an update on Ethereum smart contract
type Proposal struct {
	ContractIndex   uint64
	ContractAddress types.Address
}

// Bytes returns a unique representation of a Proposal in []byte
func (proposal *Proposal) Bytes() []byte {
	buf := new(bytes.Buffer)
	buf.Write(utils.EncodeUint64(proposal.ContractIndex))
	buf.Write(proposal.ContractAddress[:])
	return buf.Bytes()
}

func approvalKey(id *types.LikeChainID, proposalBytes []byte) []byte {
	return utils.JoinKeys([][]byte{
		updateApprovalKey,
		id.Bytes(),
		[]byte("proposal"),
		proposalBytes,
	})
}

func weightKey(proposalBytes []byte) []byte {
	return utils.JoinKeys([][]byte{
		updateWeightKey,
		proposalBytes,
	})
}

func executedKey(contractIndex uint64) []byte {
	// Not using JoinKeys since we want better interoperability with smart contract
	key := []byte("exec")
	key = append(key, utils.EncodeUint64(contractIndex)...)
	return key
}

// GetContractUpdatersWeightSum loads the weight sum of the contract updaters from state tree
func GetContractUpdatersWeightSum(state context.IImmutableState) uint64 {
	_, bs := state.ImmutableStateTree().Get(contractUpdatersWeightSumKey)
	if bs == nil {
		return 0
	}
	if len(bs) != 8 {
		log.
			WithField("data", cmn.HexBytes(bs)).
			Panic("Invalid contract updater weight raw data")
	}
	return utils.DecodeUint64(bs)
}

// GetContractUpdaters loads the contract updater list from state tree
func GetContractUpdaters(state context.IImmutableState) (updaters []Updater) {
	_, bs := state.ImmutableStateTree().Get(contractUpdatersKey)
	if bs == nil {
		return nil
	}
	err := types.AminoCodec().UnmarshalBinaryBare(bs, &updaters)
	if err != nil {
		log.
			WithField("data", cmn.HexBytes(bs)).
			WithError(err).
			Panic("Cannot unmarshal contract updaters")
	}
	return updaters
}

// SetContractUpdaters saves the contract updater list into state tree
func SetContractUpdaters(state context.IMutableState, updaters []Updater) {
	if len(updaters) == 0 {
		state.MutableStateTree().Remove(contractUpdatersKey)
		state.MutableStateTree().Remove(contractUpdatersWeightSumKey)
		return
	}
	totalWeight := uint64(0)
	for _, updater := range updaters {
		totalWeight += uint64(updater.Weight)
	}
	bs, err := types.AminoCodec().MarshalBinaryBare(updaters)
	if err != nil {
		log.
			WithField("updaters", updaters).
			WithError(err).
			Panic("Cannot marshal contract updaters")
	}
	state.MutableStateTree().Set(contractUpdatersKey, bs)
	state.MutableStateTree().Set(contractUpdatersWeightSumKey, utils.EncodeUint64(totalWeight))
}

// setUpdateApproval records a update approval into state tree
func setUpdateApproval(state context.IMutableState, updater *types.LikeChainID, proposalBytes []byte) {
	if HasApprovedUpdate(state, updater, proposalBytes) {
		log.
			WithField("proposal_bytes", cmn.HexBytes(proposalBytes)).
			Panic("Double approving contract update on the same proposal")
	}
	key := approvalKey(updater, proposalBytes)
	state.MutableStateTree().Set(key, []byte{1})
}

// HasApprovedUpdate checks if a ContractUpdater has already approved a certain contract update proposal
func HasApprovedUpdate(state context.IImmutableState, updater *types.LikeChainID, proposalBytes []byte) bool {
	key := approvalKey(updater, proposalBytes)
	_, v := state.ImmutableStateTree().Get(key)
	return v != nil
}

// IncreaseUpdateProposalWeight initializes or increments a contract update proposal's approve weight, returns the new weight
func IncreaseUpdateProposalWeight(state context.IMutableState, proposalBytes []byte, weight uint32) uint64 {
	key := weightKey(proposalBytes)
	var oldWeight uint64
	_, bs := state.MutableStateTree().Get(key)
	if bs == nil {
		oldWeight = 0
	} else {
		oldWeight = utils.DecodeUint64(bs)
	}
	newWeight := oldWeight + uint64(weight)
	state.MutableStateTree().Set(key, utils.EncodeUint64(newWeight))
	return newWeight
}

// GetUpdateProposalWeight returns the weight sum of updaters approved this proposal
func GetUpdateProposalWeight(state context.IImmutableState, proposalBytes []byte) uint64 {
	key := weightKey(proposalBytes)
	_, bs := state.ImmutableStateTree().Get(key)
	if bs == nil {
		return 0
	}
	return utils.DecodeUint64(bs)
}

// setUpdateExecution records a update execution with contract index into withdraw tree
func setUpdateExecution(state context.IMutableState, proposal *Proposal) {
	if GetUpdateExecution(state, proposal.ContractIndex) != nil {
		log.
			WithField("contract_index", proposal.ContractIndex).
			WithField("contract_addr", proposal.ContractAddress.String()).
			Panic("Double setting contract update execution on the same contract index")
	}
	key := executedKey(proposal.ContractIndex)
	state.MutableWithdrawTree().Set(key, proposal.ContractAddress[:])
}

// GetUpdateExecution returns the executed proposalBytes for a contract index
func GetUpdateExecution(state context.IImmutableState, contractIndex uint64) *types.Address {
	key := executedKey(contractIndex)
	_, contractAddrBytes := state.ImmutableWithdrawTree().Get(key)
	if contractAddrBytes == nil {
		return nil
	}
	addr, err := types.NewAddress(contractAddrBytes)
	if err != nil {
		log.
			WithField("contract_addr_bytes", cmn.HexBytes(contractAddrBytes)).
			WithError(err).
			Panic("Cannot reconstruct address from contract address bytes")
	}
	return addr
}

// GetUpdateExecutionWithProof returns the contract address and proof of an executed update
func GetUpdateExecutionWithProof(state context.IMutableState, contractIndex uint64, version int64) (*types.Address, *iavl.RangeProof) {
	key := executedKey(contractIndex)
	contractAddrBytes, proof, err := state.MutableWithdrawTree().GetVersionedWithProof(key, version)
	if err != nil || contractAddrBytes == nil {
		return nil, nil
	}
	addr, err := types.NewAddress(contractAddrBytes)
	if err != nil {
		log.
			WithField("contract_addr_bytes", cmn.HexBytes(contractAddrBytes)).
			WithError(err).
			Panic("Cannot reconstruct address from contract address bytes")
	}
	return addr, proof
}

// GetContractIndex returns the current contract index
func GetContractIndex(state context.IImmutableState) uint64 {
	_, bs := state.ImmutableStateTree().Get(currentContractIndexKey)
	if bs == nil {
		return 0
	}
	return utils.DecodeUint64(bs)
}

func increaseContractIndex(state context.IMutableState) {
	currentIndex := GetContractIndex(state)
	bs := utils.EncodeUint64(currentIndex + 1)
	state.MutableStateTree().Set(currentContractIndexKey, bs)
}
