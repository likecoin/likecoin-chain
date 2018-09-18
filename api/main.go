package main

import (
	"encoding/base64"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/likecoin/likechain/api/types"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
)

var client *rpcClient.HTTP

type registerBody struct {
	Addr string `json:"addr" binding:"required"`
	Sig  string `json:"sig" binding:"required"`
}

func register(c *gin.Context) {
	var json registerBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := types.RegisterTransaction{
		Addr: types.NewAddressFromHex(json.Addr),
		Sig:  types.NewSignatureFromHex(json.Sig),
	}

	data, err := proto.Marshal(tx.ToTransaction())
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	result, err := client.BroadcastTxCommit(data)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	if res := result.CheckTx; res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  res.Code,
			"error": res.Info,
		})
		return
	}

	res := result.DeliverTx
	if res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  res.Code,
			"error": res.Info,
		})
		return
	}

	id := types.NewLikeChainID(res.Data)
	c.JSON(http.StatusOK, gin.H{
		"id": id.ToString(),
	})
}

type transferTargetParams struct {
	Identity string `json:"identity" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Remark   string `json:"remark"`
}
type transferBody struct {
	Identity string                 `json:"identity" binding:"required"`
	To       []transferTargetParams `json:"to" binding:"required"`
	Nonce    int64                  `json:"nonce" binding:"required"`
	Fee      string                 `json:"fee" binding:"required"`
	Sig      string                 `json:"sig" binding:"required"`
}

func transfer(c *gin.Context) {
	var json transferBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := types.TransferTransaction{
		From:   types.NewIdentifier(json.Identity),
		ToList: make([]*types.TransferTransaction_TransferTarget, len(json.To)),
		Nonce:  uint64(json.Nonce),
		Fee:    types.NewBigInteger(json.Fee),
		Sig:    types.NewSignatureFromHex(json.Sig),
	}

	for i, t := range json.To {
		tx.ToList[i] = &types.TransferTransaction_TransferTarget{
			To:     types.NewIdentifier(t.Identity),
			Value:  types.NewBigInteger(t.Value),
			Remark: []byte(t.Remark),
		}
	}

	data, err := proto.Marshal(tx.ToTransaction())
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	result, err := client.BroadcastTxCommit(data)
	if res := result.CheckTx; res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  res.Code,
			"error": res.Info,
		})
		return
	}

	if res := result.DeliverTx; res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  res.Code,
			"error": res.Info,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

type accountInfoQuery struct {
	Identity string `form:"identity" binding:"required"`
}

func accountInfo(c *gin.Context) {
	var query accountInfoQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := client.ABCIQuery("account_info", []byte(query.Identity))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resQuery := result.Response
	if resQuery.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  resQuery.Code,
			"error": resQuery.Info,
		})
		return
	}

	c.Data(
		http.StatusOK,
		"application/json; charset=utf-8",
		result.Response.Value,
	)
}

type txStateQuery struct {
	TxHash string `form:"tx_hash" binding:"required"`
}

func txState(c *gin.Context) {
	var query txStateQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txHash, err := base64.StdEncoding.DecodeString(query.TxHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := client.ABCIQuery("tx_state", txHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resQuery := result.Response
	if resQuery.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  resQuery.Code,
			"error": resQuery.Info,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": string(resQuery.Value),
	})
}

type blockQuery struct {
	Height int64 `form:"height" binding:"required"`
}

func block(c *gin.Context) {
	var query blockQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := client.Block(&query.Height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result.Block,
	})
}

func main() {
	// TODO: Put in config file
	host := os.Getenv("LIKECHAIN_API_CLIENT_HOST")
	if host == "" {
		host = "localhost:26657"
	}

	port := os.Getenv("LIKECHAIN_API_PORT")
	if port == "" {
		port = "3000"
	}

	client = rpcClient.NewHTTP("tcp://"+host, "/websocket")

	router := gin.Default()

	router.POST("/register", register)
	router.POST("/transfer", transfer)

	router.GET("/account_info", accountInfo)
	router.GET("/tx_state", txState)
	router.GET("/block", block)

	router.Run(":" + port)
}
