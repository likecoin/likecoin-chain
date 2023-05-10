package keeper

import (
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

var _ types.QueryServer = Keeper{}
