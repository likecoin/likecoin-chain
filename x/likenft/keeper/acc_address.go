package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) assertBech32EqualsAccAddress(bech32 string, expectedAccAddress sdk.AccAddress) error {
	accAddress, err := sdk.AccAddressFromBech32(bech32)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !accAddress.Equals(expectedAccAddress) {
		return sdkerrors.ErrUnauthorized.Wrapf("User %s is not authorized. Expected %s.", accAddress.String(), expectedAccAddress.String())
	}
	return nil
}
