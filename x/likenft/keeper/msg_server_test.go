package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
)

func setupMsgServer(t testing.TB, dependedKeepers keepertest.LikenftDependedKeepers) (types.MsgServer, context.Context) {
	k, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, dependedKeepers)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
