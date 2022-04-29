package staking

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	originalstaking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

const FlagEnableCustomIndex = "enable-custom-index"

func AddEnableCustomIndexFlag(startCmd *cobra.Command) {
	startCmd.Flags().Bool(FlagEnableCustomIndex, false, "Enable custom index (currently for ValidatorDelegations only)")
}

func setupOriginalStakingModule(
	homePath string, appCodec codec.Codec,
	stakingKeeper *keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper,
	distrKeeper distrkeeper.Keeper, slashingKeeper slashingkeeper.Keeper,
) (*keeper.Keeper, module.AppModule) {
	stakingKeeper = stakingKeeper.SetHooks(
		types.NewMultiStakingHooks(
			distrKeeper.Hooks(),
			slashingKeeper.Hooks(),
		),
	)
	am := originalstaking.NewAppModule(appCodec, *stakingKeeper, ak, bk)
	return stakingKeeper, am
}

func setupStakingModuleWithCustomIndex(
	homePath string, appCodec codec.Codec,
	stakingKeeper *keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper,
	distrKeeper distrkeeper.Keeper, slashingKeeper slashingkeeper.Keeper,
) (*keeper.Keeper, module.AppModule) {
	stakingIndexDB, err := dbm.NewGoLevelDB("index", path.Join(homePath, "data"))
	if err != nil {
		panic(fmt.Errorf("failed to create indexing DB for staking module: %s", err))
	}
	// TODO: When to close this DB?
	// This should be handled by the finalizer set in runtime so should not be a problem.
	// But we still want to do something in code
	stakingQuerier := NewQuerier(stakingKeeper, appCodec, stakingIndexDB)
	stakingKeeper = stakingKeeper.SetHooks(
		types.NewMultiStakingHooks(
			distrKeeper.Hooks(),
			slashingKeeper.Hooks(),
			NewHooks(stakingQuerier),
		),
	)
	am := NewAppModule(appCodec, *stakingKeeper, ak, bk, stakingQuerier)
	return stakingKeeper, am
}

func SetupStakingModule(
	homePath string, appCodec codec.Codec,
	stakingKeeper *keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper,
	distrKeeper distrkeeper.Keeper, slashingKeeper slashingkeeper.Keeper,
	appOpts servertypes.AppOptions,
) (*keeper.Keeper, module.AppModule) {
	opt := appOpts.Get(FlagEnableCustomIndex)
	if shouldEnableCustomIndex, ok := opt.(bool); ok && shouldEnableCustomIndex {
		return setupStakingModuleWithCustomIndex(homePath, appCodec, stakingKeeper, ak, bk, distrKeeper, slashingKeeper)
	} else {
		return setupOriginalStakingModule(homePath, appCodec, stakingKeeper, ak, bk, distrKeeper, slashingKeeper)
	}
}
