package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
)

type registerJSON struct {
	Addr string    `json:"addr" binding:"required,eth_addr"`
	Sig  Signature `json:"sig" binding:"required"`
}

func postRegister(c *gin.Context) {
	var json registerJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := txs.RegisterTransaction{
		Addr: *types.Addr(json.Addr),
	}
	switch json.Sig.Type {
	case "eip712":
		tx.Sig = &txs.RegisterEIP712Signature{EIP712Signature: txs.SigEIP712(json.Sig.Value)}
	default:
		tx.Sig = &txs.RegisterJSONSignature{JSONSignature: txs.Sig(json.Sig.Value)}
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

	id := types.ID(res.Data)
	c.JSON(http.StatusOK, gin.H{
		"tx_hash": result.Hash,
		"id":      id.String(),
	})
}
