package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

func (k Keeper) DeductFeePerByte(ctx sdk.Context, feePayer sdk.AccAddress, bytesLength int) error {
	feePerByte := k.GetParams(ctx).FeePerByte
	amount := feePerByte.Amount.MulInt64(int64(bytesLength))
	fees := sdk.NewCoins(sdk.NewCoin(feePerByte.Denom, amount.Ceil().RoundInt()))
	if fees.IsZero() {
		return nil
	}
	acc := k.accountKeeper.GetAccount(ctx, feePayer)
	if acc == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("account %s not found", feePayer.String())
	}
	err := ante.DeductFees(k.bankKeeper, ctx, acc, fees)
	if err != nil {
		return sdkerrors.ErrInsufficientFee.Wrapf(err.Error())
	}
	return nil
}
