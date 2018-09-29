package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/types"
)

type registerJSON struct {
	Addr string `json:"addr" binding:"required,eth_addr"`
	Sig  string `json:"sig" binding:"required,eth_sig"`
}

func postRegister(c *gin.Context) {
	var json registerJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := types.RegisterTransaction{
		Addr: types.NewAddressFromHex(json.Addr),
		Sig:  types.NewSignatureFromHex(json.Sig),
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

	id := types.NewLikeChainID(res.Data)
	c.JSON(http.StatusOK, gin.H{
		"id": id.ToString(),
	})
}
