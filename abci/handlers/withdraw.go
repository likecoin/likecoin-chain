package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/errcode"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkWithdraw(state context.IImmutableState, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetWithdrawTx()
	if tx == nil {
		// TODO: log
		panic("Expect WithdrawTx but got nil")
	}

	_ = tx.From

	if !validateWithdrawTransactionFormat(tx) {
		code, info := errcode.WithdrawCheckTxInvalidFormat()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	if !validateWithdrawSignature(tx.Sig) {
		code, info := errcode.WithdrawCheckTxInvalidSignature()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	return abci.ResponseCheckTx{} // TODO
}

func deliverWithdraw(state context.IMutableState, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetWithdrawTx()
	if tx == nil {
		// TODO: log
		panic("Expect WithdrawTx but got nil")
	}

	if !validateWithdrawTransactionFormat(tx) {
		code, info := errcode.WithdrawDeliverTxInvalidFormat()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	if !validateWithdrawSignature(tx.Sig) {
		code, info := errcode.RegisterCheckTxInvalidFormat()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	return abci.ResponseDeliverTx{} // TODO
}

func validateWithdrawSignature(sig *types.Signature) bool {
	return false // TODO
}

func validateWithdrawTransactionFormat(tx *types.WithdrawTransaction) bool {
	return false // TODO
}

func withdraw(state context.IMutableState, tx *types.WithdrawTransaction) {
	// TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_WithdrawTx)(nil))
	registerCheckTxHandler(t, checkWithdraw)
	registerDeliverTxHandler(t, deliverWithdraw)
}
