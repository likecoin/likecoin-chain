package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkDeposit(context *context.Context, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetDepositTx()
	_ = tx.BlockNumber
	return abci.ResponseCheckTx{} // TODO
}

func deliverDeposit(context *context.Context, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetDepositTx()
	_ = tx.BlockNumber
	return abci.ResponseDeliverTx{} // TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_DepositTx)(nil))
	registerCheckTxHandler(t, checkDeposit)
	registerDeliverTxHandler(t, deliverDeposit)
}
