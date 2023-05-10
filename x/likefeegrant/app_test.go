package likefeegrant_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"

	"github.com/likecoin/likecoin-chain/v4/x/iscn/types"

	testutil "github.com/likecoin/likecoin-chain/v4/testutil"
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	priv3 = secp256k1.GenPrivKey()
	addr3 = sdk.AccAddress(priv3.PubKey().Address())

	fingerprint1 = "hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"
	fingerprint2 = "ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"

	stakeholder1 = types.IscnInput(`
{
	"entity": {
		"@id": "did:cosmos:5sy29r37gfxvxz21rh4r0ktpuc46pzjrmz29g45",
		"name": "Chung Wu"
	},
	"rewardProportion": 95,
	"contributionType": "http://schema.org/author"
}`)

	stakeholder2 = types.IscnInput(`
{
	"rewardProportion": 5,
	"contributionType": "http://schema.org/citation",
	"footprint": "https://en.wikipedia.org/wiki/Fibonacci_number",
	"description": "The blog post referred the matrix form of computing Fibonacci numbers."
}`)

	contentMetadata1 = types.IscnInput(`
{
	"@context": "http://schema.org/",
	"@type": "CreativeWorks",
	"title": "使用矩陣計算遞歸關係式",
	"description": "An article on computing recursive function with matrix multiplication.",
	"datePublished": "2019-04-19",
	"version": 1,
	"url": "https://nnkken.github.io/post/recursive-relation/",
	"author": "https://github.com/nnkken",
	"usageInfo": "https://creativecommons.org/licenses/by/4.0",
	"keywords": "matrix,recursion"
}`)

	contentMetadata2 = types.IscnInput(`
{
	"@context": "http://schema.org/",
	"@type": "CreativeWorks",
	"title": "another work"
}`)
)

func TestFeeGrant(t *testing.T) {
	var msg sdk.Msg
	app := testutil.SetupTestApp([]testutil.GenesisBalance{
		{addr1.String(), "1000000000000000000nanolike"},
		{addr2.String(), "1000000000000000000nanolike"},
		{addr3.String(), "1000000000000000000nanolike"},
	})

	app.NextHeader(1234567890)
	app.SetForTx()
	record := types.IscnRecord{
		RecordNotes:         "some update",
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	result := app.DeliverMsgNoError(t, msg, priv1)
	iscnId := testutil.GetIscnIdFromResult(t, result)

	record = types.IscnRecord{
		RecordNotes:         "new update",
		ContentFingerprints: []string{fingerprint1, fingerprint2},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata2,
	}
	expiration := time.Unix(2000000000, 0)
	authorization := types.NewUpdateAuthorization(iscnId.Prefix.String())
	authzMsg, err := authz.NewMsgGrant(addr1, addr2, authorization, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, authzMsg, priv1)

	balanceBefore1 := app.BankKeeper.GetAllBalances(app.Context, addr1)
	balanceBefore2 := app.BankKeeper.GetAllBalances(app.Context, addr2)

	updateMsg := types.NewMsgUpdateIscnRecord(addr1, iscnId, &record)
	msgExec := authz.NewMsgExec(addr2, []sdk.Msg{updateMsg})

	// fee allowance not granted yet, should fail
	_, err, simErr, _ := app.DeliverMsgsWithFeeGranter([]sdk.Msg{&msgExec}, priv2, addr2)
	require.NoError(t, err)
	require.Error(t, simErr)

	allowance, err := feegrant.NewAllowedMsgAllowance(&feegrant.BasicAllowance{
		SpendLimit: sdk.NewCoins(sdk.NewInt64Coin("nanolike", 1000000000000000000)),
		Expiration: &expiration,
	}, []string{authorization.MsgTypeURL()})
	require.NoError(t, err)
	feegrantMsg, err := feegrant.NewMsgGrantAllowance(allowance, addr2, addr1)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, feegrantMsg, priv2)

	_, err, simErr, deliverErr := app.DeliverMsgsWithFeeGranter([]sdk.Msg{&msgExec}, priv2, addr2)
	require.NoError(t, err)
	require.NoError(t, simErr)
	require.NoError(t, deliverErr)

	balanceAfter1 := app.BankKeeper.GetAllBalances(app.Context, addr1)
	balanceAfter2 := app.BankKeeper.GetAllBalances(app.Context, addr2)

	// addr1 grant addr2 to update iscn record
	// addr2 grant addr1 to use addr2's fee allowance for update
	// so addr2 should pay for the fee
	require.True(t, balanceBefore1.Sub(balanceAfter1...).IsZero(), "Address 1 should not pay any fee: balanceBefore1=%v, balanceAfter1=%v, diff=%v", balanceBefore1, balanceAfter1, balanceBefore1.Sub(balanceAfter1...))
	require.False(t, balanceBefore2.Sub(balanceAfter2...).IsZero(), "Address 2 should pay the ISCN fee: balanceBefore2=%v, balanceAfter2=%v, diff=%v", balanceBefore2, balanceAfter2, balanceBefore2.Sub(balanceAfter2...))
}
