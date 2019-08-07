package whitelist

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, genesisState GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, genesisState.Params)
	keeper.SetWhitelist(ctx, genesisState.Whitelist)
	return nil
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	whitelist := keeper.GetWhitelist(ctx)
	return GenesisState{
		Params:    params,
		Whitelist: whitelist,
	}
}
