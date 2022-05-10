package dual_prefix_tests

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v2/testutil"
	iscntypes "github.com/likecoin/likecoin-chain/v2/x/iscn/types"
	"github.com/stretchr/testify/require"
)

func testCreateAndTransferISCNWithBech32(t *testing.T, from string, to string) {
	app := testutil.SetupTestApp([]testutil.GenesisBalance{
		{
			Address: from,
			Coin:    "1000000000000nanolike",
		},
	})
	app.NextHeader(1234567890)
	app.SetForTx()

	fromAddr, err := sdk.AccAddressFromBech32(from)
	require.NoError(t, err)
	toAddr, err := sdk.AccAddressFromBech32(to)
	require.NoError(t, err)

	fingerprint1 := "hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"
	stakeholder1 := iscntypes.IscnInput(`
{
	"entity": {
		"@id": "did:cosmos:5sy29r37gfxvxz21rh4r0ktpuc46pzjrmz29g45",
		"name": "Chung Wu"
	},
	"rewardProportion": 95,
	"contributionType": "http://schema.org/author"
}`)

	stakeholder2 := iscntypes.IscnInput(`
{
	"rewardProportion": 5,
	"contributionType": "http://schema.org/citation",
	"footprint": "https://en.wikipedia.org/wiki/Fibonacci_number",
	"description": "The blog post referred the matrix form of computing Fibonacci numbers."
}`)

	contentMetadata1 := iscntypes.IscnInput(`
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

	// create record
	record := iscntypes.IscnRecord{
		RecordNotes:         "some update",
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []iscntypes.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	msg := iscntypes.NewMsgCreateIscnRecord(fromAddr, &record)
	result := app.DeliverMsgNoError(t, msg, priv1)
	events := result.GetEvents()

	iscnIdStrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnIdStrBytes)
	iscnId, err := iscntypes.ParseIscnId(string(iscnIdStrBytes))
	require.NoError(t, err)

	// check original ownership
	ctx := app.SetForQuery()
	resultRecord := app.IscnKeeper.GetContentIdRecord(ctx, iscnId.Prefix)
	require.True(t, resultRecord.OwnerAddress().Equals(fromAddr))

	// transfer
	app.NextHeader(1234567891)
	app.SetForTx()
	msg2 := iscntypes.NewMsgChangeIscnRecordOwnership(fromAddr, iscnId, toAddr)
	app.DeliverMsgNoError(t, msg2, priv1)

	// check new ownership
	ctx = app.SetForQuery()
	resultRecord2 := app.IscnKeeper.GetContentIdRecord(ctx, iscnId.Prefix)
	require.True(t, resultRecord2.OwnerAddress().Equals(toAddr))
}

func TestCreateAndTransferISCNFromLegacyPrefixToNew(t *testing.T) {
	testCreateAndTransferISCNWithBech32(t, legacyAddr1, newAddr2)
}

func TestCreateAndTransferISCNFromLegacyPrefixToLegacy(t *testing.T) {
	testCreateAndTransferISCNWithBech32(t, legacyAddr1, legacyAddr2)
}

func TestCreateAndTransferISCNFromNewPrefixToNew(t *testing.T) {
	testCreateAndTransferISCNWithBech32(t, newAddr1, newAddr2)
}

func TestCreateAndTransferISCNFromNewPrefixToLegacy(t *testing.T) {
	testCreateAndTransferISCNWithBech32(t, newAddr1, legacyAddr2)
}
