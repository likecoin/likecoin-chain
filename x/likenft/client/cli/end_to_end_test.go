package cli_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	nft "github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	nftcli "github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft/client/cli"
	"github.com/likecoin/likechain/testutil/network"
	iscncli "github.com/likecoin/likechain/x/iscn/client/cli"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"

	cli "github.com/likecoin/likechain/x/likenft/client/cli"
	types "github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func parseEventCreateClass(res sdk.TxResponse) types.EventNewClass {
	actualEvent := types.EventNewClass{}

ParseEventCreateClass:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likecoin.likechain.likenft.EventNewClass" {
				for _, attr := range event.Attributes {
					if attr.Key == "iscnIdPrefix" {
						actualEvent.IscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "classId" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "owner" {
						actualEvent.Owner = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventCreateClass
			}
		}
	}

	return actualEvent
}

func parseEventUpdateClass(res sdk.TxResponse) types.EventUpdateClass {
	actualEvent := types.EventUpdateClass{}

ParseEventUpdateClass:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likecoin.likechain.likenft.EventUpdateClass" {
				for _, attr := range event.Attributes {
					if attr.Key == "iscnIdPrefix" {
						actualEvent.IscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "classId" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "owner" {
						actualEvent.Owner = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventUpdateClass
			}
		}
	}

	return actualEvent
}

func parseEventMintNFT(res sdk.TxResponse) types.EventMintNFT {
	actualEvent := types.EventMintNFT{}

ParseEventMintNFT:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likecoin.likechain.likenft.EventMintNFT" {
				for _, attr := range event.Attributes {
					if attr.Key == "iscnIdPrefix" {
						actualEvent.IscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "classId" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nftId" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "owner" {
						actualEvent.Owner = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventMintNFT
			}
		}
	}

	return actualEvent
}

