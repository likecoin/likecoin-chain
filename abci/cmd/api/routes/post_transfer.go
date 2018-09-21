package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
)

type transferTargetJSON struct {
	Identity string `json:"identity" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Remark   string `json:"remark"`
}
type transferJSON struct {
	Identity string               `json:"identity" binding:"required"`
	To       []transferTargetJSON `json:"to" binding:"required"`
	Nonce    int64                `json:"nonce" binding:"required"`
	Fee      string               `json:"fee" binding:"required"`
	Sig      string               `json:"sig" binding:"required"`
}

func postTransfer(c *gin.Context) {
	var json transferJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidBigIntegerString(json.Fee) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer fee"})
		return
	}

	tx := types.TransferTransaction{
		From:   types.NewIdentifier(json.Identity),
		ToList: make([]*types.TransferTransaction_TransferTarget, len(json.To)),
		Nonce:  uint64(json.Nonce),
		Fee:    types.NewBigInteger(json.Fee),
		Sig:    types.NewSignatureFromHex(json.Sig),
	}

	for i, t := range json.To {
		if !utils.IsValidBigIntegerString(t.Value) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer value"})
			return
		}

		tx.ToList[i] = &types.TransferTransaction_TransferTarget{
			To:     types.NewIdentifier(t.Identity),
			Value:  types.NewBigInteger(t.Value),
			Remark: []byte(t.Remark),
		}
	}

	data, err := proto.Marshal(tx.ToTransaction())
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	result, err := tendermint.BroadcastTxCommit(data)
	if res := result.CheckTx; res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  res.Code,
			"error": res.Info,
		})
		return
	}

	if res := result.DeliverTx; res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  res.Code,
			"error": res.Info,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
