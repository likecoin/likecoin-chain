package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	likechain "github.com/likecoin/likechain/abci/app"
	"github.com/likecoin/likechain/abci/cmd/api/routes"
	"github.com/likecoin/likechain/abci/context"
	. "github.com/smartystreets/goconvey/convey"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpctest "github.com/tendermint/tendermint/rpc/test"
)

func request(
	router *gin.Engine,
	method, uri string,
	params map[string]interface{},
) (
	jsonRes map[string]interface{},
) {
	var req *http.Request

	switch method {
	case "GET":
		req = httptest.NewRequest(method, uri, nil)
	case "POST":
		fallthrough
	default:
		data, _ := json.Marshal(params)
		req = httptest.NewRequest(method, uri, bytes.NewReader(data))
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	resJSONBytes, _ := ioutil.ReadAll(result.Body)

	json.Unmarshal(resJSONBytes, &jsonRes)
	return jsonRes
}

func TestMain(t *testing.T) {
	Convey("Testing API", t, func() {
		mockCtx := context.NewMock()
		mockCtx.GetMutableState().SetInitialBalance(big.NewInt(100))

		app := likechain.NewLikeChainApplication(mockCtx.ApplicationContext)
		node := rpctest.StartTendermint(app)
		defer func() {
			node.Stop()
			node.Wait()
		}()

		client := rpcclient.NewLocal(node)

		router := gin.Default()
		routes.Initialize(router, client)

		uri := "/v1/register"
		res := request(router, "POST", uri, map[string]interface{}{
			"addr": "0x064b663abf9d74277a07aa7563a8a64a54de8c0a",
			"sig":  "0xb19ced763ac63a33476511ecce1df4ebd91bb9ae8b2c0d24b0a326d96c5717122ae0c9b5beacaf4560f3a2535a7673a3e567ff77f153e452907169d431c951091b",
		})
		So(res, ShouldContainKey, "id")

		likeChainID := res["id"]
		uri = "/v1/account_info?identity=" + url.QueryEscape(likeChainID.(string))
		res = request(router, "GET", uri, nil)
		So(res["id"], ShouldEqual, likeChainID)
		So(res["balance"], ShouldEqual, "100")

		status, _ := client.Status()
		uri = fmt.Sprintf("/v1/block?height=%d", status.SyncInfo.LatestBlockHeight)
		res = request(router, "GET", uri, nil)
		So(res["result"], ShouldNotBeNil)
	})
}
