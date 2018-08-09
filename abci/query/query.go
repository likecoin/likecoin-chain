package query

import (
	"github.com/likecoin/likechain/abci/context"
	abci "github.com/tendermint/tendermint/abci/types"
)

type queryHandler = func(context *context.Context, reqQuery abci.RequestQuery) abci.ResponseQuery

var queryHandlerTable = make(map[string]queryHandler)

func registerQueryHandler(path string, f queryHandler) {
	queryHandlerTable[path] = f
}

func Query(context *context.Context, reqQuery abci.RequestQuery) abci.ResponseQuery {
	f, exist := queryHandlerTable[reqQuery.Path]
	if !exist {
		return abci.ResponseQuery{} // TODO
	}
	return f(context, reqQuery)
}
