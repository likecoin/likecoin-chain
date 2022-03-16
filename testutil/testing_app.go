package testutil

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

	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	likeapp "github.com/likecoin/likechain/app"

	"github.com/likecoin/likechain/x/iscn/types"
)

const DefaultNodeHome = "/tmp/.liked-test"
const invCheckPeriod = 1

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

func SetupTestAppWithIscnGenesis(genesisBalances []GenesisBalance, iscnGenesisJson json.RawMessage) *TestingApp {
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

	if iscnGenesisJson != nil {
		genesisState[types.ModuleName] = iscnGenesisJson
	}

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

func SetupTestApp(genesisBalances []GenesisBalance) *TestingApp {
	return SetupTestAppWithIscnGenesis(genesisBalances, nil)
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
	addr := sdk.AccAddress(priv.PubKey().Address())
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
