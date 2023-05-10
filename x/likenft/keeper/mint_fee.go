package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/likecoin-chain/v4/x/likefeegrant"
)

func (k Keeper) DeductFeePerByte(ctx sdk.Context, msgSender sdk.AccAddress, bytesLength int, msg sdk.Msg) error {
	feePerByte := k.GetParams(ctx).FeePerByte
	amount := feePerByte.Amount.MulInt64(int64(bytesLength))
	fees := sdk.NewCoins(sdk.NewCoin(feePerByte.Denom, amount.Ceil().RoundInt()))
	return likefeegrant.DeductFeePerByte(
		ctx,
		k.accountKeeper, k.bankKeeper, k.feegrantKeeper,
		msgSender, fees, msg,
	)
}
