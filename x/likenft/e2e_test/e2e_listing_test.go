package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	nft "github.com/likecoin/likecoin-chain/v3/backport/cosmos-sdk/v0.46.0-rc1/x/nft"
	nftcli "github.com/likecoin/likecoin-chain/v3/backport/cosmos-sdk/v0.46.0-rc1/x/nft/client/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/testutil/network"

	cli "github.com/likecoin/likecoin-chain/v3/x/likenft/client/cli"
	types "github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestEndToEndListing(t *testing.T) {
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

	// Setup another account
	out, err := clitestutil.ExecTestCLICmd(
		ctx,
		keys.Commands(ctx.HomeDir),
		[]string{"add", "user2", "--output=json"},
	)
	require.NoError(t, err)
	var keyOut keyring.KeyOutput
	cfg.LegacyAmino.UnmarshalJSON(out.Bytes(), &keyOut)

	user2Address := keyOut.Address
	user2TxArgs := []string{
		fmt.Sprintf("--from=%s", user2Address),
		"--yes",
		"--output=json",
		fmt.Sprintf("--gas-prices=%s", cfg.MinGasPrices),
		"--broadcast-mode=block",
	}

	// Send some coins to acc2
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		bankcli.NewSendTxCmd(),
		append([]string{userAddress.String(), user2Address, fmt.Sprintf("1000000000%s", cfg.BondDenom)}, txArgs...),
	)
	require.NoError(t, err)

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

	// Create class
	out, err = clitestutil.ExecTestCLICmd(
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

	// Create Listing
	price := 123456
	expiration := time.Now().UTC().Add(30 * 24 * time.Hour)
	fullPayToRoyalty := true

	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdCreateListing(),
		append([]string{classId, nftId, fmt.Sprintf("%d", price), expiration.Format(time.RFC3339Nano), fmt.Sprintf("%t", fullPayToRoyalty)}, txArgs...),
	)
	require.NoError(t, err)

	// Check event
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualCreateEvent := parseEventCreateListing(res)
	require.Equal(t, types.EventCreateListing{
		ClassId: classId,
		NftId:   nftId,
		Seller:  userAddress.String(),
	}, actualCreateEvent)

	// Query Listing
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdListingsByNFT(),
		append([]string{classId, nftId}, queryArgs...),
	)
	require.NoError(t, err)
	listingRes := types.QueryListingsByNFTResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &listingRes)

	require.Len(t, listingRes.Listings, 1)
	require.Equal(t, types.Listing{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           userAddress.String(),
		Price:            uint64(price),
		Expiration:       expiration,
		FullPayToRoyalty: fullPayToRoyalty,
	}, listingRes.Listings[0])

	// Update Listing
	newPrice := 987654
	newExpiration := time.Now().UTC().Add(60 * 24 * time.Hour)
	newFullPayToRoyalty := false

	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdUpdateListing(),
		append([]string{classId, nftId, fmt.Sprintf("%d", newPrice), newExpiration.Format(time.RFC3339Nano), fmt.Sprintf("%t", newFullPayToRoyalty)}, txArgs...),
	)
	require.NoError(t, err)

	// Check event
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)
	actualUpdateEvent := parseEventUpdateListing(res)
	require.Equal(t, types.EventUpdateListing{
		ClassId: classId,
		NftId:   nftId,
		Seller:  userAddress.String(),
	}, actualUpdateEvent)

	// Query Listing
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdListingsByNFT(),
		append([]string{classId, nftId}, queryArgs...),
	)
	require.NoError(t, err)
	listingRes = types.QueryListingsByNFTResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &listingRes)
	require.Len(t, listingRes.Listings, 1)
	require.Equal(t, types.Listing{
		ClassId:          classId,
		NftId:            nftId,
		Seller:           userAddress.String(),
		Price:            uint64(newPrice),
		Expiration:       newExpiration,
		FullPayToRoyalty: newFullPayToRoyalty,
	}, listingRes.Listings[0])

	// Buy NFT
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdBuyNFT(),
		append([]string{classId, nftId, userAddress.String(), fmt.Sprintf("%d", newPrice)}, user2TxArgs...),
	)
	require.NoError(t, err)

	// Check event
	res = sdk.TxResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &res)

	actualBuyEvent := parseEventBuyNFT(res)
	require.Equal(t, types.EventBuyNFT{
		ClassId: classId,
		NftId:   nftId,
		Seller:  userAddress.String(),
		Buyer:   user2Address,
		Price:   uint64(newPrice),
	}, actualBuyEvent)

	// Query owner
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		nftcli.GetCmdQueryOwner(),
		append([]string{classId, nftId}, queryArgs...),
	)
	require.NoError(t, err)

	ownerRes := nft.QueryOwnerResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &ownerRes)
	require.Equal(t, nft.QueryOwnerResponse{
		Owner: user2Address,
	}, ownerRes)

	// No listing remaining
	out, err = clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdListingsByNFT(),
		append([]string{classId, nftId}, queryArgs...),
	)
	require.NoError(t, err)
	listingRes = types.QueryListingsByNFTResponse{}
	cfg.Codec.MustUnmarshalJSON(out.Bytes(), &listingRes)
	require.Len(t, listingRes.Listings, 0)
}
