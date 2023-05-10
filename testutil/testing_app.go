package testutil

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
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	likeapp "github.com/likecoin/likecoin-chain/v4/app"

	"github.com/likecoin/likecoin-chain/v4/x/iscn/types"
)

const DefaultNodeHome = "/tmp/.liked-test"
const invCheckPeriod = 1

const TestChainId = "test-chain-VhVyAD"
const ValidatorLikeAddr = "like1tf9tg46d82lm32xwq4ms7xj6xse3qu4m5n2h5h"
const ValidatorLikeValAddr = "likevaloper1tf9tg46d82lm32xwq4ms7xj6xse3qu4mzuufyy"
const ValidatorMnemonic = "celery milk ahead display high either family win pool potato plunge crunch siren table become slush bracket dust tumble talent gadget fossil wet lobster"
const GenTx = `{"body":{"messages":[{"@type":"/cosmos.staking.v1beta1.MsgCreateValidator","description":{"moniker":"asdf","identity":"","website":"","security_contact":"","details":""},"commission":{"rate":"0.100000000000000000","max_rate":"0.200000000000000000","max_change_rate":"0.010000000000000000"},"min_self_delegation":"1","delegator_address":"like1tf9tg46d82lm32xwq4ms7xj6xse3qu4m5n2h5h","validator_address":"likevaloper1tf9tg46d82lm32xwq4ms7xj6xse3qu4mzuufyy","pubkey":{"@type":"/cosmos.crypto.ed25519.PubKey","key":"omjQHY80Kp/8VAWXZ5bAf0dffldUETXPENkGo/0jFxQ="},"value":{"denom":"nanolike","amount":"100000000000000"}}],"memo":"d539239b11be3bd8716b0f5b7d80536f3f211118@localhost:26656","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AvsqtnuPtaAIaGbsqf3rs1ELvT+1ztNP5VdUdyZreWeE"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":["lArPx+wmsLj4d+/NEFB/OCaoA2RAYaDjjjQLVkNXift3KCUjzS51pp07sAURENpBK9ruy9I54sOozdswIY27JA=="]}`

type TestingApp struct {
	*likeapp.LikeApp

	txCfg   client.TxConfig
	Header  tmproto.Header
	Context sdk.Context
}

type GenesisBalance struct {
	Address string
	Coin    string
}

func SetupTestAppWithState(appState json.RawMessage, appOptions servertypes.AppOptions) *TestingApp {
	db := dbm.NewMemDB()
	encodingCfg := likeapp.MakeEncodingConfig()
	logger := log.NewTMLogger(os.Stdout)
	app := likeapp.NewLikeApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, invCheckPeriod, encodingCfg, appOptions)

	app.InitChain(
		abci.RequestInitChain{
			ChainId:         TestChainId,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   appState,
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

func SetupTestAppWithIscnGenesis(genesisBalances []GenesisBalance, iscnGenesisJson json.RawMessage) *TestingApp {
	genAccs := []authtypes.GenesisAccount{&authtypes.BaseAccount{Address: ValidatorLikeAddr}}
	balances := []banktypes.Balance{{
		Address: ValidatorLikeAddr,
		Coins:   sdk.NewCoins(sdk.NewCoin("nanolike", sdk.NewInt(1000000000000000000))),
	}}
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

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = "nanolike"
	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, nil, nil)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	crisisGenesis := crisistypes.NewGenesisState(sdk.NewInt64Coin("nanolike", 1))
	genesisState[crisistypes.ModuleName] = app.AppCodec().MustMarshalJSON(crisisGenesis)

	if iscnGenesisJson != nil {
		genesisState[types.ModuleName] = iscnGenesisJson
	}

	genesisState[genutiltypes.ModuleName] = json.RawMessage(fmt.Sprintf(`{"gen_txs":[%s]}`, GenTx))

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	return SetupTestAppWithState(stateBytes, simapp.EmptyAppOptions{})
}

func SetupTestApp(genesisBalances []GenesisBalance) *TestingApp {
	return SetupTestAppWithIscnGenesis(genesisBalances, nil)
}

func SetupTestAppWithDefaultState() *TestingApp {
	return SetupTestApp(nil)
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

func (app *TestingApp) DeliverMsgsWithFeeGranter(msgs []sdk.Msg, priv cryptotypes.PrivKey, feeGranter sdk.AccAddress) (res *sdk.Result, err error, simErr error, deliverErr error) {
	app.Header.Height = app.LastBlockHeight() + 1
	app.BeginBlock(abci.RequestBeginBlock{Header: app.Header})
	app.Context = app.BaseApp.NewContext(false, app.Header)
	chainId := ""
	addr := sdk.AccAddress(priv.PubKey().Address())
	acc := app.AccountKeeper.GetAccount(app.Context, addr)
	accNum := acc.GetAccountNumber()
	accSeq := acc.GetSequence()
	txCfg := app.txCfg
	tx, err := GenerateTx(
		app.txCfg,
		msgs,
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
		helpers.DefaultGenTxGas,
		chainId,
		[]uint64{accNum},
		[]uint64{accSeq},
		feeGranter,
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
	_, res, deliverErr = app.SimDeliver(txCfg.TxEncoder(), tx)
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()
	return res, nil, nil, deliverErr
}

func (app *TestingApp) DeliverMsgs(msgs []sdk.Msg, priv cryptotypes.PrivKey) (res *sdk.Result, err error, simErr error, deliverErr error) {
	return app.DeliverMsgsWithFeeGranter(msgs, priv, nil)
}

func (app *TestingApp) DeliverMsg(msg sdk.Msg, priv cryptotypes.PrivKey) (res *sdk.Result, err error, simErr error, deliverErr error) {
	return app.DeliverMsgs([]sdk.Msg{msg}, priv)
}

func (app *TestingApp) DeliverMsgsNoError(t *testing.T, msgs []sdk.Msg, priv cryptotypes.PrivKey) *sdk.Result {
	res, err, simErr, deliverErr := app.DeliverMsgs(msgs, priv)
	require.NoError(t, err)
	require.NoError(t, simErr)
	require.NoError(t, deliverErr)
	return res
}

func (app *TestingApp) DeliverMsgNoError(t *testing.T, msg sdk.Msg, priv cryptotypes.PrivKey) *sdk.Result {
	return app.DeliverMsgsNoError(t, []sdk.Msg{msg}, priv)
}

func (app *TestingApp) DeliverMsgSimError(t *testing.T, msg sdk.Msg, priv cryptotypes.PrivKey, errContains string, args ...interface{}) {
	_, err, simErr, _ := app.DeliverMsg(msg, priv)
	require.NoError(t, err)
	require.ErrorContains(t, simErr, errContains, args...)
}

func GetEventAttribute(events sdk.Events, typ string, attrKey []byte) []byte {
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

func GetIscnIdFromResult(t *testing.T, result *sdk.Result) types.IscnId {
	events := result.GetEvents()
	iscnIdStrBytes := GetEventAttribute(events, "iscn_record", []byte("iscn_id"))
	require.NotNil(t, iscnIdStrBytes)
	iscnId, err := types.ParseIscnId(string(iscnIdStrBytes))
	require.NoError(t, err)
	return iscnId
}

func GetClassIdFromResult(t *testing.T, result *sdk.Result) string {
	events := result.GetEvents()
	classIdEvent := GetEventAttribute(events, "likechain.likenft.v1.EventNewClass", []byte("class_id"))
	require.NotNil(t, classIdEvent)
	// strip the leading and trailing quotes
	return string(classIdEvent[1 : len(classIdEvent)-1])
}
