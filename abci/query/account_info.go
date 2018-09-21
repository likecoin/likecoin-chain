package query

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryAccountInfo(
	state context.IMutableState,
	reqQuery abci.RequestQuery,
) response.R {
	identity := string(reqQuery.Data)

	var id *types.LikeChainID
	if common.IsHexAddress(identity) {
		// Convert address to LikeChain ID
		ethAddr := common.HexToAddress(identity)
		id = account.AddressToLikeChainID(state, ethAddr)
		if id != nil {
			identity = id.ToString()
		}
	} else {
		// Decode LikeChain ID from strings
		var err error
		id, err = types.NewLikeChainIDFromString(identity)
		if err != nil {
			log.
				WithError(err).
				Debug("Unable to decode LikeChain ID from strings")
			return response.QueryInvalidIdentifier
		}
	}

	if id == nil || !account.IsLikeChainIDRegistered(state, id) {
		log.Debug(response.QueryInvalidIdentifier.Info)
		return response.QueryInvalidIdentifier
	}

	balance := account.FetchBalance(state, id.ToIdentifier())
	nextNonce := account.FetchNextNonce(state, id)

	return jsonMap{
		"id":         identity,
		"balance":    balance.String(),
		"next_nonce": nextNonce,
	}.ToResponse()
}

func init() {
	registerQueryHandler("account_info", queryAccountInfo)
}
