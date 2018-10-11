package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type addressInfoQuery struct {
	Addr string `form:"addr" binding:"required,eth_addr"`
}

func getAddressInfo(c *gin.Context) {
	var query addressInfoQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := tendermint.ABCIQuery("address_info", []byte(query.Addr))
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
