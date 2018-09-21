package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
)

type depositEventJSON struct {
	From  string `json:"from" binding:"required"`
	Value string `json:"value" binding:"required"`
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
		if !common.IsHexAddress(e.From) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender address"})
			return
		}

		if !utils.IsValidBigIntegerString(e.Value) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deposit value"})
			return
		}

		tx.Deposits[i] = &types.DepositTransaction_DepositEvent{
			FromAddr: types.NewAddressFromHex(e.From),
			Value:    types.NewBigInteger(e.Value),
		}
	}

	data, err := proto.Marshal(tx.ToTransaction())
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
