package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/types"
)

type depositEventJSON struct {
	From  string `json:"from" binding:"required,eth_addr"`
	Value string `json:"value" binding:"required,biginteger"`
}

type depositJSON struct {
	BlockNumber uint64              `json:"block" binding:"required"`
	Events      []*depositEventJSON `json:"events" binding:"required"`
}

func postDeposit(c *gin.Context) {
	var json depositJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := types.DepositTransaction{
		BlockNumber: json.BlockNumber,
		Deposits:    make([]*types.DepositTransaction_DepositEvent, len(json.Events)),
	}
	for i, e := range json.Events {
		tx.Deposits[i] = &types.DepositTransaction_DepositEvent{
			FromAddr: types.NewAddressFromHex(e.From),
			Value:    types.NewBigInteger(e.Value),
		}
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
