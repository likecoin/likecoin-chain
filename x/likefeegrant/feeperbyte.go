package likefeegrant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

func DeductFeePerByte(
	ctx sdk.Context,
	accountKeeper AccountKeeper, bankKeeper authtypes.BankKeeper, feegrantKeeper ante.FeegrantKeeper,
	msgSender sdk.AccAddress, fees sdk.Coins, msg sdk.Msg,
) error {
	if fees.IsZero() {
		return nil
	}

	feeTx := GetFeeTx(ctx)
	granter := msgSender
	if feeTx != nil {
		granter = feeTx.FeeGranter()
	}
	if granter == nil {
		granter = msgSender
	}
	if !granter.Equals(msgSender) {
		if feegrantKeeper == nil {
			return sdkerrors.ErrInvalidRequest.Wrap("error when deducting fee per byte: fee grants are not enabled")
		}
		err := feegrantKeeper.UseGrantedFees(ctx, granter, msgSender, fees, []sdk.Msg{msg})
		if err != nil {
			return sdkerrors.Wrapf(err, "error when deducting fee per byte: %s does not not allow to pay fees for %s", granter, msgSender)
		}
	}
	acc := accountKeeper.GetAccount(ctx, granter)
	if acc == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("error when deducting fee per byte: account %s not found", msgSender.String())
	}
	err := ante.DeductFees(bankKeeper, ctx, acc, fees)
	if err != nil {
		return sdkerrors.ErrInsufficientFee.Wrapf(err.Error())
	}
	return nil
}
