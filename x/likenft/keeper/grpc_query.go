package keeper

import (
	"github.com/likecoin/likechain/x/likenft/types"
)

var _ types.QueryServer = Keeper{}
