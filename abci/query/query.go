package query

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

type queryHandler = func(reqQuery abci.RequestQuery) abci.ResponseQuery

var queryHandlerTable = make(map[string]queryHandler)

func registerQueryHandler(path string, f queryHandler) {
	queryHandlerTable[path] = f
}

func Query(reqQuery abci.RequestQuery) abci.ResponseQuery {
	f, exist := queryHandlerTable[reqQuery.Path]
	if !exist {
		return abci.ResponseQuery{} // TODO
	}
	return f(reqQuery)
}
