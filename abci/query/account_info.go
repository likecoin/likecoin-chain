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
	var err error
	identity := string(reqQuery.Data)

	// Get LikeChain ID
	ethAddr := common.HexToAddress(identity)
	id := account.AddressToLikeChainID(state, ethAddr)

	if id == nil {
		log.
			WithField("addr", ethAddr.Hex()).
			Debug("Identity is not an address")

		id, err = types.NewLikeChainIDFromString(identity)
		if err != nil || !account.IsLikeChainIDRegistered(state, id) {
			log.
				WithField("identity", identity).
				Info(response.QueryInvalidIdentifier.Info)
			return response.QueryInvalidIdentifier
		}
	}

	balance := account.FetchBalance(state, id.ToIdentifier())
	nextNonce := account.FetchNextNonce(state, id)

	return jsonMap{
		"balance":   balance.String(),
		"nextNonce": nextNonce,
	}.ToResponse()
}

func init() {
	registerQueryHandler("account_info", queryAccountInfo)
}
