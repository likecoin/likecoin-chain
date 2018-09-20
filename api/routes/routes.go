package routes

import (
	"github.com/gin-gonic/gin"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
)

var tendermint *rpcClient.HTTP

// Initialize initializes routes
func Initialize(router *gin.Engine, client *rpcClient.HTTP) {
	tendermint = client

	v1 := router.Group("/v1")

	v1.POST("/register", postRegister)
	v1.POST("/transfer", postTransfer)

	v1.GET("/account_info", getAccountInfo)
	v1.GET("/tx_state", getTxState)
	v1.GET("/block", getBlock)
}
