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

	opWeightMsgCreateMintableNFT = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateMintableNFT int = 100

	opWeightMsgUpdateMintableNFT = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateMintableNFT int = 100

	opWeightMsgDeleteMintableNFT = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteMintableNFT int = 100

	opWeightMsgCreateOffer = "op_weight_msg_offer"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateOffer int = 100

	opWeightMsgUpdateOffer = "op_weight_msg_offer"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateOffer int = 100

	opWeightMsgDeleteOffer = "op_weight_msg_offer"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteOffer int = 100

	opWeightMsgCreateListing = "op_weight_msg_listing"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateListing int = 100

	opWeightMsgUpdateListing = "op_weight_msg_listing"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateListing int = 100

	opWeightMsgDeleteListing = "op_weight_msg_listing"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteListing int = 100

	opWeightMsgSellNFT = "op_weight_msg_sell_nft"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSellNFT int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	likenftGenesis := types.GenesisState{
		OfferList: []types.Offer{
			{
				ClassId: "0",
				NftId:   "0",
				Buyer:   sample.AccAddress(),
			},
			{
				ClassId: "1",
				NftId:   "1",
				Buyer:   sample.AccAddress(),
			},
		},
		ListingList: []types.Listing{
			{
				ClassId: "0",
				NftId:   "0",
				Seller:  "0",
			},
			{
				ClassId: "1",
				NftId:   "1",
				Seller:  "1",
			},
		},
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

	var weightMsgCreateMintableNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateMintableNFT, &weightMsgCreateMintableNFT, nil,
		func(_ *rand.Rand) {
			weightMsgCreateMintableNFT = defaultWeightMsgCreateMintableNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateMintableNFT,
		likenftsimulation.SimulateMsgCreateMintableNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateMintableNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateMintableNFT, &weightMsgUpdateMintableNFT, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateMintableNFT = defaultWeightMsgUpdateMintableNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateMintableNFT,
		likenftsimulation.SimulateMsgUpdateMintableNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteMintableNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteMintableNFT, &weightMsgDeleteMintableNFT, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteMintableNFT = defaultWeightMsgDeleteMintableNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteMintableNFT,
		likenftsimulation.SimulateMsgDeleteMintableNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateOffer int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateOffer, &weightMsgCreateOffer, nil,
		func(_ *rand.Rand) {
			weightMsgCreateOffer = defaultWeightMsgCreateOffer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateOffer,
		likenftsimulation.SimulateMsgCreateOffer(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateOffer int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateOffer, &weightMsgUpdateOffer, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateOffer = defaultWeightMsgUpdateOffer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateOffer,
		likenftsimulation.SimulateMsgUpdateOffer(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteOffer int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteOffer, &weightMsgDeleteOffer, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteOffer = defaultWeightMsgDeleteOffer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteOffer,
		likenftsimulation.SimulateMsgDeleteOffer(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateListing int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateListing, &weightMsgCreateListing, nil,
		func(_ *rand.Rand) {
			weightMsgCreateListing = defaultWeightMsgCreateListing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateListing,
		likenftsimulation.SimulateMsgCreateListing(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateListing int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateListing, &weightMsgUpdateListing, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateListing = defaultWeightMsgUpdateListing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateListing,
		likenftsimulation.SimulateMsgUpdateListing(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteListing int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteListing, &weightMsgDeleteListing, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteListing = defaultWeightMsgDeleteListing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteListing,
		likenftsimulation.SimulateMsgDeleteListing(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSellNFT int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSellNFT, &weightMsgSellNFT, nil,
		func(_ *rand.Rand) {
			weightMsgSellNFT = defaultWeightMsgSellNFT
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSellNFT,
		likenftsimulation.SimulateMsgSellNFT(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
