package likenft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the classesByISCN
	for _, elem := range genState.ClassesByISCNList {
		k.SetClassesByISCN(ctx, elem)
	}
	// Set all the classesByAccount
	for _, elem := range genState.ClassesByAccountList {
		k.SetClassesByAccount(ctx, elem)
	}
	// Set all the mintableNFT
	for _, elem := range genState.MintableNFTList {
		k.SetMintableNFT(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ClassesByISCNList = k.GetAllClassesByISCN(ctx)
	genesis.ClassesByAccountList = k.GetAllClassesByAccount(ctx)
	genesis.MintableNFTList = k.GetAllMintableNFT(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
