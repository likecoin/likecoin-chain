package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type blockQuery struct {
	Height int64 `form:"height" binding:"required"`
}

func getBlock(c *gin.Context) {
	var query blockQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := tendermint.Block(&query.Height)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result.Block,
	})
}
