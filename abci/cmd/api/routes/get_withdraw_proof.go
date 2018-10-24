package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type withdrawProofQuery struct {
	Identity string `form:"identity" binding:"required,identity"`
	ToAddr   string `form:"to_addr" binding:"required,eth_addr"`
	Value    string `form:"value" binding:"required,biginteger"`
	Nonce    uint64 `form:"nonce" binding:"required,min=1"`
	Fee      string `form:"fee" binding:"required,biginteger"`
	Height   int64  `form:"height" binding:"required,min=1"`
}

func getWithdrawProof(c *gin.Context) {
	var query withdrawProofQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, ok := types.NewBigIntFromString(query.Value)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdraw value"})
		return
	}

	fee, ok := types.NewBigIntFromString(query.Fee)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdraw fee"})
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
		return
	}

	idStr := accountInfo["id"].(string)

	// Get withdraw proof
	tx := txs.WithdrawTransaction{
		From:   types.IDStr(idStr),
		ToAddr: *types.Addr(query.ToAddr),
		Value:  value,
		Nonce:  query.Nonce,
		Fee:    fee,
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
