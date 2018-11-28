package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/state/contract"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
)

type contractUpdateJSON struct {
	Identity      string `json:"identity" binding:"required,identity"`
	ContractIndex uint64 `json:"contract_index" binding:"required,min=1"`
	ContractAddr  string `json:"contract_addr" binding:"required,eth_addr"`
	Nonce         uint64 `json:"nonce" binding:"required,min=1"`
	Sig           string `json:"sig" binding:"required,eth_sig"`
}

func postContractUpdate(c *gin.Context) {
	var json contractUpdateJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contractAddr, _ := types.NewAddressFromHex(json.ContractAddr)
	tx := txs.ContractUpdateTransaction{
		Proposer: types.NewIdentifier(json.Identity),
		Proposal: contract.Proposal{
			ContractIndex:   json.ContractIndex,
			ContractAddress: *contractAddr,
		},
		Nonce: json.Nonce,
		Sig:   &txs.ContractUpdateJSONSignature{JSONSignature: txs.Sig(json.Sig)},
	}

	data := txs.EncodeTx(&tx)

	result, err := tendermint.BroadcastTxCommit(data)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	if res := result.CheckTx; res.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"tx_hash": result.Hash,
			"code":    res.Code,
			"error":   res.Info,
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
