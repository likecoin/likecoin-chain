package iscn_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/likecoin/likecoin-chain/v3/x/iscn/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/iscn/types"

	testutil "github.com/likecoin/likecoin-chain/v3/testutil"
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

func TestBasicCreateAndUpdateAndChangeOwnership(t *testing.T) {
	var msg sdk.Msg
	app := testutil.SetupTestApp([]testutil.GenesisBalance{{addr1.String(), "1000000000000000000nanolike"}})

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

	events := result.GetEvents()
	iscnIdStrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnIdStrBytes)
	iscnId, err := types.ParseIscnId(string(iscnIdStrBytes))
	require.NoError(t, err)
	ipldStrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("ipld"))
	require.NotNil(t, ipldStrBytes)
	ownerStrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("owner"))
	require.NotNil(t, ownerStrBytes)
	require.Equal(t, string(ownerStrBytes), addr1.String())

	ctx := app.SetForQuery()

	idQuery := types.NewQueryRecordsByIdRequest(iscnId, 0, 0)
	idQueryRes, err := app.IscnKeeper.RecordsById(sdk.WrapSDKContext(ctx), idQuery)
	require.NoError(t, err)
	require.Equal(t, uint64(1), idQueryRes.LatestVersion)
	require.Equal(t, addr1.String(), idQueryRes.Owner)
	require.Len(t, idQueryRes.Records, 1)
	queryRecord := idQueryRes.Records[0]
	require.Equal(t, string(ipldStrBytes), queryRecord.Ipld)
	v, ok := queryRecord.Data.GetPath("@id")
	require.True(t, ok)
	require.Equal(t, iscnId.String(), v)
	v, ok = queryRecord.Data.GetPath("@type")
	require.True(t, ok)
	require.Equal(t, "Record", v)
	notes, ok := queryRecord.Data.GetPath("recordNotes")
	require.True(t, ok)
	require.Equal(t, "some update", notes)
	timestamp, ok := queryRecord.Data.GetPath("recordTimestamp")
	require.True(t, ok)
	require.Equal(t, "2009-02-13T23:31:30+00:00", timestamp)
	recordFingerprints, ok := queryRecord.Data.GetPath("contentFingerprints")
	require.True(t, ok)
	require.Len(t, recordFingerprints, 1)
	recordFingerprint1, ok := queryRecord.Data.GetPath("contentFingerprints", 0)
	require.True(t, ok)
	require.Equal(t, fingerprint1, recordFingerprint1)
	_, ok = queryRecord.Data.GetPath("recordParentIPLD")
	require.False(t, ok)
	recordStakeholders, ok := queryRecord.Data.GetPath("stakeholders")
	require.True(t, ok)
	require.Len(t, recordStakeholders, 2)
	recordStakeholder1Obj, ok := queryRecord.Data.GetPath("stakeholders", 0)
	require.True(t, ok)
	recordStakeholder1Json, err := json.Marshal(recordStakeholder1Obj)
	require.NoError(t, err)
	require.Equal(t, sdk.MustSortJSON(stakeholder1), recordStakeholder1Json)
	recordStakeholder2Obj, ok := queryRecord.Data.GetPath("stakeholders", 1)
	require.True(t, ok)
	recordStakeholder2Json, err := json.Marshal(recordStakeholder2Obj)
	require.NoError(t, err)
	require.Equal(t, sdk.MustSortJSON(stakeholder2), recordStakeholder2Json)
	recordContentMetadataObj, ok := queryRecord.Data.GetPath("contentMetadata")
	require.True(t, ok)
	recordContentMetadataJson, err := json.Marshal(recordContentMetadataObj)
	require.NoError(t, err)
	require.Equal(t, sdk.MustSortJSON(contentMetadata1), recordContentMetadataJson)

	idQuery = types.NewQueryRecordsByIdRequest(iscnId.PrefixId(), 1, 1)
	idQueryRes, err = app.IscnKeeper.RecordsById(sdk.WrapSDKContext(ctx), idQuery)
	require.NoError(t, err)
	require.Equal(t, uint64(1), idQueryRes.LatestVersion)
	require.Equal(t, addr1.String(), idQueryRes.Owner)
	require.Len(t, idQueryRes.Records, 1)
	require.Equal(t, queryRecord, idQueryRes.Records[0])

	fpQuery := types.NewQueryRecordsByFingerprintRequest(fingerprint1, 0)
	fpQueryRes, err := app.IscnKeeper.RecordsByFingerprint(sdk.WrapSDKContext(ctx), fpQuery)
	require.NoError(t, err)
	require.Len(t, fpQueryRes.Records, 1)
	require.Equal(t, queryRecord, fpQueryRes.Records[0])

	ownerQuery := types.NewQueryRecordsByOwnerRequest(addr1, 0)
	ownerQueryRes, err := app.IscnKeeper.RecordsByOwner(sdk.WrapSDKContext(ctx), ownerQuery)
	require.NoError(t, err)
	require.Len(t, ownerQueryRes.Records, 1)
	require.Equal(t, queryRecord, ownerQueryRes.Records[0])

	app.NextHeader(1234567891)
	app.SetForTx()

	record = types.IscnRecord{
		RecordNotes:         "new update",
		ContentFingerprints: []string{fingerprint1, fingerprint2},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata2,
	}
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId, &record)
	result = app.DeliverMsgNoError(t, msg, priv1)

	events = result.GetEvents()
	iscnId2StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnId2StrBytes)
	iscnId2, err := types.ParseIscnId(string(iscnId2StrBytes))
	require.NoError(t, err)
	require.Equal(t, iscnId.Prefix, iscnId2.Prefix)
	require.Equal(t, iscnId.Version+1, iscnId2.Version)
	ipld2StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("ipld"))
	require.NotNil(t, ipld2StrBytes)
	owner2StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("owner"))
	require.NotNil(t, owner2StrBytes)
	require.Equal(t, string(owner2StrBytes), addr1.String())

	ctx = app.SetForQuery()

	idQuery = types.NewQueryRecordsByIdRequest(iscnId, 0, 0)
	idQueryRes, err = app.IscnKeeper.RecordsById(sdk.WrapSDKContext(ctx), idQuery)
	require.NoError(t, err)
	require.Equal(t, uint64(2), idQueryRes.LatestVersion)
	require.Equal(t, addr1.String(), idQueryRes.Owner)
	require.Len(t, idQueryRes.Records, 1)
	require.Equal(t, queryRecord, idQueryRes.Records[0])

	idQuery = types.NewQueryRecordsByIdRequest(iscnId.PrefixId(), 1, 0)
	idQueryRes, err = app.IscnKeeper.RecordsById(sdk.WrapSDKContext(ctx), idQuery)
	require.NoError(t, err)
	require.Equal(t, uint64(2), idQueryRes.LatestVersion)
	require.Equal(t, addr1.String(), idQueryRes.Owner)
	require.Len(t, idQueryRes.Records, 2)
	require.Equal(t, queryRecord, idQueryRes.Records[0])
	queryRecord2 := idQueryRes.Records[1]

	require.Equal(t, string(ipld2StrBytes), queryRecord2.Ipld)
	v, ok = queryRecord2.Data.GetPath("@id")
	require.True(t, ok)
	require.Equal(t, iscnId2.String(), v)
	v, ok = queryRecord2.Data.GetPath("@type")
	require.True(t, ok)
	require.Equal(t, "Record", v)
	notes, ok = queryRecord2.Data.GetPath("recordNotes")
	require.True(t, ok)
	require.Equal(t, "new update", notes)
	timestamp, ok = queryRecord2.Data.GetPath("recordTimestamp")
	require.True(t, ok)
	require.Equal(t, "2009-02-13T23:31:31+00:00", timestamp)
	recordFingerprints, ok = queryRecord2.Data.GetPath("contentFingerprints")
	require.True(t, ok)
	require.Len(t, recordFingerprints, 2)
	recordFingerprint1, ok = queryRecord2.Data.GetPath("contentFingerprints", 0)
	require.True(t, ok)
	require.Equal(t, fingerprint1, recordFingerprint1)
	recordFingerprint2, ok := queryRecord2.Data.GetPath("contentFingerprints", 1)
	require.True(t, ok)
	require.Equal(t, fingerprint2, recordFingerprint2)
	recordParentIpld, ok := queryRecord2.Data.GetPath("recordParentIPLD", "/")
	require.True(t, ok)
	require.Equal(t, string(ipldStrBytes), recordParentIpld)
	recordStakeholders, ok = queryRecord.Data.GetPath("stakeholders")
	require.True(t, ok)
	require.Len(t, recordStakeholders, 2)
	recordStakeholder1Obj, ok = queryRecord.Data.GetPath("stakeholders", 0)
	require.True(t, ok)
	recordStakeholder1Json, err = json.Marshal(recordStakeholder1Obj)
	require.NoError(t, err)
	require.Equal(t, sdk.MustSortJSON(stakeholder1), recordStakeholder1Json)
	recordStakeholder2Obj, ok = queryRecord.Data.GetPath("stakeholders", 1)
	require.True(t, ok)
	recordStakeholder2Json, err = json.Marshal(recordStakeholder2Obj)
	require.NoError(t, err)
	require.Equal(t, sdk.MustSortJSON(stakeholder2), recordStakeholder2Json)
	recordContentMetadataObj, ok = queryRecord2.Data.GetPath("contentMetadata")
	require.True(t, ok)
	recordContentMetadataJson, err = json.Marshal(recordContentMetadataObj)
	require.NoError(t, err)
	require.Equal(t, sdk.MustSortJSON(contentMetadata2), recordContentMetadataJson)

	fpQuery = types.NewQueryRecordsByFingerprintRequest(fingerprint1, 0)
	fpQueryRes, err = app.IscnKeeper.RecordsByFingerprint(sdk.WrapSDKContext(ctx), fpQuery)
	require.NoError(t, err)
	require.Len(t, fpQueryRes.Records, 2)
	require.Equal(t, queryRecord, fpQueryRes.Records[0])
	require.Equal(t, queryRecord2, fpQueryRes.Records[1])

	fpQuery = types.NewQueryRecordsByFingerprintRequest(fingerprint2, 0)
	fpQueryRes, err = app.IscnKeeper.RecordsByFingerprint(sdk.WrapSDKContext(ctx), fpQuery)
	require.NoError(t, err)
	require.Len(t, fpQueryRes.Records, 1)
	require.Equal(t, queryRecord2, fpQueryRes.Records[0])

	ownerQuery = types.NewQueryRecordsByOwnerRequest(addr1, 0)
	ownerQueryRes, err = app.IscnKeeper.RecordsByOwner(sdk.WrapSDKContext(ctx), ownerQuery)
	require.NoError(t, err)
	require.Len(t, ownerQueryRes.Records, 2)
	require.Equal(t, queryRecord, ownerQueryRes.Records[0])
	require.Equal(t, queryRecord2, ownerQueryRes.Records[1])

	app.SetForTx()

	msg = types.NewMsgChangeIscnRecordOwnership(addr1, iscnId2, addr2)
	app.DeliverMsgNoError(t, msg, priv1)

	ctx = app.SetForQuery()

	idQuery = types.NewQueryRecordsByIdRequest(iscnId.PrefixId(), 1, 0)
	idQueryRes, err = app.IscnKeeper.RecordsById(sdk.WrapSDKContext(ctx), idQuery)
	require.NoError(t, err)
	require.Equal(t, uint64(2), idQueryRes.LatestVersion)
	require.Equal(t, addr2.String(), idQueryRes.Owner)
	require.Len(t, idQueryRes.Records, 2)
	require.Equal(t, queryRecord, idQueryRes.Records[0])
	require.Equal(t, queryRecord2, idQueryRes.Records[1])

	ownerQuery = types.NewQueryRecordsByOwnerRequest(addr1, 0)
	ownerQueryRes, err = app.IscnKeeper.RecordsByOwner(sdk.WrapSDKContext(ctx), ownerQuery)
	require.NoError(t, err)
	require.Len(t, ownerQueryRes.Records, 0)

	ownerQuery = types.NewQueryRecordsByOwnerRequest(addr2, 0)
	ownerQueryRes, err = app.IscnKeeper.RecordsByOwner(sdk.WrapSDKContext(ctx), ownerQuery)
	require.NoError(t, err)
	require.Len(t, ownerQueryRes.Records, 2)
	require.Equal(t, queryRecord, ownerQueryRes.Records[0])
	require.Equal(t, queryRecord2, ownerQueryRes.Records[1])

	app.SetForTx()

	msg = crisistypes.NewMsgVerifyInvariant(addr1, "iscn", "iscn-records")
	app.DeliverMsgNoError(t, msg, priv1)
}

