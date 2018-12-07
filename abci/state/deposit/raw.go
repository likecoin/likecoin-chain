package deposit

import (
	"bytes"
	"sort"

	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
)

var (
	depositApproversKey          = []byte("depositApprovers")
	depositApproversWeightSumKey = []byte("depositApproversWeightSum")

	depositWeightKey   = []byte("depositWeight")
	depositApprovalKey = []byte("depositApproval")
	depositExecutedKey = []byte("depositExecuted")

	log = logger.L
)

// Approver represents the identity and weight of a deposit approver
type Approver struct {
	ID     *types.LikeChainID
	Weight uint32
}

// Input represents a LikeCoin Transfer event on Ethereum
type Input struct {
	FromAddr types.Address
	Value    types.BigInt
}

// Validate checks if a deposit input is valid
func (input Input) Validate() bool {
	return input.Value.Int != nil && input.Value.IsWithinRange()
}

// ProposalInputs represents the inputs of a proposal
type ProposalInputs []Input

// Validate checks if deposit inputs are valid
func (inputs ProposalInputs) Validate() bool {
	for _, input := range inputs {
		if !input.Validate() {
			return false
		}
	}
	return true
}

// Len implements interface for sort
func (inputs ProposalInputs) Len() int {
	return len(inputs)
}

// Swap implements interface for sort
func (inputs ProposalInputs) Swap(i, j int) {
	inputs[i], inputs[j] = inputs[j], inputs[i]
}

// Less implements interface for sort
func (inputs ProposalInputs) Less(i, j int) bool {
	switch bytes.Compare(inputs[i].FromAddr[:], inputs[j].FromAddr[:]) {
	case -1:
		return true
	case 1:
		return false
	default:
		return inputs[i].Value.Cmp(inputs[j].Value.Int) < 0
	}
}

// Proposal represents all the LikeCoin Transfer events on a single Ethereum block
type Proposal struct {
	BlockNumber uint64
	Inputs      ProposalInputs
}

// Validate checks if a deposit proposal is valid
func (proposal Proposal) Validate() bool {
	if len(proposal.Inputs) == 0 {
		return false
	}
	return proposal.Inputs.Validate()
}

// Sort sorts the Inputs in proposal
func (proposal *Proposal) Sort() {
	sort.Sort(proposal.Inputs)
}

// Hash computes the hash of the proposal
func (proposal *Proposal) Hash() []byte {
	proposal.Sort()
	bs, err := types.AminoCodec().MarshalBinaryBare(proposal)
	if err != nil {
		log.
			WithField("proposal", proposal).
			Panic("Cannot encode deposit proposal")
	}
	return tmhash.Sum(bs)
}

func approvalKey(id *types.LikeChainID, proposalHash []byte) []byte {
	return utils.JoinKeys([][]byte{
		depositApprovalKey,
		id.Bytes(),
		[]byte("approval"),
		proposalHash,
	})
}

func weightKey(proposalHash []byte) []byte {
	return utils.JoinKeys([][]byte{
		depositWeightKey,
		proposalHash,
	})
}

func executedKey(blockNumber uint64) []byte {
	return utils.JoinKeys([][]byte{
		depositExecutedKey,
		utils.EncodeUint64(blockNumber),
	})
}

// GetDepositApproversWeightSum loads the weight sum of the deposit approvers from state tree
func GetDepositApproversWeightSum(state context.IImmutableState) uint64 {
	_, bs := state.ImmutableStateTree().Get(depositApproversWeightSumKey)
	if bs == nil {
		return 0
	}
	if len(bs) != 8 {
		log.
			WithField("data", cmn.HexBytes(bs)).
			Panic("Invalid deposit weight raw data")
	}
	return utils.DecodeUint64(bs)
}

// GetDepositApprovers loads the deposit approver list from state tree
func GetDepositApprovers(state context.IImmutableState) (approvers []Approver) {
	_, bs := state.ImmutableStateTree().Get(depositApproversKey)
	if bs == nil {
		return nil
	}
	err := types.AminoCodec().UnmarshalBinaryBare(bs, &approvers)
	if err != nil {
		log.
			WithField("data", cmn.HexBytes(bs)).
			WithError(err).
			Panic("Cannot unmarshal deposit approvers")
	}
	return approvers
}

// SetDepositApprovers saves the deposit approver list into state tree
func SetDepositApprovers(state context.IMutableState, approvers []Approver) {
	if len(approvers) == 0 {
		state.MutableStateTree().Remove(depositApproversKey)
		state.MutableStateTree().Remove(depositApproversWeightSumKey)
		return
	}
	totalWeight := uint64(0)
	for _, approver := range approvers {
		totalWeight += uint64(approver.Weight)
	}
	bs, err := types.AminoCodec().MarshalBinaryBare(approvers)
	if err != nil {
		log.
			WithField("approvers", approvers).
			WithError(err).
			Panic("Cannot marshal deposit approvers")
	}
	state.MutableStateTree().Set(depositApproversKey, bs)
	state.MutableStateTree().Set(depositApproversWeightSumKey, utils.EncodeUint64(totalWeight))
}

// setDepositApproval records a deposit approval into state tree
func setDepositApproval(state context.IMutableState, approver *types.LikeChainID, proposalHash []byte) {
	if HasApprovedDeposit(state, approver, proposalHash) {
		log.
			WithField("approver", approver.String()).
			WithField("proposal_hash", cmn.HexBytes(proposalHash)).
			Panic("Double approving the same proposal hash")
	}
	key := approvalKey(approver, proposalHash)
	state.MutableStateTree().Set(key, []byte{1})
}

// HasApprovedDeposit returns a DepositApprover's approved proposalHash for a block number
func HasApprovedDeposit(state context.IImmutableState, approver *types.LikeChainID, proposalHash []byte) bool {
	key := approvalKey(approver, proposalHash)
	_, v := state.ImmutableStateTree().Get(key)
	return v != nil
}

// IncreaseDepositProposalWeight initializes or increments a deposit proposal's approve weight, returns the new weight
func IncreaseDepositProposalWeight(state context.IMutableState, proposalHash []byte, weight uint32) uint64 {
	key := weightKey(proposalHash)
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

// GetDepositProposalWeight returns the weight sum of approvers approved this proposal
func GetDepositProposalWeight(state context.IImmutableState, proposalHash []byte) uint64 {
	key := weightKey(proposalHash)
	_, bs := state.ImmutableStateTree().Get(key)
	if bs == nil {
		return 0
	}
	return utils.DecodeUint64(bs)
}

// setDepositExecution records a deposit execution with block number into state tree
func setDepositExecution(state context.IMutableState, blockNumber uint64, proposalHash []byte) {
	if GetDepositExecution(state, blockNumber) != nil {
		log.
			WithField("block_number", blockNumber).
			WithField("proposal_hash", cmn.HexBytes(proposalHash)).
			Panic("Double setting deposit execution on the same block number")
	}
	key := executedKey(blockNumber)
	state.MutableStateTree().Set(key, proposalHash)
}

// GetDepositExecution returns the executed proposalHash for a block number
func GetDepositExecution(state context.IImmutableState, blockNumber uint64) []byte {
	key := executedKey(blockNumber)
	_, proposalHash := state.ImmutableStateTree().Get(key)
	return proposalHash
}
