package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func validateOfferExpiration(ctx sdk.Context, expireTime time.Time) error {
	if expireTime.Before(ctx.BlockTime()) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Expiration is in the past")
	}

	if expireTime.After(ctx.BlockTime().Add(types.MaxOfferDuration)) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Expiration is too far in the future. Max offer duration is %s.", types.MaxOfferDurationText)
	}

	return nil
}
