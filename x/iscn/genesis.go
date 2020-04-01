package iscn

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/iscn/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, genesisState GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, genesisState.Params)
	for _, record := range genesisState.IscnRecords {
		keeper.SetIscnRecord(ctx, record.Id, &record.Record)
	}
	return nil
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	records := []types.IscnPair{}
	keeper.IterateIscnRecords(ctx, func(id []byte, record *IscnRecord) bool {
		records = append(records, types.IscnPair{
			Id:     id,
			Record: *record,
		})
		return false
	})
	return GenesisState{
		Params:      params,
		IscnRecords: records,
	}
}
