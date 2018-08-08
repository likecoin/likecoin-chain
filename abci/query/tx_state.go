package query

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryTxState(reqQuery abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{}
}

func init() {
	registerQueryHandler("tx_state", queryTxState)
}
