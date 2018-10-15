package query

import (
	"fmt"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/transaction"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryTxState(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	txHash := utils.HashRawTx([]byte(""))

	Convey("Given a tx state query", t, func() {
		appCtx.Reset()
		reqQuery := abci.RequestQuery{
			Data: txHash,
			Path: "tx_state",
		}

		txStatusList := []types.TxStatus{
			types.TxStatusSuccess,
			types.TxStatusFail,
			types.TxStatusPending,
		}
		for _, status := range txStatusList {
			s := status.String()
			Convey(fmt.Sprintf("If it is a valid query with %s Tx", s), func() {
				transaction.SetStatus(state, txHash, status)
				res := Query(state, reqQuery)

				Convey(fmt.Sprintf("Should return code 0 and %s", s), func() {
					So(res.Code, ShouldEqual, 0)

					txStateRes, err := new(types.TxStateResponse).Unmarshal(res.Value)
					if err != nil {
						t.Error(err)
					}

					So(txStateRes.Status, ShouldEqual, s)
				})
			})
		}
	})
}
