package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type withdrawProofQuery struct {
	Identity string `form:"identity" binding:"required,identity"`
	ToAddr   string `form:"to_addr" binding:"required,eth_addr"`
	Value    string `form:"value" binding:"required,biginteger"`
	Nonce    uint64 `form:"nonce" binding:"required,min=1"`
	Fee      string `form:"fee" binding:"required,biginteger"`
	Height   int64  `form:"height" binding:"required,min=0"`
}

func getWithdrawProof(c *gin.Context) {
	var query withdrawProofQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidBigIntegerString(query.Value) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdraw value"})
		return
	}

	// Get LikeChain ID
	result, err := tendermint.ABCIQuery("account_info", []byte(query.Identity))
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

	var accountInfo gin.H
	if err := json.Unmarshal(result.Response.Value, &accountInfo); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	idStr := accountInfo["id"].(string)

	// Get withdraw proof
	tx := types.WithdrawTransaction{
		From:   types.NewIdentifier(idStr),
		ToAddr: types.NewAddressFromHex(query.ToAddr),
		Value:  types.NewBigInteger(query.Value),
		Nonce:  query.Nonce,
		Fee:    types.NewBigInteger(query.Fee),
	}

	result, err = tendermint.ABCIQueryWithOptions(
		"withdraw_proof",
		tx.Pack(),
		rpcclient.ABCIQueryOptions{Height: query.Height},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resQuery = result.Response
	if resQuery.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  resQuery.Code,
			"error": resQuery.Info,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    idStr,
		"proof": base64.StdEncoding.EncodeToString(result.Response.Value),
	})
}
