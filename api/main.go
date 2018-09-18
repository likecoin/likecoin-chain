package main

import (
	"net/http"

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
	client = rpcClient.NewHTTP("tcp://localhost:26657", "/websocket")

	router := gin.Default()

	router.POST("/register", register)

	router.GET("/account_info", accountInfo)
	router.GET("/block", block)

	router.Run(":3000")
}
