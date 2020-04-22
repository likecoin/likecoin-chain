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
	for _, author := range genesisState.Authors {
		keeper.SetAuthor(ctx, &author)
	}
	keeper.SetIscnCount(ctx, uint64(len(genesisState.IscnRecords)))
	return nil
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	authors := []Author{}
	keeper.IterateAuthors(ctx, func(_ []byte, author *Author) bool {
		authors = append(authors, *author)
		return false
	})
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
		Authors:     authors,
		IscnRecords: records,
	}
}
