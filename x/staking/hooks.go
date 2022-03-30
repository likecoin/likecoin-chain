package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type IndexingHooks struct {
	querier *Querier
}

var _ types.StakingHooks = IndexingHooks{}

func NewHooks(querier *Querier) IndexingHooks {
	return IndexingHooks{
		querier: querier,
	}
}

func (hooks IndexingHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	hooks.querier.RemoveIndex(ctx, delAddr, valAddr)
}

func (hooks IndexingHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	hooks.querier.UpdateIndex(ctx, delAddr, valAddr)
}

func (IndexingHooks) BeforeDelegationSharesModified(sdk.Context, sdk.AccAddress, sdk.ValAddress) {}
func (IndexingHooks) BeforeDelegationCreated(sdk.Context, sdk.AccAddress, sdk.ValAddress)        {}
func (IndexingHooks) AfterValidatorCreated(sdk.Context, sdk.ValAddress)                          {}
func (IndexingHooks) BeforeValidatorModified(sdk.Context, sdk.ValAddress)                        {}
func (IndexingHooks) AfterValidatorRemoved(sdk.Context, sdk.ConsAddress, sdk.ValAddress)         {}
func (IndexingHooks) AfterValidatorBonded(sdk.Context, sdk.ConsAddress, sdk.ValAddress)          {}
func (IndexingHooks) AfterValidatorBeginUnbonding(sdk.Context, sdk.ConsAddress, sdk.ValAddress)  {}
func (IndexingHooks) BeforeValidatorSlashed(sdk.Context, sdk.ValAddress, sdk.Dec)                {}

type DebugHooks struct{}

func (DebugHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	ctx.Logger().Debug(
		"ValidatorDelegations indexer: BeforeDelegationRemoved",
		"delegator_addr", delAddr.String(),
		"validator_addr", valAddr.String(),
	)
}

func (DebugHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	ctx.Logger().Debug(
		"ValidatorDelegations indexer: AfterDelegationModified",
		"delegator_addr", delAddr.String(),
		"validator_addr", valAddr.String(),
	)
}
func (DebugHooks) BeforeDelegationSharesModified(sdk.Context, sdk.AccAddress, sdk.ValAddress) {}
func (DebugHooks) BeforeDelegationCreated(sdk.Context, sdk.AccAddress, sdk.ValAddress)        {}
func (DebugHooks) AfterValidatorCreated(sdk.Context, sdk.ValAddress)                          {}
func (DebugHooks) BeforeValidatorModified(sdk.Context, sdk.ValAddress)                        {}
func (DebugHooks) AfterValidatorRemoved(sdk.Context, sdk.ConsAddress, sdk.ValAddress)         {}
func (DebugHooks) AfterValidatorBonded(sdk.Context, sdk.ConsAddress, sdk.ValAddress)          {}
func (DebugHooks) AfterValidatorBeginUnbonding(sdk.Context, sdk.ConsAddress, sdk.ValAddress)  {}
func (DebugHooks) BeforeValidatorSlashed(sdk.Context, sdk.ValAddress, sdk.Dec)                {}
