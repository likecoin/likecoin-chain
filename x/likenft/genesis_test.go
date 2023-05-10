package likenft_test

import (
	"testing"
	"time"

	"github.com/likecoin/likecoin-chain/v4/testutil"
	"github.com/likecoin/likecoin-chain/v4/testutil/nullify"
	"github.com/likecoin/likecoin-chain/v4/x/likenft"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"

	likenfttestutil "github.com/likecoin/likecoin-chain/v4/x/likenft/testutil"
)

func TestGenesis(t *testing.T) {
	accounts := testutil.CreateIncrementalAccounts(2)
	nowTime := time.Now()
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		ClassesByIscnList: []types.ClassesByISCN{
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
		BlindBoxContentList: []types.BlindBoxContent{
			{
				ClassId: "0",
				Id:      "0",
			},
			{
				ClassId: "1",
				Id:      "1",
			},
		},
		ClassRevealQueue: []types.ClassRevealQueueEntry{
			{
				RevealTime: nowTime,
				ClassId:    "0",
			},
			{
				RevealTime: nowTime,
				ClassId:    "1",
			},
		},
		OfferList: []types.Offer{
			{
				ClassId: "0",
				NftId:   "0",
				Buyer:   accounts[0].String(),
			},
			{
				ClassId: "1",
				NftId:   "1",
				Buyer:   accounts[1].String(),
			},
		},
		ListingList: []types.Listing{
			{
				ClassId: "0",
				NftId:   "0",
				Seller:  accounts[0].String(),
			},
			{
				ClassId: "1",
				NftId:   "1",
				Seller:  accounts[1].String(),
			},
		},
		OfferExpireQueue: []types.OfferExpireQueueEntry{
			{
				ExpireTime: nowTime,
				OfferKey:   []byte("0"),
			},
			{
				ExpireTime: nowTime,
				OfferKey:   []byte("1"),
			},
		},
		ListingExpireQueue: []types.ListingExpireQueueEntry{
			{
				ExpireTime: nowTime,
				ListingKey: []byte("0"),
			},
			{
				ExpireTime: nowTime,
				ListingKey: []byte("1"),
			},
		},
		RoyaltyConfigByClassList: []types.RoyaltyConfigByClass{
			{
				ClassId: "0",
			},
			{
				ClassId: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, ctrl := likenfttestutil.LikenftKeeperForBlindBoxTest(t)
	defer ctrl.Finish()

	likenft.InitGenesis(ctx, *k, genesisState)
	got := likenft.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.ClassesByIscnList, got.ClassesByIscnList)
	require.ElementsMatch(t, genesisState.ClassesByAccountList, got.ClassesByAccountList)
	require.ElementsMatch(t, genesisState.BlindBoxContentList, got.BlindBoxContentList)
	require.ElementsMatch(t, genesisState.ClassRevealQueue, got.ClassRevealQueue)
	require.ElementsMatch(t, genesisState.OfferList, got.OfferList)
	require.ElementsMatch(t, genesisState.ListingList, got.ListingList)
	require.ElementsMatch(t, genesisState.OfferExpireQueue, got.OfferExpireQueue)
	require.ElementsMatch(t, genesisState.ListingExpireQueue, got.ListingExpireQueue)
	require.ElementsMatch(t, genesisState.RoyaltyConfigByClassList, got.RoyaltyConfigByClassList)
	// this line is used by starport scaffolding # genesis/test/assert
}
