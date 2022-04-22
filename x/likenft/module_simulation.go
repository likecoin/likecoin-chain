package likenft

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/likecoin/likechain/testutil/sample"
	likenftsimulation "github.com/likecoin/likechain/x/likenft/simulation"
	"github.com/likecoin/likechain/x/likenft/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = likenftsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgNewClass = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewClass int = 100

	opWeightMsgUpdateClass = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateClass int = 100

	opWeightMsgMintNFT = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgMintNFT int = 100

	opWeightMsgBurnNFT = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBurnNFT int = 100

	opWeightMsgCreateClaimableNFT = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateClaimableNFT int = 100

	opWeightMsgUpdateClaimableNFT = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateClaimableNFT int = 100

	opWeightMsgDeleteClaimableNFT = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteClaimableNFT int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	likenftGenesis := types.GenesisState{
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&likenftGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgNewClass int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgNewClass, &weightMsgNewClass, nil,
		func(_ *rand.Rand) {
			weightMsgNewClass = defaultWeightMsgNewClass
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewClass,
		likenftsimulation.SimulateMsgNewClass(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateClass int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateClass, &weightMsgUpdateClass, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateClass = defaultWeightMsgUpdateClass
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateClass,
		likenftsimulation.SimulateMsgUpdateClass(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgMintNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgMintNFT, &weightMsgMintNFT, nil,
		func(_ *rand.Rand) {
			weightMsgMintNFT = defaultWeightMsgMintNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintNFT,
		likenftsimulation.SimulateMsgMintNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgBurnNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgBurnNFT, &weightMsgBurnNFT, nil,
		func(_ *rand.Rand) {
			weightMsgBurnNFT = defaultWeightMsgBurnNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBurnNFT,
		likenftsimulation.SimulateMsgBurnNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateClaimableNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateClaimableNFT, &weightMsgCreateClaimableNFT, nil,
		func(_ *rand.Rand) {
			weightMsgCreateClaimableNFT = defaultWeightMsgCreateClaimableNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateClaimableNFT,
		likenftsimulation.SimulateMsgCreateClaimableNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateClaimableNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateClaimableNFT, &weightMsgUpdateClaimableNFT, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateClaimableNFT = defaultWeightMsgUpdateClaimableNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateClaimableNFT,
		likenftsimulation.SimulateMsgUpdateClaimableNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteClaimableNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteClaimableNFT, &weightMsgDeleteClaimableNFT, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteClaimableNFT = defaultWeightMsgDeleteClaimableNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteClaimableNFT,
		likenftsimulation.SimulateMsgDeleteClaimableNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
