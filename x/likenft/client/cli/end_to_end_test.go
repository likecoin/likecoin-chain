package cli_test

import (
	"fmt"
	"os"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/testutil/network"
	iscncli "github.com/likecoin/likechain/x/iscn/client/cli"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/stretchr/testify/require"
)

func TestEndToEndNormal(t *testing.T) {
	tempDir := t.TempDir() // swap to ioutil for longlived files when debug
	cfg := network.DefaultConfig()

	// Override x/iscn gas fee denom to avoid the need to seed tokens

	// We do not have account addresses until network spawned
	// And do not want to spend time on seeding tester accounts
	// Fix later if using nanolike is a must for testing features

	iscnGenesis := iscntypes.GenesisState{}
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[iscntypes.StoreKey], &iscnGenesis)
	iscnGenesis.Params.FeePerByte = sdk.NewDecCoin(
		cfg.BondDenom, sdk.NewInt(iscntypes.DefaultFeePerByteAmount),
	)
	cfg.GenesisState[iscntypes.StoreKey] = cfg.Codec.MustMarshalJSON(&iscnGenesis)

	// Setup network
	net := network.New(t, cfg)
	ctx := net.Validators[0].ClientCtx
	txArgs := []string{
		fmt.Sprintf("--from=%s", net.Validators[0].Address.String()),
		"--yes",
		"--output=json",
		fmt.Sprintf("--gas-prices=%s", cfg.MinGasPrices),
		"--broadcast-mode=block",
	}
	queryArgs := []string{
		"--output=json",
	}

	// Seed input files
	createIscnFile, err := os.CreateTemp(tempDir, "create_iscn_*.json")
	require.NoError(t, err)
	require.NotNil(t, createIscnFile)
	_, err = createIscnFile.WriteString(`{
	"recordNotes": "Add IPFS fingerprint",
	"contentFingerprints": [
		"hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e",
		"ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"
	],
	"stakeholders": [
		{
			"entity": {
				"@id": "did:cosmos:5sy29r37gfxvxz21rh4r0ktpuc46pzjrmz29g45",
				"name": "Chung Wu"
			},
			"rewardProportion": 95,
			"contributionType": "http://schema.org/author"
		},
		{
			"rewardProportion": 5,
			"contributionType": "http://schema.org/citation",
			"footprint": "https://en.wikipedia.org/wiki/Fibonacci_number",
			"description": "The blog post referred the matrix form of computing Fibonacci numbers."
		}
	],
	"contentMetadata": {
		"@context": "http://schema.org/",
		"@type": "Article",
		"name": "使用矩陣計算遞歸關係式",
		"description": "An article on computing recursive function with matrix multiplication.",
		"datePublished": "2019-04-19",
		"version": 1,
		"url": "https://nnkken.github.io/post/recursive-relation/",
		"author": "https://github.com/nnkken",
		"usageInfo": "https://creativecommons.org/licenses/by/4.0",
		"keywords": "matrix,recursion"
	}
}
`)
	require.NoError(t, err)
	err = createIscnFile.Close()
	require.NoError(t, err)

	newClassFile, err := os.CreateTemp(tempDir, "new_class_*.json")
	require.NoError(t, err)
	require.NotNil(t, newClassFile)
	_, err = newClassFile.WriteString(`{
	"name": "New Class",
	"symbol": "CLS",
	"description": "Testing New Class",
	"uri": "ipfs://aabbcc",
	"uriHash": "aabbcc",
	"metadata": {
		"abc": "def"
	},
	"burnable": false
}
`)
	require.NoError(t, err)
	err = newClassFile.Close()
	require.NoError(t, err)

	updateClassFile, err := os.CreateTemp(tempDir, "update_class_*.json")
	require.NoError(t, err)
	require.NotNil(t, updateClassFile)
	_, err = updateClassFile.WriteString(`{
	"name": "Oursky Cat Photos",
	"symbol": "Meowgear",
	"description": "Photos of our beloved bosses.",
	"uri": "https://www.facebook.com/chima.fasang",
	"uriHash": "",
	"metadata": {
		"name": "Oursky Cat Photos",
		"description": "Photos of our beloved bosses.",
		"image": "ipfs://QmZu3v5qFaTrrkSJC4mz8nLoDbR5kJx1QwMUy9CZhFZjT3",
		"external_link": "https://www.facebook.com/chima.fasang"
	},
	"burnable": true
}
`)
	require.NoError(t, err)
	err = updateClassFile.Close()
	require.NoError(t, err)

	mintNftFile, err := os.CreateTemp(tempDir, "mint_nft_*.json")
	require.NoError(t, err)
	require.NotNil(t, mintNftFile)
	_, err = mintNftFile.WriteString(`{
	"uri": "ipfs://QmYXq11iygTghZeyxvTZqpDoTomaX7Vd6Cbv1wuyNxq3Fw",
	"uriHash": "QmYXq11iygTghZeyxvTZqpDoTomaX7Vd6Cbv1wuyNxq3Fw",
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

	// Create iscn
	out, err := clitestutil.ExecTestCLICmd(
		ctx,
		iscncli.NewCreateIscnTxCmd(),
		append([]string{createIscnFile.Name()}, txArgs...),
	)
	require.NoError(t, err)

	// Get iscn id prefix created
	res := sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	var iscnIdPrefix string
FindIscnIdPrefix:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "iscn_record" {
				for _, attr := range event.Attributes {
					if attr.Key == "iscn_id_prefix" {
						iscnIdPrefix = attr.Value
						break FindIscnIdPrefix
					}
				}
			}
		}
	}
	require.NotEmpty(t, iscnIdPrefix)
}
