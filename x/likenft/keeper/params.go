package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.Params{
		PriceDenom: k.PriceDenom(ctx),
	}
}

func (k Keeper) PriceDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.ParamKeyPriceDenom, &res)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
