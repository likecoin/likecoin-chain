package types_test

import (
	"testing"
	"time"

	"github.com/likecoin/likechain/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	accounts := testutil.CreateIncrementalAccounts(2)
	revealTime := time.Now()
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

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
				Params: types.DefaultParams(),
				MintableNftList: []types.MintableNFT{
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
						RevealTime: revealTime,
						ClassId:    "0",
					},
					{
						RevealTime: revealTime,
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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated classesByISCN",
			genState: &types.GenesisState{
				ClassesByIscnList: []types.ClassesByISCN{
					{
						IscnIdPrefix: "0",
					},
					{
						IscnIdPrefix: "0",
					},
				},
				Params: types.DefaultParams(),
			},
			valid: false,
		},
		{
			desc: "duplicated classesByAccount",
			genState: &types.GenesisState{
				ClassesByAccountList: []types.ClassesByAccount{
					{
						Account: accounts[0].String(),
					},
					{
						Account: accounts[0].String(),
					},
				},
				Params: types.DefaultParams(),
			},
			valid: false,
		},
		{
			desc: "empty mint price denom param",
			genState: &types.GenesisState{
				ClassesByAccountList: []types.ClassesByAccount{
					{
						Account: accounts[0].String(),
					},
					{
						Account: accounts[0].String(),
					},
				},
				Params: types.Params{
					MintPriceDenom: "",
				},
			},
			valid: false,
		},
		{
			desc: "invalid mint price denom param",
			genState: &types.GenesisState{
				ClassesByAccountList: []types.ClassesByAccount{
					{
						Account: accounts[0].String(),
					},
					{
						Account: accounts[0].String(),
					},
				},
				Params: types.Params{
					MintPriceDenom: "123456",
				},
			},
			valid: false,
		},
		{
			desc: "duplicated mintableNFT",
			genState: &types.GenesisState{
				MintableNftList: []types.MintableNFT{
					{
						ClassId: "0",
						Id:      "0",
					},
					{
						ClassId: "0",
						Id:      "0",
					},
				},
			},
		},
		{
			desc: "duplicated classRevealQueueEntry",
			genState: &types.GenesisState{
				ClassRevealQueue: []types.ClassRevealQueueEntry{
					{
						RevealTime: revealTime,
						ClassId:    "0",
					},
					{
						RevealTime: revealTime,
						ClassId:    "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated offer",
			genState: &types.GenesisState{
				OfferList: []types.Offer{
					{
						ClassId: "0",
						NftId:   "0",
						Buyer:   accounts[0].String(),
					},
					{
						ClassId: "0",
						NftId:   "0",
						Buyer:   accounts[0].String(),
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated listing",
			genState: &types.GenesisState{
				ListingList: []types.Listing{
					{
						ClassId: "0",
						NftId:   "0",
						Seller:  accounts[0].String(),
					},
					{
						ClassId: "0",
						NftId:   "0",
						Seller:  accounts[0].String(),
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
