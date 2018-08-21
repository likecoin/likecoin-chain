package query

import (
	"github.com/likecoin/likechain/abci/context"
	abci "github.com/tendermint/tendermint/abci/types"
)

type queryHandler = func(context.IImmutableState, abci.RequestQuery) abci.ResponseQuery

var queryHandlerTable = make(map[string]queryHandler)

func registerQueryHandler(path string, f queryHandler) {
	queryHandlerTable[path] = f
}

func Query(state context.IImmutableState, reqQuery abci.RequestQuery) abci.ResponseQuery {
	f, exist := queryHandlerTable[reqQuery.Path]
	if !exist {
		return abci.ResponseQuery{} // TODO
	}
	return f(state, reqQuery)
}