func parseEventBurnNFT(res sdk.TxResponse) types.EventBurnNFT {
	actualEvent := types.EventBurnNFT{}

ParseEventBurnNFT:
	for _, log := range res.Logs {
		for _, event := range log.Events {
			if event.Type == "likecoin.likechain.likenft.EventBurnNFT" {
				for _, attr := range event.Attributes {
					if attr.Key == "iscnIdPrefix" {
						actualEvent.IscnIdPrefix = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "classId" {
						actualEvent.ClassId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "nftId" {
						actualEvent.NftId = strings.Trim(attr.Value, "\"")
					}
					if attr.Key == "owner" {
						actualEvent.Owner = strings.Trim(attr.Value, "\"")
					}
				}
				break ParseEventBurnNFT
			}
		}
	}

	return actualEvent
}

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

	// Create class
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdNewClass(),
		append([]string{iscnIdPrefix, newClassFile.Name()}, txArgs...),
	)
	require.NoError(t, err)

	// Validate event emitted
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualCreateEvent := parseEventCreateClass(res)
	require.Equal(t, iscnIdPrefix, actualCreateEvent.IscnIdPrefix)
	require.NotEmpty(t, actualCreateEvent.ClassId)
	require.Equal(t, userAddress.String(), actualCreateEvent.Owner)

	// Query class
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdShowClassesByISCN(),
		append([]string{iscnIdPrefix}, queryArgs...),
	)
	require.NoError(t, err)

	// Unmarshal and check class data
	classesRes := types.QueryClassesByISCNResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &classesRes)

	require.Len(t, classesRes.Classes, 1)
	class := classesRes.Classes[0]
	require.Equal(t, "New Class", class.Name)
	require.Equal(t, "CLS", class.Symbol)
	require.Equal(t, "Testing New Class", class.Description)
	require.Equal(t, "ipfs://aabbcc", class.Uri)
	require.Equal(t, "aabbcc", class.UriHash)
	classData := types.ClassData{}
	err = classData.Unmarshal(class.Data.Value)
	require.NoError(t, err)
	expectedMetadata, err := types.JsonInput(`{
	"abc": "def"
}`).Normalize()
	require.NoError(t, err)
	actualMetadata, err := classData.Metadata.Normalize()
	require.NoError(t, err)
	require.Equal(t, expectedMetadata, actualMetadata)
	require.Equal(t, false, classData.Config.Burnable)
	require.Equal(t, iscnIdPrefix, classData.Parent.IscnIdPrefix)

	// Update class
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdUpdateClass(),
		append([]string{class.Id, updateClassFile.Name()}, txArgs...),
	)
	require.NoError(t, err)

	// Validate event emitted
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualUpdateEvent := parseEventUpdateClass(res)
	require.Equal(t, iscnIdPrefix, actualUpdateEvent.IscnIdPrefix)
	require.Equal(t, class.Id, actualUpdateEvent.ClassId)
	require.Equal(t, userAddress.String(), actualUpdateEvent.Owner)

	// Query updated class
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdShowClassesByISCN(),
		append([]string{iscnIdPrefix}, queryArgs...),
	)
	require.NoError(t, err)

	// Unmarshal and check updated class data
	updatedClassesRes := types.QueryClassesByISCNResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &updatedClassesRes)

	require.Len(t, updatedClassesRes.Classes, 1)
	updatedClass := updatedClassesRes.Classes[0]
	require.Equal(t, "Oursky Cat Photos", updatedClass.Name)
	require.Equal(t, "Meowgear", updatedClass.Symbol)
	require.Equal(t, "Photos of our beloved bosses.", updatedClass.Description)
	require.Equal(t, "https://www.facebook.com/chima.fasang", updatedClass.Uri)
	require.Equal(t, "", updatedClass.UriHash)
	updatedClassData := types.ClassData{}
	err = updatedClassData.Unmarshal(updatedClass.Data.Value)
	require.NoError(t, err)
	expectedUpdatedMetadata, err := types.JsonInput(`{
	"name": "Oursky Cat Photos",
	"description": "Photos of our beloved bosses.",
	"image": "ipfs://QmZu3v5qFaTrrkSJC4mz8nLoDbR5kJx1QwMUy9CZhFZjT3",
	"external_link": "https://www.facebook.com/chima.fasang"
}`).Normalize()
	require.NoError(t, err)
	actualUpdatedMetadata, err := updatedClassData.Metadata.Normalize()
	require.NoError(t, err)
	require.Equal(t, expectedUpdatedMetadata, actualUpdatedMetadata)
	require.Equal(t, true, updatedClassData.Config.Burnable)
	require.Equal(t, iscnIdPrefix, updatedClassData.Parent.IscnIdPrefix)

	// Mint NFT
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdMintNFT(),
		append([]string{class.Id, "token1", mintNftFile.Name()}, txArgs...),
	)
	require.NoError(t, err)

	// Validate event emitted
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualMintEvent := parseEventMintNFT(res)
	require.Equal(t, iscnIdPrefix, actualMintEvent.IscnIdPrefix)
	require.Equal(t, class.Id, actualMintEvent.ClassId)
	require.Equal(t, "token1", actualMintEvent.NftId)
	require.Equal(t, userAddress.String(), actualMintEvent.Owner)

	// Query NFT
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		nftcli.GetCmdQueryNFT(),
		append([]string{class.Id, "token1"}, queryArgs...),
	)
	require.NoError(t, err)

	// Unmarshal and check nft data
	nftRes := nft.QueryNFTResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &nftRes)

	require.Equal(t, class.Id, nftRes.Nft.ClassId)
	require.Equal(t, "token1", nftRes.Nft.Id)
	require.Equal(t, "ipfs://QmYXq11iygTghZeyxvTZqpDoTomaX7Vd6Cbv1wuyNxq3Fw", nftRes.Nft.Uri)
	require.Equal(t, "QmYXq11iygTghZeyxvTZqpDoTomaX7Vd6Cbv1wuyNxq3Fw", nftRes.Nft.UriHash)
	nftData := types.NFTData{}
	err = nftData.Unmarshal(nftRes.Nft.Data.Value)
	require.NoError(t, err)
	require.Equal(t, iscnIdPrefix, nftData.ClassParent.IscnIdPrefix)
	expectedNftMetadata, err := types.JsonInput(`{
	"name": "Sleepy Coffee #1",
	"description": "Coffee is very sleepy", 
	"image": "ipfs://QmVhp6V2JdpYftT6LnDPELWCDMkk2aHwQZ1qbWf15KRbaZ",
	"external_url": "ipfs://QmYXq11iygTghZeyxvTZqpDoTomaX7Vd6Cbv1wuyNxq3Fw"
}`).Normalize()
	require.NoError(t, err)
	actualNftMetadata, err := nftData.Metadata.Normalize()
	require.NoError(t, err)
	require.Equal(t, expectedNftMetadata, actualNftMetadata)

	// Burn NFT
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdBurnNFT(),
		append([]string{class.Id, "token1"}, txArgs...),
	)
	require.NoError(t, err)

	// Validate event emitted
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualBurnEvent := parseEventBurnNFT(res)
	require.Equal(t, iscnIdPrefix, actualBurnEvent.IscnIdPrefix)
	require.Equal(t, class.Id, actualBurnEvent.ClassId)
	require.Equal(t, "token1", actualBurnEvent.NftId)
	require.Equal(t, userAddress.String(), actualBurnEvent.Owner)

	// Check NFT is burnt
	_, err = clitestutil.ExecTestCLICmd(
		ctx,
		nftcli.GetCmdQueryNFT(),
		append([]string{class.Id, "token1"}, queryArgs...),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}
