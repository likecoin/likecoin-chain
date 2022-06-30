package testutil

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/likecoin/likecoin-chain/v3/bech32-migration/utils"
)

func AssertStakingAddressBech32(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) bool {
	ok := true
	utils.IterateStoreByPrefix(ctx, storeKey, types.ValidatorsKey, func(bz []byte) []byte {
		validator := types.MustUnmarshalValidator(cdc, bz)
		if !strings.HasPrefix(validator.OperatorAddress, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: validator.OperatorAddress", validator.OperatorAddress))
			ok = false
		}
		return bz
	})
	utils.IterateStoreByPrefix(ctx, storeKey, types.DelegationKey, func(bz []byte) []byte {
		delegation := types.MustUnmarshalDelegation(cdc, bz)
		if !strings.HasPrefix(delegation.DelegatorAddress, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: delegation.DelegatorAddress", delegation.DelegatorAddress))
			ok = false
		}
		if !strings.HasPrefix(delegation.ValidatorAddress, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: delegation.ValidatorAddress", delegation.ValidatorAddress))
			ok = false
		}
		return bz
	})
	utils.IterateStoreByPrefix(ctx, storeKey, types.RedelegationKey, func(bz []byte) []byte {
		redelegation := types.MustUnmarshalRED(cdc, bz)
		if !strings.HasPrefix(redelegation.DelegatorAddress, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: redelegation.DelegatorAddress", redelegation.DelegatorAddress))
			ok = false
		}
		if !strings.HasPrefix(redelegation.ValidatorSrcAddress, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: redelegation.ValidatorSrcAddress", redelegation.ValidatorSrcAddress))
			ok = false
		}
		if !strings.HasPrefix(redelegation.ValidatorDstAddress, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: redelegation.ValidatorDstAddress", redelegation.ValidatorDstAddress))
			ok = false
		}
		return bz
	})
	utils.IterateStoreByPrefix(ctx, storeKey, types.UnbondingDelegationKey, func(bz []byte) []byte {
		unbonding := types.MustUnmarshalUBD(cdc, bz)
		if !strings.HasPrefix(unbonding.DelegatorAddress, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: unbonding.DelegatorAddress", unbonding.DelegatorAddress))
			ok = false
		}
		if !strings.HasPrefix(unbonding.ValidatorAddress, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: unbonding.ValidatorAddress", unbonding.ValidatorAddress))
			ok = false
		}
		return types.MustMarshalUBD(cdc, unbonding)
	})
	utils.IterateStoreByPrefix(ctx, storeKey, types.HistoricalInfoKey, func(bz []byte) []byte {
		historicalInfo := types.MustUnmarshalHistoricalInfo(cdc, bz)
		for i := range historicalInfo.Valset {
			if !strings.HasPrefix(historicalInfo.Valset[i].OperatorAddress, "like") {
				ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: historicalInfo.Valset[i].OperatorAddress", historicalInfo.Valset[i].OperatorAddress))
				ok = false
			}
		}
		return cdc.MustMarshal(&historicalInfo)
	})
	return ok
}
