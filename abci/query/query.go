package query

import (
	"encoding/json"

	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/response"
	abci "github.com/tendermint/tendermint/abci/types"
)

var log = logger.L

type queryHandler = func(context.IMutableState, abci.RequestQuery) response.R

var queryHandlerTable = make(map[string]queryHandler)

func registerQueryHandler(path string, f queryHandler) {
	queryHandlerTable[path] = f
}

func Query(
	state context.IMutableState,
	reqQuery abci.RequestQuery,
) abci.ResponseQuery {
	f, exist := queryHandlerTable[reqQuery.Path]
	if !exist {
		return response.QueryPathNotExist.ToResponseQuery()
	}

	if reqQuery.Data == nil {
		log.Info(response.QueryParsingRequestError.Info)
		return response.QueryParsingRequestError.ToResponseQuery()
	}

	return f(state, reqQuery).ToResponseQuery()
}

type jsonMap map[string]interface{}

func (m jsonMap) ToResponse() response.R {
	b, err := json.Marshal(m)
	if err != nil {
		log.WithError(err).Info(response.QueryParsingResponseError.Info)
		return response.QueryParsingResponseError
	}
	return response.Success.Merge(response.R{Data: b})
}
