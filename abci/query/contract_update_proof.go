package query

import (
	"encoding/json"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/contract"
	"github.com/likecoin/likechain/abci/utils"

	abci "github.com/tendermint/tendermint/abci/types"
)

func queryContractUpdateProof(
	state context.IMutableState,
	reqQuery abci.RequestQuery,
) response.R {
	height := reqQuery.Height
	if height <= 0 {
		log.
			WithField("height", height).
			Debug("Invalid height in contract update proof query")
		return response.QueryContractUpdateProofInvalidHeight
	}

	metadata := state.GetMetadataAtHeight(height)
	if metadata == nil || metadata.WithdrawTreeVersion <= 0 {
		log.
			WithField("height", height).
			Debug("Version not found for height in contract update proof query")
		return response.QueryContractUpdateProofInvalidHeight
	}

	tree, err := state.MutableWithdrawTree().GetImmutable(metadata.WithdrawTreeVersion)
	if tree == nil || err != nil {
		log.
			WithError(err).
			WithField("height", height).
			WithField("version", metadata.WithdrawTreeVersion).
			Debug("Cannot get version in withdraw tree in contract update proof query")
		return response.QueryContractUpdateProofInvalidHeight
	}

	contractIndexBytes := reqQuery.Data
	if len(contractIndexBytes) != 8 {
		log.Debug("Invalid contractIndexBytes length in contract update proof query")
		return response.QueryContractUpdateProofNotExist
	}

	contractIndex := utils.DecodeUint64(contractIndexBytes)
	contractAddr, proof := contract.GetUpdateExecutionWithProof(state, contractIndex, metadata.WithdrawTreeVersion)
	if contractAddr == nil || proof == nil {
		return response.QueryContractUpdateProofNotExist
	}

	resJSON, err := json.Marshal(map[string]interface{}{
		"contract_address": contractAddr.String(),
		"proof":            proof,
	})

	return response.Success.Merge(response.R{
		Data: resJSON,
	})
}

func init() {
	registerQueryHandler("contract_update_proof", queryContractUpdateProof)
}
