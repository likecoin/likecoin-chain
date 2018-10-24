package query

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryAddressInfo(
	state context.IMutableState,
	reqQuery abci.RequestQuery,
) response.R {
	addrHex := string(reqQuery.Data)

	if !common.IsHexAddress(addrHex) {
	}

	addr, err := types.NewAddressFromHex(addrHex)
	if err != nil {
		log.
			WithError(err).
			Debug("Invalid query address")
		return response.QueryInvalidIdentifier
	}

	id := account.AddressToLikeChainID(state, addr)
	if id == nil {
		balance := account.FetchRawBalance(state, addr)
		return jsonMap{
			"balance": balance.String(),
		}.ToResponse()
	}

	balance := account.FetchBalance(state, id)
	nextNonce := account.FetchNextNonce(state, id)
	return jsonMap{
		"id":         id.String(),
		"balance":    balance.String(),
		"next_nonce": nextNonce,
	}.ToResponse()
}

func init() {
	registerQueryHandler("address_info", queryAddressInfo)
}
