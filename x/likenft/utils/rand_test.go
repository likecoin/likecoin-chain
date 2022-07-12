package utils_test

import (
	"encoding/hex"
	"testing"

	testkeeper "github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/utils"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestRandSeedFromLastBlockHash(t *testing.T) {
	_, ctx := testkeeper.LikenftKeeper(t)
	// hash from mainnet block 3561600
	hash1, err := hex.DecodeString("6C48857E96C33E00E067C3E9073F3A6089BB2CD05BCC5783953302E5894C38A0")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash1,
	})
	seed1 := utils.RandSeedFromLastBlock(ctx)
	require.NotZero(t, seed1)

	// hash from mainnet block 1
	hash2, err := hex.DecodeString("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
	require.NoError(t, err)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		AppHash: hash2,
	})
	seed2 := utils.RandSeedFromLastBlock(ctx)
	require.NotZero(t, seed2)

	require.NotEqual(t, seed1, seed2)
}