func TestMultipleCreateInOneTx(t *testing.T) {
	app := testutil.SetupTestApp([]testutil.GenesisBalance{{addr1.String(), "1000000000000000000nanolike"}})

	app.NextHeader(1234567890)
	app.SetForTx()
	record1 := types.IscnRecord{
		RecordNotes:         "some update",
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	record2 := types.IscnRecord{
		RecordNotes:         "another update",
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	msgs := []sdk.Msg{
		types.NewMsgCreateIscnRecord(addr1, &record1, 0),
		types.NewMsgCreateIscnRecord(addr1, &record2, 0),
	}
	app.DeliverMsgsNoError(t, msgs, priv1)
}

func TestOwnerQueryPagination(t *testing.T) {
	var msg sdk.Msg
	app := testutil.SetupTestApp([]testutil.GenesisBalance{{addr1.String(), "1000000000000000000nanolike"}})

	app.NextHeader(1234567890)
	app.SetForTx()

	// OwnerRecordsPageLimit-1 records for prefix 1, 2 records for prefix 2, 2 records for prefix 3
	record := types.IscnRecord{
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
		RecordNotes:         fmt.Sprintf("update record 1 to version %010d", 1),
	}
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	result := app.DeliverMsgNoError(t, msg, priv1)
	events := result.GetEvents()
	iscnId1StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnId1StrBytes)
	iscnId1, err := types.ParseIscnId(string(iscnId1StrBytes))
	require.NoError(t, err)
	for i := 1; i < keeper.OwnerRecordsPageLimit-1; i++ {
		iscnId1.Version = uint64(i)
		record.RecordNotes = fmt.Sprintf("update record 1 to version %010d", i+1)
		msg = types.NewMsgUpdateIscnRecord(addr1, iscnId1, &record)
		app.DeliverMsgNoError(t, msg, priv1)
	}

	record.RecordNotes = fmt.Sprintf("update record 2 to version %010d", 1)
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	result = app.DeliverMsgNoError(t, msg, priv1)
	events = result.GetEvents()
	iscnId2StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnId2StrBytes)
	iscnId2, err := types.ParseIscnId(string(iscnId2StrBytes))
	require.NoError(t, err)
	record.RecordNotes = fmt.Sprintf("update record 2 to version %010d", 2)
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId2, &record)
	app.DeliverMsgNoError(t, msg, priv1)

	record.RecordNotes = fmt.Sprintf("update record 3 to version %010d", 1)
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	result = app.DeliverMsgNoError(t, msg, priv1)
	events = result.GetEvents()
	iscnId3StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnId3StrBytes)
	iscnId3, err := types.ParseIscnId(string(iscnId3StrBytes))
	require.NoError(t, err)
	record.RecordNotes = fmt.Sprintf("update record 3 to version %010d", 2)
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId3, &record)
	app.DeliverMsgNoError(t, msg, priv1)

	ctx := app.SetForQuery()

	ownerQuery := types.NewQueryRecordsByOwnerRequest(addr1, 0)
	ownerQueryRes, err := app.IscnKeeper.RecordsByOwner(sdk.WrapSDKContext(ctx), ownerQuery)
	require.NoError(t, err)
	require.Len(t, ownerQueryRes.Records, keeper.OwnerRecordsPageLimit-1)
	require.NotZero(t, ownerQueryRes.NextSequence)
	for i, queryRecord := range ownerQueryRes.Records {
		notes, ok := queryRecord.Data.GetPath("recordNotes")
		require.True(t, ok)
		require.Equal(t, fmt.Sprintf("update record 1 to version %010d", i+1), notes)
	}

	ownerQuery = types.NewQueryRecordsByOwnerRequest(addr1, ownerQueryRes.NextSequence)
	ownerQueryRes, err = app.IscnKeeper.RecordsByOwner(sdk.WrapSDKContext(ctx), ownerQuery)
	require.NoError(t, err)
	require.Len(t, ownerQueryRes.Records, 4)
	require.Zero(t, ownerQueryRes.NextSequence)

	queryRecord := ownerQueryRes.Records[0]
	notes, ok := queryRecord.Data.GetPath("recordNotes")
	require.True(t, ok)
	require.Equal(t, fmt.Sprintf("update record 2 to version %010d", 1), notes)
	queryRecord = ownerQueryRes.Records[1]
	notes, ok = queryRecord.Data.GetPath("recordNotes")
	require.True(t, ok)
	require.Equal(t, fmt.Sprintf("update record 2 to version %010d", 2), notes)
	queryRecord = ownerQueryRes.Records[2]
	notes, ok = queryRecord.Data.GetPath("recordNotes")
	require.True(t, ok)
	require.Equal(t, fmt.Sprintf("update record 3 to version %010d", 1), notes)
	queryRecord = ownerQueryRes.Records[3]
	notes, ok = queryRecord.Data.GetPath("recordNotes")
	require.True(t, ok)
	require.Equal(t, fmt.Sprintf("update record 3 to version %010d", 2), notes)
}

