package tendermint

import (
	tmRPC "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"

	logger "github.com/likecoin/likechain/services/log"
)

var log = logger.L

// GetSignedHeader returns the signed header at the given height
func GetSignedHeader(tmClient *tmRPC.HTTP, height int64) types.SignedHeader {
	commit, err := tmClient.Commit(&height)
	if err != nil {
		panic(err)
	}
	return commit.SignedHeader
}

// GetValidators returns the current validators
func GetValidators(tmClient *tmRPC.HTTP) []types.Validator {
	rawConsensusState, err := tmClient.DumpConsensusState()
	if err != nil {
		panic(err)
	}

	jsonRes := struct {
		Validators struct {
			Validators []types.Validator `json:"validators"`
		} `json:"validators"`
	}{}

	err = AminoCodec().UnmarshalJSON(rawConsensusState.RoundState, &jsonRes)
	if err != nil {
		panic(err)
	}

	return jsonRes.Validators.Validators
}

// GetHeight returns the current height
func GetHeight(tmClient *tmRPC.HTTP) int64 {
	abciInfo, err := tmClient.ABCIInfo()
	if err != nil {
		panic(err)
	}
	return abciInfo.Response.GetLastBlockHeight()
}
