package routes

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/types"
)

type transferTargetJSON struct {
	Identity string `json:"identity" binding:"required,identity"`
	Value    string `json:"value" binding:"required,biginteger"`
	Remark   string `json:"remark" binding:"base64"`
}
type transferJSON struct {
	Identity string               `json:"identity" binding:"required,identity"`
	To       []transferTargetJSON `json:"to" binding:"required"`
	Nonce    int64                `json:"nonce" binding:"required,min=1"`
	Fee      string               `json:"fee" binding:"required,biginteger"`
	Sig      string               `json:"sig" binding:"required,eth_sig"`
}

func postTransfer(c *gin.Context) {
	var json transferJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		remark, _ := base64.StdEncoding.DecodeString(t.Remark)
		tx.ToList[i] = &types.TransferTransaction_TransferTarget{
			To:     types.NewIdentifier(t.Identity),
			Value:  types.NewBigInteger(t.Value),
			Remark: remark,
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

	if res := result.DeliverTx; res.IsErr() {
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
