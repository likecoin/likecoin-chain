package iscn

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, genesisState GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, genesisState.Params)
	// TODO: CIDs and ISCN related metadata
	return nil
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	// TODO: CIDs and ISCN related metadata
	return GenesisState{
		Params: params,
	}
}
