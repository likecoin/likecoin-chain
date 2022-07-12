package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) validateListingExpiration(ctx sdk.Context, expireTime time.Time) error {
	if expireTime.Before(ctx.BlockTime()) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Expiration is in the past")
	}

	if expireTime.After(ctx.BlockTime().Add(k.MaxListingDuration(ctx))) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Expiration is too far in the future. Max listing duration is %s.", k.MaxListingDurationText(ctx))
	}

	return nil
}
