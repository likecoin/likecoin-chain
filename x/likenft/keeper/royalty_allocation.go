package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func (k Keeper) ComputeRoyaltyAllocation(ctx sdk.Context, txnAmount uint64, fullPayToRoyalty bool, config types.RoyaltyConfig) (royaltyAmount uint64, allocations []types.RoyaltyAllocation, err error) {
	err = k.validateRoyaltyConfig(ctx, config)
	if err != nil {
		return
	}
	// max allocable amount
	rateBasisPoint := config.RateBasisPoints
	if fullPayToRoyalty {
		rateBasisPoint = 10000
	}
	allocatable := uint64(math.Floor(float64(txnAmount) / float64(10000) * float64(rateBasisPoint)))
	if allocatable <= 0 {
		return
	}
	// sum total weight
	totalWeight := uint64(0)
	for _, stakeholder := range config.Stakeholders {
		totalWeight += stakeholder.Weight
	}
	if totalWeight <= 0 {
		return
	}
	// split by weights
	for _, stakeholder := range config.Stakeholders {
		amount := uint64(math.Floor(float64(allocatable) / float64(totalWeight) * float64(stakeholder.Weight)))
		if amount > 0 {
			allocations = append(allocations, types.RoyaltyAllocation{
				Account: stakeholder.Account,
				Amount:  amount,
			})
			royaltyAmount += amount
		}
	}
	return
}
