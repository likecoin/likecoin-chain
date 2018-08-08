package query

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryWithdrawProof(reqQuery abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{}
}

func init() {
	registerQueryHandler("withdraw_proof", queryWithdrawProof)
}
