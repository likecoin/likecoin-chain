package query

import (
	"encoding/json"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
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

	metadata := state.GetMetadataAtHeight(height)
	if metadata == nil || metadata.WithdrawTreeVersion <= 0 {
		log.
			WithField("height", height).
			Debug("Version not found for height in withdraw proof query")
		return response.QueryWithdrawProofInvalidHeight
	}

	tree, err := state.MutableWithdrawTree().GetImmutable(metadata.WithdrawTreeVersion)
	if tree == nil || err != nil {
		log.
			WithError(err).
			WithField("height", height).
			WithField("version", metadata.WithdrawTreeVersion).
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
			WithField("packed_tx", cmn.HexBytes(packedTx)).
			Debug("Cannot get proof in withdraw proof query")
		return response.QueryWithdrawProofNotExist
	}

	proofJSONBytes, err := json.Marshal(proof)
	if err != nil {
		log.
			WithField("tx", packedTx).
			WithField("proof", proof).
			Panic("Cannot marshal proof as JSON in withdraw proof query")
	}

	return response.Success.Merge(response.R{
		Data: proofJSONBytes,
	})
}

func init() {
	registerQueryHandler("withdraw_proof", queryWithdrawProof)
}
