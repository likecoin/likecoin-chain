package main

import (
	"bytes"
	"encoding/base64"
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
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
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

func TestAPI(t *testing.T) {
	Convey("Testing API", t, func() {
		mockCtx := context.NewMock()

		app := likechain.NewLikeChainApplication(mockCtx.ApplicationContext)
		node := rpctest.StartTendermint(app)

		client := rpcclient.NewLocal(node)

		router := gin.Default()
		routes.Initialize(router, client)

		appHeight := int64(0)

		// Test POST /register
		// Register A account
		uri := "/v1/register"
		mockCtx.GetMutableState().SetInitialBalance(big.NewInt(100))
		sig := "0xb19ced763ac63a33476511ecce1df4ebd91bb9ae8b2c0d24b0a326d96c5717122ae0c9b5beacaf4560f3a2535a7673a3e567ff77f153e452907169d431c951091b"
		res := request(router, "POST", uri, map[string]interface{}{
			"addr": fixture.Alice.Address.Hex(),
			"sig":  sig,
		})
		So(res["error"], ShouldBeNil)
		So(res, ShouldContainKey, "id")
		appHeight++
		aliceID := res["id"].(string)

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		// Test GET /account_info
		uri = "/v1/account_info?identity=" + url.QueryEscape(aliceID)
		res = request(router, "GET", uri, nil)
		So(res["id"], ShouldEqual, aliceID)
		So(res["balance"], ShouldEqual, "100")

		// Register B account
		mockCtx.GetMutableState().SetInitialBalance(big.NewInt(200))
		uri = "/v1/register"
		sig = "0x6d8c7bb3292cab67f4814f9c2d1986430bd188b4eadf82a3fdf1e6be10f7599751985388c2a79429ee60761169e4c67e3b453daf88b637d77f87d7be68196b2c1b"
		res = request(router, "POST", uri, map[string]interface{}{
			"addr": fixture.Bob.Address.Hex(),
			"sig":  sig,
		})
		So(res["error"], ShouldBeNil)
		So(res, ShouldContainKey, "id")
		appHeight += 2
		bobID := res["id"].(string)

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/account_info?identity=" + url.QueryEscape(bobID)
		res = request(router, "GET", uri, nil)
		So(res["id"], ShouldEqual, bobID)
		So(res["balance"], ShouldEqual, "200")

		// Test GET /block
		uri = fmt.Sprintf("/v1/block?height=%d", appHeight)
		res = request(router, "GET", uri, nil)
		So(res["result"], ShouldNotBeNil)

		// Test POST /transfer
		uri = "/v1/transfer"
		sig = "0x343db6effdf722054ff57bcdad4d21b7025407557c631e7e4b1cd77411fa4c155e91c2a8fff992792d885c34411dca84d8d2eaa863ff7eecda1f3322be024d071b"
		res = request(router, "POST", uri, map[string]interface{}{
			"fee":      "0",
			"identity": fixture.Alice.Address.Hex(),
			"nonce":    1,
			"to": []map[string]interface{}{
				{
					"identity": fixture.Bob.Address.Hex(),
					"value":    "1",
				},
			},
			"sig": sig,
		})
		So(res["error"], ShouldBeNil)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/account_info?identity=" + url.QueryEscape(aliceID)
		res = request(router, "GET", uri, nil)
		So(res["balance"], ShouldEqual, "99")

		uri = "/v1/account_info?identity=" + url.QueryEscape(bobID)
		res = request(router, "GET", uri, nil)
		So(res["balance"], ShouldEqual, "201")

		// Test GET /tx_state
		rawTx, _ := (&types.TransferTransaction{
			Fee:   types.NewBigInteger("0"),
			From:  fixture.Alice.RawAddress.ToIdentifier(),
			Nonce: 1,
			Sig:   types.NewSignatureFromHex(sig),
			ToList: []*types.TransferTransaction_TransferTarget{
				&types.TransferTransaction_TransferTarget{
					To:    fixture.Bob.RawAddress.ToIdentifier(),
					Value: types.NewBigInteger("1"),
				},
			},
		}).ToTransaction().Encode()
		txHash := utils.HashRawTx(rawTx)
		uri = "/v1/tx_state?tx_hash=" + url.QueryEscape(base64.StdEncoding.EncodeToString(txHash))
		res = request(router, "GET", uri, nil)
		So(res["status"], ShouldEqual, "success")

		// Test POST /withdraw
		uri = "/v1/withdraw"
		sig = "0x9d6dca90161dcdcf5594e2070b221a6c50318e4034bbf5b25ba50402dcbe0ebb2a8fe928a28fffac5c0edb3d607b86da7df016f5ce789e7488f5fd70da37dbe61b"
		res = request(router, "POST", uri, map[string]interface{}{
			"identity": fixture.Alice.Address.Hex(),
			"nonce":    2,
			"to_addr":  types.NewZeroAddress().ToHex(),
			"value":    "1",
			"fee":      "0",
			"sig":      sig,
		})
		So(res["error"], ShouldBeNil)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		// Test GET /withdraw_proof
		uri = fmt.Sprintf(
			"/v1/withdraw_proof?identity=%s&to_addr=%s&height=%d&nonce=%d&value=%s&fee=%s",
			url.QueryEscape(aliceID),
			types.NewZeroAddress().ToHex(),
			appHeight,
			2,
			"1",
			"0",
		)
		res = request(router, "GET", uri, nil)
		So(res["proof"], ShouldNotBeNil)
	})
}
