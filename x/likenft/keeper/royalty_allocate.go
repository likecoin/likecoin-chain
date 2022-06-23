package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/likechain/x/likenft/types"
)

type RoyaltyAllocation struct {
	Account sdk.AccAddress
	Amount  uint64
}

func (k Keeper) ComputeRoyaltyAllocation(ctx sdk.Context, txnAmount uint64, config types.RoyaltyConfig) (royaltyAmount uint64, allocations []RoyaltyAllocation, err error) {
	if config.RateBasisPoints > k.MaxRoyaltyBasisPoints(ctx) {
		err = types.ErrInvalidRoyaltyConfig.Wrapf("Royalty basis points cannot be greater than %s", k.MaxRoyaltyBasisPointsText(ctx))
		return
	}
	// max allocable amount
	allocatable := uint64(math.Floor(float64(txnAmount) / float64(10000) * float64(config.RateBasisPoints)))
	// sum total weight
	totalWeight := uint64(0)
	for _, stakeholder := range config.Stakeholders {
		totalWeight += stakeholder.Weight
	}
	// split by weights
	for _, stakeholder := range config.Stakeholders {
		amount := uint64(math.Floor(float64(allocatable) / float64(totalWeight) * float64(stakeholder.Weight)))
		allocations = append(allocations, RoyaltyAllocation{
			Account: stakeholder.Account,
			Amount:  amount,
		})
		royaltyAmount += amount
	}
	return
}
