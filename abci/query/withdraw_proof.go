package query

import (
	"encoding/json"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

func queryWithdrawProof(
	state context.IMutableState,
	reqQuery abci.RequestQuery,
) response.R {
	height := reqQuery.Height
	if height <= 0 {
		log.
			WithField("height", height).
			Debug("Invalid height in withdraw proof query")
		return response.QueryWithdrawProofInvalidHeight
	}
	version := state.GetWithdrawVersionAtHeight(height)
	if version < 0 {
		log.
			WithField("height", height).
			Debug("Version not found for height in withdraw proof query")
		return response.QueryWithdrawProofInvalidHeight
	}
	tree, err := state.MutableWithdrawTree().GetImmutable(version)
	if tree == nil || err != nil {
		log.
			WithError(err).
			WithField("height", height).
			WithField("version", version).
			Debug("Cannot get version in withdraw tree in withdraw proof query")
		return response.QueryWithdrawProofInvalidHeight
	}
	packedTx := reqQuery.Data
	if packedTx == nil {
		log.Debug("Data is nil in withdraw proof query")
		return response.QueryWithdrawProofNotExist
	}
	key := crypto.Sha256(packedTx)
	value, proof, err := tree.GetWithProof(key)
	if value == nil || err != nil {
		log.
			WithError(err).
			WithField("tx", packedTx).
			Debug("Cannot get proof in withdraw proof query")
		return response.QueryWithdrawProofNotExist
	}
	valueJSONBytes, err := json.Marshal(proof)
	if err != nil {
		log.
			WithField("tx", packedTx).
			WithField("proof", proof).
			Panic("Cannot marshal proof as JSON in withdraw proof query")
	}
	return response.R{
		Code: 0,
		Data: valueJSONBytes,
	}
}

func init() {
	registerQueryHandler("withdraw_proof", queryWithdrawProof)
}
