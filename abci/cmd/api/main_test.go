package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/account"
	likechain "github.com/likecoin/likechain/abci/app"
	"github.com/likecoin/likechain/abci/cmd/api/routes"
	customvalidator "github.com/likecoin/likechain/abci/cmd/api/validator"
	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/contract"
	"github.com/likecoin/likechain/abci/state/deposit"
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

		state := mockCtx.GetMutableState()
		account.NewAccountFromID(state, Carol.ID, Carol.Address)
		account.NewAccountFromID(state, Dave.ID, Dave.Address)

		approvers := []deposit.Approver{
			{ID: Carol.ID, Weight: 33},
			{ID: Dave.ID, Weight: 67},
		}
		deposit.SetDepositApprovers(state, approvers)
		contractUpdaters := []contract.Updater{
			{ID: Carol.ID, Weight: 33},
			{ID: Dave.ID, Weight: 67},
		}
		contract.SetContractUpdaters(state, contractUpdaters)

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
			"sig": map[string]interface{}{
				"value": "",
			},
		})
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Register A account
		sig := "0xcf3a79ff76b94dd6bee6bfbbd2da201a9972b28e1ef47f4d7c66034d1aa74bf016d22572fdad933ebd867eec110234f88490f105ec5ad0af39ebc5db787b08011b"
		params := map[string]interface{}{
			"addr": Alice.Address.String(),
			"sig": map[string]interface{}{
				"type":  "eip712",
				"value": sig,
			},
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "id")
		So(res, ShouldContainKey, "tx_hash")
		appHeight++
		aliceID := res["id"].(string)
		txHashHex := res["tx_hash"].(string)

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		// Duplicated registration
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusConflict)

		//
		// Test GET /tx_state
		//
		uri = "/v1/tx_state?tx_hash=" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "success")

		uri = "/v1/tx_state?tx_hash=0x" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
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

		uri = "/v1/tx_state?tx_hash=0xABC"
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		//
		// Test GET /account_info
		//
		uri = "/v1/account_info?identity=" + url.QueryEscape(aliceID)
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["id"], ShouldEqual, aliceID)
		So(res["balance"], ShouldEqual, "0")

		// Missing params
		uri = "/v1/account_info"
		res, code = request(router, "GET", uri, nil)
		So(code, ShouldEqual, http.StatusBadRequest)
		So(res["error"], ShouldNotBeNil)

		// Invalid params
		uri = "/v1/account_info?identity=0x0000000000000000000000000000000000000000"
		res, code = request(router, "GET", uri, nil)
		So(code, ShouldEqual, http.StatusBadRequest)
		So(res["error"], ShouldNotBeNil)

		//
		// Test GET /address_info
		//
		uri = "/v1/address_info?addr=" + Alice.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["id"], ShouldEqual, aliceID)
		So(res["balance"], ShouldEqual, "0")

		// Missing params
		uri = "/v1/address_info"
		res, code = request(router, "GET", uri, nil)
		So(code, ShouldEqual, http.StatusBadRequest)
		So(res["error"], ShouldNotBeNil)

		// Invalid params
		uri = "/v1/account_info?addr=0x012345678901234567890123456789012345678" // missing one hex digit
		res, code = request(router, "GET", uri, nil)
		So(code, ShouldEqual, http.StatusBadRequest)
		So(res["error"], ShouldNotBeNil)

		// Deposit into A and B account addresses proposed by C
		uri = "/v1/deposit"
		sig = "0x254c15e8d2baf6ac11cc3d549cc94f7445c839a2a2f75ca724ebcb6dc9a498205da5773c55ff3c5cbdc3830c296ace91979deeb264407b96e52c5dc67fb4ee3a1b"
		res, code = request(router, "POST", uri, map[string]interface{}{
			"block_number": 1,
			"identity":     Carol.Address.String(),
			"inputs": []map[string]interface{}{
				{"from_addr": Alice.Address.String(), "value": "100"},
				{"from_addr": Bob.Address.String(), "value": "200"},
			},
			"nonce": 1,
			"sig": map[string]interface{}{
				"value": sig,
			},
		})
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/tx_state?tx_hash=" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "success")

		uri = "/v1/account_info?identity=" + Alice.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "0")

		uri = "/v1/address_info?addr=" + Bob.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "0")

		// Deposit into A and B account addresses proposed by D
		uri = "/v1/deposit"
		sig = "0x180fd088d9c26b7b6b62272b797cd52a572870cb2a4ffc3a7cc6f04c1b7fde3a333cbbaf113fe672264c7729dd939bfe658992e02b62445e08789b49fc4b45911c"
		res, code = request(router, "POST", uri, map[string]interface{}{
			"block_number": 1,
			"identity":     Dave.ID.String(),
			"inputs": []map[string]interface{}{
				{"from_addr": Alice.Address.String(), "value": "100"},
				{"from_addr": Bob.Address.String(), "value": "200"},
			},
			"nonce": 1,
			"sig": map[string]interface{}{
				"value": sig,
			},
		})
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/tx_state?tx_hash=" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "success")

		uri = "/v1/account_info?identity=" + Alice.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "100")

		uri = "/v1/address_info?addr=" + Bob.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "200")

		// Register B account
		uri = "/v1/register"
		sig = "0x6d8c7bb3292cab67f4814f9c2d1986430bd188b4eadf82a3fdf1e6be10f7599751985388c2a79429ee60761169e4c67e3b453daf88b637d77f87d7be68196b2c1b"
		res, code = request(router, "POST", uri, map[string]interface{}{
			"addr": Bob.Address.String(),
			"sig": map[string]interface{}{
				"value": sig,
			},
		})
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "id")
		appHeight += 2
		bobID := res["id"].(string)

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/account_info?identity=" + Bob.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["id"], ShouldEqual, bobID)
		So(res["balance"], ShouldEqual, "200")

		//
		// Test GET /block
		//
		uri = fmt.Sprintf("/v1/block?height=%d", appHeight)
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
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

		sig = "0x3cd8332511becc97ddcca750adf591a434a309331f9db77f69072dc440fa20b62496a816e6850bd1c1e5c1d17756c4f86b2e6b44d82cad813e95b6a4004798371b"
		params = map[string]interface{}{
			"fee":      "0",
			"identity": Alice.Address.String(),
			"nonce":    1,
			"outputs": []map[string]interface{}{
				{
					"identity": Bob.Address.String(),
					"value":    "1",
				},
			},
			"sig": map[string]interface{}{
				"value": sig,
			},
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		// Duplicated transfer
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusConflict)

		// Invalid logic
		params["nonce"] = 2
		params["outputs"] = []map[string]interface{}{
			{
				"identity": Bob.Address.String(),
				"value":    "999",
			},
		}
		params["sig"] = map[string]interface{}{
			"value": "0xbf61280a9930be07f0782ac4df7660a1f67d4fa3c681f9a43db36215c787cb3e12cfe342dfba1ec8a67eca7874b6839176f51c65bf238f24fc181a192fa906511c",
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldNotBeNil)
		So(res["code"], ShouldEqual, response.TransferNotEnoughBalance.Code)
		So(code, ShouldEqual, http.StatusBadRequest)

		// Invalid params
		params["sig"] = map[string]interface{}{
			"value": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		uri = "/v1/account_info?identity=" + url.QueryEscape(aliceID)
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "99")

		uri = "/v1/account_info?identity=" + url.QueryEscape(bobID)
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "201")

		//
		// Test GET /tx_state
		//
		uri = "/v1/tx_state?tx_hash=" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
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
		sig = "0x221546d3afaa5875f153a726979a90e76d8c1155abd4ed50fc888f7072c509515ada487f12f5f59a640c19735e95135b2165b1a12566171371f7b7045f7c84071c"
		params = map[string]interface{}{
			"identity": Alice.Address.String(),
			"nonce":    2,
			"to_addr":  "0x0000000000000000000000000000000000000000",
			"value":    "1",
			"fee":      "0",
			"sig": map[string]interface{}{
				"type":  "eip712",
				"value": sig,
			},
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		status, err := client.Status()
		if err != nil {
			t.Error(err)
		}
		withdrawHeight := status.SyncInfo.LatestBlockHeight

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/account_info?identity=" + Alice.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "98")

		//
		// Test GET /tx_state
		//
		uri = "/v1/tx_state?tx_hash=" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "success")

		// Duplicated transfer
		uri = "/v1/withdraw"
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
		params["sig"] = map[string]interface{}{
			"value": "0x4389b731b1d67c792cc309a52281d0ef5350973f1d7d72f38049915756fc8ca86765e505a8c4c54e5cc98473881634706ad79ae561ca959fb49da5f7193c14901b",
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldNotBeNil)
		So(res["code"], ShouldEqual, response.WithdrawNotEnoughBalance.Code)
		So(code, ShouldEqual, http.StatusBadRequest)

		//
		// Test GET /withdraw_proof
		//
		formattedQuery := "/v1/withdraw_proof?identity=%s&to_addr=%s&height=%d&nonce=%d&value=%s&fee=%s"
		uri = fmt.Sprintf(
			formattedQuery,
			url.QueryEscape(aliceID),
			"0x0000000000000000000000000000000000000000",
			withdrawHeight,
			2,
			"1",
			"0",
		)
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["proof"], ShouldNotBeNil)
		proof := res["proof"]
		So(code, ShouldEqual, http.StatusOK)

		// Using address
		uri = fmt.Sprintf(
			formattedQuery,
			Alice.Address.String(),
			"0x0000000000000000000000000000000000000000",
			withdrawHeight,
			2,
			"1",
			"0",
		)
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["proof"], ShouldNotBeNil)
		So(res["proof"], ShouldResemble, proof)
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
			"0x0000000000000000000000000000000000000000",
			-1,
			2,
			"1",
			"0",
		)
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusBadRequest)

		//
		// Test POST /contract_update
		//
		uri = "/v1/contract_update"
		sig = "0xfa14ecac48c39fa32be02bc31ea25f8c2841f6e071763fcef4edb40744b9a3a13b2acdcc497726eb8a55b8cd6eb765dff557e3505ffc53ae8004241ccb84f2ae1c"
		params = map[string]interface{}{
			"identity":       Carol.Address.String(),
			"nonce":          2,
			"contract_addr":  "0x0102030405060708091011121314151617181920",
			"contract_index": 1,
			"sig": map[string]interface{}{
				"value": sig,
			},
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		uri = "/v1/contract_update"
		sig = "0x8311a84142712d3eb1e81294400169097926a7b9e685df776c53db10b0846b3440e713aac7533d72b4d25d0afdf95e75733b8ac3b556ef7acb546cae215781761b"
		params = map[string]interface{}{
			"identity":       Dave.ID.String(),
			"nonce":          2,
			"contract_addr":  "0x0102030405060708091011121314151617181920",
			"contract_index": 1,
			"sig": map[string]interface{}{
				"value": sig,
			},
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		status, err = client.Status()
		if err != nil {
			t.Error(err)
		}
		contractUpdateHeight := status.SyncInfo.LatestBlockHeight

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		//
		// Test GET /contract_update_proof
		//
		formattedQuery = "/v1/contract_update_proof?height=%d&contract_index=1"
		uri = fmt.Sprintf(formattedQuery, contractUpdateHeight)
		res, code = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["proof"], ShouldNotBeNil)
		So(code, ShouldEqual, http.StatusOK)

		//
		// Test POST /simple_transfer
		//
		uri = "/v1/simple_transfer"
		sig = "0x623f31cc53432dd5ba38e8dde93edf71c6fc3467e52c773db40a03a461a9ac316f80dbc38717e4645efcacbf45a6b4c704153765e62140b4a46174989c40fa2d1c"
		params = map[string]interface{}{
			"identity": Alice.Address.String(),
			"to":       Bob.Address.String(),
			"value":    "1",
			"remark":   "there is no spoon",
			"fee":      "1",
			"nonce":    3,
			"sig": map[string]interface{}{
				"value": sig,
			},
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/tx_state?tx_hash=" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "success")

		uri = "/v1/account_info?identity=" + Alice.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "97")

		uri = "/v1/account_info?identity=" + Bob.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "202")

		//
		// Test POST /simple_transfer with no remark
		//
		uri = "/v1/simple_transfer"
		sig = "0x12571283ce3f744d0a448204b94024520764d5fbba538dfd0cca82e888f3df1560ef6fb97b66c668931108bc8dbf10ce6cf24fa48b5e4f5f2e627c735bb07c001c"
		params = map[string]interface{}{
			"identity": Alice.Address.String(),
			"to":       Bob.Address.String(),
			"value":    "1",
			"fee":      "1",
			"nonce":    4,
			"sig": map[string]interface{}{
				"value": sig,
			},
		}
		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/tx_state?tx_hash=" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "success")

		uri = "/v1/account_info?identity=" + Alice.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "96")

		uri = "/v1/account_info?identity=" + Bob.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "203")

		//
		// Test POST /hashed_transfer
		//
		secret := "0x1111111111111111111111111111111111111111111111111111111111111111"
		commit := "0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc"
		uri = "/v1/hashed_transfer"
		sig = "0xe3675fdc2d6d68e156f6b2b860fdf60c3fb7af59e1fdabb2c71791f4ebf352f645005b4f599d6461792d88ffc42a50b815e1a781c91bb09d5d5a97389695c8341b"
		params = map[string]interface{}{
			"identity":    Bob.Address.String(),
			"to":          Alice.Address.String(),
			"value":       "2",
			"hash_commit": commit,
			"expiry":      999999999999,
			"fee":         "0",
			"nonce":       1,
			"sig": map[string]interface{}{
				"value": sig,
			},
		}

		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		htlcTxHash := res["tx_hash"].(string)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/tx_state?tx_hash=" + htlcTxHash
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "pending")

		uri = "/v1/account_info?identity=" + Alice.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "96")

		uri = "/v1/account_info?identity=" + Bob.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "201")

		//
		// Test POST /claim_hashed_transfer
		//
		uri = "/v1/claim_hashed_transfer"
		sig = "0xa276348dbf18a4f144008c7d382493bf970401a935061e97433fb6cb3c918da52f440736c8299b6644933b632e0cc292a374385316fa728bba2933f394a120f91b"
		params = map[string]interface{}{
			"identity":     Alice.Address.String(),
			"htlc_tx_hash": htlcTxHash,
			"secret":       secret,
			"nonce":        5,
			"sig": map[string]interface{}{
				"value": sig,
			},
		}

		res, code = request(router, "POST", uri, params)
		So(res["error"], ShouldBeNil)
		So(code, ShouldEqual, http.StatusOK)
		So(res, ShouldContainKey, "tx_hash")
		txHashHex = res["tx_hash"].(string)
		appHeight += 2

		if err := rpcclient.WaitForHeight(client, appHeight, nil); err != nil {
			t.Error(err)
		}

		uri = "/v1/tx_state?tx_hash=" + txHashHex
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "success")

		uri = "/v1/tx_state?tx_hash=" + htlcTxHash
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["status"], ShouldEqual, "success")

		uri = "/v1/account_info?identity=" + Alice.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "98")

		uri = "/v1/account_info?identity=" + Bob.Address.String()
		res, _ = request(router, "GET", uri, nil)
		So(res["error"], ShouldBeNil)
		So(res["balance"], ShouldEqual, "201")
	})
}
