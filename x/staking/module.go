package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	_ module.AppModule = AppModule{}
)

type AppModule struct {
	staking.AppModule
	querier *Querier
	keeper  *keeper.Keeper
}

func NewAppModule(
	cdc codec.Codec, k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper,
	querier *Querier,
) AppModule {
	oldAppModule := staking.NewAppModule(cdc, k, ak, bk)
	return AppModule{
		AppModule: oldAppModule,
		querier:   querier,
		keeper:    &k,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(*am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.querier)

	m := keeper.NewMigrator(*am.keeper)
	cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	am.querier.BeginWriteIndex(ctx)
	am.AppModule.BeginBlock(ctx, req)
}

func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	res := am.AppModule.EndBlock(ctx, req)
	am.querier.CommitWriteIndex(ctx)
	return res
}
