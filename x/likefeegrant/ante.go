package likefeegrant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type contextKey string

const FeeTxKey contextKey = "fee-tx-key"

var _ sdk.AnteDecorator = FeeTxContextDecorator{}

type FeeTxContextDecorator struct{}

func NewFeeTxContextDecorator() FeeTxContextDecorator {
	return FeeTxContextDecorator{}
}

func (FeeTxContextDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}
	newCtx := ctx.WithValue(FeeTxKey, feeTx)
	return next(newCtx, tx, simulate)
}

func GetFeeTx(ctx sdk.Context) sdk.FeeTx {
	value := ctx.Value(FeeTxKey)
	if value == nil {
		return nil
	}
	return value.(sdk.FeeTx)
}
