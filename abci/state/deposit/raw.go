package deposit

import (
	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	cmn "github.com/tendermint/tendermint/libs/common"
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

// Proposal represents all the LikeCoin Transfer events on a single Ethereum block
type Proposal struct {
	BlockNumber uint64
	Inputs      []Input
}

var (
	depositApproversKey          = []byte("depositApprovers")
	depositApproversWeightSumKey = []byte("depositApproversWeightSum")

	depositProposalKey = []byte("deposit")
	depositWeightKey   = []byte("depositWeight")
	depositApprovalKey = []byte("depositApproval")
	depositExecutedKey = []byte("depositExecuted")

	log = logger.L
)

func approvalKey(id *types.LikeChainID, blockNumber uint64) []byte {
	return utils.JoinKeys([][]byte{
		depositApprovalKey,
		id.Bytes(),
		[]byte("block"),
		utils.EncodeUint64(blockNumber),
	})
}

func proposalKey(txHash []byte) []byte {
	return utils.JoinKeys([][]byte{
		depositProposalKey,
		txHash,
	})
}

func weightKey(txHash []byte) []byte {
	return utils.JoinKeys([][]byte{
		depositWeightKey,
		txHash,
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

// SetDepositProposal records a deposit proposal into state tree
func SetDepositProposal(state context.IMutableState, txHash []byte, proposal Proposal) {
	bs, err := types.AminoCodec().MarshalBinaryBare(proposal)
	if err != nil {
		log.
			WithField("proposal", proposal).
			WithError(err).
			Panic("Unable to marshal deposit proposal")
	}
	key := proposalKey(txHash)
	state.MutableStateTree().Set(key, bs)
}

// GetDepositProposal loads a deposit proposal from state tree
func GetDepositProposal(state context.IImmutableState, txHash []byte) *Proposal {
	key := proposalKey(txHash)
	_, bs := state.ImmutableStateTree().Get(key)
	if bs == nil {
		return nil
	}
	result := Proposal{}
	err := types.AminoCodec().UnmarshalBinaryBare(bs, &result)
	if err != nil {
		log.
			WithField("data", cmn.HexBytes(bs)).
			WithError(err).
			Panic("Unable to unmarshal deposit proposal")
	}
	return &result
}

// SetDepositApproval records a deposit approval into state tree
func SetDepositApproval(state context.IMutableState, approver *types.LikeChainID, blockNumber uint64, txHash []byte) {
	key := approvalKey(approver, blockNumber)
	state.MutableStateTree().Set(key, txHash)
}

// GetDepositApproval returns a DepositApprover's approved txHash for a block number
func GetDepositApproval(state context.IImmutableState, approver *types.LikeChainID, blockNumber uint64) []byte {
	key := approvalKey(approver, blockNumber)
	_, txHash := state.ImmutableStateTree().Get(key)
	return txHash
}

// IncreaseDepositProposalWeight initializes or increments a deposit proposal's approve weight, returns the new weight
func IncreaseDepositProposalWeight(state context.IMutableState, txHash []byte, weight uint32) uint64 {
	key := weightKey(txHash)
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
func GetDepositProposalWeight(state context.IImmutableState, txHash []byte) uint64 {
	key := weightKey(txHash)
	_, bs := state.ImmutableStateTree().Get(key)
	if bs == nil {
		return 0
	}
	return utils.DecodeUint64(bs)
}

// SetDepositExecution records a deposit execution with block number into state tree
func SetDepositExecution(state context.IMutableState, blockNumber uint64, txHash []byte) {
	key := executedKey(blockNumber)
	state.MutableStateTree().Set(key, txHash)
}

// GetDepositExecution returns the executed txHash for a block number
func GetDepositExecution(state context.IImmutableState, blockNumber uint64) []byte {
	key := executedKey(blockNumber)
	_, txHash := state.ImmutableStateTree().Get(key)
	return txHash
}
