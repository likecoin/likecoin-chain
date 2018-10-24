package query

import (
	"encoding/base64"
	"encoding/json"
	"math/big"

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

	// First try decode LikeChainID from strings
	id, err := types.NewLikeChainIDFromString(identity)
	if err != nil {
		// LikeChainID failed, try address
		addr, err := types.NewAddressFromHex(identity)
		if err != nil {
			return response.QueryInvalidIdentifier
		}
		id = account.AddressToLikeChainID(state, addr)
		if id == nil {
			return response.QueryInvalidIdentifier
		}
		identity = id.String()
	}

	if !account.IsLikeChainIDRegistered(state, id) {
		log.Debug(response.QueryInvalidIdentifier.Info)
		return response.QueryInvalidIdentifier
	}

	balance := account.FetchBalance(state, id)
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

// AccountInfoRes represents response data of account_info query
type AccountInfoRes struct {
	ID        []byte
	Balance   *big.Int
	NextNonce uint64
}

// GetAccountInfoRes transforms the raw byte response from account_info query back to AccountInfoRes structure
func GetAccountInfoRes(data []byte) *AccountInfoRes {
	accountInfo := struct {
		ID        string `json:"id"`
		NextNonce uint64 `json:"next_nonce"`
		Balance   string `json:"balance"`
	}{}
	err := json.Unmarshal(data, &accountInfo)
	if err != nil {
		return nil
	}
	balance, succ := new(big.Int).SetString(accountInfo.Balance, 10)
	if !succ {
		return nil
	}
	result := AccountInfoRes{
		Balance:   balance,
		NextNonce: accountInfo.NextNonce,
	}
	idBytes, err := base64.StdEncoding.DecodeString(accountInfo.ID)
	if err == nil {
		result.ID = idBytes
	}
	return &result
}
