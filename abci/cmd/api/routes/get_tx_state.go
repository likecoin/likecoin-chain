package routes

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

type txStateQuery struct {
	TxHash string `form:"tx_hash" binding:"required,hex"`
}

func getTxState(c *gin.Context) {
	var query txStateQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txHashHex := query.TxHash
	if txHashHex[0:2] == "0x" {
		txHashHex = txHashHex[2:]
	}
	txHash, _ := hex.DecodeString(txHashHex)

	result, err := tendermint.ABCIQuery("tx_state", txHash)
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
