package main_test

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
	customvalidator "github.com/likecoin/likechain/abci/cmd/api/validator"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
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
	statusCode int,
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
	return jsonRes, result.StatusCode
}

func TestAPI(t *testing.T) {
	Convey("Testing API", t, func() {
		mockCtx := context.NewMock()

		app := likechain.NewLikeChainApplication(mockCtx.ApplicationContext)
		node := rpctest.StartTendermint(app)

		client := rpcclient.NewLocal(node)

		router := gin.Default()

		customvalidator.Bind()

		routes.Initialize(router, client)

		appHeight := int64(0)

		// Test invalid path
		uri := "/404"
		res, code := request(router, "GET", uri, nil)
		So(code, ShouldEqual, http.StatusNotFound)

		//
		// Test POST /register
		//
		uri = "/v1/register"
		// Missing params
		res, code = request(router, "POST", uri, map[string]interface{}{})
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Invalid params
		res, code = request(router, "POST", uri, map[string]interface{}{
			"addr": "",
			"sig":  "",
		})
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Register A account
		mockCtx.SetInitialBalance(big.NewInt(100))
		sig := "0xb19ced763ac63a33476511ecce1df4ebd91bb9ae8b2c0d24b0a326d96c5717122ae0c9b5beacaf4560f3a2535a7673a3e567ff77f153e452907169d431c951091b"
		params := map[string]interface{}{
			"addr": fixture.Alice.Address.Hex(),
			"sig":  sig,
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "id")
		appHeight++
		aliceID := res["id"].(string)

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		// Duplicated registration
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusConflict)

		//
		// Test GET /account_info
		//
		uri = "/v1/account_info?identity=" + url.QueryEscape(aliceID)
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["id"], ShouldEqual, aliceID)
		So(res["balance"], ShouldEqual, "100")

		// Missing params
		uri = "/v1/account_info"
		res, code = request(router, "GET", uri, nil)
		So(code, ShouldEqual, http.StatusBadRequest)
		So(res["error"], ShouldNotBeNil)

		// Invalid params
		uri = "/v1/account_info?identity=" + types.NewZeroAddress().ToHex()
		res, code = request(router, "GET", uri, nil)
		So(code, ShouldEqual, http.StatusBadRequest)
		So(res["error"], ShouldNotBeNil)

		// Register B account
		mockCtx.SetInitialBalance(big.NewInt(200))
		uri = "/v1/register"
		sig = "0x6d8c7bb3292cab67f4814f9c2d1986430bd188b4eadf82a3fdf1e6be10f7599751985388c2a79429ee60761169e4c67e3b453daf88b637d77f87d7be68196b2c1b"
		res, code = request(router, "POST", uri, map[string]interface{}{
			"addr": fixture.Bob.Address.Hex(),
			"sig":  sig,
		})
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "id")
		appHeight += 2
		bobID := res["id"].(string)

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/account_info?identity=" + fixture.Bob.Address.Hex()
		res, _ = request(router, "GET", uri, nil)
		So(res["id"], ShouldEqual, bobID)
		So(res["balance"], ShouldEqual, "200")

		//
		// Test GET /block
		//
		uri = fmt.Sprintf("/v1/block?height=%d", appHeight)
		res, code = request(router, "GET", uri, nil)
		So(res["result"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusOK)

		// Missing params
		uri = "/v1/block"
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Invalid height
		uri = "/v1/block?height=-1"
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		uri = "/v1/block?height=999"
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		//
		// Test POST /transfer
		//
		uri = "/v1/transfer"

		// Missing params
		res, code = request(router, "POST", uri, map[string]interface{}{})
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		sig = "0x343db6effdf722054ff57bcdad4d21b7025407557c631e7e4b1cd77411fa4c155e91c2a8fff992792d885c34411dca84d8d2eaa863ff7eecda1f3322be024d071b"
		params = map[string]interface{}{
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
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		// Duplicated transfer
		res, code = request(router, "POST", uri, params)
		So(code, ShouldEqual, http.StatusConflict)
		So(res["error"], ShouldNotBeNil)

		// Invalid logic
		params["nonce"] = 2
		params["to"] = []map[string]interface{}{
			{
				"identity": fixture.Bob.Address.Hex(),
				"value":    "999",
			},
		}
		params["sig"] = "0x1ef7d4812b55645c4aea8ac1a16354278827f39722dd3c4f38f23dda004795db121d8d4341613cdffed28f29418cf7df46e55f21380daab5711d3ffb8f8ee0771c"
		res, code = request(router, "POST", uri, params)
		So(res["code"], ShouldEqual, response.TransferCheckTxNotEnoughBalance.Code)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Invalid params
		params["sig"] = types.NewZeroSignature().ToHex()
		res, code = request(router, "POST", uri, params)
		So(code, ShouldEqual, http.StatusBadRequest)
		So(res["error"], ShouldNotBeNil)

		uri = "/v1/account_info?identity=" + url.QueryEscape(aliceID)
		res, _ = request(router, "GET", uri, nil)
		So(res["balance"], ShouldEqual, "99")

		uri = "/v1/account_info?identity=" + url.QueryEscape(bobID)
		res, _ = request(router, "GET", uri, nil)
		So(res["balance"], ShouldEqual, "201")

		//
		// Test GET /tx_state
		//
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
		res, _ = request(router, "GET", uri, nil)
		So(res["status"], ShouldEqual, "success")

		// Missing params
		uri = "/v1/tx_state"
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Invalid params
		uri = "/v1/tx_state?tx_hash=123ABC"
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		//
		// Test POST /withdraw
		//
		uri = "/v1/withdraw"
		sig = "0x9d6dca90161dcdcf5594e2070b221a6c50318e4034bbf5b25ba50402dcbe0ebb2a8fe928a28fffac5c0edb3d607b86da7df016f5ce789e7488f5fd70da37dbe61b"
		params = map[string]interface{}{
			"identity": fixture.Alice.Address.Hex(),
			"nonce":    2,
			"to_addr":  types.NewZeroAddress().ToHex(),
			"value":    "1",
			"fee":      "0",
			"sig":      sig,
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		// Duplicated transfer
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusConflict)

		// Missing params
		res, code = request(router, "POST", uri, map[string]interface{}{})
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Invalid params
		params["fee"] = "-1"
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Invalid logic
		params["fee"] = "0"
		params["nonce"] = 3
		params["value"] = "999"
		params["sig"] = "0x4389b731b1d67c792cc309a52281d0ef5350973f1d7d72f38049915756fc8ca86765e505a8c4c54e5cc98473881634706ad79ae561ca959fb49da5f7193c14901b"
		res, code = request(router, "POST", uri, params)
		So(res["code"], ShouldEqual, response.WithdrawCheckTxNotEnoughBalance.Code)
		So(code, ShouldEqual, http.StatusBadRequest)

		//
		// Test GET /withdraw_proof
		//
		formattedQuery := "/v1/withdraw_proof?identity=%s&to_addr=%s&height=%d&nonce=%d&value=%s&fee=%s"
		uri = fmt.Sprintf(
			formattedQuery,
			url.QueryEscape(aliceID),
			types.NewZeroAddress().ToHex(),
			appHeight,
			2,
			"1",
			"0",
		)
		res, code = request(router, "GET", uri, nil)
		So(res["proof"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusOK)

		// Using address
		uri = fmt.Sprintf(
			formattedQuery,
			fixture.Alice.Address.Hex(),
			types.NewZeroAddress().ToHex(),
			appHeight,
			2,
			"1",
			"0",
		)
		res, code = request(router, "GET", uri, nil)
		So(res["proof"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusOK)

		// Missing params
		uri = "/v1/withdraw_proof"
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Invalid params
		uri = fmt.Sprintf(
			formattedQuery,
			url.QueryEscape(aliceID),
			types.NewZeroAddress().ToHex(),
			-1,
			2,
			"1",
			"0",
		)
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)
	})
}
