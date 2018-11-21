package routes

import (
	"github.com/gin-gonic/gin"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var tendermint rpcclient.Client

// Initialize initializes routes
func Initialize(router *gin.Engine, client rpcclient.Client) {
	tendermint = client

	v1 := router.Group("/v1")

	v1.POST("/register", postRegister)
	v1.POST("/deposit", postDeposit)
	v1.POST("/transfer", postTransfer)
	v1.POST("/withdraw", postWithdraw)
	v1.POST("/hashed_transfer", postHashedTransfer)
	v1.POST("/claim_hashed_transfer", postClaimHashedTransfer)
	v1.POST("/simple_transfer", postSimpleTransfer)

	v1.GET("/account_info", getAccountInfo)
	v1.GET("/address_info", getAddressInfo)
	v1.GET("/tx_state", getTxState)
	v1.GET("/withdraw_proof", getWithdrawProof)
	v1.GET("/block", getBlock)
}
