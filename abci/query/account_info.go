package query

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryAccountInfo(reqQuery abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{}
}

func init() {
	registerQueryHandler("account_info", queryAccountInfo)
}
