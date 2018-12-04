package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/state/htlc"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
)

type hashedTransferJSON struct {
	Identity   string    `json:"identity" binding:"required,identity"`
	To         string    `json:"to" binding:"required,identity"`
	Value      string    `json:"value" binding:"required,biginteger"`
	Fee        string    `json:"fee" binding:"required,biginteger"`
	HashCommit string    `json:"hash_commit" binding:"required,bytes32"`
	Expiry     int64     `json:"expiry" binding:"required"`
	Nonce      int64     `json:"nonce" binding:"required,min=1"`
	Sig        Signature `json:"sig" binding:"required"`
}

func postHashedTransfer(c *gin.Context) {
	var json hashedTransferJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, ok := types.NewBigIntFromString(json.Value)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value"})
		return
	}

	fee, ok := types.NewBigIntFromString(json.Value)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fee"})
		return
	}

	commitSlice, err := utils.Hex2Bytes(json.HashCommit)
	if err != nil || len(commitSlice) != 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hash commit"})
		return
	}
	commit := [32]byte{}
	copy(commit[:], commitSlice)

	tx := txs.HashedTransferTransaction{
		HashedTransfer: htlc.HashedTransfer{
			From:       types.NewIdentifier(json.Identity),
			To:         types.NewIdentifier(json.To),
			Value:      value,
			HashCommit: commit,
			Expiry:     json.Expiry,
		},
		Fee:   fee,
		Nonce: uint64(json.Nonce),
	}
	switch json.Sig.Type {
	case "eip712":
		tx.Sig = &txs.HashedTransferEIP712Signature{EIP712Signature: txs.SigEIP712(json.Sig.Value)}
	default:
		tx.Sig = &txs.HashedTransferJSONSignature{JSONSignature: txs.Sig(json.Sig.Value)}
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
