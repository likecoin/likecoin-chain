package iscn_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	likeapp "github.com/likecoin/likechain/app"

	"github.com/likecoin/likechain/x/iscn/keeper"
	"github.com/likecoin/likechain/x/iscn/types"
)

// TODO: seems a useless param just for creating the app, but not sure if there's a better way to handle
const DefaultNodeHome = "/tmp/.liked-test"
const invCheckPeriod = 1

type TestingApp struct {
	*likeapp.LikeApp

	txCfg   client.TxConfig
	Header  tmproto.Header
	Context sdk.Context
}

type genesisBalance struct {
	Address string
	Coin    string
}

func SetupTestApp(genesisBalances []genesisBalance) *TestingApp {
	genAccs := []authtypes.GenesisAccount{}
	balances := []banktypes.Balance{}
	for _, balance := range genesisBalances {
		addr := balance.Address
		genAccs = append(genAccs, &authtypes.BaseAccount{Address: addr})
		coin, err := sdk.ParseCoinNormalized(balance.Coin)
		if err != nil {
			panic(err)
		}
		balance := banktypes.Balance{Address: addr, Coins: sdk.NewCoins(coin)}
		balances = append(balances, balance)
	}
	db := dbm.NewMemDB()
	encodingCfg := likeapp.MakeEncodingConfig()
	logger := log.NewTMLogger(os.Stdout)
	app := likeapp.NewLikeApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, invCheckPeriod, encodingCfg, simapp.EmptyAppOptions{})
	genesisState := likeapp.ModuleBasics.DefaultGenesis(encodingCfg.Marshaler)
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		totalSupply = totalSupply.Add(b.Coins...)
	}

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	crisisGenesis := crisistypes.NewGenesisState(sdk.NewInt64Coin("nanolike", 1))
	genesisState[crisistypes.ModuleName] = app.AppCodec().MustMarshalJSON(crisisGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	app.Commit()

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	return &TestingApp{
		LikeApp: app,
		txCfg:   simapp.MakeTestEncodingConfig().TxConfig,
		Header:  header,
		Context: app.BaseApp.NewContext(false, header),
	}
}

func (app *TestingApp) NextHeader(unixTimestamp int64) {
	app.Header = tmproto.Header{
		Time: time.Unix(unixTimestamp, 0),
	}
}

func (app *TestingApp) SetForQuery() sdk.Context {
	app.Header.Height = app.LastBlockHeight() + 1
	app.BeginBlock(abci.RequestBeginBlock{Header: app.Header})
	app.Context = app.BaseApp.NewContext(false, app.Header)
	return app.Context
}

func (app *TestingApp) SetForTx() {
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()
}

func (app *TestingApp) DeliverMsgs(msgs []sdk.Msg, priv cryptotypes.PrivKey) (res *sdk.Result, err error, simErr error, deliverErr error) {
	app.Header.Height = app.LastBlockHeight() + 1
	app.BeginBlock(abci.RequestBeginBlock{Header: app.Header})
	app.Context = app.BaseApp.NewContext(false, app.Header)
	chainId := ""
	addr := sdk.AccAddress(priv1.PubKey().Address())
	acc := app.AccountKeeper.GetAccount(app.Context, addr)
	accNum := acc.GetAccountNumber()
	accSeq := acc.GetSequence()
	txCfg := app.txCfg
	tx, err := helpers.GenTx(
		app.txCfg,
		msgs,
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
		helpers.DefaultGenTxGas,
		chainId,
		[]uint64{accNum},
		[]uint64{accSeq},
		priv,
	)
	if err != nil {
		return nil, err, nil, nil
	}
	txBytes, err := txCfg.TxEncoder()(tx)
	if err != nil {
		return nil, err, nil, nil
	}
	_, _, simErr = app.Simulate(txBytes)
	if simErr != nil {
		return nil, nil, simErr, nil
	}
	_, res, deliverErr = app.Deliver(txCfg.TxEncoder(), tx)
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()
	return res, nil, nil, deliverErr
}

func (app *TestingApp) DeliverMsg(msg sdk.Msg, priv cryptotypes.PrivKey) (res *sdk.Result, err error, simErr error, deliverErr error) {
	return app.DeliverMsgs([]sdk.Msg{msg}, priv)
}

func (app *TestingApp) DeliverMsgNoError(t *testing.T, msg sdk.Msg, priv cryptotypes.PrivKey) *sdk.Result {
	res, err, simErr, deliverErr := app.DeliverMsgs([]sdk.Msg{msg}, priv)
	require.NoError(t, err)
	require.NoError(t, simErr)
	require.NoError(t, deliverErr)
	return res
}

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())

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

