package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type checkTxHandler = func(context.IImmutableState, *types.Transaction) abci.ResponseCheckTx
type deliverTxHandler = func(context.IMutableState, *types.Transaction) abci.ResponseDeliverTx

var checkTxHandlerTable = make(map[reflect.Type]checkTxHandler)
var deliverTxHandlerTable = make(map[reflect.Type]deliverTxHandler)

func registerCheckTxHandler(t reflect.Type, f checkTxHandler) {
	checkTxHandlerTable[t] = f
}

func registerDeliverTxHandler(t reflect.Type, f deliverTxHandler) {
	deliverTxHandlerTable[t] = f
}

func CheckTx(state context.IImmutableState, tx *types.Transaction) abci.ResponseCheckTx {
	t := reflect.TypeOf(tx.GetTx())
	f, exist := checkTxHandlerTable[t]
	if !exist {
		return abci.ResponseCheckTx{} // TODO
	}
	return f(state, tx)
}

func DeliverTx(state context.IMutableState, tx *types.Transaction) abci.ResponseDeliverTx {
	t := reflect.TypeOf(tx.GetTx())
	f, exist := deliverTxHandlerTable[t]
	if !exist {
		return abci.ResponseDeliverTx{} // TODO
	}
	return f(state, tx)
}
