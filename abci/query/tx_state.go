package query

import (
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/transaction"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryTxState(
	state context.IMutableState,
	reqQuery abci.RequestQuery,
) response.R {
	txHash := reqQuery.Data
	txStatus := transaction.GetStatus(state, txHash)

	data, err := (&types.TxStateResponse{
		Status: txStatus.String(),
	}).Marshal()
	if err != nil {
		log.WithError(err).Debug("Unable to parse tx state response to JSON")
		return response.QueryParsingResponseError
	}

	return response.Success.Merge(response.R{
		Data: data,
	})
}

func init() {
	registerQueryHandler("tx_state", queryTxState)
}
