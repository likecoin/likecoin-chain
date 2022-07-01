package e2e_test

import (
	"fmt"
	"os"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/testutil/network"

	cli "github.com/likecoin/likecoin-chain/v3/x/likenft/client/cli"
	types "github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestEndToEndRoyaltyConfig(t *testing.T) {
	tempDir := t.TempDir() // swap to ioutil for longlived files when debug
	cfg := network.DefaultConfig()

	// We do not have account addresses until network spawned
	// And do not want to spend time on seeding tester accounts
	// Fix later if using nanolike is a must for testing features

	likenftGenesis := types.GenesisState{}
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[types.StoreKey], &likenftGenesis)
	likenftGenesis.Params.FeePerByte = sdk.NewDecCoin(
		cfg.BondDenom, sdk.NewInt(types.DefaultFeePerByteAmount),
	)
	likenftGenesis.Params.PriceDenom = cfg.BondDenom
	cfg.GenesisState[types.StoreKey] = cfg.Codec.MustMarshalJSON(&likenftGenesis)

	// Setup network
	net := network.New(t, cfg)
	ctx := net.Validators[0].ClientCtx
	userAddress := net.Validators[0].Address
	txArgs := []string{
		fmt.Sprintf("--from=%s", userAddress.String()),
		"--yes",
		"--output=json",
		fmt.Sprintf("--gas-prices=%s", cfg.MinGasPrices),
		"--broadcast-mode=block",
	}
	queryArgs := []string{
		"--output=json",
	}

	// Seed input files
	newClassFile, err := os.CreateTemp(tempDir, "new_class_*.json")
	require.NoError(t, err)
	require.NotNil(t, newClassFile)
	_, err = newClassFile.WriteString(`{
	"name": "Oursky Cat Photos",
	"symbol": "Meowgear",
	"description": "Photos of our beloved bosses.",
	"uri": "https://www.facebook.com/chima.fasang",
	"uri_hash": "",
	"metadata": {
		"name": "Oursky Cat Photos",
		"description": "Photos of our beloved bosses.",
		"image": "ipfs://QmZu3v5qFaTrrkSJC4mz8nLoDbR5kJx1QwMUy9CZhFZjT3",
		"external_link": "https://www.facebook.com/chima.fasang"
	},
	"config": {
		"burnable": true
	}
}
`)
	require.NoError(t, err)
	err = newClassFile.Close()
	require.NoError(t, err)

	mintNftFile, err := os.CreateTemp(tempDir, "mint_nft_*.json")
	require.NoError(t, err)
	require.NotNil(t, mintNftFile)
	_, err = mintNftFile.WriteString(`{
	"uri": "ipfs://QmYXq11iygTghZeyxvTZqpDoTomaX7Vd6Cbv1wuyNxq3Fw",
	"uri_hash": "QmYXq11iygTghZeyxvTZqpDoTomaX7Vd6Cbv1wuyNxq3Fw",
	"metadata": {
		"name": "Sleepy Coffee #1",
		"description": "Coffee is very sleepy", 
		"image": "ipfs://QmVhp6V2JdpYftT6LnDPELWCDMkk2aHwQZ1qbWf15KRbaZ",
		"external_url": "ipfs://QmYXq11iygTghZeyxvTZqpDoTomaX7Vd6Cbv1wuyNxq3Fw"
	}
}`)
	require.NoError(t, err)
	err = mintNftFile.Close()
	require.NoError(t, err)

	createRoyaltyConfigFile, err := os.CreateTemp(tempDir, "create_royalty_config_*.json")
	require.NoError(t, err)
	require.NotNil(t, createRoyaltyConfigFile)
	_, err = createRoyaltyConfigFile.WriteString(`{
		"rate_basis_points": 123,
		"stakeholders": [
			{
				"account": "cosmos172nhdqasd2t9e8vvqw4cxfnnutt98q7elzluk9",
				"weight": 1
			}
		]
}`)
	require.NoError(t, err)
	err = createRoyaltyConfigFile.Close()
	require.NoError(t, err)

	updateRoyaltyConfigFile, err := os.CreateTemp(tempDir, "update_royalty_config_*.json")
	require.NoError(t, err)
	require.NotNil(t, updateRoyaltyConfigFile)
	_, err = updateRoyaltyConfigFile.WriteString(`{
		"rate_basis_points": 1000,
		"stakeholders": [
			{
				"account": "cosmos172nhdqasd2t9e8vvqw4cxfnnutt98q7elzluk9",
				"weight": 1
			},
			{
				"account": "cosmos1r623mw6k77g6s3t67fy3042u9nshdl49fgvtex",
				"weight": 2
			}
		]
}`)
	require.NoError(t, err)
	err = updateRoyaltyConfigFile.Close()
	require.NoError(t, err)

	seedUser1, err := sdk.AccAddressFromBech32("cosmos172nhdqasd2t9e8vvqw4cxfnnutt98q7elzluk9")
	require.NoError(t, err)
	seedUser2, err := sdk.AccAddressFromBech32("cosmos1r623mw6k77g6s3t67fy3042u9nshdl49fgvtex")
	require.NoError(t, err)

	// Create class
	out, err := clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdNewClass(),
		append([]string{"--account", newClassFile.Name()}, txArgs...),
	)
	require.NoError(t, err)

	res := sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualCreateClassEvent := parseEventCreateClass(res)
	require.NotEmpty(t, actualCreateClassEvent.ClassId)
	require.Empty(t, actualCreateClassEvent.ParentIscnIdPrefix)
	require.Equal(t, userAddress.String(), actualCreateClassEvent.ParentAccount)

	classId := actualCreateClassEvent.ClassId

	// Mint NFT
	nftId := "nft1"
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdMintNFT(),
		append([]string{classId, "--id", nftId, "--input", mintNftFile.Name()}, txArgs...),
	)
	require.NoError(t, err)

	// Create config
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdCreateRoyaltyConfig(),
		append([]string{classId, createRoyaltyConfigFile.Name()}, txArgs...),
	)
	require.NoError(t, err)

	// Check event
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualCreateEvent := parseEventCreateRoyaltyConfig(res)
	require.Equal(t, types.EventCreateRoyaltyConfig{
		ClassId: classId,
	}, actualCreateEvent)

	// Query config
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdShowRoyaltyConfig(),
		append([]string{classId}, queryArgs...),
	)
	require.NoError(t, err)
	configRes := types.QueryRoyaltyConfigResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &configRes)

	require.Equal(t, types.RoyaltyConfig{
		RateBasisPoints: uint64(123),
		Stakeholders: []types.RoyaltyStakeholder{
			{
				Account: seedUser1,
				Weight:  1,
			},
		},
	}, configRes.RoyaltyConfig)

	// Update config
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdUpdateRoyaltyConfig(),
		append([]string{classId, updateRoyaltyConfigFile.Name()}, txArgs...),
	)
	require.NoError(t, err)

	// Check event
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualUpdateEvent := parseEventUpdateRoyaltyConfig(res)
	require.Equal(t, types.EventUpdateRoyaltyConfig{
		ClassId: classId,
	}, actualUpdateEvent)

	// Query config
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdShowRoyaltyConfig(),
		append([]string{classId}, queryArgs...),
	)
	require.NoError(t, err)
	configRes = types.QueryRoyaltyConfigResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &configRes)

	require.Equal(t, types.RoyaltyConfig{
		RateBasisPoints: uint64(1000),
		Stakeholders: []types.RoyaltyStakeholder{
			{
				Account: seedUser1,
				Weight:  1,
			},
			{
				Account: seedUser2,
				Weight:  2,
			},
		},
	}, configRes.RoyaltyConfig)

	// Delete config
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdDeleteRoyaltyConfig(),
		append([]string{classId}, txArgs...),
	)
	require.NoError(t, err)

	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualDeleteEvent := parseEventDeleteRoyaltyConfig(res)
	require.Equal(t, types.EventDeleteRoyaltyConfig{
		ClassId: classId,
	}, actualDeleteEvent)

	// Query config
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdShowRoyaltyConfig(),
		append([]string{classId}, queryArgs...),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}
