package likenft_test

import (
	"testing"

	"github.com/likecoin/likechain/testutil"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"

	likenfttestutil "github.com/likecoin/likechain/x/likenft/testutil"
)

func TestGenesis(t *testing.T) {
	accounts := testutil.CreateIncrementalAccounts(2)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		ClassesByISCNList: []types.ClassesByISCN{
			{
				IscnIdPrefix: "0",
			},
			{
				IscnIdPrefix: "1",
			},
		},
		ClassesByAccountList: []types.ClassesByAccount{
			{
				Account: accounts[0].String(),
			},
			{
				Account: accounts[1].String(),
			},
		},
		ClaimableNFTList: []types.ClaimableNFT{
			{
				ClassId: "0",
				Id:      "0",
			},
			{
				ClassId: "1",
				Id:      "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, ctrl := likenfttestutil.LikenftKeeperForClaimableTest(t)
	defer ctrl.Finish()

	likenft.InitGenesis(ctx, *k, genesisState)
	got := likenft.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.ClassesByISCNList, got.ClassesByISCNList)
	require.ElementsMatch(t, genesisState.ClassesByAccountList, got.ClassesByAccountList)
	require.ElementsMatch(t, genesisState.ClaimableNFTList, got.ClaimableNFTList)
	// this line is used by starport scaffolding # genesis/test/assert
}
