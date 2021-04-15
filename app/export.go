package app

import (
	"log"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// export the state of likechain for a genesis file
func (app *LikeApp) ExportAppStateAndValidators(forZeroHeight bool, jailAllowedAddrs []string,
) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	height := app.LastBlockHeight() + 1

	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}

	genState := app.mm.ExportGenesis(ctx, app.appCodec)
	appState, err := codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	validators, err := staking.WriteValidators(ctx, app.stakingKeeper)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, nil
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
//      in favour of export at a block height
func (app *LikeApp) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) {
	applyAllowedAddrs := false

	// Check if there is a whitelist
	if len(jailAllowedAddrs) > 0 {
		applyAllowedAddrs = true
	}

	allowedAddrsMap := make(map[string]bool)

	for _, addr := range jailAllowedAddrs {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			log.Fatal(err)
		}
		allowedAddrsMap[addr] = true
	}

	/* Just to be safe, assert the invariants on current state. */
	app.crisisKeeper.AssertInvariants(ctx)

	/* Handle fee distribution state. */

	// withdraw all validator commission
	app.stakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		_, _ = app.distrKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		return false
	})

	// withdraw all delegator rewards
	dels := app.stakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		_, err := app.distrKeeper.WithdrawDelegationRewards(ctx, delegation.GetDelegatorAddr(), delegation.GetValidatorAddr())
		if err != nil {
			panic(err)
		}
	}

	// clear validator slash events
	app.distrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	app.distrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators
	app.stakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps := app.distrKeeper.GetValidatorOutstandingRewardsCoins(ctx, val.GetOperator())
		feePool := app.distrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.distrKeeper.SetFeePool(ctx, feePool)

		app.distrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	// reinitialize all delegations
	for _, del := range dels {
		app.distrKeeper.Hooks().BeforeDelegationCreated(
			ctx, del.GetDelegatorAddr(), del.GetValidatorAddr(),
		)
		app.distrKeeper.Hooks().AfterDelegationModified(
			ctx, del.GetDelegatorAddr(), del.GetValidatorAddr(),
		)
	}

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	app.stakingKeeper.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		app.stakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	// iterate through unbonding delegations, reset creation height
	app.stakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		app.stakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.keys[stakingtypes.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])
		validator, found := app.stakingKeeper.GetValidator(ctx, addr)
		if !found {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
			validator.Jailed = true
		}

		app.stakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)

	/* Handle slashing state. */

	// reset start height on signing infos
	app.slashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			app.slashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
}
