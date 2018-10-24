package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
)

type depositEventJSON struct {
	From  string `json:"from" binding:"required,eth_addr"`
	Value string `json:"value" binding:"required,biginteger"`
}

type depositJSON struct {
	BlockNumber uint64              `json:"block" binding:"required"`
	Inputs      []*depositEventJSON `json:"inputs" binding:"required"`
}

func postDeposit(c *gin.Context) {
	var json depositJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := txs.DepositTransaction{
		BlockNumber: json.BlockNumber,
		Inputs:      make([]txs.DepositInput, len(json.Inputs)),
	}
	for i, e := range json.Inputs {
		value, ok := types.NewBigIntFromString(e.Value)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input value"})
			return
		}
		tx.Inputs[i] = txs.DepositInput{
			FromAddr: *types.Addr(e.From),
			Value:    value,
		}
	}

	data := txs.EncodeTx(&tx)

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
