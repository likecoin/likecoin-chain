package keeper_test

import (
	"testing"

	testkeeper "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.LikenftKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
