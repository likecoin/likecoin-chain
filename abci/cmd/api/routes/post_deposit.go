package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
)

type depositInputJSON struct {
	FromAddr string `json:"from_addr" binding:"required,eth_addr"`
	Value    string `json:"value" binding:"required,biginteger"`
}

type depositJSON struct {
	Identity    string             `json:"identity" binding:"required,identity"`
	BlockNumber uint64             `json:"block_number" binding:"required"`
	Inputs      []depositInputJSON `json:"inputs" binding:"required"`
	Nonce       int64              `json:"nonce" binding:"required,min=1"`
	Sig         string             `json:"sig" binding:"required,eth_sig"`
}

func postDeposit(c *gin.Context) {
	var json depositJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := txs.DepositTransaction{
		Proposer: types.NewIdentifier(json.Identity),
		Proposal: deposit.Proposal{
			BlockNumber: json.BlockNumber,
			Inputs:      make([]deposit.Input, 0, len(json.Inputs)),
		},
		Nonce: uint64(json.Nonce),
		Sig:   &txs.DepositJSONSignature{JSONSignature: txs.Sig(json.Sig)},
	}
	for _, input := range json.Inputs {
		value, ok := types.NewBigIntFromString(input.Value)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input value"})
			return
		}
		tx.Proposal.Inputs = append(tx.Proposal.Inputs, deposit.Input{
			FromAddr: *types.Addr(input.FromAddr),
			Value:    value,
		})
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
