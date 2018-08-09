package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkTransfer(context *context.Context, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetTransferTx()
	_ = tx.From
	return abci.ResponseCheckTx{} // TODO
}

func deliverTransfer(context *context.Context, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetTransferTx()
	_ = tx.From
	return abci.ResponseDeliverTx{} // TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_TransferTx)(nil))
	registerCheckTxHandler(t, checkTransfer)
	registerDeliverTxHandler(t, deliverTransfer)
}
