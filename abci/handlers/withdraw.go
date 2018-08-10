package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkWithdraw(ctx context.Context, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetWithdrawTx()
	_ = tx.From
	return abci.ResponseCheckTx{} // TODO
}

func deliverWithdraw(ctx context.Context, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetWithdrawTx()

	if !validateWithdrawSignature(tx.Sig) {
		panic("Invalid signature")
	}

	if !validateWithdrawTransaction(tx) {
		panic("Invalid WithdrawTransaction in WithdrawTx")
	}

	return abci.ResponseDeliverTx{} // TODO
}

func validateWithdrawSignature(sig *types.Signature) bool {
	return false // TODO
}

func validateWithdrawTransaction(tx *types.WithdrawTransaction) bool {
	return false // TODO
}

func withdraw(ctx context.Context, tx *types.WithdrawTransaction) {
	// TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_WithdrawTx)(nil))
	registerCheckTxHandler(t, checkWithdraw)
	registerDeliverTxHandler(t, deliverWithdraw)
}
