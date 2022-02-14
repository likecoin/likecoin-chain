package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.LikenftKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