func TestFingerprintQueryPagination(t *testing.T) {
	var msg sdk.Msg
	app := testutil.SetupTestApp([]testutil.GenesisBalance{{addr1.String(), "1000000000000000000nanolike"}})

	app.NextHeader(1234567890)
	app.SetForTx()

	record := types.IscnRecord{
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	for i := 0; i < 2*keeper.FingerprintRecordsPageLimit-1; i++ {
		record.RecordNotes = fmt.Sprintf("record %010d", i)
		msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
		app.DeliverMsgNoError(t, msg, priv1)
	}

	ctx := app.SetForQuery()

	fpQuery := types.NewQueryRecordsByFingerprintRequest(fingerprint1, 0)
	fpQueryRes, err := app.IscnKeeper.RecordsByFingerprint(sdk.WrapSDKContext(ctx), fpQuery)
	require.NoError(t, err)
	require.Len(t, fpQueryRes.Records, keeper.FingerprintRecordsPageLimit)
	require.NotZero(t, fpQueryRes.NextSequence)
	for i, queryRecord := range fpQueryRes.Records {
		notes, ok := queryRecord.Data.GetPath("recordNotes")
		require.True(t, ok)
		require.Equal(t, fmt.Sprintf("record %010d", i), notes)
	}

	fpQuery = types.NewQueryRecordsByFingerprintRequest(fingerprint1, fpQueryRes.NextSequence)
	fpQueryRes, err = app.IscnKeeper.RecordsByFingerprint(sdk.WrapSDKContext(ctx), fpQuery)
	require.NoError(t, err)
	require.Len(t, fpQueryRes.Records, keeper.FingerprintRecordsPageLimit-1)
	require.Zero(t, fpQueryRes.NextSequence)
	for i, queryRecord := range fpQueryRes.Records {
		notes, ok := queryRecord.Data.GetPath("recordNotes")
		require.True(t, ok)
		require.Equal(t, fmt.Sprintf("record %010d", i+keeper.FingerprintRecordsPageLimit), notes)
	}
}

func TestFailureCases(t *testing.T) {
	var msg sdk.Msg
	var record types.IscnRecord
	app := testutil.SetupTestApp([]testutil.GenesisBalance{
		{addr1.String(), "1000000000000000000nanolike"},
		{addr2.String(), "1nanolike"},
		{addr3.String(), "1000000000000000000nanolike"},
	})

	goodRecord := func() types.IscnRecord {
		return types.IscnRecord{
			ContentFingerprints: []string{fingerprint1},
			Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
			ContentMetadata:     contentMetadata1,
		}
	}

	app.NextHeader(1234567890)
	app.SetForTx()

	// Test for MsgCreateIscnRecord

	// ensure everything works fine when no modification is made
	record = goodRecord()
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	res := app.DeliverMsgNoError(t, msg, priv1)
	iscnId, err := types.ParseIscnId(string(testutil.GetEventAttribute(res.GetEvents(), "iscn_record", []byte("iscn_id"))))
	require.NoError(t, err)

	// wrong sender address checksum
	record = goodRecord()
	msg = &types.MsgCreateIscnRecord{"cosmos1ww3qews2y5jxe8apw2zt8stqqrcu2tptejfwag", record, 0}
	_, err, simErr, _ := app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrInvalidAddress))

	// wrong sender address prefix
	record = goodRecord()
	msg = &types.MsgCreateIscnRecord{"iaa1nr4zjtg87mtgvf2zetvmny8htuxplsduyc0h9f", record, 0}
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrInvalidAddress))

	// invalid fingerprint
	record = goodRecord()
	record.ContentFingerprints[0] = ""
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnRecord))

	// invalid stakeholder
	record = goodRecord()
	record.Stakeholders[0] = types.IscnInput(``)
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnRecord))

	// invalid content metadata
	record = goodRecord()
	record.ContentMetadata = types.IscnInput(``)
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnRecord))

	// balance not enough for ISCN fee
	record = goodRecord()
	msg = types.NewMsgCreateIscnRecord(addr2, &record, 0)
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrDeductIscnFee))

	// Test for MsgUpdateIscnRecord

	// ensure everything works fine when no modification is made
	record = goodRecord()
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId, &record)
	res = app.DeliverMsgNoError(t, msg, priv1)
	iscnId2, err := types.ParseIscnId(string(testutil.GetEventAttribute(res.GetEvents(), "iscn_record", []byte("iscn_id"))))
	require.NoError(t, err)

	// wrong sender address checksum
	record = goodRecord()
	msg = &types.MsgUpdateIscnRecord{"cosmos1ww3qews2y5jxe8apw2zt8stqqrcu2tptejfwag", iscnId2.String(), record}
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrInvalidAddress))

	// wrong sender address prefix
	record = goodRecord()
	msg = &types.MsgUpdateIscnRecord{"iaa1nr4zjtg87mtgvf2zetvmny8htuxplsduyc0h9f", iscnId2.String(), record}
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrInvalidAddress))

	// invalid ISCN ID format
	record = goodRecord()
	msg = &types.MsgUpdateIscnRecord{addr1.String(), iscnId2.String()[1:], record}
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnId))

	// not owner
	record = goodRecord()
	msg = types.NewMsgUpdateIscnRecord(addr2, iscnId2, &record)
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrUnauthorized))

	// invalid version
	record = goodRecord()
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId, &record)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnVersion))

	// non-existing ISCN ID
	invalidIscnId, err := types.ParseIscnId("iscn://a/b/1")
	require.NoError(t, err)
	record = goodRecord()
	msg = types.NewMsgUpdateIscnRecord(addr1, invalidIscnId, &record)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrRecordNotFound))

	// existing ISCN ID prefix with future version
	iscnId3 := iscnId
	iscnId3.Version = 3
	record = goodRecord()
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId3, &record)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrRecordNotFound))

	// invalid fingerprint
	record = goodRecord()
	record.ContentFingerprints[0] = ""
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId2, &record)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnRecord))

	// invalid stakeholder
	record = goodRecord()
	record.Stakeholders[0] = types.IscnInput(``)
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId2, &record)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnRecord))

	// invalid content metadata
	record = goodRecord()
	record.ContentMetadata = types.IscnInput(``)
	msg = types.NewMsgUpdateIscnRecord(addr1, iscnId2, &record)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnRecord))

	// balance not enough for ISCN fee
	// also test for success case for MsgChangeIscnRecordOwnership
	msg = types.NewMsgChangeIscnRecordOwnership(addr1, iscnId2, addr2)
	app.DeliverMsgNoError(t, msg, priv1)
	record = goodRecord()
	msg = types.NewMsgUpdateIscnRecord(addr2, iscnId2, &record)
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrDeductIscnFee))

	// Test for MsgChangeIscnRecordOwnership

	// wrong sender address checksum
	msg = &types.MsgChangeIscnRecordOwnership{"cosmos1ww3qews2y5jxe8apw2zt8stqqrcu2tptejfwag", iscnId2.String(), addr3.String()}
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrInvalidAddress))

	// wrong sender address prefix
	msg = &types.MsgChangeIscnRecordOwnership{"iaa1nr4zjtg87mtgvf2zetvmny8htuxplsduyc0h9f", iscnId2.String(), addr3.String()}
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrInvalidAddress))

	// wrong new owner address checksum
	msg = &types.MsgChangeIscnRecordOwnership{addr2.String(), iscnId2.String(), "cosmos1ww3qews2y5jxe8apw2zt8stqqrcu2tptejfwag"}
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrInvalidAddress))

	// wrong sender address prefix
	msg = &types.MsgChangeIscnRecordOwnership{addr2.String(), iscnId2.String(), "iaa1nr4zjtg87mtgvf2zetvmny8htuxplsduyc0h9f"}
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrInvalidAddress))

	// non-owner
	msg = types.NewMsgChangeIscnRecordOwnership(addr1, iscnId2, addr3)
	_, err, simErr, _ = app.DeliverMsg(msg, priv1)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, sdkerrors.ErrUnauthorized))

	// non-existing ISCN ID
	msg = types.NewMsgChangeIscnRecordOwnership(addr2, invalidIscnId, addr3)
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrRecordNotFound))

	// previous version
	msg = types.NewMsgChangeIscnRecordOwnership(addr2, iscnId, addr3)
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnVersion))

	// future version
	msg = types.NewMsgChangeIscnRecordOwnership(addr2, iscnId3, addr3)
	_, err, simErr, _ = app.DeliverMsg(msg, priv2)
	require.NoError(t, err)
	require.Error(t, simErr)
	require.True(t, errors.Is(simErr, types.ErrInvalidIscnVersion))
}

