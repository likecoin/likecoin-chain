package types_test

import (
	"testing"

	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
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
						Account: "0",
					},
					{
						Account: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated classesByISCN",
			genState: &types.GenesisState{
				ClassesByISCNList: []types.ClassesByISCN{
					{
						IscnIdPrefix: "0",
					},
					{
						IscnIdPrefix: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated classesByAccount",
			genState: &types.GenesisState{
				ClassesByAccountList: []types.ClassesByAccount{
					{
						Account: "0",
					},
					{
						Account: "0",
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
