package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/utils"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type contractUpdateProofQuery struct {
	ContractIndex uint64 `form:"contract_index" binding:"required,min=1"`
	Height        int64  `form:"height" binding:"required,min=1"`
}

func getContractUpdateProof(c *gin.Context) {
	var query contractUpdateProofQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := tendermint.ABCIQueryWithOptions(
		"contract_update_proof",
		utils.EncodeUint64(query.ContractIndex),
		rpcclient.ABCIQueryOptions{Height: query.Height},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resQuery := result.Response
	if resQuery.IsErr() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  resQuery.Code,
			"error": resQuery.Info,
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(result.Response.Value))
}
