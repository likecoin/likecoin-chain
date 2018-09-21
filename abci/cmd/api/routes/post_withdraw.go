package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/likecoin/likechain/abci/types"
)

type withdrawJSON struct {
	Identity string `json:"identity"`
	ToAddr   string `json:"to_addr"`
	Value    string `json:"value"`
	Nonce    uint64 `json:"nonce"`
	Fee      string `json:"fee"`
	Sig      string `json:"sig"`
}

func postWithdraw(c *gin.Context) {
	var json withdrawJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := types.WithdrawTransaction{
		From:   types.NewIdentifier(json.Identity),
		ToAddr: types.NewAddressFromHex(json.ToAddr),
		Value:  types.NewBigInteger(json.Value),
		Nonce:  json.Nonce,
		Fee:    types.NewBigInteger(json.Fee),
		Sig:    types.NewSignatureFromHex(json.Sig),
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
