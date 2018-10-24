package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/txs"
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

	value, ok := types.NewBigIntFromString(json.Value)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdraw value"})
		return
	}

	fee, ok := types.NewBigIntFromString(json.Fee)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdraw fee"})
		return
	}

	tx := txs.WithdrawTransaction{
		From:   types.NewIdentifier(json.Identity),
		ToAddr: *types.Addr(json.ToAddr),
		Value:  value,
		Nonce:  json.Nonce,
		Fee:    fee,
		Sig:    &txs.WithdrawJSONSignature{JSONSignature: txs.Sig(json.Sig)},
	}

	data := txs.EncodeTx(&tx)

	result, err := tendermint.BroadcastTxCommit(data)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	if res := result.CheckTx; res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"tx_hash": result.Hash,
			"code":    res.Code,
			"error":   res.Info,
		})
		return
	}

	res := result.DeliverTx
	if res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"tx_hash": result.Hash,
			"code":    res.Code,
			"error":   res.Info,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tx_hash": result.Hash,
	})
}
