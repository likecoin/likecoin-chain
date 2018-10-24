package query

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/txstatus"
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

		Convey("If it is a query to a non-existing Tx", func() {
			res := Query(state, reqQuery)
			code := response.QueryTxNotExist.Code

			Convey(fmt.Sprintf("Should return code %d", code), func() {
				So(res.Code, ShouldEqual, code)
			})
		})

		txStatusList := []txstatus.TxStatus{
			txstatus.TxStatusSuccess,
			txstatus.TxStatusFail,
			txstatus.TxStatusPending,
		}
		for _, status := range txStatusList {
			s := status.String()
			Convey(fmt.Sprintf("If it is a valid query with %s Tx", s), func() {
				txstatus.SetStatus(state, txHash, status)
				res := Query(state, reqQuery)

				Convey(fmt.Sprintf("Should return code 0 and %s", s), func() {
					So(res.Code, ShouldEqual, 0)
					jsonRes := make(map[string]interface{})
					err := json.Unmarshal(res.Value, &jsonRes)
					So(err, ShouldBeNil)
					So(jsonRes, ShouldContainKey, "status")
					So(jsonRes["status"], ShouldEqual, s)
				})
			})
		}
	})
}
