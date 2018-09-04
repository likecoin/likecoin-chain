package query

import (
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryWithdrawProof(
	state context.IImmutableState,
	reqQuery abci.RequestQuery,
) response.R {
	return response.R{}
}

func init() {
	registerQueryHandler("withdraw_proof", queryWithdrawProof)
}
