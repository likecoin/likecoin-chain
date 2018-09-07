package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/response"

	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var log = logger.L

type checkTxHandler = func(context.IImmutableState, *types.Transaction) response.R
type deliverTxHandler = func(context.IMutableState, *types.Transaction, []byte) response.R

var checkTxHandlerTable = make(map[reflect.Type]checkTxHandler)
var deliverTxHandlerTable = make(map[reflect.Type]deliverTxHandler)

func RegisterCheckTxHandler(t reflect.Type, f checkTxHandler) {
	checkTxHandlerTable[t] = f
}

func RegisterDeliverTxHandler(t reflect.Type, f deliverTxHandler) {
	deliverTxHandlerTable[t] = f
}

func CheckTx(state context.IImmutableState, tx *types.Transaction) abci.ResponseCheckTx {
	t := reflect.TypeOf(tx.GetTx())
	f, exist := checkTxHandlerTable[t]
	if !exist {
		return abci.ResponseCheckTx{} // TODO
	}
	return f(state, tx).ToResponseCheckTx()
}

func DeliverTx(state context.IMutableState, tx *types.Transaction, txHash []byte) abci.ResponseDeliverTx {
	t := reflect.TypeOf(tx.GetTx())
	f, exist := deliverTxHandlerTable[t]
	if !exist {
		return abci.ResponseDeliverTx{} // TODO
	}
	return f(state, tx, txHash).ToResponseDeliverTx()
}
