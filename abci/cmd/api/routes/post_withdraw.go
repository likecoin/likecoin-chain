package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/types"
)

type withdrawJSON struct {
	Identity string `json:"identity" binding:"required,identity"`
	ToAddr   string `json:"to_addr" binding:"required,eth_addr"`
	Value    string `json:"value" binding:"required,biginteger"`
	Nonce    uint64 `json:"nonce" binding:"required,min=1"`
	Fee      string `json:"fee" binding:"required,biginteger"`
	Sig      string `json:"sig" binding:"required,eth_sig"`
}

func postWithdraw(c *gin.Context) {
	var json withdrawJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := types.WithdrawTransaction{
		From:   types.NewIdentifier(json.Identity),
		ToAddr: types.NewAddressFromHex(json.ToAddr),
		Value:  types.NewBigInteger(json.Value),
		Nonce:  json.Nonce,
		Fee:    types.NewBigInteger(json.Fee),
		Sig:    types.NewSignatureFromHex(json.Sig),
	}

	data, err := tx.ToTransaction().Encode()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	result, err := tendermint.BroadcastTxCommit(data)
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

	c.JSON(http.StatusOK, gin.H{})
}
