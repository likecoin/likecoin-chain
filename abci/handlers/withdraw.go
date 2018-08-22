package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

func checkWithdraw(state context.IImmutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetWithdrawTx()
	if tx == nil {
		// TODO: log
		panic("Expect WithdrawTx but got nil")
	}

	_ = tx.From

	if !validateWithdrawTransactionFormat(tx) {
		return response.WithdrawCheckTxInvalidFormat
	}

	if !validateWithdrawSignature(tx.Sig) {
		return response.WithdrawCheckTxInvalidSignature
	}

	return response.Success // TODO
}

func deliverWithdraw(state context.IMutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetWithdrawTx()
	if tx == nil {
		// TODO: log
		panic("Expect WithdrawTx but got nil")
	}

	if !validateWithdrawTransactionFormat(tx) {
		return response.WithdrawDeliverTxInvalidFormat
	}

	if !validateWithdrawSignature(tx.Sig) {
		return response.RegisterCheckTxInvalidFormat
	}

	return response.Success // TODO
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