func TestSimulation(t *testing.T) {
	const seedCount = 10
	const txCount = 100

	goodRecord := func() types.IscnRecord {
		return types.IscnRecord{
			ContentFingerprints: []string{fingerprint1},
			Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
			ContentMetadata:     contentMetadata1,
		}
	}

	testWithRand := func(r *rand.Rand) {
		prefixArr := []string{}
		contentIdRecordMap := map[string]types.ContentIdRecord{}
		notesMap := map[string]string{}
		keys := []struct {
			PrivKey cryptotypes.PrivKey
			Address sdk.AccAddress
		}{
			{priv1, addr1},
			{priv2, addr2},
			{priv3, addr3},
		}
		addrToPrivKey := map[string]cryptotypes.PrivKey{
			addr1.String(): priv1,
			addr2.String(): priv2,
			addr3.String(): priv3,
		}

		doRandomTx := func(r *rand.Rand, app *testutil.TestingApp) {
			x := r.Intn(100)
			if x < 50 || len(contentIdRecordMap) == 0 {
				key := keys[r.Intn(len(keys))]
				privKey := key.PrivKey
				addr := key.Address
				record := goodRecord()
				notes := fmt.Sprintf("create notes %d", r.Int63())
				record.RecordNotes = notes
				msg := types.NewMsgCreateIscnRecord(addr, &record, 0)
				res := app.DeliverMsgNoError(t, msg, privKey)
				iscnId, err := types.ParseIscnId(string(testutil.GetEventAttribute(res.GetEvents(), "iscn_record", []byte("iscn_id"))))
				require.NoError(t, err)
				prefix := iscnId.Prefix.String()
				prefixArr = append(prefixArr, prefix)
				contentIdRecordMap[prefix] = types.ContentIdRecord{
					OwnerAddressBytes: addr.Bytes(),
					LatestVersion:     1,
				}
				notesMap[iscnId.String()] = notes
			} else {
				prefix := prefixArr[r.Intn(len(contentIdRecordMap))]
				iscnId, err := types.ParseIscnId(prefix)
				require.NoError(t, err)
				contentIdRecord := contentIdRecordMap[prefix]
				owner := contentIdRecord.OwnerAddress()
				privKey := addrToPrivKey[owner.String()]
				iscnId.Version = contentIdRecord.LatestVersion
				if x < 80 {
					record := goodRecord()
					notes := fmt.Sprintf("update notes %d", r.Int63())
					record.RecordNotes = notes
					msg := types.NewMsgUpdateIscnRecord(owner, iscnId, &record)
					app.DeliverMsgNoError(t, msg, privKey)
					contentIdRecord.LatestVersion++
					contentIdRecordMap[prefix] = contentIdRecord
					iscnId.Version++
					notesMap[iscnId.String()] = notes
				} else {
					newOwner := keys[r.Intn(len(keys))].Address
					msg := types.NewMsgChangeIscnRecordOwnership(owner, iscnId, newOwner)
					app.DeliverMsgNoError(t, msg, privKey)
					contentIdRecord.OwnerAddressBytes = newOwner.Bytes()
					contentIdRecordMap[prefix] = contentIdRecord
				}
			}
		}

		verifyState := func(app *testutil.TestingApp) {
			ctx := app.Context
			for prefix, contentIdRecord := range contentIdRecordMap {
				prefixIscnId, err := types.ParseIscnId(prefix)
				require.NoError(t, err)
				query := types.NewQueryRecordsByIdRequest(prefixIscnId.PrefixId(), 1, 0)
				res, err := app.IscnKeeper.RecordsById(sdk.WrapSDKContext(ctx), query)
				require.NoError(t, err)
				require.Equal(t, contentIdRecord.LatestVersion, res.LatestVersion)
				require.Equal(t, contentIdRecord.OwnerAddress().String(), res.Owner)
				require.Len(t, res.Records, int(contentIdRecord.LatestVersion))
				for i, record := range res.Records {
					iscnIdAny, ok := record.Data.GetPath("@id")
					require.True(t, ok)
					iscnIdStr, ok := iscnIdAny.(string)
					require.True(t, ok)
					iscnId, err := types.ParseIscnId(string(iscnIdStr))
					require.NoError(t, err)
					require.Equal(t, uint64(i+1), iscnId.Version)
					notes, ok := record.Data.GetPath("recordNotes")
					require.True(t, ok)
					require.Equal(t, notesMap[iscnId.String()], notes)
				}
			}
		}

		genesisBalances := []testutil.GenesisBalance{
			{addr1.String(), "1000000000000000000nanolike"},
			{addr2.String(), "1000000000000000000nanolike"},
			{addr3.String(), "1000000000000000000nanolike"},
		}
		app := testutil.SetupTestApp(genesisBalances)
		for i := 0; i < txCount; i++ {
			doRandomTx(r, app)
		}
		ctx := app.SetForQuery()
		verifyState(app)

		iscnGenesis := app.IscnKeeper.ExportGenesis(ctx)
		iscnGenesisJson := app.AppCodec().MustMarshalJSON(iscnGenesis)
		app = testutil.SetupTestAppWithIscnGenesis(genesisBalances, iscnGenesisJson)
		app.SetForQuery()
		verifyState(app)

		app.SetForTx()
		for i := 0; i < txCount; i++ {
			doRandomTx(r, app)
		}

		app.SetForQuery()
		verifyState(app)
	}

	for seed := int64(0); seed < seedCount; seed++ {
		r := rand.New(rand.NewSource(seed))
		testWithRand(r)
	}
}

