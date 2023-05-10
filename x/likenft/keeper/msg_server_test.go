package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/likecoin/likecoin-chain/v4/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func setupMsgServer(t testing.TB, dependedKeepers keepertest.LikenftDependedKeepers) (types.MsgServer, context.Context, *keeper.Keeper) {
	k, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, dependedKeepers)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx), k
}
