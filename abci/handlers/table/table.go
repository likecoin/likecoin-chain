package table

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

type checkTxHandler = func(
	context.IImmutableState,
	*types.Transaction,
) response.R

type deliverTxHandler = func(
	context.IMutableState,
	*types.Transaction,
	[]byte,
) response.R

var checkTxHandlerTable = make(map[reflect.Type]checkTxHandler)
var deliverTxHandlerTable = make(map[reflect.Type]deliverTxHandler)

func getType(t interface{}) reflect.Type {
	return reflect.TypeOf(t)
}

func getTypeFromTx(tx *types.Transaction) reflect.Type {
	return reflect.TypeOf(tx.GetTx())
}

// RegisterCheckTxHandler registers a CheckTx handler for a type
func RegisterCheckTxHandler(t reflect.Type, f checkTxHandler) {
	checkTxHandlerTable[t] = f
}

// GetCheckTxHandlerFromTx retrieves a CheckTx handler for a tx
func GetCheckTxHandlerFromTx(tx *types.Transaction) (
	reflect.Type,
	checkTxHandler,
	bool,
) {
	_type := getTypeFromTx(tx)
	handler, ok := checkTxHandlerTable[_type]
	return _type, handler, ok
}

// RegisterDeliverTxHandler registers a CheckTx handler for a type
func RegisterDeliverTxHandler(t reflect.Type, f deliverTxHandler) {
	deliverTxHandlerTable[t] = f
}

// GetDeliverTxHandlerFromTx retrieves a Deliver handler for a tx
func GetDeliverTxHandlerFromTx(tx *types.Transaction) (
	reflect.Type,
	deliverTxHandler,
	bool,
) {
	_type := getTypeFromTx(tx)
	handler, ok := deliverTxHandlerTable[_type]
	return _type, handler, ok
}