func TestUpdateAuthorization(t *testing.T) {
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
	msg, err := authz.NewMsgGrant(addr1, addr2, types.NewUpdateAuthorization(iscnId.Prefix.String()), &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msg, priv1)

	updateMsg := types.NewMsgUpdateIscnRecord(addr1, iscnId, &record)
	msgExec := authz.NewMsgExec(addr2, []sdk.Msg{updateMsg})
	msg = &msgExec
	result = app.DeliverMsgNoError(t, msg, priv2)

	events := result.GetEvents()
	iscnId2StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnId2StrBytes)
	iscnId2, err := types.ParseIscnId(string(iscnId2StrBytes))
	require.NoError(t, err)
	require.Equal(t, iscnId.Prefix, iscnId2.Prefix)
	require.Equal(t, iscnId.Version+1, iscnId2.Version)
	ipld2StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("ipld"))
	require.NotNil(t, ipld2StrBytes)
	owner2StrBytes := testutil.GetEventAttribute(events, "iscn_record", []byte("owner"))
	require.NotNil(t, owner2StrBytes)
	require.Equal(t, string(owner2StrBytes), addr1.String())

	updateMsg = types.NewMsgUpdateIscnRecord(addr1, iscnId2, &record)
	msgExec = authz.NewMsgExec(addr3, []sdk.Msg{updateMsg})
	msg = &msgExec
	_, _, simErr, _ := app.DeliverMsg(msg, priv3)
	require.ErrorContains(t, simErr, "authorization not found")

	record = types.IscnRecord{
		RecordNotes:         "another record",
		ContentFingerprints: []string{fingerprint1, fingerprint2},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata2,
	}
	msg = types.NewMsgCreateIscnRecord(addr1, &record, 0)
	result = app.DeliverMsgNoError(t, msg, priv1)
	iscnId3 := testutil.GetIscnIdFromResult(t, result)

	record = types.IscnRecord{
		RecordNotes:         "another record updated",
		ContentFingerprints: []string{fingerprint1, fingerprint2},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata2,
	}

	updateMsg = types.NewMsgUpdateIscnRecord(addr1, iscnId3, &record)
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{updateMsg})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv2)
	require.ErrorContains(t, simErr, "ISCN ID prefix mismatch")
}
