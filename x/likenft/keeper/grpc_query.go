package keeper

import (
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

var _ types.QueryServer = Keeper{}
