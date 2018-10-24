package query

import (
	"encoding/json"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/txstatus"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryTxState(
	state context.IMutableState,
	reqQuery abci.RequestQuery,
) response.R {
	txHash := reqQuery.Data
	txStatus := txstatus.GetStatus(state, txHash)

	if txStatus == txstatus.TxStatusNotSet {
		return response.QueryTxNotExist
	}

	return jsonMap{
		"status": txStatus.String(),
	}.ToResponse()
}

func init() {
	registerQueryHandler("tx_state", queryTxState)
}

// TxStateRes represents response data of tx_state query
type TxStateRes struct {
	Status string `json:"status"`
}

// GetTxStateRes transforms the raw byte response from tx_state query back to GetTxStateRes structure
func GetTxStateRes(data []byte) *TxStateRes {
	result := TxStateRes{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}
	return &result
}
