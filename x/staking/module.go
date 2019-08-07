package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/likecoin/likechain/x/whitelist"
)

var (
	_ module.AppModule = AppModule{}
)

type AppModule struct {
	staking.AppModule
	whitelistKeeper whitelist.Keeper
}

func NewAppModule(stakingKeeper staking.Keeper, distrKeeper stakingTypes.DistributionKeeper,
	accKeeper stakingTypes.AccountKeeper, supplyKeeper stakingTypes.SupplyKeeper,
	whitelistKeeper whitelist.Keeper) AppModule {
	return AppModule{
		AppModule:       staking.NewAppModule(stakingKeeper, distrKeeper, accKeeper, supplyKeeper),
		whitelistKeeper: whitelistKeeper,
	}
}

// proxy handler which intercepts MsgCreateValidator
func (am AppModule) NewHandler() sdk.Handler {
	stakingHandler := am.AppModule.NewHandler()
	wrappedHandler := whitelist.WrapStakingHandler(am.whitelistKeeper, stakingHandler)
	return wrappedHandler
}
