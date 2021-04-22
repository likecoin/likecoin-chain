package iscn_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	likeapp "github.com/likecoin/likechain/app"

	"github.com/likecoin/likechain/x/iscn/types"
)

// TODO: seems a useless param just for creating the app, but not sure if there's a better way to handle
const DefaultNodeHome = "/tmp/.liked-test"
const invCheckPeriod = 1

func SetupTestApp(genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *likeapp.LikeApp {
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
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1}})

	return app
}

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())

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

func TestMsgCreateIscn(t *testing.T) {
	var header tmproto.Header
	txGen := simapp.MakeTestEncodingConfig().TxConfig

	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}
	genAccs := []authtypes.GenesisAccount{acc}
	balance := banktypes.Balance{Address: acc.Address, Coins: sdk.NewCoins(sdk.NewInt64Coin("nanolike", 1e18))}
	app := SetupTestApp(genAccs, balance)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	authAcc := app.AccountKeeper.GetAccount(ctx, addr1)
	initialSeq := authAcc.GetSequence()

	record := types.IscnRecord{
		RecordNotes:         "some update",
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []types.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	msg1 := types.NewMsgCreateIscnRecord(addr1, &record)
	header = tmproto.Header{
		Time:   time.Unix(1234567890, 0),
		Height: app.LastBlockHeight() + 1,
	}
	accNumber := authAcc.GetAccountNumber()
	_, result, err := simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{msg1}, "", []uint64{accNumber}, []uint64{initialSeq}, true, true, priv1)
	require.NoError(t, err)

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

	query := types.NewQueryRecordsByIdRequest(iscnId, 0, 0)
	queryRes, err := app.IscnKeeper.RecordsById(sdk.WrapSDKContext(ctx), query)
	require.NoError(t, err)
	require.Equal(t, uint64(1), queryRes.LatestVersion)
	require.Equal(t, addr1.String(), queryRes.Owner)
	require.Len(t, queryRes.Records, 1)
	queryRecord := queryRes.Records[0]
	require.Equal(t, string(ipldStrBytes), queryRecord.Ipld)
	v, ok := queryRecord.Data.GetPath("@id")
	require.True(t, ok)
	require.Equal(t, iscnId.String(), v)
	timestamp, ok := queryRecord.Data.GetPath("recordTimestamp")
	require.True(t, ok)
	require.Equal(t, "2009-02-13T23:31:30+00:00", timestamp)
	recordFingerprint1, ok := queryRecord.Data.GetPath("contentFingerprints", 0)
	require.True(t, ok)
	require.Equal(t, fingerprint1, recordFingerprint1)
	_, ok = queryRecord.Data.GetPath("recordParentIPLD")
	require.False(t, ok)
	recordContentMetadataMap, ok := queryRecord.Data.GetPath("contentMetadata")
	require.True(t, ok)
	recordContentMetadataJson, err := json.Marshal(recordContentMetadataMap)
	require.NoError(t, err)
	require.Equal(t, sdk.MustSortJSON(contentMetadata1), recordContentMetadataJson)

	msg2 := crisistypes.NewMsgVerifyInvariant(addr1, "iscn", "iscn-records")
	header = tmproto.Header{
		Time:   header.Time.Add(6 * time.Second),
		Height: app.LastBlockHeight() + 1,
	}
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{msg2}, "", []uint64{accNumber}, []uint64{initialSeq + 1}, true, true, priv1)
	require.NoError(t, err)
}
