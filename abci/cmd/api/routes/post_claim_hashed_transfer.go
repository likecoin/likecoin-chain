package routes

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
)

type claimHashedTrasnferJSON struct {
	Identity   string    `json:"identity" binding:"required,identity"`
	HTLCTxHash string    `json:"htlc_tx_hash" binding:"required,bytes32"`
	Secret     string    `json:"secret" binding:"required"`
	Nonce      int64     `json:"nonce" binding:"required,min=1"`
	Sig        Signature `json:"sig" binding:"required"`
}

func postClaimHashedTransfer(c *gin.Context) {
	var json claimHashedTrasnferJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	htlcTxHashHex := json.HTLCTxHash
	if htlcTxHashHex[0:2] == "0x" {
		htlcTxHashHex = htlcTxHashHex[2:]
	}
	htlcTxHash, _ := hex.DecodeString(htlcTxHashHex)

	var secret []byte
	if len(json.Secret) > 2 && json.Secret[0:2] == "0x" {
		json.Secret = json.Secret[2:]
	}
	if len(json.Secret) == 0 {
		secret = nil
	} else if len(json.Secret) == 64 {
		var err error
		secret, err = hex.DecodeString(json.Secret)
		if err != nil || len(secret) != 32 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid secret"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid secret"})
		return
	}

	tx := txs.ClaimHashedTransferTransaction{
		From:       types.NewIdentifier(json.Identity),
		HTLCTxHash: htlcTxHash,
		Secret:     secret,
		Nonce:      uint64(json.Nonce),
	}
	switch json.Sig.Type {
	case "eip712":
		tx.Sig = &txs.ClaimHashedTransferEIP712Signature{EIP712Signature: txs.SigEIP712(json.Sig.Value)}
	default:
		tx.Sig = &txs.ClaimHashedTransferJSONSignature{JSONSignature: txs.Sig(json.Sig.Value)}
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
