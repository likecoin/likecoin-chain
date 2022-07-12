package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func (k Keeper) validateRoyaltyConfig(ctx sdk.Context, config types.RoyaltyConfig) error {
	if config.RateBasisPoints > k.MaxRoyaltyBasisPoints(ctx) {
		return types.ErrInvalidRoyaltyConfig.Wrapf("Royalty basis points cannot be greater than %s", k.MaxRoyaltyBasisPointsText(ctx))
	}
	return nil
}

func (k Keeper) validateRoyaltyConfigInput(ctx sdk.Context, input types.RoyaltyConfigInput) error {
	if input.RateBasisPoints > k.MaxRoyaltyBasisPoints(ctx) {
		return types.ErrInvalidRoyaltyConfig.Wrapf("Royalty basis points cannot be greater than %s", k.MaxRoyaltyBasisPointsText(ctx))
	}
	for _, stakeholder := range input.Stakeholders {
		_, err := sdk.AccAddressFromBech32(stakeholder.Account)
		if err != nil {
			return types.ErrInvalidRoyaltyConfig.Wrapf("Stakeholder address %s is invalid", stakeholder.Account)
		}
	}
	return nil
}