func getEventAttribute(events sdk.Events, typ string, attrKey []byte) []byte {
	for _, e := range events {
		if e.Type == typ {
			for _, attr := range e.Attributes {
				if bytes.Equal(attr.Key, attrKey) {
					return attr.Value
				}
			}
		}
	}
	return nil
}

func TestBasicCreateAndUpdateAndChangeOwnership(t *testing.T) {
	var msg sdk.Msg
	app := SetupTestApp([]genesisBalance{{addr1.String(), "1000000000000000000nanolike"}})

	app.NextHeader(1234567890)
	app.SetForTx()
	record := types.IscnRecord{
		RecordNotes:         "some update",
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	msg = types.NewMsgCreateIscnRecord(addr1, &record)
	result := app.DeliverMsgNoError(t, msg, priv1)

	events := result.GetEvents()
	iscnIdStrBytes := getEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnIdStrBytes)
	iscnId, err := types.ParseIscnId(string(iscnIdStrBytes))
	require.NoError(t, err)
	ipldStrBytes := getEventAttribute(events, "iscn_record", []byte("ipld"))
	require.NotNil(t, ipldStrBytes)
	ownerStrBytes := getEventAttribute(events, "iscn_record", []byte("owner"))
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
	iscnId2StrBytes := getEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnId2StrBytes)
	iscnId2, err := types.ParseIscnId(string(iscnId2StrBytes))
	require.NoError(t, err)
	require.Equal(t, iscnId.Prefix, iscnId2.Prefix)
	require.Equal(t, iscnId.Version+1, iscnId2.Version)
	ipld2StrBytes := getEventAttribute(events, "iscn_record", []byte("ipld"))
	require.NotNil(t, ipld2StrBytes)
	owner2StrBytes := getEventAttribute(events, "iscn_record", []byte("owner"))
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

	app.SetForTx()

	msg = crisistypes.NewMsgVerifyInvariant(addr1, "iscn", "iscn-records")
	app.DeliverMsgNoError(t, msg, priv1)
}

func TestFingerprintQueryPagination(t *testing.T) {
	var msg sdk.Msg
	app := SetupTestApp([]genesisBalance{{addr1.String(), "1000000000000000000nanolike"}})

	app.NextHeader(1234567890)
	app.SetForTx()

	record := types.IscnRecord{
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	for i := 0; i < 2*keeper.FingerprintPageLimit-1; i++ {
		record.RecordNotes = fmt.Sprintf("record %010d", i)
		msg = types.NewMsgCreateIscnRecord(addr1, &record)
		app.DeliverMsgNoError(t, msg, priv1)
	}

	ctx := app.SetForQuery()

	fpQuery := types.NewQueryRecordsByFingerprintRequest(fingerprint1, 0)
	fpQueryRes, err := app.IscnKeeper.RecordsByFingerprint(sdk.WrapSDKContext(ctx), fpQuery)
	require.NoError(t, err)
	require.Len(t, fpQueryRes.Records, keeper.FingerprintPageLimit)
	require.NotZero(t, fpQueryRes.NextSequence)
	for i, queryRecord := range fpQueryRes.Records {
		notes, ok := queryRecord.Data.GetPath("recordNotes")
		require.True(t, ok)
		require.Equal(t, fmt.Sprintf("record %010d", i), notes)
	}

	fpQuery = types.NewQueryRecordsByFingerprintRequest(fingerprint1, fpQueryRes.NextSequence)
	fpQueryRes, err = app.IscnKeeper.RecordsByFingerprint(sdk.WrapSDKContext(ctx), fpQuery)
	require.NoError(t, err)
	require.Len(t, fpQueryRes.Records, keeper.FingerprintPageLimit-1)
	require.Zero(t, fpQueryRes.NextSequence)
	for i, queryRecord := range fpQueryRes.Records {
		notes, ok := queryRecord.Data.GetPath("recordNotes")
		require.True(t, ok)
		require.Equal(t, fmt.Sprintf("record %010d", i+keeper.FingerprintPageLimit), notes)
	}
}

func TestFailureCases(t *testing.T) {
	// TODO
}

func TestUpdateOwnership(t *testing.T) {
	// TODO
}

func TestDeductFund(t *testing.T) {
	// TODO
}

func TestExportAndImport(t *testing.T) {
	// TODO
}
