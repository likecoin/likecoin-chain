package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkWithdraw(rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetWithdrawTx()
	_ = tx.From
	return abci.ResponseCheckTx{} // TODO
}

func deliverWithdraw(rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetWithdrawTx()
	_ = tx.From
	return abci.ResponseDeliverTx{} // TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_WithdrawTx)(nil))
	registerCheckTxHandler(t, checkWithdraw)
	registerDeliverTxHandler(t, deliverWithdraw)
}
