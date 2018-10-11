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
		log.Debug("Invalid address strings")
		return response.QueryInvalidIdentifier
	}

	ethAddr := common.HexToAddress(addrHex)
	id := account.AddressToLikeChainID(state, ethAddr)

	if id == nil {
		identifier := types.NewAddressFromHex(addrHex).ToIdentifier()
		balance := account.FetchRawBalance(state, identifier)
		return jsonMap{
			"balance": balance.String(),
		}.ToResponse()
	}

	identifier := id.ToIdentifier()
	balance := account.FetchBalance(state, identifier)
	nextNonce := account.FetchNextNonce(state, id)
	return jsonMap{
		"id":         identifier.ToString(),
		"balance":    balance.String(),
		"next_nonce": nextNonce,
	}.ToResponse()
}

func init() {
	registerQueryHandler("address_info", queryAddressInfo)
}
