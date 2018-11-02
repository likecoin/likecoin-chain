package routes

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
)

type depositApprovalJSON struct {
	Identity      string `json:"identity" binding:"required,identity"`
	DepositTxHash string `json:"deposit_tx_hash" binding:"required,txHash"`
	Nonce         int64  `json:"nonce" binding:"required,min=1"`
	Sig           string `json:"sig" binding:"required,eth_sig"`
}

func postDepositApproval(c *gin.Context) {
	var json depositApprovalJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	depositTxHashHex := json.DepositTxHash
	if depositTxHashHex[0:2] == "0x" {
		depositTxHashHex = depositTxHashHex[2:]
	}
	depositTxHash, _ := hex.DecodeString(depositTxHashHex)

	tx := txs.DepositApprovalTransaction{
		Approver:      types.NewIdentifier(json.Identity),
		DepositTxHash: depositTxHash,
		Nonce:         uint64(json.Nonce),
		Sig:           &txs.DepositApprovalJSONSignature{JSONSignature: txs.Sig(json.Sig)},
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
