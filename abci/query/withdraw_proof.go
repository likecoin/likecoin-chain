package query

import (
	"github.com/likecoin/likechain/abci/context"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryWithdrawProof(state context.IImmutableState, reqQuery abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{}
}

func init() {
	registerQueryHandler("withdraw_proof", queryWithdrawProof)
}
