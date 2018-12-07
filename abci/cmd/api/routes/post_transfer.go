package routes

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
)

type transferTargetJSON struct {
	Identity string `json:"identity" binding:"required,identity"`
	Value    string `json:"value" binding:"required,biginteger"`
	Remark   string `json:"remark" binding:"base64"`
}
type transferJSON struct {
	Identity string               `json:"identity" binding:"required,identity"`
	Outputs  []transferTargetJSON `json:"outputs" binding:"required"`
	Nonce    int64                `json:"nonce" binding:"required,min=1"`
	Fee      string               `json:"fee" binding:"required,biginteger"`
	Sig      Signature            `json:"sig" binding:"required"`
}

func postTransfer(c *gin.Context) {
	var json transferJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fee, ok := types.NewBigIntFromString(json.Fee)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer fee"})
		return
	}

	tx := txs.TransferTransaction{
		From:    types.NewIdentifier(json.Identity),
		Outputs: make([]txs.TransferOutput, len(json.Outputs)),
		Nonce:   uint64(json.Nonce),
		Fee:     fee,
	}
	switch json.Sig.Type {
	case "eip712":
		c.JSON(http.StatusBadRequest, gin.H{"error": "EIP-712 signature not supported for Transfer transaction"})
		return
	default:
		tx.Sig = &txs.TransferJSONSignature{JSONSignature: txs.Sig(json.Sig.Value)}
	}

	for i, t := range json.Outputs {
		value, ok := types.NewBigIntFromString(t.Value)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer value"})
			return
		}
		remark, _ := base64.StdEncoding.DecodeString(t.Remark)
		tx.Outputs[i] = txs.TransferOutput{
			To:     types.NewIdentifier(t.Identity),
			Value:  value,
			Remark: remark,
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
