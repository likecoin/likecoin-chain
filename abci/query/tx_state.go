package query

import (
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/handlers/transfer"
	"github.com/likecoin/likechain/abci/response"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryTxState(
	state context.IMutableState,
	reqQuery abci.RequestQuery,
) response.R {
	txHash := reqQuery.Data
	txStatus := transfer.GetStatus(state, txHash)

	return response.Success.Merge(response.R{
		Data: []byte(txStatus.String()),
	})
}

func init() {
	registerQueryHandler("tx_state", queryTxState)
}
