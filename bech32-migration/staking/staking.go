package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/likecoin/likechain/bech32-migration/utils"
)

func MigrateAddressBech32(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) {
	ctx.Logger().Info("Migration of address bech32 for staking module begin")
	validatorCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.ValidatorsKey, func(bz []byte) []byte {
		validator := types.MustUnmarshalValidator(cdc, bz)
		validator.OperatorAddress = utils.ConvertValAddr(validator.OperatorAddress)
		validatorCount++
		return types.MustMarshalValidator(cdc, &validator)
	})
	delegationCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.DelegationKey, func(bz []byte) []byte {
		delegation := types.MustUnmarshalDelegation(cdc, bz)
		delegation.DelegatorAddress = utils.ConvertAccAddr(delegation.DelegatorAddress)
		delegation.ValidatorAddress = utils.ConvertValAddr(delegation.ValidatorAddress)
		delegationCount++
		return types.MustMarshalDelegation(cdc, delegation)
	})
	redelegationCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.RedelegationKey, func(bz []byte) []byte {
		redelegation := types.MustUnmarshalRED(cdc, bz)
		redelegation.DelegatorAddress = utils.ConvertAccAddr(redelegation.DelegatorAddress)
		redelegation.ValidatorSrcAddress = utils.ConvertValAddr(redelegation.ValidatorSrcAddress)
		redelegation.ValidatorDstAddress = utils.ConvertValAddr(redelegation.ValidatorDstAddress)
		redelegationCount++
		return types.MustMarshalRED(cdc, redelegation)
	})
	unbondingDelegationCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.UnbondingDelegationKey, func(bz []byte) []byte {
		unbonding := types.MustUnmarshalUBD(cdc, bz)
		unbonding.DelegatorAddress = utils.ConvertAccAddr(unbonding.DelegatorAddress)
		unbonding.ValidatorAddress = utils.ConvertValAddr(unbonding.ValidatorAddress)
		unbondingDelegationCount++
		return types.MustMarshalUBD(cdc, unbonding)
	})
	// Not migrating historical info, since it's weird to change "history"
	// historicalInfoCount := uint64(0)
	// utils.IterateStoreByPrefix(ctx, storeKey, types.HistoricalInfoKey, func(bz []byte) []byte {
	// 	historicalInfo := types.MustUnmarshalHistoricalInfo(cdc, bz)
	// 	for i := range historicalInfo.Valset {
	// 		historicalInfo.Valset[i].OperatorAddress = utils.ConvertValAddr(historicalInfo.Valset[i].OperatorAddress)
	// 	}
	// 	historicalInfoCount++
	// 	return cdc.MustMarshal(&historicalInfo)
	// })
	ctx.Logger().Info(
		"Migration of address bech32 for staking module done",
		"validator_count", validatorCount,
		"delegation_count", delegationCount,
		"redelegation_count", redelegationCount,
		"unbonding_delegation_count", unbondingDelegationCount,
		// "historical_info_count", historicalInfoCount,
	)
}
