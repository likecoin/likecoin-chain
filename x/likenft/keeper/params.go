package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.Params{
		PriceDenom:             k.PriceDenom(ctx),
		FeePerByte:             k.FeePerByte(ctx),
		MaxOfferDurationDays:   k.maxOfferDurationDays(ctx),
		MaxListingDurationDays: k.maxListingDurationDays(ctx),
		MaxRoyaltyBasisPoints:  k.MaxRoyaltyBasisPoints(ctx),
	}
}

func (k Keeper) PriceDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.ParamKeyPriceDenom, &res)
	return
}

func (k Keeper) FeePerByte(ctx sdk.Context) (res sdk.DecCoin) {
	k.paramstore.Get(ctx, types.ParamKeyFeePerByte, &res)
	return
}

func (k Keeper) maxOfferDurationDays(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.ParamKeyMaxOfferDurationDays, &res)
	return
}

func (k Keeper) MaxOfferDuration(ctx sdk.Context) time.Duration {
	days := k.maxOfferDurationDays(ctx)
	return time.Duration(days) * 24 * time.Hour
}

func (k Keeper) MaxOfferDurationText(ctx sdk.Context) string {
	days := k.maxOfferDurationDays(ctx)
	return fmt.Sprintf("%d days", days)
}

func (k Keeper) maxListingDurationDays(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.ParamKeyMaxListingDurationDays, &res)
	return
}

func (k Keeper) MaxListingDuration(ctx sdk.Context) time.Duration {
	days := k.maxListingDurationDays(ctx)
	return time.Duration(days) * 24 * time.Hour
}

func (k Keeper) MaxListingDurationText(ctx sdk.Context) string {
	days := k.maxListingDurationDays(ctx)
	return fmt.Sprintf("%d days", days)
}

func (k Keeper) MaxRoyaltyBasisPoints(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.ParamKeyMaxRoyaltyBasisPoints, &res)
	return
}

func (k Keeper) MaxRoyaltyBasisPointsText(ctx sdk.Context) string {
	points := k.MaxRoyaltyBasisPoints(ctx)
	return fmt.Sprintf("%d (%.2f%%)", points, float64(points)/100)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
