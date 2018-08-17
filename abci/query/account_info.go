package query

import (
	"github.com/likecoin/likechain/abci/context"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryAccountInfo(ctx context.ImmutableContext, reqQuery abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{}
}

func init() {
	registerQueryHandler("account_info", queryAccountInfo)
}
