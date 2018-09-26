package routes

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

type txStateQuery struct {
	TxHash string `form:"tx_hash" binding:"required"`
}

func getTxState(c *gin.Context) {
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
